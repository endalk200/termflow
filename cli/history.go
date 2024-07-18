package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)

func LoadShellHistory() ([]string, error) {
	// HISTFILE_LOCATION, ok := os.LookupEnv("HISTFILE")
	// if !ok {
	// 	log.Panicf("HISTFILE environment variable is not set")
	// }

	HISTFILE_LOCATION := "/Users/endalk200/.zsh_history"

	history := []string{}

	file, err := os.Open(HISTFILE_LOCATION)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while trying to open %s history file: %v", strings.Map(unicode.ToLower, HISTFILE_LOCATION), err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		history = append(history, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("something went wrong while trying to read %s history file: %v", strings.Map(unicode.ToLower, HISTFILE_LOCATION), err)
	}

	return history, nil
}

// Function to convert Unix timestamp to datetime string
func timestampToDatetime(timestamp string) string {
	// Example: timestamp is "1720979874"
	// Assuming timestamp is in seconds, convert to datetime string
	// Replace this with your preferred conversion logic based on your requirements
	// Here we simply use the timestamp as-is for illustration
	return timestamp
}
