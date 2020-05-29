package unit

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"istio.io/api/networking/v1alpha3"
)

// StringMatchExtended copies istio StringMatchExtended definition and extends it to add helper methods.
type StringMatchExtended struct {
	*v1alpha3.StringMatch
}

// IsEmpty returns true if struct is empty
func (sm StringMatchExtended) IsEmpty() bool {
	return reflect.DeepEqual(sm, StringMatchExtended{})
}

// Match will return true if a given string matches a StringMatch object.
func (sm *StringMatchExtended) Match(s string) bool {
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

func matchRequest(input parser.Input, httpMatchRequest *v1alpha3.HTTPMatchRequest) bool {
	authority := &StringMatchExtended{httpMatchRequest.Authority}
	uri := &StringMatchExtended{httpMatchRequest.Uri}
	method := &StringMatchExtended{httpMatchRequest.Method}
	// _ = &StringMatchExtended{httpMatchRequest.Headers} // TODO

	return uri.Match(input.URI) && authority.Match(input.Authority) && method.Match(input.Method)
}
