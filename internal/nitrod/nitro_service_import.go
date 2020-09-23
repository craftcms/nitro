package nitrod

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *NitroService) ImportDatabase(stream NitroService_ImportDatabaseServer) error {
	var container string
	var database string
	var compressed bool

	// create a temp file
	file, err := ioutil.TempFile(os.TempDir(), "nitro-database-import-")
	if err != nil {
		s.logger.Println(fmt.Errorf("error creating the temp file %v", err))
		return status.Errorf(codes.FailedPrecondition, "Unable to create a temp file, err:", err.Error())
	}
	s.logger.Println("Created temporary file for database import", file.Name())

	// handle the file streaming requests
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "unable to create the stream", err)
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

	// if the file is compressed, extract it and we are done
	if compressed {
		s.logger.Println("The file is compressed, extracting now")

		// open the recently saved file
		f, err := os.Open(file.Name())
		if err != nil {
			s.logger.Println("error opening the database import file: ", file.Name())
			return status.Errorf(codes.Unknown, "error opening the database import file", err.Error())
		}

		// create the gzip reader
		reader, err := gzip.NewReader(f)
		if err == io.EOF {
			s.logger.Println("got eof in the compressed reader", err)
			return nil
		}
		if err != nil {
			s.logger.Println("error creating the gzip reader", err.Error())
			return status.Errorf(codes.Unknown, "error reading the compressed file. %w", err.Error())
		}
		reader.Multistream(true)

		// create a new compressed file to extract into
		compressedFile, err := ioutil.TempFile(os.TempDir(), "nitro-compressed-database-")
		if err != nil {
			s.logger.Println(fmt.Errorf("error creating the compressed temp file %v", err))
			return status.Errorf(codes.FailedPrecondition, "Unable to create a compressed temp file, err:", err.Error())
		}
		s.logger.Println("Created temporary file for compressed database import", compressedFile.Name())

		if _, err := io.Copy(compressedFile, reader); err != nil {
			s.logger.Println("error copying the compressed file. %w", err)
			return err
		}

		if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
			return status.Errorf(codes.Internal, "unable to send the response %v", err)
		}

		return nil
	}

	if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
		return status.Errorf(codes.Internal, "unable to send the response %v", err)
	}

	return nil
}
