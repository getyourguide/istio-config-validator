package parser

import (
	"fmt"
	"os"

	networking "istio.io/api/networking/v1alpha3"
	v1alpha3 "istio.io/client-go/pkg/apis/networking/v1alpha3"
	"istio.io/istio/pilot/pkg/config/kube/crd"
	"istio.io/istio/pkg/config/schema/collections"
	"istio.io/istio/pkg/config/schema/gvk"
)

func ParseVirtualServices(files []string, strict bool) ([]*v1alpha3.VirtualService, error) {
	out := []*v1alpha3.VirtualService{}
	for _, file := range files {
		fileContent, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("reading file '%s' failed: %w", file, err)
		}
		configs, _, err := crd.ParseInputs(string(fileContent))
		if err != nil {
			return nil, fmt.Errorf("failed to parse CRD %q: %w", file, err)
		}
		for _, c := range configs {
			if c.Meta.GroupVersionKind != gvk.VirtualService {
				continue
			}
			if strict {
				schema, exists := collections.Pilot.FindByGroupVersionAliasesKind(gvk.VirtualService)
				if !exists {
					return nil, fmt.Errorf("failed to find schema for VirtualService")
				}
				warn, err := schema.ValidateConfig(c)
				if err != nil {
					return nil, fmt.Errorf("failed to validate VirtualService: %w", err)
				}
				if warn != nil {
					return nil, fmt.Errorf("failed to validate VirtualService: %w", warn)
				}
			}
			spec, ok := c.Spec.(*networking.VirtualService)
			if !ok {
				return nil, fmt.Errorf("failed to convert spec to VirtualService: %w", err)
			}
			vs := &v1alpha3.VirtualService{
				ObjectMeta: c.ToObjectMeta(),
				Spec:       *spec, //nolint as deep copying mess up with reflect.DeepEqual comparison.
			}
			out = append(out, vs)
		}
	}
	return out, nil
}
