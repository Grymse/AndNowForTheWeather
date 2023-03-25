package main

import (
	"os"
	"encoding/csv"
)

/**
 * Write the given data to the given file
 * @param data The data to write. This should be in CSV format
 * @param filename The filename to write to (Do not include .csv, this is added automatically)
 * @return The error, if any
*/
func writeRawCSVFile(data []byte, filename string) (err error) {
	return os.WriteFile(getFilePath(filename, "csv"), data, 0644);
}

/**
 * Read the given CSV file and parse it into a 2D array of strings.
 * @param filename The filename to read from (Do not include .csv, this is added automatically)
 * @return The data in the file
 * @return The error, if any
*/
func readCSVFile(filename string) ([][]string, error) {

	result := make([][]string, 0);

	bytes, err := os.Open(getFilePath(filename, "csv"));
	if err != nil {
		return result, err;
	}
	reader := csv.NewReader(bytes);
	
	result, err = reader.ReadAll();
	if err != nil {
		return result, err;
	}

	return result, nil
}

/**
 * Get the file path for the given area and file type
 * @param area The area to get the file path for
 * @param fileType The file type to get the file path for
 * @return The file path
*/
func getFilePath(area string, fileType string) string {
	return "./data/" + area + "." + fileType
}