package nitrod

import (
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"

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
	file, err := s.createFile(os.TempDir(), "nitro-db-upload")
	if err != nil {
		s.logger.Println("Error creating a temp file for the upload:", err.Error())
		return status.Errorf(codes.Internal, "Unable creating a temp file for the upload")
	}

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

		// open the recently saved file
		f, err := os.Open(file.Name())
		if err != nil {
			s.logger.Println("error opening the database import file: ", file.Name())
			return status.Errorf(codes.Unknown, "error opening the database import file: err: %s", err.Error())
		}

		// create the gzip reader
		switch options.CompressionType {
		case "gz":
			reader, err := gzip.NewReader(f)
			if err != nil {
				s.logger.Println("error creating the gzip reader", err.Error())
				return status.Errorf(codes.Unknown, "error reading the compressed file. %s", err.Error())
			}

			reader.Multistream(true)

			defer reader.Close()
		default:
			reader, err := zip.OpenReader(f.Name())
			if err != nil {
				return err
			}

			defer reader.Close()
		}

		compressedFile, err := s.createFile(os.TempDir(), "nitro-compressed-db-")
		if err != nil {
			s.logger.Println("error creating the compressed db file:", err.Error())
			return err
		}

		if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
			return status.Errorf(codes.Internal, "unable to send the response %v", err)
		}
		options.File = compressedFile.Name()
	}

	// import the database
	if _, err := s.importDatabase(options); err != nil {
		s.logger.Printf("Error importing database: %s\n", err)
		return err
	}

	if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
		return status.Errorf(codes.Internal, "unable to send the response %v", err)
	}

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
			s.logger.Printf("Created the database %q\n", opts.Database)
		}

		// copy the file into the containers tmp dir
		if output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf("docker cp %s %s:/tmp", opts.File, opts.Container)}); err != nil {
			s.logger.Println()
			return string(output), err
		}
		s.logger.Printf("Copied file %q into container %q\n", opts.File, opts.Container)

		s.logger.Printf("Beginning import of file %q", opts.File)

		// if we are skipping create, it has the use statement and no database name
		if opts.CreateDatabase {
			opts.Database = "emptydatabase"
		}

		// import the database
		output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf("docker exec -i %q mysql -unitro -pnitro %s < %s", opts.Container, opts.Database, opts.File)})
		if err != nil {
			s.logger.Println(string(output))
			return "", err
		}
	default:
		output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf(scripts.FmtDockerPostgresCreateDatabase, opts.Container, opts.Database)})
		if err != nil {
			s.logger.Println(string(output))
			return "", err
		}
		s.logger.Printf("created database %q for engine %q", opts.Database, opts.Container)

		output, err = s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf(scripts.FmtDockerPostgresImportDatabase, opts.Container, opts.Database, opts.File)})
		if err != nil {
			s.logger.Println(string(output))
			return "", err
		}
	}

	s.logger.Printf("Imported database %q into %q", opts.Database, opts.Container)

	return fmt.Sprintf("Imported database %q into %q", opts.Database, opts.Container), nil
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
