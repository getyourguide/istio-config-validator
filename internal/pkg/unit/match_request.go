package unit

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"istio.io/api/networking/v1alpha3"
)

// ExtendedStringMatch copies istio ExtendedStringMatch definition and extends it to add helper methods.
type ExtendedStringMatch struct {
	*v1alpha3.StringMatch
}

// IsEmpty returns true if struct is empty
func (sm ExtendedStringMatch) IsEmpty() bool {
	return reflect.DeepEqual(sm, ExtendedStringMatch{})
}

// Match will return true if a given string matches a StringMatch object.
func (sm *ExtendedStringMatch) Match(s string) bool {
	if sm.IsEmpty() {
		return true
	}

	switch {
	case sm.GetExact() != "":
		return sm.GetExact() == s
	case sm.GetPrefix() != "":
		return strings.HasPrefix(s, sm.GetPrefix())
	case sm.GetRegex() != "":
		r, err := regexp.Compile(sm.GetRegex())
		if err != nil {
			return false
		}
		return r.MatchString(s)
	}
	return false
}

// matchRequest takes an Input and evaluates against a HTTPMatchRequest block. It replicates
// Istio VirtualService semantic returning true when ALL conditions within the block are true.
// TODO: Add support for all fields within a match block. The ones supported today are:
// Authority, Uri, Method and Headers.
func matchRequest(input parser.Input, httpMatchRequest *v1alpha3.HTTPMatchRequest) bool {
	authority := &ExtendedStringMatch{httpMatchRequest.Authority}
	uri := &ExtendedStringMatch{httpMatchRequest.Uri}
	method := &ExtendedStringMatch{httpMatchRequest.Method}

	for headerName, sm := range httpMatchRequest.Headers {
		if _, ok := input.Headers[headerName]; !ok {
			continue
		}
		header := &ExtendedStringMatch{sm}
		if !header.Match(input.Headers[headerName]) {
			return false
		}
	}

	return uri.Match(input.URI) && authority.Match(input.Authority) && method.Match(input.Method)
}
