package nitrod

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/craftcms/nitro/internal/scripts"
)

func (s *NitroService) ImportDatabase(stream NitroService_ImportDatabaseServer) error {
	var container string
	var database string
	var compressed bool
	var isMySQL bool

	// create a temp file
	file, err := s.createFile(os.TempDir(), "nitro-db-upload")
	if err != nil {
		s.logger.Println("Error creating a temp file for the upload:", err.Error())
		return status.Errorf(codes.Internal, "Unable creating a temp file for the upload")
	}

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
		if container == "" {
			container = req.GetContainer()
		}
		if database == "" {
			database = req.GetDatabase()
		}
		if (compressed == false) && (req.GetCompressed()) {
			compressed = req.GetCompressed()
		}

		// write the backup content into the temp file
		_, err = file.Write(req.GetData())
		if err != nil {
			return status.Errorf(codes.Internal, "unable to write the backup to the temp file")
		}
	}

	s.logger.Printf("Importing database %q for container %q\n", database, container)

	// check the database engine
	if strings.Contains(container, "mysql") {
		isMySQL = true
	}

	// if the file is compressed, extract it and we are done
	if compressed {
		s.logger.Println("The file is compressed, extracting now")

		// open the recently saved file
		f, err := os.Open(file.Name())
		if err != nil {
			s.logger.Println("error opening the database import file: ", file.Name())
			return status.Errorf(codes.Unknown, "error opening the database import file: err: %s", err.Error())
		}

		// create the gzip reader
		reader, err := gzip.NewReader(f)
		if err != nil {
			s.logger.Println("error creating the gzip reader", err.Error())
			return status.Errorf(codes.Unknown, "error reading the compressed file. %s", err.Error())
		}
		reader.Multistream(true)

		compressedFile, err := s.createFile(os.TempDir(), "nitro-compressed-db-")
		if err != nil {
			s.logger.Println("error creating the compress db file:", err.Error())
			return err
		}

		// import the database
		if err := s.importDatabase(isMySQL, container, database, compressedFile.Name()); err != nil {
			return err
		}

		if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
			return status.Errorf(codes.Internal, "unable to send the response %v", err)
		}

		return nil
	}

	// import the database
	if err := s.importDatabase(isMySQL, container, database, file.Name()); err != nil {
		s.logger.Printf("Error imported the database: %s\n", err)
		return err
	}

	if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
		return status.Errorf(codes.Internal, "unable to send the response %v", err)
	}

	return nil
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

func (s *NitroService) importDatabase(mysql bool, container, database, file string) error {
	if mysql {
		// create the mysql database
		if output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf(scripts.FmtDockerMysqlCreateDatabaseIfNotExists, container, database)}); err != nil {
			s.logger.Println(string(output))
			return err
		}
		s.logger.Printf("Created the database %q\n", database)

		// copy the file into the containers tmp dir
		if output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf("docker cp %s %s:/", file, container)}); err != nil {
			s.logger.Println(string(output))
			return err
		}
		s.logger.Printf("Copied the file %q into the %q\n", file, container)

		s.logger.Printf("Beginning import of file %q", file)

		// remove the /tmp prefix since mysql defaults to /
		sp := strings.Split(file, "/")
		f := sp[len(sp)-1]

		// import the database
		output, err := s.command.Run("/bin/bash", []string{"-c", fmt.Sprintf("docker exec -i %q mysql -unitro -pnitro %s < /%s", container, database, f)})
		if err != nil {
			s.logger.Println(string(output))
			return err
		}

		s.logger.Println(string(output))

		s.logger.Printf("Imported the database %q into %q", database, container)

		return nil
	}

	s.logger.Println("creating postgres database")

	//switch  {
	//case tr:
	//if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerPostgresCreateDatabase, containerName, databaseName)); err != nil {
	//	fmt.Println(output)
	//	return err
	//}
	//
	//fmt.Println("Created database", databaseName)
	//
	//if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerPostgresImportDatabase, containerName, databaseName, fileFullPath)); err != nil {
	//	fmt.Println(output)
	//	return err
	//}
	//default:
	//if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerMysqlCreateDatabaseIfNotExists, containerName, databaseName)); err != nil {
	//	fmt.Println(output)
	//	return err
	//}
	//
	//fmt.Println("Created database", databaseName)
	//
	//if output, err := script.Run(false, fmt.Sprintf(scripts.FmtDockerMysqlImportDatabase, fileFullPath, containerName, databaseName)); err != nil {
	//	fmt.Println(output)
	//	return err
	//}
	//}
	return errors.New("postgres import is not implemented")
}
