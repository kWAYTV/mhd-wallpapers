package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type WallpaperData struct {
	Version int                               `json:"version"`
	Data    map[string]map[string]interface{} `json:"data"`
}

func main() {
	// Read the JSON file
	jsonData, err := os.ReadFile("wallpapers.json")
	if err != nil {
		fmt.Println("Error reading JSON file:", err)
		return
	}

	// Parse the JSON data
	var wallpaperData WallpaperData
	err = json.Unmarshal(jsonData, &wallpaperData)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Create output directory
	outputDir := "wallpapers"
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// Use a WaitGroup to wait for all downloads to complete
	var wg sync.WaitGroup

	// Iterate through the data and download wallpapers
	for _, item := range wallpaperData.Data {
		for key, value := range item {
			if url, ok := value.(string); ok && (strings.HasPrefix(key, "dhd") || strings.HasPrefix(key, "dsd")) {
				wg.Add(1)
				go downloadWallpaper(url, outputDir, &wg)
			}
		}
	}

	// Wait for all downloads to complete
	wg.Wait()
	fmt.Println("All downloads completed!")
}

func downloadWallpaper(url, outputDir string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Extract filename from URL
	filename := filepath.Base(url)
	filename = strings.Split(filename, "?")[0] // Remove query parameters

	// Create the full path for the output file
	outputPath := filepath.Join(outputDir, filename)

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error downloading %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Create the output file
	out, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("Error creating file %s: %v\n", outputPath, err)
		return
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("Error writing to file %s: %v\n", outputPath, err)
		return
	}

	fmt.Printf("Downloaded: %s\n", filename)
}