package unit

import (
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"istio.io/api/networking/v1alpha3"
)

func Test_Something(t *testing.T) {
	httpMatchRequest := &v1alpha3.HTTPMatchRequest{
		Uri: &v1alpha3.StringMatch{
			MatchType: &v1alpha3.StringMatch_Exact{
				Exact: "/",
			},
		},
	}
	input := parser.Input{Authority: "www.example.com", URI: "/"}
	matchRequest(input, httpMatchRequest)
}

func Test_matchRequest(t *testing.T) {
	type args struct {
		input            parser.Input
		httpMatchRequest *v1alpha3.HTTPMatchRequest
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "single match exact (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/exac", Method: "GET"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Exact{
						Exact: "/exac",
					}}}},
		want: true,
	}, {
		name: "single match exact (false)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/exac", Method: "GET"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Exact{
						Exact: "/exac/",
					}}}},
		want: false,
	}, {
		name: "single match prefix (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					}}}},
		want: true,
	}, {
		name: "single match prefix (false)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/not-prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					}}}},
		want: false,
	}, {
		name: "single match regex (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/regex/test", Method: "POST"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Regex{
						Regex: "/reg.+?(/)",
					}}}},
		want: true,
	}, {
		name: "single match regex (false)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/not-regex/test", Method: "PATCH"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Regex{
						Regex: "/reg(/)",
					}}}},
		want: false,
	}, {
		name: "multiple match exact, prefix and regex (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Authority: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Regex{
						Regex: "(www.)example.com",
					}},
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					}},
				Method: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Exact{
						Exact: "GET",
					}}}},
		want: true,
	}, {
		name: "multiple match exact, prefix and regex (false)",
		args: args{
			input: parser.Input{Authority: "www.example.co", URI: "/prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &v1alpha3.HTTPMatchRequest{
				Authority: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Regex{
						Regex: "(www.)example.com",
					}},
				Uri: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					}},
				Method: &v1alpha3.StringMatch{
					MatchType: &v1alpha3.StringMatch_Exact{
						Exact: "GET",
					}}}},
		want: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := matchRequest(tt.args.input, tt.args.httpMatchRequest); got != tt.want {
				t.Errorf("matchRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
