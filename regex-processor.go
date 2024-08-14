package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func main() {
	var choice, path, format, regexPattern, outputPath string

	// Step 1: Ask for file or directory using numbered options
	fmt.Println("Choose an option:")
	fmt.Println("1. Process a file")
	fmt.Println("2. Process a directory")
	fmt.Scanln(&choice)

	// Step 2: Ask for the file or directory path based on the user's choice
	switch choice {
	case "1":
		fmt.Println("Enter the file path:")
		fmt.Scanln(&path)
	case "2":
		fmt.Println("Enter the directory path:")
		fmt.Scanln(&path)
		// Step 2b: Get file format if the user chose directory
		fmt.Println("Enter the file format (e.g., '.txt', '.log'):")
		fmt.Scanln(&format)
	default:
		fmt.Println("Invalid choice.")
		return
	}

	// Step 3: Get regex pattern
	fmt.Println("Enter the regex pattern to search for:")
	fmt.Scanln(&regexPattern)
	re, err := regexp.Compile(regexPattern)
	if err != nil {
		fmt.Println("Invalid regex pattern:", err)
		return
	}

	// Step 4: Get the output file path
	fmt.Println("Enter the file path to save the result:")
	fmt.Scanln(&outputPath)

	// Step 5: Process files
	results := make(map[string]int)
	if choice == "1" {
		count, err := processFile(path, re)
		if err != nil {
			fmt.Println("Error processing file:", err)
			return
		}
		results[path] = count
	} else if choice == "2" {
		files, err := ioutil.ReadDir(path)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			return
		}
		for _, file := range files {
			if filepath.Ext(file.Name()) == format {
				filePath := filepath.Join(path, file.Name())
				count, err := processFile(filePath, re)
				if err != nil {
					fmt.Println("Error processing file:", err)
					return
				}
				results[filePath] = count
			}
		}
	}

	// Step 6: Save results to the output file
	err = saveResults(outputPath, results)
	if err != nil {
		fmt.Println("Error saving results:", err)
		return
	}

	fmt.Println("Results saved successfully!")
}

// processFile reads the file and counts lines matching the regex pattern
func processFile(filePath string, re *regexp.Regexp) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		if re.MatchString(line) {
			count++
		}
	}
	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

// saveResults writes the results to the specified output file
func saveResults(outputPath string, results map[string]int) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	for fileName, count := range results {
		line := fmt.Sprintf("%s: %d\n", fileName, count)
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}

	return nil
}
