package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

const baseURL = "http://localhost:8080/api/v1"

const SPLIT_LINE = "==========================================="

func clearChannels() {
	// Create client
	client := &http.Client{}

	// Create request
	req, err := http.NewRequest("DELETE", baseURL+"/channels", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	// Read Response Body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(respBody))
}

func readNews() {
	resp, err := http.Get(baseURL + "/news")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func readJobs() {
	resp, err := http.Get(baseURL + "/jobs")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func readChannels() {
	resp, err := http.Get(baseURL + "/channels")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func addChannel(url string) {
	payload := map[string]string{"link": url}
	data, _ := json.Marshal(payload)
	resp, err := http.Post(baseURL+"/channels", "application/json", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func askModel(query string) {
	resp, err := http.Get(baseURL + "/query/" + url.QueryEscape(query))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:", string(body))
}

func main() {
	rssURL := flag.String("url", "", "RSS feed URL")
	reset := flag.Bool("reset", false, "Clear all channels")
	query := flag.String("query", "", "Query to last news")
	showChannels := flag.Bool("show-channels", false, "Show all channels")
	showNews := flag.Bool("show-news", false, "Show all news")
	showJobs := flag.Bool("show-jobs", false, "Show all jobs")
	flag.Parse()

	if *rssURL == "" && !*showChannels && !*reset && !*showNews && !*showJobs && *query == "" {
		log.Fatal("Please specify either --url or --show-channels or --show-news or --reset parameter")
	}

	if *reset {
		clearChannels()
	}

	if *rssURL != "" {
		addChannel(*rssURL)
	}

	if *showNews {
		readNews()
	}

	if *showChannels {
		readChannels()
	}

	if *showJobs {
		readJobs()
	}

	log.Printf("Query %s", *query)

	if *query != "" {
		askModel(*query)
	}
}
