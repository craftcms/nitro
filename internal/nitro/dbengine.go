package nitro

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func FindEngineByDump(fileAbsPath string) string {

	file, err := os.Open(fileAbsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	foundEngine := ""
	scanner := bufio.NewScanner(file)
	line := 1

	for scanner.Scan() {
		fmt.Println(scanner.Text())
		_, matchesPostgres := checkSubstrings(scanner.Text(), "PostgreSQL", "pg_dump")
		_, matchesMySql := checkSubstrings(scanner.Text(), "MySQL", "mysqldump", "MariaDB")

		if matchesMySql > 0 {
			foundEngine = "mysql"
		}

		if matchesPostgres > 0 {
			foundEngine = "postgres"
		}

		// Stop looking after 50 lines or when we have found the engine
		if line > 50 || foundEngine != "" {
			break
		}

		line++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return foundEngine
}

func checkSubstrings(str string, subs ...string) (bool, int) {

	matches := 0
	isCompleteMatch := true

	for _, sub := range subs {
		if strings.Contains(str, sub) {
			matches += 1
		} else {
			isCompleteMatch = false
		}
	}

	return isCompleteMatch, matches
}
