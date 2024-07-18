package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"
)

type ShellHistoryEntry struct {
	Timestamp  string
	ReturnType string
	Command    string
}

func LoadShellHistory() ([]ShellHistoryEntry, error) {
	HISTFILE_LOCATION := "/Users//.zsh_history"

	history := []ShellHistoryEntry{}

	file, err := os.Open(HISTFILE_LOCATION)
	if err != nil {
		return nil, fmt.Errorf("something went wrong while trying to open %s history file: %v", strings.Map(unicode.ToLower, HISTFILE_LOCATION), err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Ignore lines without the expected format
		if !strings.HasPrefix(line, ": ") {
			continue
		}

		// Parse the line into timestamp, return type, and command
		parts := strings.SplitN(line, ";", 2)
		if len(parts) != 2 {
			continue
		}

		infoParts := strings.Split(parts[0], ":")
		if len(infoParts) < 3 {
			continue
		}

		var timestamp int64
		var returnType int

		fmt.Sscanf(infoParts[1], "%d", &timestamp)
		fmt.Sscanf(infoParts[2], "%d", &returnType)

		timeFormatted := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		var formattedReturnType string
		if returnType == 0 {
			formattedReturnType = "OK"
		} else {
			formattedReturnType = "ERROR"
		}

		command := parts[1]

		history = append(history, ShellHistoryEntry{
			Timestamp:  timeFormatted,
			ReturnType: formattedReturnType,
			Command:    command,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("something went wrong while trying to read %s history file: %v", strings.Map(unicode.ToLower, HISTFILE_LOCATION), err)
	}

	return history, nil
}
