package unit

import (
	"reflect"
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"github.com/stretchr/testify/require"
	networking "istio.io/api/networking/v1"
	v1 "istio.io/client-go/pkg/apis/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestRun(t *testing.T) {
	testcasefiles := []string{"../../../examples/virtualservice_test.yml"}
	configfiles := []string{"../../../examples/virtualservice.yml"}
	var strict bool
	_, _, err := Run(testcasefiles, configfiles, strict)
	require.NoError(t, err)
}

func TestRunDelegate(t *testing.T) {
	testcasefiles := []string{"../../../examples/virtualservice_delegate_test.yml"}
	configfiles := []string{"../../../examples/delegate_virtualservice.yml"}
	var strict bool
	_, _, err := Run(testcasefiles, configfiles, strict)
	require.NoError(t, err)
}

func TestGetRoute(t *testing.T) {
	type args struct {
		input           parser.Input
		virtualServices []*v1.VirtualService
		checkHosts      bool
	}
	tests := []struct {
		name    string
		args    args
		want    *networking.HTTPRoute
		wantErr bool
	}{
		{
			name: "no host match, empty destination",
			args: args{
				input: parser.Input{Authority: "www.example.com", URI: "/"},
				virtualServices: []*v1.VirtualService{{
					Spec: networking.VirtualService{
						Hosts: []string{"www.another-example.com"},
						Http: []*networking.HTTPRoute{{
							Match: []*networking.HTTPMatchRequest{{
								Uri: &networking.StringMatch{
									MatchType: &networking.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}},
				checkHosts: true,
			},
			want:    &networking.HTTPRoute{},
			wantErr: false,
		},
		{
			name: "match a fallback destination",
			args: args{
				input: parser.Input{Authority: "www.example.com", URI: "/path-to-fallback"},
				virtualServices: []*v1.VirtualService{{
					Spec: networking.VirtualService{
						Hosts: []string{"www.example.com"},
						Http: []*networking.HTTPRoute{{
							Match: []*networking.HTTPMatchRequest{{
								Uri: &networking.StringMatch{
									MatchType: &networking.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}, {
							Route: []*networking.HTTPRouteDestination{{
								Destination: &networking.Destination{
									Host: "fallback.fallback.svc.cluster.local",
								},
							}},
						}},
					},
				}},
				checkHosts: true,
			},
			want: &networking.HTTPRoute{
				Route: []*networking.HTTPRouteDestination{{
					Destination: &networking.Destination{
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
				virtualServices: []*v1.VirtualService{{
					Spec: networking.VirtualService{
						Hosts: []string{"www.notmatch.com"},
						Http: []*networking.HTTPRoute{{
							Route: []*networking.HTTPRouteDestination{{
								Destination: &networking.Destination{
									Host: "notmatch.notmatch.svc.cluster.local",
								},
							}},
							Match: []*networking.HTTPMatchRequest{{
								Uri: &networking.StringMatch{
									MatchType: &networking.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}, {
					Spec: networking.VirtualService{
						Hosts: []string{"www.match.com"},
						Http: []*networking.HTTPRoute{{
							Route: []*networking.HTTPRouteDestination{{
								Destination: &networking.Destination{
									Host: "match.match.svc.cluster.local",
								},
							}},
							Match: []*networking.HTTPMatchRequest{{
								Uri: &networking.StringMatch{
									MatchType: &networking.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}},
				checkHosts: true,
			},
			want: &networking.HTTPRoute{
				Route: []*networking.HTTPRouteDestination{{
					Destination: &networking.Destination{
						Host: "match.match.svc.cluster.local",
					},
				}}, Match: []*networking.HTTPMatchRequest{{
					Uri: &networking.StringMatch{
						MatchType: &networking.StringMatch_Exact{
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
				virtualServices: []*v1.VirtualService{{
					Spec: networking.VirtualService{
						Hosts: []string{"www.match.com"},
						Http: []*networking.HTTPRoute{{
							Route: []*networking.HTTPRouteDestination{{
								Destination: &networking.Destination{
									Host: "match.match.svc.cluster.local",
								},
							}},
							Match: []*networking.HTTPMatchRequest{{
								Uri: &networking.StringMatch{
									MatchType: &networking.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}},
				checkHosts: true,
			},
			want: &networking.HTTPRoute{
				Route: []*networking.HTTPRouteDestination{{
					Destination: &networking.Destination{
						Host: "match.match.svc.cluster.local",
					},
				}},
				Match: []*networking.HTTPMatchRequest{{
					Uri: &networking.StringMatch{
						MatchType: &networking.StringMatch_Exact{
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
				virtualServices: []*v1.VirtualService{{
					Spec: networking.VirtualService{
						Http: []*networking.HTTPRoute{{
							Route: []*networking.HTTPRouteDestination{{
								Destination: &networking.Destination{
									Host: "match.match.svc.cluster.local",
								},
							}},
							Match: []*networking.HTTPMatchRequest{{
								Uri: &networking.StringMatch{
									MatchType: &networking.StringMatch_Exact{
										Exact: "/",
									},
								},
							}},
						}},
					},
				}},
				checkHosts: false,
			},
			want: &networking.HTTPRoute{
				Route: []*networking.HTTPRouteDestination{{
					Destination: &networking.Destination{
						Host: "match.match.svc.cluster.local",
					},
				}}, Match: []*networking.HTTPMatchRequest{{
					Uri: &networking.StringMatch{
						MatchType: &networking.StringMatch_Exact{
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
		delegate        *networking.Delegate
		virtualServices []*v1.VirtualService
	}
	tests := []struct {
		name    string
		args    args
		want    *v1.VirtualService
		wantErr bool
	}{
		{
			name: "match",
			args: args{
				delegate: &networking.Delegate{
					Name: "delegate",
				},
				virtualServices: []*v1.VirtualService{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "delegate",
						Namespace: "default",
					},
				}},
			},
			want: &v1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "delegate",
					Namespace: "default",
				},
			},
			wantErr: false,
		}, {
			name: "match with namespace",
			args: args{
				delegate: &networking.Delegate{
					Name:      "delegate",
					Namespace: "test",
				},
				virtualServices: []*v1.VirtualService{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "delegate",
						Namespace: "default",
					},
				}, {
					ObjectMeta: metav1.ObjectMeta{
						Name:      "delegate",
						Namespace: "test",
					},
				}},
			},
			want: &v1.VirtualService{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "delegate",
					Namespace: "test",
				},
			},
			wantErr: false,
		}, {
			name: "no match",
			args: args{
				delegate: &networking.Delegate{
					Name: "delegate",
				},
				virtualServices: []*v1.VirtualService{{
					ObjectMeta: metav1.ObjectMeta{
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
				delegate: &networking.Delegate{
					Name:      "delegate",
					Namespace: "production",
				},
				virtualServices: []*v1.VirtualService{{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "delegate",
						Namespace: "default",
					},
				}, {
					ObjectMeta: metav1.ObjectMeta{
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
