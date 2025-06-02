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
	rss_link := c.FormValue("rss_link")
	feed, _ := h.reader.ReadFeed(rss_link)
	// name := feed.Title
	link := feed.Link

	//save in db the relevant content
	// author, link, description etc
	var feedq Feed
	h.db.Get(&feed, "SELECT * FROM feeds")

	fmt.Println("feeds", feedq.Name)

	return c.JSON(http.StatusOK, link)
}

func (h *Handler) RegisterFeedsRoutes(e *echo.Echo) {
	e.POST("/add_feed", h.addFeed)
}
