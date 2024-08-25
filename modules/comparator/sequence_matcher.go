package comparator

type SequenceMatcher struct{}

func NewSequenceMatcher() *SequenceMatcher {
	return &SequenceMatcher{}
}

// Compare calculates the similarity between two sequences of tokens using an optimized Longest Common Subsequence (LCS) algorithm.
func (sm *SequenceMatcher) Compare(tokens1, tokens2 []string) float64 {
	m := len(tokens1)
	n := len(tokens2)

	// Create a single slice to store lengths of longest common subsequence.
	lcs := make([]int, n+1)

	for i := 1; i <= m; i++ {
		prev := 0
		for j := 1; j <= n; j++ {
			temp := lcs[j]
			if tokens1[i-1] == tokens2[j-1] {
				lcs[j] = prev + 1
			} else {
				lcs[j] = max(lcs[j], lcs[j-1])
			}
			prev = temp
		}
	}

	// The length of the longest common subsequence.
	lcsLength := lcs[n]

	// Calculate the similarity score.
	similarity := float64(2*lcsLength) / float64(m+n)
	return similarity
}
