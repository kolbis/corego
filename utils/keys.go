package utils

import "strings"

// Keys used to have a common way of buiding keys such as cache keys
type Keys struct {
	Sep string
}

// NewKeys creates a new instance of Keys
func NewKeys() Keys {
	return Keys{
		Sep: "-",
	}
}

// Build a key
func (k Keys) Build(prefix string, name string, params ...string) string {
	var key string

	if prefix != "" {
		key = prefix + k.Sep
	}

	key += name

	if params != nil && len(params) > 0 {
		key += "-" + strings.Join(params, k.Sep)
	}

	return strings.ToLower(key)
}
