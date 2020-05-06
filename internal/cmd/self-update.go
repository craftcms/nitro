package cmd

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var selfUpdateCommand = &cobra.Command{
	Use:   "self-update",
	Short: "Update Nitro to the latest",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println(runtime.GOOS)
		fileUrl := "https://raw.githubusercontent.com/craftcms/nitro/master/install.sh"

		tempFile, _ := os.Getwd()
		tempFile += "/temp_nitro_update.sh"

		// localFile := filepath.Join(tempFolder, "temp_nitro_update.sh")

		if err := DownloadFile(tempFile, fileUrl); err != nil {
			return err
		}
		fmt.Println("successfully downloaded file to "+tempFile)
		if err := os.Chmod(tempFile, 0777); err != nil {
			return err
		}
//		test := exec.Command()
//		output, err := test.StdoutPipe()
//		if err != nil {
//			return err
//		}

//		if err := test.Start(); err != nil {
//			fmt.Println(err)
//		}
//		fmt.Println(output)
//		return nil
		ch := make(chan string)
		go func() {
			//if err := RunCommandCh(ch, "\r\n", "C:\\Windows\\system32\\cmd.exe", "/c", "\"\"C:\\Program Files\\Git\\bin\\sh.exe\" --login -i -- D:\\dev\\nitro\\temp_nitro_update.sh\""); err != nil {
			if err := RunCommandCh(ch, "\r\n", tempFile); err != nil {
				log.Fatal(err)
			}
		}()

		for v := range ch {
			fmt.Println(v)
		}

		return nil
	},
}

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
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
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

	return cmd.Wait()
}
