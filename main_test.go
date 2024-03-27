package main

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestReadFileUrls(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test-urls.txt")
	if err != nil {
		t.Fatalf("Error creating temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	// Write sample URLs to the temporary file
	urls := []string{
		"https://google.com",
		"https://go.dev",
	}
	for _, u := range urls {
		_, err := tmpFile.WriteString(u + "\n")
		if err != nil {
			t.Fatalf("Error writing to temporary file: %v", err)
		}
	}

	// Read Urls from the temporary file using the function being tested
	readUrls, err := readFileUrls(tmpFile.Name())
	if err != nil {
		t.Fatalf("Error reading URLs from file: %v", err)
	}

	// Check if the read Urls match the expected Urls
	for i, u := range urls {
		if readUrls[i] != u {
			t.Errorf("Expected URL: %s, Got: %s", u, readUrls[i])
		}
	}
}

func TestGetHttpResponse(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mock HTML Content"))
	}))
	defer mockServer.Close()

	// Call the function being tested
	_, err := getHttpResponse(mockServer.URL)
	if err != nil {
		t.Fatalf("Error getting HTTP response: %v", err)
	}
}

func TestSaveResult(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mock HTML Content"))
	}))
	defer mockServer.Close()

	// Call the function being tested
	err := saveResult(mockServer.URL, []byte("Mock HTML Content"))
	if err != nil {
		t.Fatalf("Error saving result: %v", err)
	}

	parseUrl, err := url.Parse(mockServer.URL)
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}
	filePath := filepath.Join("result", parseUrl.Hostname()+".html")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Check if the file is created and contains the expected content
	expectedContent := "Mock HTML Content"
	if string(content) != expectedContent {
		t.Errorf("Expected content: %s, Got: %s", expectedContent, string(content))
	}
}

func TestSaveCrawlResult(t *testing.T) {
	// Create a mock server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Mock HTML Content"))
	}))
	defer mockServer.Close()

	// Call the function being tested
	err := saveCrawlResult(mockServer.URL)
	if err != nil {
		t.Fatalf("Error crawling and saving: %v", err)
	}

	parseUrl, err := url.Parse(mockServer.URL)
	if err != nil {
		t.Fatalf("Error parsing URL: %v", err)
	}
	filePath := filepath.Join("result", parseUrl.Hostname()+".html")
	content, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Error reading file: %v", err)
	}

	// Check if the file is created and contains the expected content
	expectedContent := "Mock HTML Content"
	if string(content) != expectedContent {
		t.Errorf("Expected content: %s, Got: %s", expectedContent, string(content))
	}
}
