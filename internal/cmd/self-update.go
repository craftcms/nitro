package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	selfUpdateCommand = &cobra.Command{
		Use:   "self-update",
		Short: "Update nitro to the latest",
		RunE: func(cmd *cobra.Command, args []string) error {
			fileUrl := "https://raw.githubusercontent.com/pixelandtonic/nitro/develop/get.sh"

			tempFolder := os.TempDir()

			localFile := filepath.Join(tempFolder, "get.sh")

			if err := DownloadFile(localFile, fileUrl); err != nil {
				panic(err)
			}

			err1 := os.Chmod(localFile, 0777)
			if err1 != nil {
				log.Println(err1)
			}

			ch := make(chan string)
			go func() {
				err := RunCommandCh(ch, "\r\n", localFile)
				if err != nil {
					log.Fatal(err)
				}
			}()

			for v := range ch {
				fmt.Println(v)
			}


			defer os.Remove(tempFolder)

			return nil
		},
	}
)

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}


// RunCommandCh runs an arbitrary command and streams the output to a channel.
func RunCommandCh(stdoutCh chan<- string, cutset string, command string, flags ...string) error {
	cmd := exec.Command(command, flags...)

	output, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("RunCommand: cmd.StdoutPipe(): %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("RunCommand: cmd.Start(): %v", err)
	}

	go func() {
		defer close(stdoutCh)
		for {
			buf := make([]byte, 1024)
			n, err := output.Read(buf)
			if err != nil {
				if err != io.EOF {
					log.Fatal(err)
				}
				if n == 0 {
					break
				}
			}
			text := strings.TrimSpace(string(buf[:n]))
			for {
				// Take the index of any of the given cutset
				n := strings.IndexAny(text, cutset)
				if n == -1 {
					// If not found, but still have data, send it
					if len(text) > 0 {
						stdoutCh <- text
					}
					break
				}
				// Send data up to the found cutset
				stdoutCh <- text[:n]
				// If cutset is last element, stop there.
				if n == len(text) {
					break
				}
				// Shift the text and start again.
				text = text[n+1:]
			}
		}
	}()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("RunCommand: cmd.Wait(): %v", err)
	}
	return nil
}