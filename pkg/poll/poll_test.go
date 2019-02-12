package poll

import "testing"

func TestSplitParameters(t *testing.T) {
	input := `"My Topic's" Option1 Option2 "Option 3"`

	output := SplitParameters(input)

	if output[0] != "My Topic's" {
		t.Errorf("First item in slice is incorrect, got: %s, want: %s.", output[0], "My Topic's")
	}
	if output[1] != "Option1" {
		t.Errorf("Second item in slice is incorrect, got: %s, want: %s.", output[1], "Option1")
	}
	if output[2] != "Option2" {
		t.Errorf("Thrid item in slice is incorrect, got: %s, want: %s.", output[2], "Option2")
	}
	if output[3] != "Option 3" {
		t.Errorf("Third item in slice is incorrect, got: %s, want: %s.", output[3], "Option 3")
	}
}
