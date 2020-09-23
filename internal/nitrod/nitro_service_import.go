package nitrod

import (
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

	// create a temp file
	file, err := ioutil.TempFile(os.TempDir(), "nitro-database-import-")
	if err != nil {
		s.logger.Println(fmt.Errorf("error creating the temp file %v", err))
		return status.Errorf(codes.FailedPrecondition, "Unable to create a temp file, err:", err.Error())
	}
	s.logger.Println("Created temporary file for database import", file.Name())

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "unable to create the stream", err)
		}

		container = req.GetContainer()
		database = req.GetDatabase()

		// get the backup
		backup := req.GetData()

		// write the backup content into the temp file
		_, err = file.Write(backup)
		if err != nil {
			return status.Errorf(codes.Internal, "unable to write the backup to the temp file")
		}
	}

	s.logger.Printf("Creating database %s for container %s\n", database, container)

	if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
		return status.Errorf(codes.Internal, "unable to send the response %v", err)
	}

	return nil
}
