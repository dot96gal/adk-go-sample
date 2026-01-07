package agents

import (
	"github.com/dot96gal/adk-go-sample/internal/tools"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
)

func NewTechBlogAgent(model model.LLM) (agent.Agent, error) {
	techBlogTool, err := functiontool.New(
		functiontool.Config{
			Name:        "GetTechBlogEntry",
			Description: "テックブログのエントリーを取得します",
		},
		tools.GetTechBlogEntry,
	)
	if err != nil {
		return nil, err
	}

	description := `テックブログのエントリーを教えてくれるエージェントです。`
	instruction := `あなたはテックブログのエントリーを教えてくれる有能なエージェントです。
ツールを利用してテックブログのエントリーを取得し、ユーザーに提供してください。

テックブログのエントリーを取得するには、"GetTechBlogEntry"ツールを使用してください。
ユーザに提供する情報は下記のMarkdownのフォーマットに従ってください。

# テックブログ
- [タイトル](URL)
`

	techBlogAgent, err := llmagent.New(llmagent.Config{
		Name:        "tech_blog_agent",
		Model:       model,
		Description: description,
		Instruction: instruction,
		Tools: []tool.Tool{
			techBlogTool,
		},
	})
	if err != nil {
		return nil, err
	}

	return techBlogAgent, nil
}
