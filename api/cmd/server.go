package main

import (
	"net/http"
	"rei-api"

	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
)

func addFeed(c echo.Context) error {
	rss_link := c.FormValue("rss_link")
	reader := rei.NewParser()
	feed, _ := reader.ReadFeed(rss_link)
	// name := feed.Title
	link := feed.Link

	//save in db the relevant content
	// author, link, description etc
	return c.JSON(http.StatusOK, link)
}

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.GET("/feeds", func(c echo.Context) error {
		// it should get custom feeds from a user
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURL("https://www.gmkonan.dev/rss.xml")
		feeds := feed.Items
		return c.JSON(http.StatusOK, feeds)
	})

	// add feed that user inputed
	// parse the rss and save relevant content
	e.POST("/add_feed", addFeed)

	// get all feeds from user
	e.Logger.Fatal(e.Start(":1323"))
}
