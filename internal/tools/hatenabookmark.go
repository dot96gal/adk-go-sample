package tools

import (
	"fmt"

	"github.com/mmcdole/gofeed"
	"google.golang.org/adk/tool"
)

type GetHatenaBookmarkEntryArgs struct{}

type GetHatenaBookmarkEntryResponse struct {
	Entries []Entry `json:"entries"`
}

func GetHatenaBookmarkEntry(ctx tool.Context, _ GetHatenaBookmarkEntryArgs) (GetHatenaBookmarkEntryResponse, error) {
	const url = "https://b.hatena.ne.jp/hotentry.rss"

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		return GetHatenaBookmarkEntryResponse{}, fmt.Errorf("can not parse feed: %s", url)
	}

	entries := make([]Entry, 0, len(feed.Items))
	for _, item := range feed.Items {
		entries = append(entries, Entry{
			Title: item.Title,
			URL:   item.Link,
		})
	}

	return GetHatenaBookmarkEntryResponse{
		Entries: entries,
	}, nil
}
