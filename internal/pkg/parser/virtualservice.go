package parser

import (
	"fmt"
	"os"

	networking "istio.io/api/networking/v1"
	v1 "istio.io/client-go/pkg/apis/networking/v1"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pkg/config/schema/gvk"
)

func ParseVirtualServices(files []string) ([]*v1.VirtualService, error) {
	out := []*v1.VirtualService{}
	for _, file := range files {
		fileContent, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading file %q failed: %w", file, err)
		}
		configs, _, err := crd.ParseInputs(string(fileContent))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CRD %q: %w", file, err)
		}
		for _, c := range configs {
			if c.Meta.GroupVersionKind != gvk.VirtualService {
				continue
			}
			spec, ok := c.Spec.(*networking.VirtualService)
			if !ok {
				return nil, fmt.Errorf("failed to convert spec in %q to VirtualService: %w", file, err)
			}
			vs := &v1.VirtualService{
				ObjectMeta: c.ToObjectMeta(),
				Spec:       *spec, //nolint as deep copying mess up with reflect.DeepEqual comparison.
			}
			out = append(out, vs)
		}
	}
	return out, nil
}
