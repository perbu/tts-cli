package feed

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path"
	"testing"
)

func TestFeed_Scan(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.Chdir(tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	mAI := &MockAI{
		logger: slog.Default(),
	}
	feed := New(mAI, debugLogger())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Create a few files
	err = os.WriteFile(path.Join(tmpDir, "file1.txt"), []byte("This is the first test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = os.WriteFile(path.Join(tmpDir, "file2.txt"), []byte("This is another test"), 0644)
	if err != nil {
		t.Fatal(err)
	}
	err = feed.Scan(ctx, tmpDir)
	if err != nil {
		t.Fatal(err)
	}
	rss, err := feed.GenerateRSS()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(rss)
}

type MockAI struct {
	logger *slog.Logger
}

func (m *MockAI) Summary(ctx context.Context, content string) (string, error) {
	m.logger.Info("summary called", "content", content)
	return "very short summary", nil
}

func (m *MockAI) Illustration(ctx context.Context, content string) ([]byte, error) {
	m.logger.Info("illustration called", "content", content)
	return nil, nil
}

func (m *MockAI) Speech(ctx context.Context, content string) (io.ReadCloser, error) {
	m.logger.Info("speech called", "content", content)
	buf := bytes.NewBuffer([]byte("this is kinda speechy bytes"))
	return io.NopCloser(buf), nil
}

func debugLogger() *slog.Logger {
	lh := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	return slog.New(lh)
}
