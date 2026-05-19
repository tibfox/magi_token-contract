// Package tokencore holds pure, host-independent primitives used by the
// token contract. It deliberately imports nothing from the wasm `sdk`
// package so it can be exercised by the standard Go toolchain (`go test`),
// which makes the TOKEN-02 (decimals encoding) and TOKEN-01 (Init auth
// precondition) fixes behaviorally testable rather than source-gated.
package tokencore

// EncodeDecimals serializes a token-decimals value as a single byte string.
//
// TOKEN-02: it must NOT use `string(d)`. In Go, `string(uint8)` is a rune
// conversion: values >= 128 are encoded as a 2-byte UTF-8 sequence (e.g.
// `string(uint8(200))` == "\xc3\x88"), so a reader that takes the first
// byte recovers the wrong value and the token's reported decimals are
// permanently corrupted. Wrapping in a []byte forces a true byte cast.
func EncodeDecimals(d uint8) string {
	return string([]byte{d})
}

// DecodeDecimals reverses EncodeDecimals. An empty string (unset state)
// decodes to 0, matching the contract's "no decimals stored" semantics.
func DecodeDecimals(s string) uint8 {
	if s == "" {
		return 0
	}
	return s[0]
}

// Init auth precondition messages (kept identical to the strings the
// contract previously aborted with, so behavior is unchanged).
const (
	MsgCallerRequired = "Caller required"
	MsgOwnerRequired  = "Contract owner required"
	MsgNotOwner       = "Only contract owner can initialize"
)

// CheckInitAuth validates the owner/caller precondition for Init.
//
// TOKEN-01: the previous code nil-checked `caller` but then dereferenced
// `owner` (`*caller != *owner`) without a nil check, so an environment
// where the host omits `contract.owner` (returns a nil *string) caused a
// nil-pointer panic instead of a clean abort. This function nil-checks
// BOTH pointers before any dereference and returns the abort message to
// use; ok == true means the caller is the owner and Init may proceed.
func CheckInitAuth(owner, caller *string) (msg string, ok bool) {
	if caller == nil {
		return MsgCallerRequired, false
	}
	if owner == nil {
		return MsgOwnerRequired, false
	}
	if *caller != *owner {
		return MsgNotOwner, false
	}
	return "", true
}
