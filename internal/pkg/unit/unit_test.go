package unit

import (
	"fmt"
	"os"
	"path"
	"reflect"
	"strings"
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"gopkg.in/yaml.v3"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/pkg/log"
)

func TestRun(t *testing.T) {
	testcasefiles := []string{"../../../examples/virtualservice_test.yml"}
	configfiles := []string{"../../../examples/virtualservice.yml"}

	_, _, err := Run(testcasefiles, configfiles)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRoute(t *testing.T) {
	type args struct {
		input           parser.Input
		virtualServices []*v1alpha3.VirtualService
	}
	tests := []struct {
		name    string
		args    args
		want    *networkingv1alpha3.HTTPRoute
		wantErr bool
	}{
		{
			name: "no host match, empty destination",
			args: args{
				input: parser.Input{Authority: "www.exemple.com", URI: "/"},
				virtualServices: []*v1alpha3.VirtualService{{
					Spec: networkingv1alpha3.VirtualService{
						Hosts: []string{"www.another-example.com"},
						Http: []*networkingv1alpha3.HTTPRoute{{
							Match: []*networkingv1alpha3.HTTPMatchRequest{{
								Uri: &networkingv1alpha3.StringMatch{
									MatchType: &networkingv1alpha3.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}},
			},
			want:    &networkingv1alpha3.HTTPRoute{},
			wantErr: false,
		},
		{
			name: "match a fallback destination",
			args: args{
				input: parser.Input{Authority: "www.example.com", URI: "/path-to-fallback"},
				virtualServices: []*v1alpha3.VirtualService{{
					Spec: networkingv1alpha3.VirtualService{
						Hosts: []string{"www.example.com"},
						Http: []*networkingv1alpha3.HTTPRoute{{
							Match: []*networkingv1alpha3.HTTPMatchRequest{{
								Uri: &networkingv1alpha3.StringMatch{
									MatchType: &networkingv1alpha3.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}, {
							Route: []*networkingv1alpha3.HTTPRouteDestination{{
								Destination: &networkingv1alpha3.Destination{
									Host: "fallback.fallback.svc.cluster.local",
								},
							}},
						}},
					},
				}},
			},
			want: &networkingv1alpha3.HTTPRoute{
				Route: []*networkingv1alpha3.HTTPRouteDestination{{
					Destination: &networkingv1alpha3.Destination{
						Host: "fallback.fallback.svc.cluster.local",
					},
				}},
			},
			wantErr: false,
		},
		{
			name: "match single destination, multiple virtualservices",
			args: args{
				input: parser.Input{Authority: "www.match.com", URI: "/"},
				virtualServices: []*v1alpha3.VirtualService{{
					Spec: networkingv1alpha3.VirtualService{
						Hosts: []string{"www.notmatch.com"},
						Http: []*networkingv1alpha3.HTTPRoute{{
							Route: []*networkingv1alpha3.HTTPRouteDestination{{
								Destination: &networkingv1alpha3.Destination{
									Host: "notmatch.notmatch.svc.cluster.local",
								},
							}},
							Match: []*networkingv1alpha3.HTTPMatchRequest{{
								Uri: &networkingv1alpha3.StringMatch{
									MatchType: &networkingv1alpha3.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}, {
					Spec: networkingv1alpha3.VirtualService{
						Hosts: []string{"www.match.com"},
						Http: []*networkingv1alpha3.HTTPRoute{{
							Route: []*networkingv1alpha3.HTTPRouteDestination{{
								Destination: &networkingv1alpha3.Destination{
									Host: "match.match.svc.cluster.local",
								},
							}},
							Match: []*networkingv1alpha3.HTTPMatchRequest{{
								Uri: &networkingv1alpha3.StringMatch{
									MatchType: &networkingv1alpha3.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}},
			},
			want: &networkingv1alpha3.HTTPRoute{
				Route: []*networkingv1alpha3.HTTPRouteDestination{{
					Destination: &networkingv1alpha3.Destination{
						Host: "match.match.svc.cluster.local",
					},
				}}, Match: []*networkingv1alpha3.HTTPMatchRequest{{
					Uri: &networkingv1alpha3.StringMatch{
						MatchType: &networkingv1alpha3.StringMatch_Exact{
							Exact: "/",
						},
					},
				}},
			},
			wantErr: false,
		},
		{
			name: "match and assert rewrite and destination",
			args: args{
				input: parser.Input{Authority: "www.match.com", URI: "/"},
				virtualServices: []*v1alpha3.VirtualService{{
					Spec: networkingv1alpha3.VirtualService{
						Hosts: []string{"www.match.com"},
						Http: []*networkingv1alpha3.HTTPRoute{{
							Route: []*networkingv1alpha3.HTTPRouteDestination{{
								Destination: &networkingv1alpha3.Destination{
									Host: "match.match.svc.cluster.local",
								},
							}},
							Match: []*networkingv1alpha3.HTTPMatchRequest{{
								Uri: &networkingv1alpha3.StringMatch{
									MatchType: &networkingv1alpha3.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}},
			},
			want: &networkingv1alpha3.HTTPRoute{
				Route: []*networkingv1alpha3.HTTPRouteDestination{{
					Destination: &networkingv1alpha3.Destination{
						Host: "match.match.svc.cluster.local",
					},
				}},
				Match: []*networkingv1alpha3.HTTPMatchRequest{{
					Uri: &networkingv1alpha3.StringMatch{
						MatchType: &networkingv1alpha3.StringMatch_Exact{
							Exact: "/",
						},
					},
				}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRoute(tt.args.input, tt.args.virtualServices)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRoute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

type yamlFormatter struct{}

var yamlFormat Formatter = yamlFormatter{}

func (y yamlFormatter) Pass(input parser.Input) string {
	b, err := yaml.Marshal(input)
	if err != nil {
		log.Warn(fmt.Sprintf("failed to parse input to YAML, fellback to default format: %v", err))
		return defaultFormatter.Pass(input)
	}
	return fmt.Sprintf("Pass, input:\n%s", string(b))
}

func (y yamlFormatter) Fail(input parser.Input) string {
	b, err := yaml.Marshal(input)
	if err != nil {
		log.Warn(fmt.Sprintf("failed to parse input to YAML, fellback to default format: %v", err))
		return defaultFormatter.Fail(input)
	}
	return fmt.Sprintf("Fail, input:\n%s", string(b))
}

func TestFromFile(t *testing.T) {
	workdir, err := os.Getwd()
	if err != nil {
		t.Error(err)
	}
	testCasesRoot := path.Join(workdir, "testcases")
	entries, err := os.ReadDir(testCasesRoot)
	if err != nil {
		t.Error(err)
	}

	for _, entry := range entries {
		testRoot := path.Join(testCasesRoot, entry.Name())
		tests := []string{path.Join(testRoot, "tests.yml")}
		services := []string{path.Join(testRoot, "service.yml")}
		parsed, err := parser.New(tests, services)
		if err != nil {
			t.Error(err)
		}

		t.Run(entry.Name(), func(t *testing.T) {
			for _, testCase := range parsed.TestCases {
				runner := &TestRunner{
					TestCases:       []*parser.TestCase{testCase},
					VirtualServices: parsed.VirtualServices,
					Format:          yamlFormat,
				}
				t.Run(testCase.Description, func(t *testing.T) {
					summary, details, err := runner.Run()
					t.Log(strings.Join(summary, "\n"))
					t.Log(strings.Join(details, "\n"))
					if err != nil {
						t.Error(err)
					}
				})
			}
		})

	}
}
