# tts-cli

This is a simple CLI wrapper around the OpenAI Text-to-Speech API. There where already a few of these, but they were
written in Python. Dependency management in Python is a nightmare, so I decided to write my own in Go.

## Installation

```bash
go install github.com/perbu/tts-cli@latest
```

## Usage

```bash
tts-cli [-debug] <input file>
```

This will read the input file, send it to the OpenAI API, and write the resulting audio to an AAC file with the same
name as the input file, but with ".aac" appended.

