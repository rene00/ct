package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
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

func getBoolFromValue(value string) (string, error) {
	for k, v := range map[string][]string{
		"0": []string{"n", "no"},
		"1": []string{"y", "yes"},
	} {
		for _, i := range v {
			if strings.ToLower(value) == i {
				return k, nil
			}
		}
	}

	_, err := strconv.ParseBool(value)
	if err != nil {
		return "", fmt.Errorf("Value is not a bool")
	}

	return value, nil
}
