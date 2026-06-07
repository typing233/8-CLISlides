package executor

import (
	"strings"
	"testing"
)

func TestPreprocess(t *testing.T) {
	input := "before\n\n~~~\necho hello\n~~~\n\nafter"
	result := Preprocess(input)
	if !strings.Contains(result, "hello") {
		t.Errorf("expected preprocessed output to contain 'hello', got: %s", result)
	}
	if strings.Contains(result, "~~~") {
		t.Errorf("expected ~~~ to be replaced, got: %s", result)
	}
}

func TestExecuteCodeBlock(t *testing.T) {
	out, err := ExecuteCodeBlock("echo test123", "bash")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "test123") {
		t.Errorf("expected output to contain 'test123', got: %s", out)
	}
}
