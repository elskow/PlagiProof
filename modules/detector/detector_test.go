package detector

import (
	"testing"
	"time"
)

func Test_removeCommentsAndNamespace(t *testing.T) {
	code := `#include <iostream> 
			using namespace std;
			cout << "Hello, World!" << endl;`
	expected := `cout << "Hello, World!" << endl;`
	result, err := removeCommentsAndNamespace(code)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %v but got %v", expected, result)
	}

	const iterations = 100
	var totalElapsed time.Duration

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		removeCommentsAndNamespace(code)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
	}

	averageElapsed := totalElapsed / iterations
	t.Logf("Average time taken: %v", averageElapsed)
}

func Test_extractTokens(t *testing.T) {
	code := `// its a comment
	/* another comment */
	cout << "Hello, World!" << endl;`
	expectedTokens := []string{"cout", "<<", "\"Hello, World!\"", "<<", "endl", ";"}
	tokens, err := extractTokens(code)
	if err != nil {
		t.Errorf("Expected no error but got %v", err)
	}

	if len(tokens) != len(expectedTokens) {
		t.Errorf("Expected %v but got %v", expectedTokens, tokens)
	}

	for i, token := range tokens {
		if token != expectedTokens[i] {
			t.Errorf("Expected token %v but got %v", expectedTokens[i], token)
		}
	}

	const iterations = 100
	var totalElapsed time.Duration

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		extractTokens(code)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
	}

	averageElapsed := totalElapsed / iterations
	t.Logf("Average time taken: %v", averageElapsed)
}

func Test_generateFingerprints(t *testing.T) {
	code := `// its a comment
	/* another comment */
	cout << "Hello, World!" << endl;`
	token, _ := extractTokens(code)

	const iterations = 100
	var totalElapsed time.Duration

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		_, err := generateFingerprints(token)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
	}

	averageElapsed := totalElapsed / iterations
	t.Logf("Average time taken: %v", averageElapsed)
}

func Test_runDetector(t *testing.T) {
	code := `// its a comment
	/* another comment */
	cout << "Hello, World!" << endl;`

	const iterations = 100
	var totalElapsed time.Duration

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		_, err := Run(code)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
	}

	averageElapsed := totalElapsed / iterations
	t.Logf("Average time taken: %v", averageElapsed)
}
