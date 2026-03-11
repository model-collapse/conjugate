//go:build !cgo_analysis
// +build !cgo_analysis

package diagon

import (
	"regexp"
	"strings"
	"unicode"
)

// Analyzer performs text analysis, converting text into tokens.
// This is a pure Go implementation that works without CGO.
type Analyzer struct {
	name      string
	tokenizer func(string) []string
	filters   []func(string) string
	stopWords map[string]bool
}

// Token represents a single analyzed token with position information.
type Token struct {
	Text        string
	Position    int
	StartOffset int
	EndOffset   int
	Type        string
}

// Common English stop words
var defaultStopWords = map[string]bool{
	"a": true, "an": true, "and": true, "are": true, "as": true, "at": true,
	"be": true, "but": true, "by": true, "for": true, "if": true, "in": true,
	"into": true, "is": true, "it": true, "no": true, "not": true, "of": true,
	"on": true, "or": true, "such": true, "that": true, "the": true, "their": true,
	"then": true, "there": true, "these": true, "they": true, "this": true,
	"to": true, "was": true, "will": true, "with": true,
}

// Standard tokenizer - splits on non-letter/digit boundaries
var standardTokenizerRegex = regexp.MustCompile(`[\p{L}\p{N}]+`)

func standardTokenizer(text string) []string {
	return standardTokenizerRegex.FindAllString(text, -1)
}

// Whitespace tokenizer - splits on whitespace only
func whitespaceTokenizer(text string) []string {
	return strings.Fields(text)
}

// Keyword tokenizer - returns entire input as single token
func keywordTokenizer(text string) []string {
	if text == "" {
		return nil
	}
	return []string{text}
}

// Lowercase filter
func lowercaseFilter(s string) string {
	return strings.ToLower(s)
}

// ASCII folding filter - converts accented characters to ASCII equivalents
func asciiFoldingFilter(s string) string {
	var result strings.Builder
	for _, r := range s {
		switch {
		case r >= 'À' && r <= 'Å':
			result.WriteRune('A')
		case r >= 'à' && r <= 'å':
			result.WriteRune('a')
		case r == 'Ç':
			result.WriteRune('C')
		case r == 'ç':
			result.WriteRune('c')
		case r >= 'È' && r <= 'Ë':
			result.WriteRune('E')
		case r >= 'è' && r <= 'ë':
			result.WriteRune('e')
		case r >= 'Ì' && r <= 'Ï':
			result.WriteRune('I')
		case r >= 'ì' && r <= 'ï':
			result.WriteRune('i')
		case r == 'Ñ':
			result.WriteRune('N')
		case r == 'ñ':
			result.WriteRune('n')
		case r >= 'Ò' && r <= 'Ö':
			result.WriteRune('O')
		case r >= 'ò' && r <= 'ö':
			result.WriteRune('o')
		case r >= 'Ù' && r <= 'Ü':
			result.WriteRune('U')
		case r >= 'ù' && r <= 'ü':
			result.WriteRune('u')
		case r == 'Ý' || r == 'Ÿ':
			result.WriteRune('Y')
		case r == 'ý' || r == 'ÿ':
			result.WriteRune('y')
		default:
			result.WriteRune(r)
		}
	}
	return result.String()
}

// NewStandardAnalyzer creates a standard analyzer (standard tokenizer + lowercase + stop).
func NewStandardAnalyzer() (*Analyzer, error) {
	return &Analyzer{
		name:      "standard",
		tokenizer: standardTokenizer,
		filters:   []func(string) string{lowercaseFilter},
		stopWords: defaultStopWords,
	}, nil
}

// NewSimpleAnalyzer creates a simple analyzer (whitespace tokenizer + lowercase).
func NewSimpleAnalyzer() (*Analyzer, error) {
	return &Analyzer{
		name:      "simple",
		tokenizer: whitespaceTokenizer,
		filters:   []func(string) string{lowercaseFilter},
	}, nil
}

// NewWhitespaceAnalyzer creates a whitespace analyzer.
func NewWhitespaceAnalyzer() (*Analyzer, error) {
	return &Analyzer{
		name:      "whitespace",
		tokenizer: whitespaceTokenizer,
	}, nil
}

// NewKeywordAnalyzer creates a keyword analyzer (no tokenization).
func NewKeywordAnalyzer() (*Analyzer, error) {
	return &Analyzer{
		name:      "keyword",
		tokenizer: keywordTokenizer,
	}, nil
}

// NewChineseAnalyzer creates a Chinese analyzer.
// Note: This pure Go implementation uses simple Unicode-based segmentation.
// For production Chinese analysis, consider using a proper segmentation library.
func NewChineseAnalyzer(dictPath string) (*Analyzer, error) {
	// Simple Chinese tokenizer - splits on CJK character boundaries
	chineseTokenizer := func(text string) []string {
		var tokens []string
		var currentToken strings.Builder

		for _, r := range text {
			if unicode.Is(unicode.Han, r) {
				// CJK characters are individual tokens
				if currentToken.Len() > 0 {
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
				tokens = append(tokens, string(r))
			} else if unicode.IsLetter(r) || unicode.IsNumber(r) {
				currentToken.WriteRune(r)
			} else {
				if currentToken.Len() > 0 {
					tokens = append(tokens, currentToken.String())
					currentToken.Reset()
				}
			}
		}
		if currentToken.Len() > 0 {
			tokens = append(tokens, currentToken.String())
		}
		return tokens
	}

	return &Analyzer{
		name:      "chinese",
		tokenizer: chineseTokenizer,
		filters:   []func(string) string{lowercaseFilter},
	}, nil
}

// NewEnglishAnalyzer creates an English analyzer (standard + lowercase + ascii folding + stop).
func NewEnglishAnalyzer() (*Analyzer, error) {
	return &Analyzer{
		name:      "english",
		tokenizer: standardTokenizer,
		filters:   []func(string) string{lowercaseFilter, asciiFoldingFilter},
		stopWords: defaultStopWords,
	}, nil
}

// NewMultilingualAnalyzer creates a multilingual analyzer (standard + lowercase + ascii folding).
func NewMultilingualAnalyzer() (*Analyzer, error) {
	return &Analyzer{
		name:      "multilingual",
		tokenizer: standardTokenizer,
		filters:   []func(string) string{lowercaseFilter, asciiFoldingFilter},
	}, nil
}

// NewSearchAnalyzer creates a search analyzer optimized for queries.
func NewSearchAnalyzer() (*Analyzer, error) {
	return &Analyzer{
		name:      "search",
		tokenizer: standardTokenizer,
		filters:   []func(string) string{lowercaseFilter},
	}, nil
}

// NewAnalyzer creates an analyzer by name.
func NewAnalyzer(name string) (*Analyzer, error) {
	switch name {
	case "standard":
		return NewStandardAnalyzer()
	case "simple":
		return NewSimpleAnalyzer()
	case "whitespace":
		return NewWhitespaceAnalyzer()
	case "keyword":
		return NewKeywordAnalyzer()
	case "chinese":
		return NewChineseAnalyzer("")
	case "english":
		return NewEnglishAnalyzer()
	case "multilingual":
		return NewMultilingualAnalyzer()
	case "search":
		return NewSearchAnalyzer()
	default:
		// Default to standard analyzer for unknown types
		return NewStandardAnalyzer()
	}
}

// Close destroys the analyzer and frees resources.
func (a *Analyzer) Close() {
	// No-op for pure Go implementation
}

// Analyze analyzes text and returns tokens.
func (a *Analyzer) Analyze(text string) ([]Token, error) {
	if a.tokenizer == nil {
		return nil, nil
	}

	rawTokens := a.tokenizer(text)
	tokens := make([]Token, 0, len(rawTokens))

	position := 0
	offset := 0

	for _, raw := range rawTokens {
		// Apply filters
		filtered := raw
		for _, filter := range a.filters {
			filtered = filter(filtered)
		}

		// Skip stop words
		if a.stopWords != nil && a.stopWords[filtered] {
			offset += len(raw) + 1 // +1 for space
			continue
		}

		// Skip empty tokens
		if filtered == "" {
			offset += len(raw) + 1
			continue
		}

		startOffset := strings.Index(text[offset:], raw)
		if startOffset >= 0 {
			startOffset += offset
		} else {
			startOffset = offset
		}

		tokens = append(tokens, Token{
			Text:        filtered,
			Position:    position,
			StartOffset: startOffset,
			EndOffset:   startOffset + len(raw),
			Type:        "word",
		})

		position++
		offset = startOffset + len(raw)
	}

	return tokens, nil
}

// AnalyzeToStrings is a convenience method that returns just the token text.
func (a *Analyzer) AnalyzeToStrings(text string) ([]string, error) {
	tokens, err := a.Analyze(text)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(tokens))
	for i, token := range tokens {
		result[i] = token.Text
	}
	return result, nil
}

// Name returns the analyzer name.
func (a *Analyzer) Name() string {
	return a.name
}

// Description returns the analyzer description.
func (a *Analyzer) Description() string {
	switch a.name {
	case "standard":
		return "Standard analyzer with lowercase and stop word filters"
	case "simple":
		return "Simple analyzer with whitespace tokenization and lowercase"
	case "whitespace":
		return "Whitespace analyzer - splits on whitespace only"
	case "keyword":
		return "Keyword analyzer - no tokenization"
	case "chinese":
		return "Chinese analyzer with character-based segmentation"
	case "english":
		return "English analyzer with lowercase, ASCII folding, and stop words"
	case "multilingual":
		return "Multilingual analyzer with lowercase and ASCII folding"
	case "search":
		return "Search analyzer optimized for queries"
	default:
		return "Text analyzer"
	}
}

// ClearError clears the last error.
func ClearError() {
	// No-op for pure Go implementation
}
