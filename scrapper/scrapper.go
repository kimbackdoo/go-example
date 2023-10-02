package scrapper

import (
	"os"
	"strings"

	"github.com/labstack/echo/v4"
)

// func Example() {
// 	crawlAsync()
// 	fmt.Println()
// 	crawl()
// }

func handleHome(c echo.Context) error {
	return c.File("templates/home.html")
}

func handleScrap(c echo.Context) error {
	defer os.Remove("jobs.async.csv")

	jobGroup := strings.ToUpper(cleanString(c.FormValue("jobGroup")))
	crawlAsync(jobGroup)
	return c.Attachment("jobs.csv", "jobs.csv")
}

func Bootstrap() {
	e := echo.New()
	e.GET("/", handleHome)
	e.POST("/scrap", handleScrap)
	e.Logger.Fatal(e.Start(":1323"))
}
