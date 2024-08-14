package detector

import (
	"fmt"
	"testing"
)

func Test_removeCommentsAndNamespace(t *testing.T) {
	code := `#include <iostream> 
			using namespace std;
			cout << "Hello, World!" << endl;`
	expected := `cout << "Hello, World!" << endl;`
	result := removeCommentsAndNamespace(code)
	if result != expected {
		t.Errorf("Expected %q but got %q", expected, result)
	}
}

func Test_getToken(t *testing.T) {
	code := `// its a comment
			/* another comment */
			cout << "Hello, World!" << endl;`
	tokens := getToken(code)
	expectedTokens := []string{"cout", "<<", "\"Hello, World!\"", "<<", "endl", ";"}
	if fmt.Sprintf("%q", tokens) != fmt.Sprintf("%q", expectedTokens) {
		t.Errorf("Expected %q but got %q", expectedTokens, tokens)
	}
}
