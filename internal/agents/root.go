package agents

import (
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/tool"
)

func NewRootAgent(model model.LLM) (agent.Agent, error) {
	hatenaBookmarkAgent, err := NewHatenaBookmarkAgent(model)
	if err != nil {
		return nil, err
	}

	techBlogAgent, err := NewTechBlogAgent(model)
	if err != nil {
		return nil, err
	}

	description := `ユーザとおしゃべりするエージェントです。`
	instruction := `あなたはユーザとおしゃべりするエージェントです。

はてなブックマークのホットエントリーを取得する場合は"hatena_bookmark_agent"に処理を依頼してください。
テックブログのエントリーを取得する場合は"tech_blog_agent"に処理を依頼してください。
`

	rootAgent, err := llmagent.New(llmagent.Config{
		Name:        "root_agent",
		Model:       model,
		Description: description,
		Instruction: instruction,
		Tools:       []tool.Tool{},
		SubAgents: []agent.Agent{
			hatenaBookmarkAgent,
			techBlogAgent,
		},
	})
	if err != nil {
		return nil, err
	}

	return rootAgent, nil
}
