package urlchecker

import (
	"fmt"
	"net/http"
)

type urlInfo struct {
	url    string
	status string
}

func hitUrl(url string, channel chan<- urlInfo) {
	res, err := http.Get(url)
	if err != nil || res.StatusCode >= 400 {
		channel <- urlInfo{url: url, status: "FAILED"}
	}
	channel <- urlInfo{url: url, status: "OK"}
}

func Example() {
	urls := []string{
		"https://www.naver.com",
		"https://www.google.com",
		"https://www.facebook.com",
		"https://www.instagram.com",
		"https://www.amazon.com",
		"https://www.youtube.com",
		"https://www.netflix.com",
		"https://www.kakao.com",
		"https://www.daum.net",
		"https://www.tistory.com",
		"https://www.nate.com",
		"https://www.yahoo.com",
		"https://www.bing.com",
		"https://www.wikipedia.org",
		"https://www.reddit.com",
	}

	channel := make(chan urlInfo)
	results := make(map[string]string)
	for _, url := range urls {
		go hitUrl(url, channel)
	}

	for i := 0; i < len(urls); i++ {
		res := <-channel
		results[res.url] = res.status
	}

	for url, result := range results {
		fmt.Println(url, result)
	}
}
