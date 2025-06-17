package rei

func Sync(url string, reader *Reader) {
	reader.ReadFeed(url)
}
