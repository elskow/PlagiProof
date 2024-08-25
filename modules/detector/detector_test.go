package detector

import (
	"testing"
	"time"
)

func Test_removeCommentsAndNamespace(t *testing.T) {
	detector := NewCodeDetector()
	tests := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Basic removal",
			code: `#include <iostream> 
			using namespace std;
			cout << "Hello, World!" << endl;`,
			expected: `cout << "Hello, World!" << endl;`,
		},
		{
			name:     "No comments or namespace",
			code:     `cout << "Hello, World!" << endl;`,
			expected: `cout << "Hello, World!" << endl;`,
		},
		{
			name: "Only comments",
			code: `// its a comment
			/* another comment */
			cout << "Hello, World!" << endl;`,
			expected: `cout << "Hello, World!" << endl;`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := detector.removeCommentsAndNamespace(tt.code)
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %v but got %v", tt.expected, result)
			}
		})
	}

	const iterations = 100
	var totalElapsed time.Duration
	code := `#include <iostream> 
			using namespace std;
			cout << "Hello, World!" << endl;`

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		detector.removeCommentsAndNamespace(code)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
	}

	averageElapsed := totalElapsed / iterations
	t.Logf("Average time taken: %v", averageElapsed)
}

func Test_extractTokens(t *testing.T) {
	detector := NewCodeDetector()
	tests := []struct {
		name           string
		code           string
		expectedTokens []string
	}{
		{
			name: "Basic extraction",
			code: `// its a comment
			/* another comment */
			cout << "Hello, World!" << endl;`,
			expectedTokens: []string{"cout", "<<", "\"Hello, World!\"", "<<", "endl", ";"},
		},
		{
			name:           "No comments",
			code:           `cout << "Hello, World!" << endl;`,
			expectedTokens: []string{"cout", "<<", "\"Hello, World!\"", "<<", "endl", ";"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens, err := detector.extractTokens(tt.code)
			if err != nil {
				t.Errorf("Expected no error but got %v", err)
			}

			if len(tokens) != len(tt.expectedTokens) {
				t.Errorf("Expected %v but got %v", tt.expectedTokens, tokens)
			}

			for i, token := range tokens {
				if token != tt.expectedTokens[i] {
					t.Errorf("Expected token %v but got %v", tt.expectedTokens[i], token)
				}
			}
		})
	}

	const iterations = 100
	var totalElapsed time.Duration
	code := `// its a comment
	/* another comment */
	cout << "Hello, World!" << endl;`

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		detector.extractTokens(code)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
	}

	averageElapsed := totalElapsed / iterations
	t.Logf("Average time taken: %v", averageElapsed)
}

func Test_generateFingerprints(t *testing.T) {
	detector := NewCodeDetector()
	code := `// its a comment
	/* another comment */
	cout << "Hello, World!" << endl;`
	tokens, _ := detector.extractTokens(code)

	const iterations = 100
	var totalElapsed time.Duration

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		_, err := detector.generateFingerprints(tokens)
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
	detector := NewCodeDetector()
	code := `// its a comment
	/* another comment */
	cout << "Hello, World!" << endl;`

	const iterations = 100
	var totalElapsed time.Duration

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		_, err := detector.Run(code)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}
	}

	averageElapsed := totalElapsed / iterations
	t.Logf("Average time taken: %v", averageElapsed)
}
