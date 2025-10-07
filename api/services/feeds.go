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
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type Post struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Link        string    `json:"link" db:"link"`
	FeedId      uuid.UUID `json:"feedId" db:"feed_id"`
	FeedName    string    `json:"feedName" db:"feed_name"`
	IsRead      bool      `json:"isRead" db:"is_read"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
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

	if readFeedErr != nil || feed == nil {
		fmt.Println("Error reading feed:", readFeedErr)
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to fetch RSS feed. Please check the URL is correct and accessible.",
		})
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

	// Add all existing posts from the feed
	for _, item := range feed.Items {
		_, err := h.db.NamedExec(`INSERT INTO posts (id, name, description, link, feed_id, is_read, created_at, updated_at)
			VALUES (:id, :name, :description, :link, :feed_id, :is_read, :createdAt, :updatedAt)`,
			map[string]interface{}{
				"id":          uuid.New(),
				"name":        item.Title,
				"description": item.Description,
				"link":        item.Link,
				"feed_id":     feedUuid,
				"is_read":     false,
				"createdAt":   time.Now(),
				"updatedAt":   time.Now(),
			})

		if err != nil {
			fmt.Println("Error inserting post:", err)
			// Continue with other posts even if one fails
		}
	}

	return c.JSON(http.StatusCreated, feed.Link)
}

func (h *Handler) listFeeds(c echo.Context) error {
	var feeds []Feed
	err := h.db.Select(&feeds, "SELECT * FROM feeds ORDER BY created_at DESC")

	if err != nil {
		fmt.Println("Error listing feeds:", err)
		return c.JSON(http.StatusInternalServerError, "Error listing feeds")
	}

	return c.JSON(http.StatusOK, feeds)
}

func (h *Handler) listPosts(c echo.Context) error {
	var posts []Post
	err := h.db.Select(&posts, `
		SELECT p.*, f.name as feed_name
		FROM posts p
		JOIN feeds f ON p.feed_id = f.id
		ORDER BY p.created_at DESC
	`)

	if err != nil {
		fmt.Println("Error listing posts:", err)
		return c.JSON(http.StatusInternalServerError, "Error listing posts")
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) listUnreadPosts(c echo.Context) error {
	var posts []Post
	err := h.db.Select(&posts, `
		SELECT p.*, f.name as feed_name
		FROM posts p
		JOIN feeds f ON p.feed_id = f.id
		WHERE p.is_read = false
		ORDER BY p.created_at DESC
	`)

	if err != nil {
		fmt.Println("Error listing unread posts:", err)
		return c.JSON(http.StatusInternalServerError, "Error listing unread posts")
	}

	return c.JSON(http.StatusOK, posts)
}

func (h *Handler) markPostAsRead(c echo.Context) error {
	postId := c.Param("id")

	_, err := h.db.Exec("UPDATE posts SET is_read = true, updated_at = $1 WHERE id = $2", time.Now(), postId)

	if err != nil {
		fmt.Println("Error marking post as read:", err)
		return c.JSON(http.StatusInternalServerError, "Error marking post as read")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

func (h *Handler) deleteFeed(c echo.Context) error {
	feedId := c.Param("id")

	// Delete all posts from this feed first
	_, err := h.db.Exec("DELETE FROM posts WHERE feed_id = $1", feedId)
	if err != nil {
		fmt.Println("Error deleting posts:", err)
		return c.JSON(http.StatusInternalServerError, "Error deleting posts")
	}

	// Then delete the feed
	_, err = h.db.Exec("DELETE FROM feeds WHERE id = $1", feedId)
	if err != nil {
		fmt.Println("Error deleting feed:", err)
		return c.JSON(http.StatusInternalServerError, "Error deleting feed")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "success"})
}

// SyncFeed fetches new posts from a feed
func (h *Handler) SyncFeed(feedId uuid.UUID, rssLink string) error {
	feed, err := h.reader.ReadFeed(rssLink)
	if err != nil {
		fmt.Println("Error reading feed during sync:", err)
		return err
	}

	for _, item := range feed.Items {
		// Check if post already exists
		var exists int
		checkErr := h.db.Get(&exists, "SELECT 1 FROM posts WHERE link = $1 LIMIT 1", item.Link)

		if checkErr == sql.ErrNoRows {
			// Post doesn't exist, insert it
			_, insertErr := h.db.NamedExec(`INSERT INTO posts (id, name, description, link, feed_id, is_read, created_at, updated_at)
				VALUES (:id, :name, :description, :link, :feed_id, :is_read, :createdAt, :updatedAt)`,
				map[string]interface{}{
					"id":          uuid.New(),
					"name":        item.Title,
					"description": item.Description,
					"link":        item.Link,
					"feed_id":     feedId,
					"is_read":     false,
					"createdAt":   time.Now(),
					"updatedAt":   time.Now(),
				})

			if insertErr != nil {
				fmt.Println("Error inserting new post during sync:", insertErr)
			}
		}
	}

	// Update feed's updated_at timestamp
	h.db.Exec("UPDATE feeds SET updated_at = $1 WHERE id = $2", time.Now(), feedId)

	return nil
}

// SyncAllFeeds syncs all feeds in the database
func (h *Handler) SyncAllFeeds() {
	var feeds []Feed
	err := h.db.Select(&feeds, "SELECT * FROM feeds")

	if err != nil {
		fmt.Println("Error fetching feeds for sync:", err)
		return
	}

	for _, feed := range feeds {
		fmt.Printf("Syncing feed: %s (%s)\n", feed.Name, feed.RssLink)
		h.SyncFeed(feed.Id, feed.RssLink)
	}

	fmt.Println("Sync completed for all feeds")
}

func (h *Handler) RegisterFeedsRoutes(e *echo.Echo) {
	e.POST("/add_feed", h.addFeed)
	e.GET("/list_feeds", h.listFeeds)
	e.GET("/list_posts", h.listPosts)
	e.GET("/list_unread_posts", h.listUnreadPosts)
	e.PUT("/posts/:id/read", h.markPostAsRead)
	e.DELETE("/feeds/:id", h.deleteFeed)
}
