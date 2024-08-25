package detector

import (
	"fmt"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"hash"
	"hash/fnv"
	"regexp"
	"strings"
	"sync"
)

const ngramSize = 3

var (
	commentNamespacePattern = regexp.MustCompile(`(#.*?\n)|(using\s+namespace\s+std\s*;)|(std::)|(#include\s*<.*?>)`)
	skipTokenTypes          = map[chroma.TokenType]bool{
		chroma.Text:               true,
		chroma.Error:              true,
		chroma.Comment:            true,
		chroma.CommentSpecial:     true,
		chroma.CommentSingle:      true,
		chroma.CommentMultiline:   true,
		chroma.CommentPreproc:     true,
		chroma.CommentPreprocFile: true,
		chroma.CommentHashbang:    true,
	}
	hashPool = sync.Pool{
		New: func() interface{} {
			return fnv.New32a()
		},
	}
	lexer = chroma.Coalesce(lexers.Get("cpp"))
)

// Run the detector on the given code.
func Run(code string) ([]string, error) {
	cleanedCode, err := removeCommentsAndNamespace(code)
	if err != nil {
		return nil, fmt.Errorf("failed to remove comments and namespace: %v", err)
	}

	tokens, err := extractTokens(cleanedCode)
	if err != nil {
		return nil, fmt.Errorf("failed to extract tokens: %v", err)
	}

	fingerprints, err := generateFingerprints(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to generate fingerprints: %v", err)
	}

	return fingerprints, nil
}

// generateFingerprints generates FNV-1a fingerprints for the given code.
func generateFingerprints(tokens []string) ([]string, error) {
	ngrams := generateNgrams(tokens, ngramSize)

	fingerprints := make([]string, len(ngrams))
	var wg sync.WaitGroup
	for i, ngram := range ngrams {
		wg.Add(1)
		go func(i int, ngram string) {
			defer wg.Done()
			hashed := hashPool.Get().(hash.Hash32)
			defer hashPool.Put(hashed)
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

// removeCommentsAndNamespace removes comments and namespace declarations from the given code.
func removeCommentsAndNamespace(code string) (string, error) {
	cleanedCode := commentNamespacePattern.ReplaceAllString(code, "")
	return strings.TrimSpace(cleanedCode), nil
}

// extractTokens extracts tokens from the given code.
func extractTokens(code string) ([]string, error) {
	if lexer == nil {
		return nil, fmt.Errorf("lexer not found")
	}

	tokenStream, err := lexer.Tokenise(nil, code)
	if err != nil {
		return nil, fmt.Errorf("failed to tokenize code: %v", err)
	}

	tokens := make([]string, 0, len(code)/5) // Estimate initial capacity
	for _, token := range tokenStream.Tokens() {
		if skipTokenTypes[token.Type] {
			continue
		}
		tokens = append(tokens, token.String())
	}

	return tokens, nil
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
