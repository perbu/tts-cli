package split

import (
	"fmt"
	"regexp"
)

// SplitText splits the input text into chunks of at most maxLength bytes. The text will be split at paragraph
// boundaries, if possible. If the paragraph is too long, it will be split at sentence boundaries. If the
// sentence is too long, it will return an error as the text is considered invalid for this operation.
func SplitText(maxLength int, input string) ([]string, error) {
	// If the input is short enough, return it as is:
	if len(input) <= maxLength {
		return []string{input}, nil
	}
	// Split at paragraph boundaries:
	paragraphs, err := SplitParagraphs(maxLength, input)
	if err != nil {
		return nil, fmt.Errorf("splitting paragraphs: %w", err)
	}
	for _, p := range paragraphs {
		if len(p) > maxLength {
			return nil, fmt.Errorf("chunk too long: %d bytes", len(p))
		}
	}
	return paragraphs, nil
}

// SplitParagraphs splits the input text into paragraphs. A paragraph is defined as a sequence of
// characters that ends with two or more newlines. The newlines will be included in the output
// except for the last paragraph.
func SplitParagraphs(maxLength int, input string) ([]string, error) {
	if len(input) <= maxLength {
		return []string{input}, nil
	}
	re := regexp.MustCompile(`\n\n+`) // Match two or more newlines

	paragraphs := re.Split(input, -1)
	result := make([]string, 0, len(paragraphs))

	for i := 0; i < len(paragraphs); i++ {
		p := paragraphs[i]
		if len(p) > maxLength {
			sentences, err := SplitSentences(maxLength, p)
			if err != nil {
				return nil, fmt.Errorf("splitting sentences: %w", err)
			}
			result = append(result, sentences...)
		} else {
			// Try to merge with subsequent paragraphs
			for i+1 < len(paragraphs) && len(p+"\n\n"+paragraphs[i+1]) <= maxLength {
				p += "\n" + paragraphs[i+1]
				i++
			}
			result = append(result, p)
		}
	}
	return result, nil
}

// SplitSentences splits the input text into sentences. A sentence is defined as a sequence of
// characters that ends with a period, exclamation mark, or question mark, followed by a space or end of input.
// The ending punctuation/space will be included in the output.
func SplitSentences(maxLength int, input string) ([]string, error) {
	re := regexp.MustCompile(`(\.|\!|\?)\s+`)
	sentences := re.Split(input, -1)
	result := make([]string, 0, len(sentences))

	for i := 0; i < len(sentences); i++ {
		s := sentences[i]
		if i < len(sentences)-1 {
			s += string(re.Find([]byte(input[len(s):])))
		}
		if len(s) > maxLength {
			return nil, fmt.Errorf("sentence too long: %d bytes", len(s))
		}
		// Try to merge with subsequent sentences
		for i+1 < len(sentences) && len(s+sentences[i+1]) <= maxLength {
			s += sentences[i+1]
			if i+1 < len(sentences)-1 {
				s += string(re.Find([]byte(input[len(s):])))
			}
			i++
		}
		result = append(result, s)
	}

	return result, nil
}
