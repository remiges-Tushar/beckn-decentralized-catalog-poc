package crypto

import (
	"github.com/gibson042/canonicaljson-go"
)

func CanonicalizeJSON(v any) ([]byte, error) {
	return canonicaljson.Marshal(v)
}
