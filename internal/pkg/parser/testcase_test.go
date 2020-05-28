package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseTestCases(t *testing.T) {
	expectedTestCases := []*TestCase{
		{Description: "happy path users"},
		{Description: "Partner service only accepts GET or OPTIONS"},
	}
	configuration := &Configuration{
		RootDir: "../../../examples/",
	}
	parser, err := New(configuration)
	if err != nil {
		t.Errorf("error getting test cases %v", err)
	}
	if len(parser.TestCases) == 0 {
		t.Error("test cases are empty")
	}

	for _, expected := range expectedTestCases {
		testPass := false
		for _, out := range parser.TestCases {
			if expected.Description == out.Description {
				testPass = true
			}
		}
		if !testPass {
			t.Errorf("could not find expected description:'%v'", expected.Description)
		}

	}
}

func TestUnfoldRequest(t *testing.T) {
	testCases := []struct {
		In  Request
		Out []Input
	}{{
		Request{
			Authority: []string{"www.example.com", "example.com"},
			Method:    []string{"GET", "OPTIONS"},
		},
		[]Input{{
			Authority: "www.example.com",
			Method:    "GET",
		}, {
			Authority: "www.example.com",
			Method:    "OPTIONS",
		}, {
			Authority: "example.com",
			Method:    "GET",
		}, {
			Authority: "example.com",
			Method:    "OPTIONS",
		}},
	}}

	for _, testCase := range testCases {
		assert.ElementsMatch(t, testCase.Out, testCase.In.Unfold())
	}

}
