package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/mmcdole/gofeed"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
	"google.golang.org/genai"
)

type getHatenaBookmarkEntryArgs struct{}

type getHatenaBookmarkEntryResponse struct {
	Entries []Entry `json:"entries"`
}

type Entry struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func getHatenaBookmarkEntry(ctx tool.Context, _ getHatenaBookmarkEntryArgs) (getHatenaBookmarkEntryResponse, error) {
	const url = "https://b.hatena.ne.jp/hotentry.rss"

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		return getHatenaBookmarkEntryResponse{}, fmt.Errorf("can not parse feed: %s", url)
	}

	entries := make([]Entry, 0, len(feed.Items))
	for _, item := range feed.Items {
		entries = append(entries, Entry{
			Title: item.Title,
			URL:   item.Link,
		})
	}

	return getHatenaBookmarkEntryResponse{
		Entries: entries,
	}, nil
}

type getTechBlogEntryArgs struct{}

type getTechBlogEntryResponse struct {
	Entries []Entry `json:"entries"`
}

func getTechBlogEntry(ctx tool.Context, _ getTechBlogEntryArgs) (getTechBlogEntryResponse, error) {
	const url = "https://yamadashy.github.io/tech-blog-rss-feed/feeds/rss.xml"

	parser := gofeed.NewParser()
	feed, err := parser.ParseURL(url)
	if err != nil {
		return getTechBlogEntryResponse{}, fmt.Errorf("can not parse feed: %s", url)
	}

	entries := make([]Entry, 0, len(feed.Items))
	for _, item := range feed.Items {
		entries = append(entries, Entry{
			Title: item.Title,
			URL:   item.Link,
		})
	}

	return getTechBlogEntryResponse{
		Entries: entries,
	}, nil

}

func main() {
	ctx := context.Background()

	modelName := os.Getenv("GEMINI_MODEL")
	apikey := os.Getenv("GOOGLE_API_KEY")

	model, err := gemini.NewModel(
		ctx,
		modelName,
		&genai.ClientConfig{
			APIKey: apikey,
		},
	)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	hatenaBookmarkTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_hatena_bookmark_entry",
			Description: "Retrieves hatena bookmark entries",
		},
		getHatenaBookmarkEntry,
	)
	if err != nil {
		log.Fatal(err)
	}

	techBlogTool, err := functiontool.New(
		functiontool.Config{
			Name:        "get_tech_blog_entry",
			Description: "Retrieves tech blog entries",
		},
		getTechBlogEntry,
	)
	if err != nil {
		log.Fatal(err)
	}

	description := `ニュースを教えてくれるエージェントです。`
	instruction := `あなたはニュースを教えてくれる有能なアシスタントです。
ツールを利用してニュースを取得し、ユーザーに提供してください。

はてなブックマークのホットエントリーを取得するには、"get_hatena_bookmark_entry"ツールを使用してください。
はてなブックマークのホットエントリーを取得したら、タイトルから技術系の記事かどうかを判定し、技術系の記事のみをユーザに提供してください。

技術系ブログのエントリーを取得するには、"get_tech_blog_entry"ツールを使用してください。
技術系ブログのエントリーを取得したら、タイトルから技術系の記事かどうかを判定し、技術系の記事のみをユーザに提供してください。


ユーザに提供するニュースは下記のMarkdownのフォーマットに従ってください。

# はてなブックマーク
- [タイトル](URL)

# 技術系ブログ
- [タイトル](URL)
`

	hatenaBookmarkAgent, err := llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       model,
		Description: description,
		Instruction: instruction,
		Tools: []tool.Tool{
			// geminitool.GoogleSearch{},
			hatenaBookmarkTool,
			techBlogTool,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(hatenaBookmarkAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
