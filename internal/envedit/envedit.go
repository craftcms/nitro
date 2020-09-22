package envedit

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

// Editer is an interface that is responsible for seting environment variables and
// saving the files
type Editer interface {
	Set(env, val string)
	Save() error
}

type FileEditer struct {
	file *os.File
}

func New(file string) (*FileEditer, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	return &FileEditer{
		file: f,
	}, nil
}

func (e *FileEditer) Set(env, val string) {
	var lines []string

	// scan the files and change them
	scanner := bufio.NewScanner(e.file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), env) {
			sp := strings.Split(scanner.Text(), "=")

			switch len(sp) {
			case 2:
				sp[len(sp)-1] = val
			}

			l := strings.Join(sp, "=")
			fmt.Println(l)
			lines = append(lines, l)
		} else {
			lines = append(lines, scanner.Text())
		}
	}

	w := bufio.NewWriter(e.file)
	for _, line := range lines {
		_, err := fmt.Fprintln(w, line)
		if err != nil {
			fmt.Println(err)
		}
	}

	if err := w.Flush(); err != nil {
		fmt.Println(err)
	}
}

func (e *FileEditer) Save() error {
	return errors.New("not implemented")
}

func Set(file, env, val string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// read the file and only change what we are looking for
	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), env) {
			sp := strings.Split(scanner.Text(), "=")
			switch len(sp) {
			case 2:
				sp[2] = val
			}

			lines = append(lines, strings.Join(sp, "="))
		} else {
			lines = append(lines, scanner.Text())
		}
	}

	return nil
}
