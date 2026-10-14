package ipso

import "net/http"

// Configuration controls how a registry is created by New and NewContext, and how it talks to the OMA API.
// The zero value creates an empty registry that makes no requests.
type Configuration struct {
	InitOnNew      bool     // fetch the registry from the OMA API on creation
	SkipInitErrors bool     // skip objects whose definition cannot be fetched or parsed instead of failing
	Sanitize       bool     // sanitize object and resource descriptions on creation
	Sanitizer      []string // strings cleaned from descriptions by sanitizing, DefaultSanitizer is used when empty, see (*Reg).Sanitize
	// HTTPClient performs all OMA API requests, the objects list and every object definition.
	// A shared client with a 5 second request timeout is used when nil; a custom client applies its own Timeout only.
	HTTPClient *http.Client
	// BaseURL is the OMA API root, the objects list is requested from BaseURL + "/Object".
	// Object definitions are fetched from the ObjectLink URLs of that list, not from BaseURL.
	// The official API (https://www.openmobilealliance.org/api/lwm2m/v1) is used when empty.
	BaseURL string
}

// DefaultConfiguration returns a Configuration that fetches the registry from the OMA API, skips objects that fail
// to load, and sanitizes the result with DefaultSanitizer.
func DefaultConfiguration() Configuration {
	return Configuration{
		InitOnNew:      true,
		SkipInitErrors: true, // true by default, because OMA API returns some objects without ObjectLink filled in,
		// which makes it impossible to initialize objects and resources with strict validation
		Sanitize: true,
	}
}

// DefaultSanitizer returns an array with default sanitizer strings: line breaks, tabs, backslashes, bullets
// and double spaces. The whitespace ones are replaced by a single space, the others are removed.
func DefaultSanitizer() []string {
	return []string{
		"\n", "\\\n", "\t", "\\t", "\\", "•", "  ",
	}
}
