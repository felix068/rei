package rei

import (
	"fmt"

	"github.com/mmcdole/gofeed"
)

type Reader struct {
	feedReader *gofeed.Parser
}

func NewParser() Reader {
	return Reader{gofeed.NewParser()}
}

func (r Reader) ReadFeed(url string) (*gofeed.Feed, error) {
	feed, err := r.feedReader.ParseURL(url)

	if err != nil {
		fmt.Println("Error on read feed", url, err)
	}

	return feed, nil

}
