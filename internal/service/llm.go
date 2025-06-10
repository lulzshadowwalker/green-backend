package service

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
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

const (
	// Maximum number of sensor readings to include in the prompt
	maxReadings = 10
	// Rough estimate of max tokens for the prompt (leaving room for response)
	maxPromptTokens = 2000
	// Average tokens per character (rough estimate)
	tokensPerChar = 0.25
)

func NewLLMService(readingsStore SensorReadingsStore, apiKey string) LLMService {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	// TEMPORARY: Skip SSL verification for testing
	// TODO: Remove this in production once CA certificates are properly configured
	config := openai.DefaultConfig(apiKey)
	config.HTTPClient = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	client := openai.NewClientWithConfig(config)

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

	// Limit and optimize readings for the prompt
	limitedReadings := limitReadings(readings)
	prompt := buildPrompt(plant, limitedReadings)

	// Log token usage for monitoring
	estimatedTokens := int(float64(len(prompt)) * tokensPerChar)
	log.Printf("LLM request: plant=%s, original_readings=%d, limited_readings=%d, estimated_tokens=%d",
		plant, len(readings), len(limitedReadings), estimatedTokens)

	req := openai.ChatCompletionRequest{
		Model:     "gpt-3.5-turbo",
		MaxTokens: 1000, // Limit response tokens to control costs
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are an expert greenhouse assistant. Given the following sensor readings and plant type, provide actionable advice for optimal plant health. Be concise and practical. Keep in mind, you are providing this advice to a simple farmer who is likely not to be very technical. Keep the language friendly and easy to understand without sacrificing accuracy. Also, keep in mind that you cannot use rich text formatting in your responses.",
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

	if len(readings) == 0 {
		b.WriteString("No recent sensor readings available.\n")
	} else {
		b.WriteString(fmt.Sprintf("Recent sensor readings (%d readings):\n", len(readings)))
		for _, r := range readings {
			b.WriteString(fmt.Sprintf("- %s: %.1f (%s)\n",
				r.SensorType, r.Value, r.Timestamp.Format("15:04")))
		}
	}

	b.WriteString("\nWhat advice do you have for optimal care of this plant, given these readings?")
	return b.String()
}

// limitReadings reduces the number of readings to prevent API token limits
func limitReadings(readings []internal.SensorReading) []internal.SensorReading {
	if len(readings) == 0 {
		return readings
	}

	// Sort readings by timestamp (most recent first)
	sort.Slice(readings, func(i, j int) bool {
		return readings[i].Timestamp.After(readings[j].Timestamp)
	})

	// Group readings by sensor type to ensure we have recent data for each sensor
	sensorGroups := make(map[string][]internal.SensorReading)
	for _, reading := range readings {
		sensorGroups[reading.SensorType] = append(sensorGroups[reading.SensorType], reading)
	}

	var limitedReadings []internal.SensorReading
	readingsPerSensor := maxReadings / len(sensorGroups)
	if readingsPerSensor < 1 {
		readingsPerSensor = 1
	}

	// Take the most recent readings from each sensor type
	for _, sensorReadings := range sensorGroups {
		count := readingsPerSensor
		if count > len(sensorReadings) {
			count = len(sensorReadings)
		}
		limitedReadings = append(limitedReadings, sensorReadings[:count]...)
	}

	// Sort the final result by timestamp again
	sort.Slice(limitedReadings, func(i, j int) bool {
		return limitedReadings[i].Timestamp.After(limitedReadings[j].Timestamp)
	})

	// Further limit if still too many
	if len(limitedReadings) > maxReadings {
		limitedReadings = limitedReadings[:maxReadings]
	}

	// Check if the prompt would be too long
	testPrompt := buildPrompt("test", limitedReadings)
	estimatedTokens := int(float64(len(testPrompt)) * tokensPerChar)

	// If still too long, reduce further by taking every nth reading
	originalCount := len(limitedReadings)
	for estimatedTokens > maxPromptTokens && len(limitedReadings) > 10 {
		// Take every other reading
		var reduced []internal.SensorReading
		for i := 0; i < len(limitedReadings); i += 2 {
			reduced = append(reduced, limitedReadings[i])
		}
		limitedReadings = reduced
		testPrompt = buildPrompt("test", limitedReadings)
		estimatedTokens = int(float64(len(testPrompt)) * tokensPerChar)
	}

	// Log if we had to reduce further due to token limits
	if len(limitedReadings) < originalCount {
		log.Printf("Further reduced readings due to token limit: %d -> %d (estimated tokens: %d)",
			originalCount, len(limitedReadings), estimatedTokens)
	}

	return limitedReadings
}
