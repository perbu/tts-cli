package ai

import (
	"context"
	"fmt"
	"github.com/perbu/tts-cli/split"
	"github.com/sashabaranov/go-openai"
	"io"
	"net/http"
	"time"
)

type AI struct {
	client     *openai.Client
	httpClient *http.Client
}

func New(client *openai.Client) *AI {
	hc := &http.Client{
		Timeout: 10 * time.Second,
	}
	return &AI{client: client, httpClient: hc}
}

// Speech uses the OpenAI API to convert the input text to speech. The OpenAI API limits requests to
// 4096 bytes of text at a time. If the input text is longer than that, it will be split into multiple
// requests. The text will be split at paragraph boundaries, if possible. If the paragraph is too long, it
// will be split at sentence boundaries. If the sentence is too long, it will fail as the text is obviously
// garbage.
func (ai AI) Speech(ctx context.Context, content string) (io.ReadCloser, error) {
	inputs, err := split.SplitText(4000, content)
	if err != nil {
		return nil, fmt.Errorf("splitText: %w", err)
	}

	var pr, pw = io.Pipe()
	go func() {
		defer pw.Close()
		for _, chunk := range inputs {
			sp, err := ai.client.CreateSpeech(ctx, openai.CreateSpeechRequest{
				Model: openai.TTSModel1HD,
				Input: chunk,
				Voice: openai.VoiceAlloy,
				Speed: 1,
			})
			// Check if the context was cancelled
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			_, err = io.Copy(pw, sp.ReadCloser)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
		}
	}()
	return pr, nil
}

func (ai AI) Summary(ctx context.Context, content string) (string, error) {
	resp, err := ai.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Summarize the following text into one or two sentences. No more than 50 words in total.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: content,
				},
			},
		},
	)
	if err != nil {
		return "", fmt.Errorf("openai.CreateChatCompletion: %w", err)
	}
	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}
	return resp.Choices[0].Message.Content, nil
}

// Illustration generates an illustration for the document. It returns a byte slice, containing a PNG image.
func (ai AI) Illustration(ctx context.Context, content string) ([]byte, error) {
	resp, err := ai.client.CreateImage(ctx, openai.ImageRequest{
		Model:  openai.CreateImageModelDallE2,
		Prompt: "An illustration for a podcast with covering the following:\n" + content,
		Size:   openai.CreateImageSize512x512,
		N:      1,
	})
	if err != nil {
		return nil, fmt.Errorf("openai.CreateImage: %w", err)
	}
	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("no data in response")
	}
	contentUrl := resp.Data[0].URL
	// get the content of the URL
	resp2, err := ai.httpClient.Get(contentUrl)
	if err != nil {
		return nil, fmt.Errorf("httpClient.Get(content): %w", err)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("bad status code(payload): %d", resp2.StatusCode)
	}
	payload, err := io.ReadAll(resp2.Body)
	if err != nil {
		return nil, fmt.Errorf("io.ReadAll: %w", err)
	}
	return payload, nil

}
