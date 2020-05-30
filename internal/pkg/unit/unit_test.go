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
		name    string
		args    args
		want    []*networkingv1alpha3.HTTPRouteDestination
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
					}}},
			},
			want:    []*networkingv1alpha3.HTTPRouteDestination{},
			wantErr: false,
		}, {
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
							Route: []*networkingv1alpha3.HTTPRouteDestination{{
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
			want: []*networkingv1alpha3.HTTPRouteDestination{{
				Destination: &networkingv1alpha3.Destination{
					Host: "match.match.svc.cluster.local",
				}}},
			wantErr: false,
		}, {
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
							}}}, {
							Route: []*networkingv1alpha3.HTTPRouteDestination{{
								Destination: &networkingv1alpha3.Destination{
									Host: "fallback.fallback.svc.cluster.local",
								}}},
						}},
					}}},
			},
			want: []*networkingv1alpha3.HTTPRouteDestination{{
				Destination: &networkingv1alpha3.Destination{
					Host: "fallback.fallback.svc.cluster.local",
				}}},
			wantErr: false,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDestination(tt.args.input, tt.args.virtualServices)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDestination() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDestination() = %v, want %v", got, tt.want)
			}
		})
	}
}
