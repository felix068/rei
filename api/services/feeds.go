package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"rei-api"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	db     *sqlx.DB
	reader *rei.Reader
}
type RssLink struct {
	Link string `json:"link"`
}

type Feed struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Link        string    `json:"link" db:"link"`
	RssLink     string    `json:"rssLink" db:"rss_link"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   string    `json:"updatedAt" db:"updated_at"`
}

func FeedHandler(dbClient *sqlx.DB, reader *rei.Reader) *Handler {
	return &Handler{
		db:     dbClient,
		reader: reader,
	}
}

// create posts model and table and test

func (h *Handler) addFeed(c echo.Context) error {
	rssInput := new(RssLink)
	c.Bind(rssInput)

	fmt.Println(rssInput.Link)
	feed, readFeedErr := h.reader.ReadFeed(rssInput.Link)

	if readFeedErr != nil {
		fmt.Println(readFeedErr)
		return c.JSON(http.StatusInternalServerError, readFeedErr)
	}

	feedUuid := uuid.New()

	checkIfRssLinkExists := `SELECT 1 FROM feeds WHERE rss_link = $1 LIMIT 1`
	row := h.db.QueryRow(checkIfRssLinkExists, rssInput.Link)

	fmt.Println(&row)
	var temp int // We just need something to scan into
	err := row.Scan(&temp)

	if err != nil && err != sql.ErrNoRows {
		return c.JSON(http.StatusInternalServerError, "database error checking RSS link existence")
	}

	if err == nil {
		return c.JSON(http.StatusConflict, "Duplicate feed, this feed is already on db")
	}

	fmt.Println(feed.FeedLink)
	// The sqlx library converts the named parameters (:name) to PostgreSQL's numbered format ($1, $2, etc.) behind the scenes.
	insertFeed := `INSERT INTO feeds (id, name, description, link, rss_link, created_at, updated_at) 
	VALUES (:id, :name, :description, :link, :rss_link, :createdAt, :updatedAt)`

	_, insertErr := h.db.NamedExec(insertFeed, map[string]interface{}{
		"id":          feedUuid,
		"name":        feed.Title,
		"description": feed.Description,
		"link":        feed.Link,
		"rss_link":    rssInput.Link,
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	})

	if insertErr != nil {
		fmt.Println(insertErr)
		return c.JSON(http.StatusInternalServerError, insertErr)
	}

	fmt.Println(feed.Items[0])
	insertPost := `INSERT INTO posts (id, name, description, link, feed_id, created_at, updated_at) 
	VALUES (:id, :name, :description, :link, :feed_id, :createdAt, :updatedAt)`

	// add posts that already exist
	for _, item := range feed.Items {
		_, err := h.db.NamedExec(insertPost, map[string]interface{}{
			"id":          uuid.New(),
			"name":        item.Title,
			"description": item.Description,
			"link":        item.Link,
			"feed_id":     feedUuid,
			"createdAt":   time.Now(),
			"updatedAt":   time.Now(),
		})

		if err != nil {
			fmt.Println(err)
			return c.JSON(http.StatusInternalServerError, err)
		}
	}

	// rei.Sync(rssInput.Link, h.reader)

	return c.JSON(http.StatusCreated, feed.Link)
}

func (h *Handler) listFeeds(c echo.Context) error {
	// wip
	var feeds Feed
	h.db.Get(&feeds, "SELECT * FROM feeds")

	fmt.Println("feeds", feeds.Name)

	return c.JSON(http.StatusCreated, feeds)

}

func (h *Handler) RegisterFeedsRoutes(e *echo.Echo) {
	e.POST("/add_feed", h.addFeed)
	e.GET("/list_feeds", h.listFeeds)
}
