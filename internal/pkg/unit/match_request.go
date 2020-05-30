package unit

import (
	"fmt"
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
func (sm *ExtendedStringMatch) Match(s string) (bool, error) {
	if sm.IsEmpty() {
		return true, nil
	}

	switch {
	case sm.GetExact() != "":
		return sm.GetExact() == s, nil
	case sm.GetPrefix() != "":
		return strings.HasPrefix(s, sm.GetPrefix()), nil
	case sm.GetRegex() != "":
		r, err := regexp.Compile(sm.GetRegex())
		if err != nil {
			return false, fmt.Errorf("could not compile regex %s: %v", sm.GetRegex(), err)
		}
		return r.MatchString(s), nil
	}
	return false, nil
}

// matchRequest takes an Input and evaluates against a HTTPMatchRequest block. It replicates
// Istio VirtualService semantic returning true when ALL conditions within the block are true.
// TODO: Add support for all fields within a match block. The ones supported today are:
// Authority, Uri, Method and Headers.
func matchRequest(input parser.Input, httpMatchRequest *v1alpha3.HTTPMatchRequest) (bool, error) {
	authority := &ExtendedStringMatch{httpMatchRequest.Authority}
	uri := &ExtendedStringMatch{httpMatchRequest.Uri}
	method := &ExtendedStringMatch{httpMatchRequest.Method}

	for headerName, sm := range httpMatchRequest.Headers {
		if _, ok := input.Headers[headerName]; !ok {
			continue
		}
		header := &ExtendedStringMatch{sm}
		match, err := header.Match(input.Headers[headerName])
		if err != nil {
			return false, err
		}
		if !match {
			return false, nil
		}
	}

	uriMatch, err := uri.Match(input.URI)
	if err != nil {
		return false, err
	}
	authorityMatch, err := authority.Match(input.Authority)
	if err != nil {
		return false, err
	}
	methodMatch, err := method.Match(input.Method)
	if err != nil {
		return false, err
	}
	return authorityMatch && uriMatch && methodMatch, nil
}
