package split

import (
	"reflect"
	"testing"
)

func TestSplitSentences(t *testing.T) {
	tests := []struct {
		name      string
		maxLength int
		input     string
		want      []string
		wantErr   bool
	}{
		{
			name:      "Basic split and merge",
			maxLength: 50,
			input:     "This is a sentence. This is another one.",
			want:      []string{"This is a sentence. This is another one."},
			wantErr:   false,
		},
		{
			name:      "Basic, can't merge",
			maxLength: 50,
			input:     "This is a sentence that is longer than 50 chars. This is another sentence.",
			want:      []string{"This is a sentence that is longer than 50 chars. ", "This is another sentence."},
			wantErr:   false,
		},
		{
			name:      "Single sentence",
			maxLength: 100,
			input:     "This is a single sentence that should not be split.",
			want:      []string{"This is a single sentence that should not be split."},
			wantErr:   false,
		},
		{
			name:      "Multiple punctuation",
			maxLength: 50,
			input:     "Hello! How are you? I'm fine. Thanks for asking.",
			want:      []string{"Hello! How are you? I'm fine. Thanks for asking."},
			wantErr:   false,
		},
		{
			name:      "Sentence too long",
			maxLength: 10,
			input:     "This sentence is definitely too long to fit within the maximum length.",
			wantErr:   true,
		},
		{
			name:      "Empty input",
			maxLength: 50,
			input:     "",
			want:      []string{""},
			wantErr:   false,
		},
		{
			name:      "Respect maxLength when merging",
			maxLength: 6,
			input:     "One. Two. 3. Four. Five. Six.",
			want:      []string{"One. ", "Two. 3. ", "Four. ", "Five. ", "Six."},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitSentences(tt.maxLength, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitSentences() error = '%v', wantErr '%v'", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitSentences() = '%#v', want '%#v'", got, tt.want)
			}
		})
	}
}

func TestSplitParagraphs(t *testing.T) {
	tests := []struct {
		name      string
		maxLength int
		input     string
		want      []string
		wantErr   bool
	}{
		{
			name:      "Basic split",
			maxLength: 50,
			input:     "This is paragraph one.\n\nThis is paragraph two.",
			want:      []string{"This is paragraph one.\n\nThis is paragraph two."},
			wantErr:   false,
		},
		{
			name:      "Single paragraph",
			maxLength: 100,
			input:     "This is a single paragraph that should not be split.",
			want:      []string{"This is a single paragraph that should not be split."},
			wantErr:   false,
		},
		{
			name:      "Merge short paragraphs",
			maxLength: 50,
			input:     "Short para 1.\n\nShort para 2.\n\nLonger paragraph that won't be merged.",
			want:      []string{"Short para 1.\nShort para 2.", "Longer paragraph that won't be merged."},
			wantErr:   false,
		},
		{
			name:      "Long paragraph split into sentences",
			maxLength: 35,
			input:     "This is a long paragraph. It should be split into sentences. Like this.",
			want:      []string{"This is a long paragraph. ", "It should be split into sentences. ", "Like this."},
			wantErr:   false,
		},
		{
			name:      "Empty input",
			maxLength: 50,
			input:     "",
			want:      []string{""},
			wantErr:   false,
		},
		{
			name:      "Input shorter than maxLength",
			maxLength: 100,
			input:     "Short input.",
			want:      []string{"Short input."},
			wantErr:   false,
		},
		{
			name:      "Multiple newlines",
			maxLength: 50,
			input:     "Para 1.\n\n\nPara 2.\n\n\n\nPara 3.",
			want:      []string{"Para 1.\n\n\nPara 2.\n\n\n\nPara 3."},
			wantErr:   false,
		},
		{
			name:      "Very long sentence",
			maxLength: 10,
			input:     "This sentence is definitely too long to fit within the maximum length.",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitParagraphs(tt.maxLength, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitParagraphs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				// If we expect an error, we don't need to check the output
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitParagraphs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestSplitText(t *testing.T) {
	tests := []struct {
		name      string
		maxLength int
		input     string
		want      []string
		wantErr   bool
	}{
		{
			name:      "Short input",
			maxLength: 50,
			input:     "This is a short input.",
			want:      []string{"This is a short input."},
			wantErr:   false,
		},
		{
			name:      "Multiple paragraphs",
			maxLength: 50,
			input:     "This is paragraph one.\n\nThis is paragraph two.\n\nThis is paragraph three.",
			want:      []string{"This is paragraph one.\nThis is paragraph two.", "This is paragraph three."},
			wantErr:   false,
		},
		{
			name:      "Long paragraph split into sentences",
			maxLength: 37,
			input:     "This is a long paragraph. It should be split into sentences. Like this.",
			want:      []string{"This is a long paragraph. ", "It should be split into sentences. ", "Like this."},
			wantErr:   false,
		},
		{
			name:      "Empty input",
			maxLength: 50,
			input:     "",
			want:      []string{""},
			wantErr:   false,
		},
		{
			name:      "Input exactly maxLength",
			maxLength: 19,
			input:     "Exactly maxLength.",
			want:      []string{"Exactly maxLength."},
			wantErr:   false,
		},
		{
			name:      "Very long sentence",
			maxLength: 10,
			input:     "This sentence is definitely too long to fit within the maximum length.",
			wantErr:   true,
		},
		{
			name:      "Mixed short and long paragraphs",
			maxLength: 48,
			input:     "Short.\n\nThis is a longer paragraph that needs splitting.\n\nAnother short one.",
			want:      []string{"Short.", "This is a longer paragraph that needs splitting.", "Another short one."},
			wantErr:   false,
		},
		{
			name:      "Unicode characters",
			maxLength: 27,
			input:     "こんにちは。世界。\n\nHello. World.",
			want:      []string{"こんにちは。世界。", "Hello. World."},
			wantErr:   false,
		},
		{
			name:      "Negative maxLength",
			maxLength: -1,
			input:     "Any input",
			wantErr:   true,
		},
		{
			name:      "Overflow",
			maxLength: 10,
			input:     "This sentence is too long and will make it all fail",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SplitText(tt.maxLength, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SplitText() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				// If we expect an error, we don't need to check the output
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitText() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
