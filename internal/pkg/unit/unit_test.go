package unit

import (
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"istio.io/api/networking/v1alpha3"
)

func TestRun(t *testing.T) {
	configuration := &Configuration{
		RootDir: "../../../examples/",
	}

	err := Run(configuration)
	if err != nil {
		t.Error(err)
	}
}

func Test_inputMatch(t *testing.T) {
	type args struct {
		input           parser.Input
		virtualServices []*v1alpha3.VirtualService
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{
		name: "happy path",
		args: args{
			input: parser.Input{Authority: "www.exemple.com", URI: "/"},
			virtualServices: []*v1alpha3.VirtualService{{
				Hosts: []string{"www.example.com"},
				Http: []*v1alpha3.HTTPRoute{{
					Match: []*v1alpha3.HTTPMatchRequest{{
						Uri: &v1alpha3.StringMatch{
							MatchType: &v1alpha3.StringMatch_Exact{
								Exact: "/",
							},
						},
					}},
				}},
			}},
		},
		want: true,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := inputMatch(tt.args.input, tt.args.virtualServices); got != tt.want {
				t.Errorf("inputMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
