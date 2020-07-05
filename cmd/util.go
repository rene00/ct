package cmd

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Get value from user input.
func getValueFromConsole(value, valueText string) (string, error) {
	if value == "" {
		reader := bufio.NewReader(os.Stdin)
		if valueText == "" {
			valueText = "Value:"
		}
		fmt.Print(fmt.Sprintf("%s ", valueText))
		userInput, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		userInput = strings.TrimSuffix(userInput, "\n")
		if err != nil {
			return "", err
		}
		if userInput == "" {
			return "", errors.New("No user input")
		}
		return userInput, nil
	}
	return value, nil
}

// Check if string exists in slice.
func stringInSlice(str string, slice []string) bool {
	exists := false
	for _, v := range slice {
		if v == str {
			exists = true
			break
		}
	}
	return exists
}

func parseTimestamp(timestamp string) (time.Time, error) {
	parsedTimestamp, err := time.Parse("2006-01-02", timestamp)
	if err != nil {
		return parsedTimestamp, err
	}
	return parsedTimestamp, nil
}
