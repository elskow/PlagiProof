package comparator

import (
	"github.com/elskow/PlagiProof/modules/detector"
	"os"
	"testing"
	"time"
)

func readCodeFromFile(filePath string) (string, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func TestCodeComparator_Compare(t *testing.T) {

	t.Run("Same code", func(t *testing.T) {
		code1, err := readCodeFromFile("../../example/diff1.cpp")
		if err != nil {
			t.Fatalf("Failed to read code1: %v", err)
		}

		code2, err := readCodeFromFile("../../example/diff2.cpp")
		if err != nil {
			t.Fatalf("Failed to read code2: %v", err)
		}

		c := NewCodeComparator(detector.NewCodeDetector())
		similarity, err := c.Compare(code1, code2)

		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}

		if similarity > 0.5 {
			t.Errorf("Expected similarity >= 0.5 but got %v", similarity)
		}

		t.Logf("Similarity: %v", similarity)
	})

	t.Run("Different code", func(t *testing.T) {
		code1, err := readCodeFromFile("../../example/same1.cpp")
		if err != nil {
			t.Fatalf("Failed to read code1: %v", err)
		}

		code2, err := readCodeFromFile("../../example/same2.cpp")
		if err != nil {
			t.Fatalf("Failed to read code2: %v", err)
		}

		c := NewCodeComparator(detector.NewCodeDetector())
		similarity, err := c.Compare(code1, code2)

		if err != nil {
			t.Errorf("Expected no error but got %v", err)
		}

		if similarity <= 0.5 {
			t.Errorf("Expected similarity < 0.5 but got %v", similarity)
		}

		t.Logf("Similarity: %v", similarity)
	})
}

func BenchmarkCodeComparator_Compare(b *testing.B) {
	code1, err := readCodeFromFile("../../example/diff1.cpp")
	if err != nil {
		b.Fatalf("Failed to read code1: %v", err)
	}

	code2, err := readCodeFromFile("../../example/diff2.cpp")
	if err != nil {
		b.Fatalf("Failed to read code2: %v", err)
	}

	c := NewCodeComparator(detector.NewCodeDetector())

	const iterations = 100
	var totalElapsed time.Duration

	for i := 0; i < iterations; i++ {
		begin := time.Now()
		_, err := c.Compare(code1, code2)
		elapsed := time.Since(begin)
		totalElapsed += elapsed
		if err != nil {
			b.Errorf("Expected no error but got %v", err)
		}
	}

	averageElapsed := totalElapsed / iterations
	b.Logf("Average time taken: %v", averageElapsed)
}
