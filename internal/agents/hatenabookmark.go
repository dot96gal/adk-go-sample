package agents

import (
	"github.com/dot96gal/adk-go-sample/internal/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

func NewHatenaBookmarkAgent(model model.LLM) (agent.Agent, error) {
	hatenaBookmarkTool, err := functiontool.New(
		functiontool.Config{
			Name:        "GetHatenaBookmarkEntry",
			Description: "はてなブックマークのホットエントリーを取得します",
		},
		tools.GetHatenaBookmarkEntry,
	)
	if err != nil {
		return nil, err
	}

	description := `はてなブックマークのホットエントリーを教えてくれるエージェントです。`
	instruction := `あなたははてなブックマークのホットエントリーを教えてくれる有能なエージェントです。
ツールを利用してはてなブックマークのホットエントリーを取得し、ユーザーに提供してください。

はてなブックマークのホットエントリーを取得するには、"GetHatenaBookmarkEntry"ツールを使用してください。
ユーザに提供する情報は下記のMarkdownのフォーマットに従ってください。

# はてなブックマーク
- [タイトル](URL)
`

	hatenaBookmarkAgent, err := llmagent.New(llmagent.Config{
		Name:        "hatena_bookmark_agent",
		Model:       model,
		Description: description,
		Instruction: instruction,
		Tools: []tool.Tool{
			hatenaBookmarkTool,
		},
	})
	if err != nil {
		return nil, err
	}

	return hatenaBookmarkAgent, nil
}
