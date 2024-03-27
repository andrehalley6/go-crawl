package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	if len(os.Args) != 3 || os.Args[1] != "crawl" {
		fmt.Println("Usage: go run main.go crawl <url>")
		os.Exit(1)
	}

	inputFile := os.Args[2]

	urls, err := readFileUrls(inputFile)
	if err != nil {
		fmt.Printf("Error reading URLs from file: %v\n", err)
		os.Exit(1)
	}

	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(wg *sync.WaitGroup, url string) {
			// err := saveCrawlResult(url)
			// if err != nil {
			// 	fmt.Printf("Error crawling and saving: %v\n", err)
			// }
			resp, err := getHttpResponse(url)
			if err != nil {
				fmt.Printf("Error getting HTTP response: %v\n", err)
			}

			err = saveResult(url, resp)
			if err != nil {
				fmt.Printf("Error saving result: %v\n", err)
			}
			wg.Done()
		}(&wg, url)
	}

	wg.Wait()
}

func readFileUrls(inputFile string) ([]string, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return urls, nil
}

func getHttpResponse(rawUrl string) ([]byte, error) {
	resp, err := http.Get(rawUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func saveResult(rawUrl string, body []byte) error {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	filePath := filepath.Join("result", parsedUrl.Hostname()+".html")
	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}
	return nil
}

func saveCrawlResult(rawUrl string) error {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return err
	}

	filePath := filepath.Join("result", parsedUrl.Hostname()+".html")

	resp, err := http.Get(rawUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, body, 0644)
	if err != nil {
		return err
	}
	return nil
}
