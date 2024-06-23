package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/perbu/tts-cli/tts"
	"github.com/sashabaranov/go-openai"
	"io"
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
	outputFlag := flag.String("o", "", "Output file")
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
	inputFileName := flag.Arg(0)
	// Output file is the same as the input file, but with a .aac extension:
	if *outputFlag == "" {
		*outputFlag = inputFileName + ".aac"
	}

	// Slurp the input file:
	input, err := os.ReadFile(inputFileName)
	if err != nil {
		return fmt.Errorf("os.ReadFile: %w", err)
	}
	if *debugFlag {
		fmt.Fprintf(stdout, "read %d bytes from %s\n", len(input), inputFileName)
	}

	apiKey := getEnvStr(env, apiKeyEnvVar, "")
	if apiKey == "" {
		return fmt.Errorf("'%s' is required", apiKeyEnvVar)
	}
	c := openai.NewClient(apiKey)

	// Open the output for writing:
	output, err := os.OpenFile(*outputFlag, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}
	defer output.Close()
	if *debugFlag {
		fmt.Fprintf(stdout, "opened %s for writing\n", *outputFlag)
	}
	ttsOutput, err := tts.Speech(ctx, c, string(input), *debugFlag)
	if err != nil {
		return fmt.Errorf("tts.Speech: %w", err)
	}
	// The below could be replaced with io.Copy(output, ttsOutput), but we want to print debug output:
	for {
		buf := make([]byte, 65536)
		n, err := ttsOutput.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("ttsOutput.Read: %w", err)
		}
		if n == 0 {
			break
		}
		n, err = output.Write(buf[:n])
		if err != nil {
			return fmt.Errorf("output.Write: %w", err)
		}
		if *debugFlag {
			fmt.Fprintf(stdout, "wrote %d bytes to %s\n", n, *outputFlag)
		}
	}
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
