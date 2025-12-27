package main

import (
	"context"
	"log"
	"os"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
	"google.golang.org/adk/tool"
	"google.golang.org/adk/tool/geminitool"
	"google.golang.org/genai"
)

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

	description := "特定の都市の現在の時刻を教えてくれるエージェントです。"
	instruction := "あなたは特定の都市の現在の時刻を教えてくれる有能なアシスタントです。"

	timeAgent, err := llmagent.New(llmagent.Config{
		Name:        "hello_time_agent",
		Model:       model,
		Description: description,
		Instruction: instruction,
		Tools: []tool.Tool{
			geminitool.GoogleSearch{},
		},
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(timeAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
