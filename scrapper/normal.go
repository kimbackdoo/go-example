package scrapper

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func saveFile(data []JobData, fileName string) {
	file, err := os.Create(fileName)
	checkErr(err)

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headers := []string{"제목", "회사", "기술", "경력", "지역", "URL"}
	checkErr(writer.Write(headers))

	for _, job := range data {
		row := []string{job.title, job.company, strings.Join(job.skills, ", "), job.level, job.location, job.url}
		checkErr(writer.Write(row))
		fmt.Println("CSV 파일에 저장함", job.company)
	}

	fmt.Println("전체: ", len(data), "jobs")
}

func scrap(page int) []JobData {
	pageURL := pageBaseURL + fmt.Sprint(page)
	fmt.Println("Requesting", pageURL)

	res := httpGet(pageURL)
	defer res.Body.Close()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	checkErr(err)

	var jobs []JobData
	doc.Find(".css-mao678 > li").Each(func(_ int, element *goquery.Selection) {
		jobs = append(jobs, jobData(element))
	})

	return jobs
}

func crawl() {
	totalPages := totalPages()
	fmt.Println("총", totalPages, "페이지를 확인함")

	startTime := time.Now()

	var jobs []JobData
	for page := 1; page <= totalPages; page++ {
		jobs = append(jobs, scrap(page)...)
	}

	saveFile(jobs, "jobs.csv")

	gap := time.Since(startTime).Seconds()
	fmt.Println("실행 시간:", gap, "초")
}
