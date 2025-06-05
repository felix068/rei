package services

import (
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
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   string    `json:"updatedAt" db:"updated_at"`
}

func FeedHandler(dbClient *sqlx.DB, reader *rei.Reader) *Handler {
	return &Handler{
		db:     dbClient,
		reader: reader,
	}
}

func (h *Handler) addFeed(c echo.Context) error {
	rssInput := new(RssLink)
	c.Bind(rssInput)

	fmt.Println(rssInput.Link)
	feed, _ := h.reader.ReadFeed(rssInput.Link)
	uuid := uuid.New()

	// The sqlx library converts the named parameters (:name) to PostgreSQL's numbered format ($1, $2, etc.) behind the scenes.
	insertFeed := `INSERT INTO feeds (id, name, description, link, created_at, updated_at) 
	VALUES (:id, :name, :description, :link, :createdAt, :updatedAt)`

	_, err := h.db.NamedExec(insertFeed, map[string]interface{}{
		"id":          uuid,
		"name":        feed.Title,
		"description": feed.Description,
		"link":        feed.Link,
		"createdAt":   time.Now(),
		"updatedAt":   time.Now(),
	})
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusInternalServerError, err)
	}

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
