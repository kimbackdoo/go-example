package scrapper

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func saveFileRowAsync(writer *csv.Writer, job JobData, channel chan<- bool) {
	row := []string{job.title, job.company, strings.Join(job.skills, ", "), job.level, job.location, job.url}
	checkErr(writer.Write(row))
	channel <- true
	fmt.Println("CSV 파일에 저장함", job.company)
}

func saveFileAsync(data []JobData, fileName string) {
	file, err := os.Create(fileName)
	checkErr(err)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"제목", "회사", "기술", "경력", "지역", "URL"}
	checkErr(writer.Write(headers))

	fileChannel := make(chan bool)
	for _, job := range data {
		go saveFileRowAsync(writer, job, fileChannel)
	}

	for i := 0; i < len(data); i++ {
		<-fileChannel
	}

	fmt.Println("전체: ", len(data), "jobs")
}

func jobDataAsync(element *goquery.Selection, channel chan<- JobData) {
	channel <- jobData(element)
}

func scrapAsync(jobGroup string, page int, pageChannel chan<- []JobData) {
	pageURL := pageBaseURL + fmt.Sprint(page) + "&jobGroup=" + jobGroup
	fmt.Println("Requesting", pageURL)

	res := httpGet(pageURL)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	jobCardChannel := make(chan JobData)
	jobCard := doc.Find(".css-mao678 > li")
	jobCard.Each(func(_ int, element *goquery.Selection) {
		go jobDataAsync(element, jobCardChannel)
	})

	var jobs []JobData
	for i := 0; i < jobCard.Length(); i++ {
		jobs = append(jobs, <-jobCardChannel)
	}

	pageChannel <- jobs
}

func crawlAsync(jobGroup string) {
	totalPages := totalPages()
	fmt.Println("총", totalPages, "페이지를 확인함")

	startTime := time.Now()

	channel := make(chan []JobData)
	for page := 1; page <= totalPages; page++ {
		go scrapAsync(jobGroup, page, channel)
	}

	var jobs []JobData
	for i := 0; i < totalPages; i++ {
		jobs = append(jobs, <-channel...)
	}

	saveFileAsync(jobs, "jobs.async.csv")

	gap := time.Since(startTime).Seconds()
	fmt.Println("Async 실행 시간:", gap, "초")
}
