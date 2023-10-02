package scrapper

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type APIResponse struct {
	Data APIData
}

type APIData struct {
	TotalPage int
}

type JobData struct {
	title    string
	company  string
	skills   []string
	level    string
	location string
	url      string
}

var (
	baseURL     = "https://www.rallit.com"
	pageBaseURL = baseURL + "?pageNumber="
	apiBaseURL  = baseURL + "/api/v1/position?pageSize=20&pageNumber="
)

func cleanString(str string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(str)), " ")
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func checkStatusCode(res *http.Response) {
	if res.StatusCode != 200 {
		log.Fatalln("Request failed with Status:", res.StatusCode)
	}
}

func httpGet(url string) *http.Response {
	res, err := http.Get(url)
	checkErr(err)
	checkStatusCode(res)

	return res
}

func httpGetData(url string) APIData {
	res := httpGet(url)
	defer res.Body.Close()

	data, err := io.ReadAll(res.Body)
	checkErr(err)

	var apiResponse APIResponse
	checkErr(json.Unmarshal(data, &apiResponse))

	return apiResponse.Data
}

func totalPages() int {
	return httpGetData(apiBaseURL + "1").TotalPage
}

func jobData(element *goquery.Selection) JobData {
	title := element.Find(".summary__title").Text()
	company := element.Find(".summary__company-name").Text()
	skills := []string{}
	element.Find(".css-13kyeyo").Each(func(_ int, element *goquery.Selection) {
		skills = append(skills, element.Text())
	})
	level := element.Find(".css-oz575p").Text()
	location := element.Find(".css-oz575p").Text()
	href, _ := element.Find("a").Attr("href")

	return JobData{
		title:    title,
		company:  company,
		skills:   skills,
		level:    level,
		location: location,
		url:      baseURL + href,
	}
}
