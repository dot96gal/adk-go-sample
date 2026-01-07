package main

import (
	"context"
	"log"
	"os"

	"github.com/dot96gal/adk-go-sample/internal/agents"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	"google.golang.org/adk/model/gemini"
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

	rootAgent, err := agents.NewRootAgent(model)
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(rootAgent),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
