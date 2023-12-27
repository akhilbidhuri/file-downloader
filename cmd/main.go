package main

import (
	"fmt"
	"os"

	"github.com/akhilbidhuri/file-downloader/internal/download"
	"github.com/akhilbidhuri/file-downloader/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <file_url> [output_file_path]")
		return
	}

	url := os.Args[1]
	if !utils.ValidateURL(url) {
		fmt.Println("Invalid URL provided!")
		return
	}
	outputPath := ""
	if len(os.Args) > 2 {
		outputPath = os.Args[2]
	}

	err := download.Process(url, outputPath)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
