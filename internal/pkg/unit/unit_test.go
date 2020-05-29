package unit

import (
	"reflect"
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
)

func TestRun(t *testing.T) {
	testcasefiles := []string{"../../../examples/virtualservice_test.yml"}
	configfiles := []string{"../../../examples/virtualservice.yml"}

	err := Run(testcasefiles, configfiles)
	if err != nil {
		t.Error(err)
	}
}

func TestGetDestination(t *testing.T) {
	type args struct {
		input           parser.Input
		virtualServices []*v1alpha3.VirtualService
	}
	tests := []struct {
		name string
		args args
		want []*networkingv1alpha3.HTTPRouteDestination
	}{{
		name: "no match, empty destination",
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
				}}},
		},
		want: []*networkingv1alpha3.HTTPRouteDestination{},
	}, {
		name: "match single destination, multiple virtualservices",
		args: args{
			input: parser.Input{Authority: "www.match.com", URI: "/"},
			virtualServices: []*v1alpha3.VirtualService{{
				Spec: networkingv1alpha3.VirtualService{
					Hosts: []string{"www.notmatch.com"},
					Http: []*networkingv1alpha3.HTTPRoute{{
						Route: []*networkingv1alpha3.HTTPRouteDestination{
							&networkingv1alpha3.HTTPRouteDestination{
								Destination: &networkingv1alpha3.Destination{
									Host: "notmatch.notmatch.svc.cluster.local",
								}}},
						Match: []*networkingv1alpha3.HTTPMatchRequest{{
							Uri: &networkingv1alpha3.StringMatch{
								MatchType: &networkingv1alpha3.StringMatch_Exact{
									Exact: "/",
								},
							},
						}},
					}},
				}}, {
				Spec: networkingv1alpha3.VirtualService{
					Hosts: []string{"www.match.com"},
					Http: []*networkingv1alpha3.HTTPRoute{{
						Route: []*networkingv1alpha3.HTTPRouteDestination{
							&networkingv1alpha3.HTTPRouteDestination{
								Destination: &networkingv1alpha3.Destination{
									Host: "match.match.svc.cluster.local",
								}}},
						Match: []*networkingv1alpha3.HTTPMatchRequest{{
							Uri: &networkingv1alpha3.StringMatch{
								MatchType: &networkingv1alpha3.StringMatch_Exact{
									Exact: "/",
								},
							},
						}},
					}},
				}}},
		},
		want: []*networkingv1alpha3.HTTPRouteDestination{
			&networkingv1alpha3.HTTPRouteDestination{
				Destination: &networkingv1alpha3.Destination{
					Host: "match.match.svc.cluster.local",
				}}},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetDestination(tt.args.input, tt.args.virtualServices); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDestination() = %v, want %v", got, tt.want)
			}
		})
	}
}
