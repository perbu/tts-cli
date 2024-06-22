package tts

import (
	"context"
	"fmt"
	"github.com/perbu/tts-cli/split"
	"github.com/sashabaranov/go-openai"
	"io"
)

// Speech uses the OpenAI API to convert the input text to speech. The OpenAI API limits requests to
// 4096 bytes of text at a time. If the input text is longer than that, it will be split into multiple
// requests. The text will be split at paragraph boundaries, if possible. If the paragraph is too long, it
// will be split at sentence boundaries. If the sentence is too long, it will fail as the text is obviously
// garbage.
func Speech(ctx context.Context, c *openai.Client, input string, debug bool) (io.ReadCloser, error) {
	inputs, err := split.SplitText(4000, input)
	if err != nil {
		return nil, fmt.Errorf("splitText: %w", err)
	}
	if debug {
		fmt.Println("Split into", len(inputs), "chunks")
	}
	// make a io.Pipe to stream the output
	var pr, pw = io.Pipe()
	go func() {
		defer pw.Close()
		for i, chunk := range inputs {
			if debug {
				fmt.Println("Chunk", i, ":", len(chunk), "bytes")
			}
			sp, err := c.CreateSpeech(ctx, openai.CreateSpeechRequest{
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
		}
	}()
	return pr, nil
}
