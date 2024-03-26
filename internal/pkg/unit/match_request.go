package unit

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"istio.io/api/networking/v1alpha3"
)

// MapMatcher checks map parameters, which requires executing multiple ExtendedStringMatch checks
type MapMatcher struct {
	matchers map[string]*v1alpha3.StringMatch
}

func (mm MapMatcher) Match(input map[string]string) (bool, error) {
	for param, matcher := range mm.matchers {
		value, exists := input[param]
		if !exists {
			return false, nil
		}

		if matcher.GetMatchType() == nil {
			// match type exact: "", and {} checks for presence in query and headers respectively
			continue
		}

		extendedMatcher := &ExtendedStringMatch{StringMatch: matcher}
		matched, err := extendedMatcher.Match(value)
		if err != nil {
			return false, err
		}

		if !matched {
			return false, nil
		}
	}

	return true, nil
}

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

	switch m := sm.GetMatchType().(type) {
	case *v1alpha3.StringMatch_Exact:
		return m.Exact == s, nil
	case *v1alpha3.StringMatch_Prefix:
		return strings.HasPrefix(s, sm.GetPrefix()), nil
	case *v1alpha3.StringMatch_Regex:
		// The rule will not match if only a subsequence of the string matches the regex.
		// https://www.envoyproxy.io/docs/envoy/latest/api-v3/config/route/v3/route_components.proto#envoy-v3-api-field-config-route-v3-routematch-safe-regex
		r, err := regexp.Compile("^" + sm.GetRegex() + "$")
		if err != nil {
			return false, fmt.Errorf("could not compile regex %s: %v", sm.GetRegex(), err)
		}
		return r.MatchString(s), nil
	default:
		return false, fmt.Errorf("unknown matcher type %T", sm.GetMatchType())
	}
}

// matchRequest takes an Input and evaluates against a HTTPMatchRequest block. It replicates
// Istio VirtualService semantic returning true when ALL conditions within the block are true.
// TODO: Add support for all fields within a match block. The ones supported today are:
// Authority, Uri, Method, Headers, Scheme, and QueryParams.
func matchRequest(input parser.Input, httpMatchRequest *v1alpha3.HTTPMatchRequest) (bool, error) {
	uri := &ExtendedStringMatch{httpMatchRequest.Uri}
	if matched, err := uri.Match(input.URI); !matched || err != nil {
		return false, err
	}

	authority := &ExtendedStringMatch{httpMatchRequest.Authority}
	if matched, err := authority.Match(input.Authority); !matched || err != nil {
		return false, err
	}

	method := &ExtendedStringMatch{httpMatchRequest.Method}
	if matched, err := method.Match(input.Method); !matched || err != nil {
		return false, err
	}

	scheme := &ExtendedStringMatch{httpMatchRequest.Scheme}
	if matched, err := scheme.Match(input.Scheme); !matched || err != nil {
		return false, err
	}

	headers := &MapMatcher{matchers: httpMatchRequest.Headers}
	if matched, err := headers.Match(input.Headers); !matched || err != nil {
		return false, err
	}

	query := &MapMatcher{matchers: httpMatchRequest.QueryParams}
	if matched, err := query.Match(input.QueryParams); !matched || err != nil {
		return false, err
	}

	return true, nil
}
