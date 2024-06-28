# AI Podcaster tool

This a tool to generate a podcast from a bunch of text files. It uses the OpenAI Text-to-Speech API to 
generate the audio, ChatGPT to generate a summary and Dall-E to generate an image for the podcast.

## Installation

```bash
go install github.com/perbu/tts-cli@latest
```

## Usage

```bash
tts-cli [-debug] <input file>
```

This will read the input file, send it to the OpenAI API, and write the resulting audio to an mp3 file with the same
name as the input file, but with ".mp3" appended.

## Makefile for mp3 generation

```makefile
# Find all .txt files in the source directory
TXT_FILES := $(wildcard source/*.txt)

# Generate the list of target .aac files in the current directory
AUDIO_FILES := $(notdir $(TXT_FILES:.txt=.mp3))

# Default target
all: $(AUDIO_FILES)

# Rule to create .aac file from .txt file
%.mp3: source/%.txt
	tts-cli -o $@ $<

# Clean target to remove all generated audio files
clean:
	rm -f *.aac

.PHONY: all clean
```

