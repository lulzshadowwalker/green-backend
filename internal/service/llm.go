package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/lulzshadowwalker/green-backend/internal"
	"github.com/sashabaranov/go-openai"
)

type LLMService interface {
	StreamPlantAdvice(ctx context.Context, plant string, w io.Writer) error
}

type llmService struct {
	readingsStore SensorReadingsStore
	openaiClient  *openai.Client
}

func NewLLMService(readingsStore SensorReadingsStore, apiKey string) LLMService {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	client := openai.NewClient(apiKey)

	return &llmService{
		readingsStore: readingsStore,
		openaiClient:  client,
	}
}

func (s *llmService) StreamPlantAdvice(ctx context.Context, plant string, w io.Writer) error {

	since := time.Now().Add(-6 * time.Hour)
	readings, err := s.readingsStore.GetSensorReadingsSince(ctx, since)
	if err != nil {
		return fmt.Errorf("failed to fetch sensor readings: %w", err)
	}

	prompt := buildPrompt(plant, readings)

	req := openai.ChatCompletionRequest{
		Model: "gpt-3.5-turbo",
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert greenhouse assistant. Given the following sensor readings and plant type, provide actionable advice for optimal plant health. Be concise and practical.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
		Stream: true,
	}

	stream, err := s.openaiClient.CreateChatCompletionStream(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to create OpenAI stream: %w", err)
	}
	defer stream.Close()



	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error receiving from OpenAI stream: %w", err)
		}

		for _, choice := range response.Choices {
			content := choice.Delta.Content
			if content != "" {
				_, err := fmt.Fprint(w, content)
				if err != nil {
					return fmt.Errorf("error writing response: %w", err)
				}
			}
		}
	}



	return nil
}

func buildPrompt(plant string, readings []internal.SensorReading) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("Plant: %s\n", plant))
	b.WriteString("Sensor readings from the past 6 hours:\n")
	for _, r := range readings {
		b.WriteString(fmt.Sprintf("- [%s] %s: %.2f\n", r.Timestamp.Format(time.RFC3339), r.SensorType, r.Value))
	}
	b.WriteString("\nWhat advice do you have for optimal care of this plant, given these readings?")
	return b.String()
}


