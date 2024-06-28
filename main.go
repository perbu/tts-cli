package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/perbu/tts-cli/ai"
	"github.com/perbu/tts-cli/feed"
	"github.com/sashabaranov/go-openai"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//go:embed .version
var embeddedVersion string

const (
	apiKeyEnvVar = "OPENAI_API_KEY"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	_ = godotenv.Load()
	if err := run(ctx, os.Stdout, os.Args, os.Environ()); err != nil {
		fmt.Println("run error: ", err)
		os.Exit(1)
	}
	fmt.Println("clean exit")
}

func run(ctx context.Context, stdout *os.File, args []string, env []string) error {
	debugFlag := flag.Bool("d", false, "Enable debug output")
	versionFlag := flag.Bool("v", false, "Print version and exit")
	if err := flag.CommandLine.Parse(args[1:]); err != nil {
		return fmt.Errorf("flag.CommandLine.Parse: %w", err)
	}
	if *versionFlag {
		fmt.Println(embeddedVersion)
		return nil
	}
	if *debugFlag {
		fmt.Fprintln(stdout, "debug output enabled")
	}
	// Input file:
	if flag.NArg() != 1 {
		return fmt.Errorf("usage: %s <input-file>", args[0])
	}
	workingDir := flag.Arg(0)
	if workingDir == "" {
		workingDir = "."
	}
	logger := makeLogger(*debugFlag)
	apiKey := getEnvStr(env, apiKeyEnvVar, "")
	if apiKey == "" {
		return fmt.Errorf("'%s' is required", apiKeyEnvVar)
	}
	c := ai.New(openai.NewClient(apiKey))
	// Create a new FeedManager
	fm := feed.New(c, logger)
	// Scan the directory for txt files and create missing elements
	err := fm.Scan(ctx, workingDir)
	if err != nil {
		return fmt.Errorf("fm.Scan: %w", err)
	}
	// Generate the RSS feed
	rss, err := fm.GenerateRSS()
	if err != nil {
		return fmt.Errorf("fm.GenerateRSS: %w", err)
	}
	fmt.Println(rss)
	return nil

}

func getEnvStr(env []string, key, defaultValue string) string {
	for _, e := range env {
		pair := strings.Split(e, "=")
		if len(pair) != 2 {
			continue
		}
		if pair[0] == key {
			return pair[1]
		}
	}
	return defaultValue
}

func makeLogger(debug bool) *slog.Logger {
	level := slog.LevelInfo
	if debug {
		level = slog.LevelDebug
	}
	fh := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})
	return slog.New(fh)
}
