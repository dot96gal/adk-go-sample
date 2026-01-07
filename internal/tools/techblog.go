package tools

import (
	"fmt"

	"github.com/mmcdole/gofeed"
	"google.golang.org/adk/tool"
)

type GetTechBlogEntryArgs struct{}

type GetTechBlogEntryResponse struct {
	Entries []Entry `json:"entries"`
}

func GetTechBlogEntry(ctx tool.Context, _ GetTechBlogEntryArgs) (GetTechBlogEntryResponse, error) {
	const url = "https://yamadashy.github.io/tech-blog-rss-feed/feeds/rss.xml"

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		return GetTechBlogEntryResponse{}, fmt.Errorf("can not parse feed: %s", url)
	}

	entries := make([]Entry, 0, len(feed.Items))
	for _, item := range feed.Items {
		entries = append(entries, Entry{
			Title: item.Title,
			URL:   item.Link,
		})
	}

	return GetTechBlogEntryResponse{
		Entries: entries,
	}, nil

}
