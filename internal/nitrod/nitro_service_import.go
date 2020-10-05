package nitrod

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/internal/scripts"
)

type DatabaseImportOptions struct {
	Engine          string
	Database        string
	Container       string
	File            string
	Compressed      bool
	CompressionType string
	CreateDatabase  bool
}

func (s *NitroService) ImportDatabase(stream NitroService_ImportDatabaseServer) error {
	options := DatabaseImportOptions{}

	// create a temp file
	file, err := s.createFile(os.TempDir(), "nitro-db-upload-")
	if err != nil {
		s.logger.Println("Error creating a temp file for the upload:", err.Error())
		return status.Errorf(codes.Internal, "Unable creating a temp file for the upload")
	}
	defer file.Close()

	options.File = file.Name()

	// handle the file streaming requests
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "unable to create the stream: %s", err.Error())
		}

		// set the variables for later use only if they are not empty
		if options.Engine == "" {
			options.Engine = req.GetEngine()
		}
		if options.Container == "" {
			options.Container = req.GetContainer()
		}
		if options.Database == "" {
			options.Database = req.GetDatabase()
		}
		if (options.Compressed == false) && (req.GetCompressed()) {
			options.Compressed = req.GetCompressed()
		}
		if options.CompressionType == "" && req.GetCompressionType() != "" {
			options.CompressionType = req.GetCompressionType()
		}
		if (options.CreateDatabase == false) && (req.GetCreateDatabase()) {
			options.CreateDatabase = req.GetCreateDatabase()
		}

		// write the backup content into the temp file
		_, err = file.Write(req.GetData())
		if err != nil {
			return status.Errorf(codes.Internal, "unable to write the backup to the temp file")
		}
	}

	// if the file is compressed, extract it and we are done
	if options.Compressed {
		s.logger.Println("The file is compressed, extracting now")

		uncompressedFile, err := s.createFile(os.TempDir(), "nitro-compressed-db-")
		if err != nil {
			s.logger.Println("error creating the compressed db file:", err.Error())
			return err
		}
		defer uncompressedFile.Close()

		// create the gzip reader
		switch options.CompressionType {
		case "gz":
			s.logger.Println("Using gzip to open to file", file.Name())
			f, err := os.Open(file.Name())
			if err != nil {
				s.logger.Println("error reopening the file for the gzip reader")
				return status.Errorf(codes.Unknown, "error reopening the file for the gzip reader. %s", err.Error())
			}
			reader, err := gzip.NewReader(f)
			if err != nil {
				s.logger.Println("error creating the gzip reader", err.Error())
				return status.Errorf(codes.Unknown, "error reading the compressed file. %s", err.Error())
			}

			for {
				reader.Multistream(false)
				// copy contents to the new file
				if _, err := io.Copy(uncompressedFile, reader); err != nil {
					return err
				}

				err = reader.Reset(f)
				if err == io.EOF {
					break
				}
				if err != nil {
					return err
				}
			}

			if err := reader.Close(); err != nil {
				s.logger.Println("error closing the gzip reader, err:", err.Error())
				return err
			}

			options.File = uncompressedFile.Name()
		default:
			s.logger.Println("Using zip to open to file", file.Name())
			reader, err := zip.OpenReader(file.Name())
			if err != nil {
				return err
			}

			for _, z := range reader.File {

				// only accept sql files, ignoring MACOS specific files
				if strings.Contains(z.Name, ".sql") && strings.Contains(z.Name, "MACOSX") == false {
					b, err := z.Open()
					if err != nil {
						return err
					}

					// copy the contents into the uncompressed temp file
					if _, err := io.Copy(uncompressedFile, b); err != nil {
						return err
					}
				}
			}

			if err := reader.Close(); err != nil {
				s.logger.Println("error closing the zip reader, err:", err.Error())
				return err
			}

			options.File = uncompressedFile.Name()
		}
	}

	// import the database
	if _, err := s.importDatabase(options); err != nil {
		s.logger.Printf("Error importing database: %s\n", err)
		return err
	}

	if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
		return status.Errorf(codes.Internal, "unable to send the response %v", err)
	}

	// remove the temp file to save space
	if err := os.Remove(options.File); err != nil {
		s.logger.Println("error removing temp file:", options.File)
	}
	s.logger.Println("removed temp file:", options.File)

	return nil
}

func (s *NitroService) importDatabase(opts DatabaseImportOptions) (string, error) {
	switch opts.Engine {
	case "mysql":
		// should we skip creating the database?
		if opts.CreateDatabase == false {
			if output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf(scripts.FmtDockerMysqlCreateDatabaseIfNotExists, opts.Container, opts.Database)}); err != nil {
				s.logger.Println(string(output))
				return string(output), err
			}
			s.logger.Printf("Created the MySQL database %q\n", opts.Database)
		}

		s.logger.Printf("Beginning MySQL import of file %q", opts.File)

		// if we are skipping create, it has the use statement and no database name
		if opts.CreateDatabase {
			opts.Database = "emptydatabase"
		}

		// import the database
		output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf(`docker exec -i %q mysql -unitro -pnitro %s < %s`, opts.Container, opts.Database, opts.File)})
		if err != nil {
			s.logger.Println("Error importing the MySQL database:", string(output))
			return "", err
		}
	default:
		output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf(scripts.FmtDockerPostgresCreateDatabase, opts.Container, opts.Database)})
		if err != nil {
			s.logger.Println("Error creating the PostgreSQL database:", string(output))
			return "", err
		}
		s.logger.Printf("created PostgreSQL database %q for engine %q", opts.Database, opts.Container)

		output, err = s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf(scripts.FmtDockerPostgresImportDatabase, opts.Container, opts.Database, opts.File)})
		if err != nil {
			s.logger.Println("Error importing PostgreSQL database:", string(output))
			return "", err
		}
	}

	s.logger.Printf("Imported %s database %q into %q", opts.Engine, opts.Database, opts.Container)

	return fmt.Sprintf("Imported %s database %q into %q", opts.Engine, opts.Database, opts.Container), nil
}

func (s *NitroService) createFile(dir string, pattern string) (*os.File, error) {
	f, err := ioutil.TempFile(dir, pattern)
	if err != nil {
		return nil, err
	}

	// TODO remove this
	// set the owner of the file to the ubuntu user for importing
	if err := os.Chown(f.Name(), 1000, 1000); err != nil {
		return nil, err
	}

	s.logger.Println("Created file", f.Name())

	return f, nil
}
