package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/mmcdole/gofeed"
)

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
	e.Logger.Fatal(e.Start(":1323"))
}
