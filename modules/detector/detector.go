package detector

import (
	"fmt"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/elskow/PlagiProof/constants"
	"hash"
	"hash/fnv"
	"regexp"
	"strings"
	"sync"
)

// Detector interface defines the methods for code detection.
type Detector interface {
	Run(code string) ([]string, error)
	removeCommentsAndNamespace(code string) (string, error)
	extractTokens(code string) ([]string, error)
	generateFingerprints(tokens []string) ([]string, error)
}

// CodeDetector struct implements the Detector interface.
type CodeDetector struct {
	commentNamespacePattern *regexp.Regexp
	skipTokenTypes          map[chroma.TokenType]bool
	hashPool                sync.Pool
	lexer                   chroma.Lexer
}

// NewCodeDetector creates a new instance of CodeDetector.
func NewCodeDetector() *CodeDetector {
	return &CodeDetector{
		commentNamespacePattern: regexp.MustCompile(`(?m)^\s*//.*$|(?s)/\*.*?\*/|(?m)^\s*using\s+namespace\s+\w+;\s*$|(?m)^\s*#include\s*<.*?>\s*$`),
		skipTokenTypes: map[chroma.TokenType]bool{
			chroma.Text:               true,
			chroma.Error:              true,
			chroma.Comment:            true,
			chroma.CommentSpecial:     true,
			chroma.CommentSingle:      true,
			chroma.CommentMultiline:   true,
			chroma.CommentPreproc:     true,
			chroma.CommentPreprocFile: true,
			chroma.CommentHashbang:    true,
		},
		hashPool: sync.Pool{
			New: func() interface{} {
				return fnv.New32a()
			},
		},
		lexer: chroma.Coalesce(lexers.Get("cpp")),
	}
}

// Run the detector on the given code.
func (d *CodeDetector) Run(code string) ([]string, error) {
	cleanedCode, err := d.removeCommentsAndNamespace(code)
	if err != nil {
		return nil, fmt.Errorf("failed to remove comments and namespace: %v", err)
	}

	tokens, err := d.extractTokens(cleanedCode)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tokens: %v", err)
	}

	fingerprints, err := d.generateFingerprints(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to generate fingerprints: %v", err)
	}

	return fingerprints, nil
}

// removeCommentsAndNamespace removes comments and the 'using namespace' directive from the code.
func (d *CodeDetector) removeCommentsAndNamespace(code string) (string, error) {
	code = d.commentNamespacePattern.ReplaceAllString(code, "")

	// Trim leading and trailing whitespace from each line
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		lines[i] = strings.TrimSpace(line)
	}
	code = strings.Join(lines, "\n")

	// Remove empty lines
	reEmptyLines := regexp.MustCompile(`(?m)^\s*$[\r\n]*`)
	code = reEmptyLines.ReplaceAllString(code, "")

	return code, nil
}

// extractTokens extracts tokens from the given code.
func (d *CodeDetector) extractTokens(code string) ([]string, error) {
	if d.lexer == nil {
		return nil, fmt.Errorf("lexer not found")
	}

	tokenStream, err := d.lexer.Tokenise(nil, code)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize code: %v", err)
	}

	tokens := make([]string, 0, len(code)/5) // Estimate initial capacity
	for _, token := range tokenStream.Tokens() {
		if d.skipTokenTypes[token.Type] {
			continue
		}
		tokens = append(tokens, token.String())
	}

	return tokens, nil
}

// generateFingerprints generates FNV-1a fingerprints for the given code.
func (d *CodeDetector) generateFingerprints(tokens []string) ([]string, error) {
	ngrams := generateNgrams(tokens, constants.Ngrams)

	fingerprints := make([]string, len(ngrams))
	var wg sync.WaitGroup
	for i, ngram := range ngrams {
		wg.Add(1)
		go func(i int, ngram string) {
			defer wg.Done()
			hashed := d.hashPool.Get().(hash.Hash32)
			defer d.hashPool.Put(hashed)
			hashed.Reset()
			_, err := hashed.Write([]byte(ngram))
			if err != nil {
				return
			}
			fingerprints[i] = fmt.Sprintf("%x", hashed.Sum32())
		}(i, ngram)
	}
	wg.Wait()

	return fingerprints, nil
}

// generateNgrams generates n-grams of the given size from the given tokens.
func generateNgrams(tokens []string, size int) []string {
	if len(tokens) < size {
		return []string{}
	}
	ngrams := make([]string, 0, len(tokens)-size+1)
	var sb strings.Builder
	for i := 0; i < len(tokens)-size+1; i++ {
		sb.Reset()
		for j := 0; j < size; j++ {
			if j > 0 {
				sb.WriteString(" ")
			}
			sb.WriteString(tokens[i+j])
		}
		ngrams = append(ngrams, sb.String())
	}
	return ngrams
}
