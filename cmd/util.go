package cmd

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

// Insert the metric into the metric table.
func setMetricID(db *sql.DB, metric Metric) error {
	var sqlStmt string

	sqlStmt = `
	INSERT INTO metric
		(
			id,
			name
		)
		VALUES
		(
			NULL,
			?
		)
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(metric.Name); err != nil {
		return err
	}
	return nil
}

// Get the metric ID. If the metric does not exist within the metric table, getMetricID will call setMetricID to create the metric within the table and then call itself to return the metric ID.
func getMetricID(db *sql.DB, metric Metric) (int, error) {
	var sqlStmt string
	var metricID int

	sqlStmt = `
	SELECT id
		FROM metric
		WHERE name = ?
	`
	stmt, err := db.Prepare(sqlStmt)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	if err = stmt.QueryRow(metric.Name).Scan(&metricID); err != nil {
		if err = setMetricID(db, metric); err != nil {
			return 0, err
		}
	}
	if metricID == 0 {
		return getMetricID(db, metric)
	}
	return metricID, nil
}

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

func getMetricConfig(db *sql.DB, metric Metric) (MetricConfig, error) {
	var sqlStmt string
	metricConfig := MetricConfig{}

	sqlStmt = `
	SELECT opt, val
		FROM CONFIG
		WHERE metric_id = ?
	`
	rows, err := db.Query(sqlStmt, metric.ID)
	if err != nil {
		return metricConfig, err
	}
	defer rows.Close()

	for rows.Next() {
		var option string
		var value string
		if err = rows.Scan(&option, &value); err != nil {
			return metricConfig, err
		}
		switch option {
		case "frequency":
			metricConfig.Frequency = value
		case "value_text":
			metricConfig.ValueText = value
		case "data_type":
			metricConfig.DataType = value
		default:
			return metricConfig, errors.New("Unsupported metric config")
		}
	}
	return metricConfig, nil
}

func parseTimestamp(timestamp string) (time.Time, error) {
	parsedTimestamp, err := time.Parse("2006-01-02", timestamp)
	if err != nil {
		return parsedTimestamp, err
	}
	return parsedTimestamp, nil
}
