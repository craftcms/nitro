package cmd

import (
	"fmt"
	"github.com/craftcms/nitro/internal/helpers"
	"io"
	"io/ioutil"
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
	Short: "Update Nitro",
	RunE: func(cmd *cobra.Command, args []string) error {
		currentOs := runtime.GOOS
		fileUrl := "https://raw.githubusercontent.com/craftcms/nitro/master/install.sh"

		if currentOs == "windows" {
			fmt.Println("Self updating is not available in Windows, please visit https://github.com/craftcms/nitro to update")
			return nil
		}

		tempDirectory, _ := os.Getwd()
		tempShell := tempDirectory + "/temp_nitro_update.sh"
		tempBat := tempDirectory + "/temp_nitro_update.bat"

		if err := DownloadFile(tempShell, fileUrl); err != nil {
			return err
		}

		if err := os.Chmod(tempShell, 0777); err != nil {
			return err
		}

		ch := make(chan string)
		go func() {
			if currentOs == "windows" {

				batOutput := []byte("sh.exe --login -i -- " + tempShell + " update")
				if err := ioutil.WriteFile(tempBat, batOutput, 0644); err != nil {
					log.Fatal(err)
				}

				if err := RunCommandCh(ch, "\r\n", os.ExpandEnv("$ComSpec"), "/c", tempBat); err != nil {
					log.Fatal(err)
				}
			} else {
				if err := RunCommandCh(ch, "\r\n", tempShell, "update"); err != nil {
					log.Fatal(err)
				}
			}
		}()

		for v := range ch {
			fmt.Println(v)
		}

		if helpers.FileExists(tempBat) {
			os.Remove(tempBat)
		}

		return os.Remove(tempShell)
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
