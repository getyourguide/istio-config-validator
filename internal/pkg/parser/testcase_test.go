package parser

import (
	"testing"
)

func TestParseTestCases(t *testing.T) {
	expectedTestCases := []*TestCase{
		{Description: "happy path users"},
		{Description: "Partner service only accepts GET and OPTIONS"},
	}
	configuration := &Configuration{
		RootDir: "../../../examples/",
	}
	outTestCases, err := New(configuration)
	if err != nil {
		t.Errorf("error getting test cases %v", err)
	}
	if len(outTestCases) == 0 {
		t.Error("test cases are empty")
	}

	for _, expected := range expectedTestCases {
		testPass := false
		for _, out := range outTestCases {
			if expected.Description == out.Description {
				testPass = true
			}
		}
		if !testPass {
			t.Errorf("could not find expected description:'%v'", expected.Description)
		}

	}

}
