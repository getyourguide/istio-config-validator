package unit

import (
	"reflect"
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	networkingv1alpha3 "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRun(t *testing.T) {
	testcasefiles := []string{"../../../examples/virtualservice_test.yml"}
	configfiles := []string{"../../../examples/virtualservice.yml"}

	_, _, err := Run(testcasefiles, configfiles)
	if err != nil {
		t.Error(err)
	}
}

func TestRunDelegate(t *testing.T) {
	testcasefiles := []string{"../../../examples/virtualservice_delegate_test.yml"}
	configfiles := []string{"../../../examples/delegate_virtualservice.yml"}

	_, _, err := Run(testcasefiles, configfiles)
	if err != nil {
		t.Error(err)
	}
}

func TestGetRoute(t *testing.T) {
	type args struct {
		input           parser.Input
		virtualServices []*v1alpha3.VirtualService
		checkHosts      bool
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
				checkHosts: true,
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
				checkHosts: true,
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
				checkHosts: true,
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
				checkHosts: true,
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
		{
			name: "match virtualservice with no hosts",
			args: args{
				input: parser.Input{Authority: "www.match.com", URI: "/"},
				virtualServices: []*v1alpha3.VirtualService{{
					Spec: networkingv1alpha3.VirtualService{
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
				checkHosts: false,
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRoute(tt.args.input, tt.args.virtualServices, tt.args.checkHosts)
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

func TestGetDelegatedVirtualService(t *testing.T) {
	type args struct {
		delegate        *networkingv1alpha3.Delegate
		virtualServices []*v1alpha3.VirtualService
	}
	tests := []struct {
		name    string
		args    args
		want    *v1alpha3.VirtualService
		wantErr bool
	}{
		{
			name: "match",
			args: args{
				delegate: &networkingv1alpha3.Delegate{
					Name: "delegate",
				},
				virtualServices: []*v1alpha3.VirtualService{{
					ObjectMeta: v1.ObjectMeta{
						Name:      "delegate",
						Namespace: "default",
					},
				}},
			},
			want: &v1alpha3.VirtualService{
				ObjectMeta: v1.ObjectMeta{
					Name:      "delegate",
					Namespace: "default",
				},
			},
			wantErr: false,
		}, {
			name: "match with namespace",
			args: args{
				delegate: &networkingv1alpha3.Delegate{
					Name:      "delegate",
					Namespace: "test",
				},
				virtualServices: []*v1alpha3.VirtualService{{
					ObjectMeta: v1.ObjectMeta{
						Name:      "delegate",
						Namespace: "default",
					},
				}, {
					ObjectMeta: v1.ObjectMeta{
						Name:      "delegate",
						Namespace: "test",
					},
				}},
			},
			want: &v1alpha3.VirtualService{
				ObjectMeta: v1.ObjectMeta{
					Name:      "delegate",
					Namespace: "test",
				},
			},
			wantErr: false,
		}, {
			name: "no match",
			args: args{
				delegate: &networkingv1alpha3.Delegate{
					Name: "delegate",
				},
				virtualServices: []*v1alpha3.VirtualService{{
					ObjectMeta: v1.ObjectMeta{
						Name:      "delegate-abc",
						Namespace: "default",
					},
				}},
			},
			want:    nil,
			wantErr: true,
		}, {
			name: "no match with namespace",
			args: args{
				delegate: &networkingv1alpha3.Delegate{
					Name:      "delegate",
					Namespace: "production",
				},
				virtualServices: []*v1alpha3.VirtualService{{
					ObjectMeta: v1.ObjectMeta{
						Name:      "delegate",
						Namespace: "default",
					},
				}, {
					ObjectMeta: v1.ObjectMeta{
						Name:      "delegate",
						Namespace: "test",
					},
				}},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDelegatedVirtualService(tt.args.delegate, tt.args.virtualServices)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDelegatedVirtualService() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDelegatedVirtualService() = %v, want %v", got, tt.want)
			}
		})
	}
}
