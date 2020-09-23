package nitrod

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxSize = 1 << 20

func (s *NitroService) ImportDatabase(stream NitroService_ImportDatabaseServer) error {
	var container string
	var database string

	buffer := bytes.Buffer{}
	backupSize := 0

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			s.logger.Println("size of backup is:", backupSize)
			break
		}
		if err != nil {
			return status.Errorf(codes.Internal, "unable to create the stream", err)
		}

		container = req.GetContainer()
		database = req.GetDatabase()

		// get the backup
		backup := req.GetData()
		size := len(backup)

		backupSize += size
		if size > maxSize {
			return status.Errorf(codes.InvalidArgument, "the backup size is too large")
		}

		// write the backup content into the buffer
		_, err = buffer.Write(backup)
		if err == bufio.ErrBufferFull {
			return err
		}
		if err != nil {
			return status.Errorf(codes.Internal, "error writing to the buffer", err.Error())
		}
	}

	s.logger.Println("out of the for loop")

	// create a temp file
	tempFile, err := ioutil.TempFile(os.TempDir(), "nitro-database-import-")
	if err != nil {
		s.logger.Println(fmt.Errorf("error creating the temp file %v", err))
		return status.Errorf(codes.FailedPrecondition, "Unable to create a temp file, err:", err.Error())
	}

	s.logger.Println("Created temporary file for database import", tempFile.Name())

	// save the data into the file
	_, err = tempFile.Write(buffer.Bytes())
	if err != nil {
		return status.Errorf(codes.Internal, "unable to write the backup to the temp file")
	}

	s.logger.Printf("Creating database %s for container %s\n", database, container)

	if err := stream.SendAndClose(&ServiceResponse{Message: "Successfully imported the database"}); err != nil {
		return status.Errorf(codes.Internal, "unable to send the response %v", err)
	}

	return nil
}
