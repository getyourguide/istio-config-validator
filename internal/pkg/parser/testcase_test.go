package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnknownField(t *testing.T) {
	testsFiles := []string{"testdata/invalid_test.yml"}
	_, err := ParseTestCases(testsFiles, true)
	require.ErrorContains(t, err, "json: unknown field")

	_, err = ParseTestCases(testsFiles, false)
	require.NoError(t, err)
}

func TestParseTestCases(t *testing.T) {
	expectedTestCases := []*TestCase{
		{Description: "happy path users"},
		{Description: "Partner service only accepts GET or OPTIONS"},
	}
	testcasefiles := []string{"../../../examples/virtualservice_test.yml"}
	testCases, err := ParseTestCases(testcasefiles, false)
	require.NoError(t, err)
	require.NotEmpty(t, testCases)

	for _, expected := range expectedTestCases {
		testPass := false
		for _, out := range testCases {
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
		Name  string
		In    Request
		Out   []Input
		Error error
	}{
		{
			"single authority, method and URI",
			Request{
				Authority: []string{"www.example.com"},
				Method:    []string{"GET"},
				URI:       []string{"/"},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
				},
			},
			nil,
		},
		{
			"single authority, method and URI with headers",
			Request{
				Authority: []string{"www.example.com"},
				Method:    []string{"GET"},
				URI:       []string{"/"},
				Headers: map[string]string{
					"Cookie": "namnamnamnamnamnam",
					"x-y-z":  "X-Y-Z",
				},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
			},
			nil,
		},
		{
			"single authority, method and multiple URIs",
			Request{
				Authority: []string{"www.example.com"},
				Method:    []string{"GET"},
				URI:       []string{"/", "/healthz", "/.well-known/foo"},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/healthz",
				},
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/.well-known/foo",
				},
			},
			nil,
		},
		{
			"multiple authorities and single method and URI",
			Request{
				Authority: []string{"www.example.com", "example.com", "foo.bar"},
				Method:    []string{"GET"},
				URI:       []string{"/"},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "foo.bar",
					Method:    "GET",
					URI:       "/",
				},
			},
			nil,
		},
		{
			"single authority and multiple methods and single URI",
			Request{
				Authority: []string{"www.example.com"},
				Method:    []string{"GET", "POST", "PUT"},
				URI:       []string{"/"},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "POST",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "PUT",
					URI:       "/",
				},
			},
			nil,
		},
		{
			"multiple authorities and multiple methods and multiple URIs with headers",
			Request{
				Authority: []string{"www.example.com", "example.com"},
				Method:    []string{"GET", "POST", "PUT"},
				URI:       []string{"/", "/healthz"},
				Headers: map[string]string{
					"Cookie": "namnamnamnamnamnam",
					"x-y-z":  "X-Y-Z",
				},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "www.example.com",
					Method:    "POST",
					URI:       "/",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "www.example.com",
					Method:    "PUT",
					URI:       "/",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "example.com",
					Method:    "GET",
					URI:       "/",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "example.com",
					Method:    "POST",
					URI:       "/",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "example.com",
					Method:    "PUT",
					URI:       "/",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/healthz",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "www.example.com",
					Method:    "POST",
					URI:       "/healthz",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "www.example.com",
					Method:    "PUT",
					URI:       "/healthz",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "example.com",
					Method:    "GET",
					URI:       "/healthz",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "example.com",
					Method:    "POST",
					URI:       "/healthz",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
				{
					Authority: "example.com",
					Method:    "PUT",
					URI:       "/healthz",
					Headers: map[string]string{
						"Cookie": "namnamnamnamnamnam",
						"x-y-z":  "X-Y-Z",
					},
				},
			},
			nil,
		},
		{
			"multiple authorities and multiple methods and multiple URIs",
			Request{
				Authority: []string{"www.example.com", "example.com"},
				Method:    []string{"GET", "POST", "PUT"},
				URI:       []string{"/", "/healthz"},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "POST",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "PUT",
					URI:       "/",
				},
				{
					Authority: "example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "example.com",
					Method:    "POST",
					URI:       "/",
				},
				{
					Authority: "example.com",
					Method:    "PUT",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/healthz",
				},
				{
					Authority: "www.example.com",
					Method:    "POST",
					URI:       "/healthz",
				},
				{
					Authority: "www.example.com",
					Method:    "PUT",
					URI:       "/healthz",
				},
				{
					Authority: "example.com",
					Method:    "GET",
					URI:       "/healthz",
				},
				{
					Authority: "example.com",
					Method:    "POST",
					URI:       "/healthz",
				},
				{
					Authority: "example.com",
					Method:    "PUT",
					URI:       "/healthz",
				},
			},
			nil,
		},
		{
			"multiple authorities and multiple methods and single URI",
			Request{
				Authority: []string{"www.example.com", "example.com"},
				Method:    []string{"GET", "POST", "PUT"},
				URI:       []string{"/"},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "POST",
					URI:       "/",
				},
				{
					Authority: "www.example.com",
					Method:    "PUT",
					URI:       "/",
				},
				{
					Authority: "example.com",
					Method:    "GET",
					URI:       "/",
				},
				{
					Authority: "example.com",
					Method:    "POST",
					URI:       "/",
				},
				{
					Authority: "example.com",
					Method:    "PUT",
					URI:       "/",
				},
			},
			nil,
		},
		{
			"empty authority list",
			Request{
				Authority: []string{},
				Method:    []string{"GET", "OPTIONS"},
				URI:       []string{"/"},
			},
			[]Input{},
			ErrEmptyAuthorityList,
		},
		{
			"empty method list",
			Request{
				Authority: []string{"www.example.com", "www.example.net"},
				Method:    []string{},
				URI:       []string{"/"},
			},
			[]Input{},
			ErrEmptyMethodList,
		},
		{
			"empty URI list",
			Request{
				Authority: []string{"www.example.com"},
				Method:    []string{"GET"},
				URI:       []string{},
			},
			[]Input{},
			ErrEmptyURIList,
		},
		{
			"empty authority list and method list",
			Request{
				Authority: []string{},
				Method:    []string{},
				URI:       []string{"/"},
			},
			[]Input{},
			ErrEmptyAuthorityList,
		},
		{
			"empty authority list and URI list",
			Request{
				Authority: []string{},
				Method:    []string{"GET"},
				URI:       []string{},
			},
			[]Input{},
			ErrEmptyAuthorityList,
		},
		{
			"empty authority list and method list and URI list",
			Request{
				Authority: []string{},
				Method:    []string{},
				URI:       []string{},
			},
			[]Input{},
			ErrEmptyAuthorityList,
		},
		{
			"empty method list and URI list",
			Request{
				Authority: []string{"www.example.com"},
				Method:    []string{},
				URI:       []string{},
			},
			[]Input{},
			ErrEmptyMethodList,
		},
		{
			"query parameters should be removed",
			Request{
				Authority: []string{"www.example.com"},
				Method:    []string{"POST"},
				URI:       []string{"/reseller?partner_id=12344"},
			},
			[]Input{
				{
					Authority: "www.example.com",
					Method:    "POST",
					URI:       "/reseller",
				},
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.Name, func(t *testing.T) {
			got, gotErr := testCase.In.Unfold()
			if gotErr != testCase.Error {
				t.Errorf("expected err=%v, got err=%v", testCase.Error, gotErr)
			}
			assert.ElementsMatch(t, testCase.Out, got)
		})
	}
}
