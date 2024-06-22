package tts

import (
	"context"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/sashabaranov/go-openai"
	"io"
)

func Speech(ctx context.Context, c *openai.Client, input string, debug bool) (io.ReadCloser, error) {
	info := whatlanggo.Detect(input)
	if debug {
		fmt.Printf("Detected language: %v\n", info.Lang)
	}
	sp, err := c.CreateSpeech(ctx, openai.CreateSpeechRequest{
		Model: openai.TTSModel1HD,
		Input: input,
		Voice: openai.VoiceAlloy,
		Speed: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("openai.CreateSpeech: %w", err)
	}
	return sp.ReadCloser, nil
}
