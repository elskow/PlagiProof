package comparator

import (
	"github.com/elskow/PlagiProof/constants"
	"github.com/elskow/PlagiProof/modules/detector"
)

type Comparator interface {
	Compare(code1, code2 string) (float64, error)
	IsSimilar(code1, code2 string) (bool, error)
}

type CodeComparator struct {
	detector        detector.Detector
	sequenceMatcher *SequenceMatcher
}

func NewCodeComparator(detector detector.Detector) *CodeComparator {
	return &CodeComparator{
		detector:        detector,
		sequenceMatcher: NewSequenceMatcher(),
	}
}

func (c *CodeComparator) Compare(code1, code2 string) (float64, error) {
	tokens1, err := c.detector.Run(code1)
	if err != nil {
		return 0, err
	}

	tokens2, err := c.detector.Run(code2)
	if err != nil {
		return 0, err
	}

	return c.sequenceMatcher.Compare(tokens1, tokens2), nil
}

func (c *CodeComparator) IsSimilar(code1, code2 string) (bool, error) {
	similarity, err := c.Compare(code1, code2)
	if err != nil {
		return false, err
	}

	return similarity >= constants.ThresPlagiarism, nil
}
