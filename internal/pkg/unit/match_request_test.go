package unit

import (
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"istio.io/api/networking/v1alpha3"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
)

func Test_matchRequest(t *testing.T) {
	type args struct {
		input            parser.Input
		httpMatchRequest *v1alpha3.HTTPMatchRequest
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr bool
	}{{
		name: "no match conditions should always match",
		args: args{
			input:            parser.Input{Authority: "www.example.com", URI: "/", Method: "GET"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{}},
		want:    true,
		wantErr: false,
	}, {
		name: "single match exact (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/exac", Method: "GET"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Exact{
						Exact: "/exac",
					}}}},
		want:    true,
		wantErr: false,
	}, {
		name: "single match exact (false)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/exac", Method: "GET"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Exact{
						Exact: "/exac/",
					}}}},
		want:    false,
		wantErr: false,
	}, {
		name: "single match prefix (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					}}}},
		want:    true,
		wantErr: false,
	}, {
		name: "single match prefix (false)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/not-prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					}}}},
		want:    false,
		wantErr: false,
	}, {
		name: "single match regex (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/regex/test", Method: "POST"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Regex{
						Regex: "/reg.+?(/)",
					}}}},
		want:    true,
		wantErr: false,
	}, {
		name: "single match regex (false)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/not-regex/test", Method: "PATCH"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Regex{
						Regex: "/reg(/)",
					}}}},
		want:    false,
		wantErr: false,
	}, {
		name: "multiple match exact, prefix and regex (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Authority: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Regex{
						Regex: "(www.)example.com",
					},
				},
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					},
				},
				Method: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Exact{
						Exact: "GET",
					}}}},
		want:    true,
		wantErr: false,
	}, {
		name: "multiple match exact, prefix and regex (false)",
		args: args{
			input: parser.Input{Authority: "www.example.co", URI: "/prefix/anotherpath", Method: "GET"},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Authority: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Regex{
						Regex: "(www.)example.com",
					},
				},
				Uri: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Prefix{
						Prefix: "/prefix",
					},
				},
				Method: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Exact{
						Exact: "GET",
					}}}},
		want:    false,
		wantErr: false,
	}, {
		name: "multiple match in headers, regex, prefix, exact (true)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/", Method: "GET", Headers: map[string]string{
				"x-header-exact":  "exact",
				"x-header-prefix": "prefix-something",
				"x-header-regex":  "capture-this-regex",
			}},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Authority: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Regex{
						Regex: "(www.)example.com",
					},
				},
				Headers: map[string]*networkingv1alpha3.StringMatch{
					"x-header-prefix": {
						MatchType: &networkingv1alpha3.StringMatch_Prefix{
							Prefix: "prefix-",
						},
					},
					"x-header-exact": {
						MatchType: &networkingv1alpha3.StringMatch_Exact{
							Exact: "exact",
						},
					},
					"x-header-regex": {
						MatchType: &networkingv1alpha3.StringMatch_Regex{
							Regex: ".+?-this-.+?",
						}}},
			}},
		want:    true,
		wantErr: false,
	}, {
		name: "multiple match in headers, regex, prefix, exact (false)",
		args: args{
			input: parser.Input{Authority: "www.example.com", URI: "/", Method: "GET", Headers: map[string]string{
				"x-header-exact":  "exact",
				"x-header-prefix": "prefix-something",
				"x-header-regex":  "capture-this-regex",
			}},
			httpMatchRequest: &networkingv1alpha3.HTTPMatchRequest{
				Authority: &networkingv1alpha3.StringMatch{
					MatchType: &networkingv1alpha3.StringMatch_Regex{
						Regex: "(www.)example.com",
					},
				},
				Headers: map[string]*networkingv1alpha3.StringMatch{
					"x-header-prefix": {
						MatchType: &networkingv1alpha3.StringMatch_Prefix{
							Prefix: "not-prefix-",
						},
					},
					"x-header-exact": {
						MatchType: &networkingv1alpha3.StringMatch_Exact{
							Exact: "exact",
						},
					},
					"x-header-regex": {
						MatchType: &networkingv1alpha3.StringMatch_Regex{
							Regex: ".+?-this-.+?",
						}}},
			}},
		want:    false,
		wantErr: false,
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := matchRequest(tt.args.input, tt.args.httpMatchRequest)
			if (err != nil) != tt.wantErr {
				t.Errorf("matchRequest() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("matchRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
