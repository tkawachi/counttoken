package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pkoukk/tiktoken-go"
)

func main() {
	path := flag.String("path", "", "Path to the file or directory")
	model := flag.String("model", "gpt-4", "Model name")
	flag.Parse()

	if *path == "" {
		log.Fatal("Please provide a valid path")
	}

	fileInfo, err := os.Stat(*path)
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		totalTokenCount, err := countTokensInDirectory(*model, *path)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Total Token count: %d\n", totalTokenCount)
	} else {
		fileTokenCount, err := countTokensInFile(*model, *path)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Token count: %d\n", fileTokenCount)
	}
}

func countTokensInDirectory(model, directoryPath string) (int, error) {
	totalTokenCount := 0

	err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			fileTokenCount, err := countTokensInFile(model, path)
			if err != nil {
				log.Printf("Error counting tokens in file %s: %v", path, err)
				return nil
			}

			fmt.Printf("File: %s, Token count: %d\n", path, fileTokenCount)
			totalTokenCount += fileTokenCount
		}

		return nil
	})

	if err != nil {
		return 0, err
	}

	return totalTokenCount, nil
}

func countTokensInFile(model, filePath string) (int, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	tkm, err := tiktoken.EncodingForModel(model)
	if err != nil {
		return 0, err
	}

	// Read whole file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return 0, err
	}
	token := tkm.Encode(string(data), nil, nil)
	tokenCount := len(token)

	return tokenCount, nil
}
