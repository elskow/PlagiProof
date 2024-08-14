package detector

import (
	"crypto/sha256"
	"encoding/hex"
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"log"
	"regexp"
	"strings"
)

func removeCommentsAndNamespace(code string) string {
	code = regexp.MustCompile(`#.*?\n`).ReplaceAllString(code, "")
	code = regexp.MustCompile(`using\s+namespace\s+std\s*;`).ReplaceAllString(code, "")
	code = regexp.MustCompile(`std::`).ReplaceAllString(code, "")
	code = regexp.MustCompile(`#include\s*<.*?>`).ReplaceAllString(code, "")
	code = strings.TrimSpace(code)
	return code
}

func getToken(code string) []string {
	lexer := lexers.Get("cpp")
	if lexer == nil {
		log.Fatal("lexer not found")
	}

	lexer = chroma.Coalesce(lexer)
	token, err := lexer.Tokenise(nil, code)

	if err != nil {
		log.Fatal(err)
	}

	var tokens []string

	shouldSkippedTokenTypes := map[chroma.TokenType]bool{
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

	for _, t := range token.Tokens() {
		if shouldSkippedTokenTypes[t.Type] {
			continue
		}

		tokens = append(tokens, t.Value)
	}

	return tokens
}

func ngrams(tokens []string, n int) []string {
	var ngrams []string
	for i := 0; i < len(tokens)-n+1; i++ {
		ngrams = append(ngrams, strings.Join(tokens[i:i+n], " "))
	}
	return ngrams
}

func generateFingerprints(code string) []string {
	code = removeCommentsAndNamespace(code)
	tokens := getToken(code)
	n := 3
	ngrams := ngrams(tokens, n)

	var fingerprints []string
	for _, ngram := range ngrams {
		hash := sha256.Sum256([]byte(ngram))
		fingerprints = append(fingerprints, hex.EncodeToString(hash[:]))
	}

	return fingerprints
}
