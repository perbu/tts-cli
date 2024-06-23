# tts-cli

This is a simple CLI wrapper around the OpenAI Text-to-Speech API. There where already a few of these, but they were
written in Python. Dependency management in Python is a nightmare, so I decided to write my own in Go.

I use this, in addition to `make` and `dropcaster` to generate a podcast out of a bunch of text-files so I can "read" 
them while biking or walking.

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

