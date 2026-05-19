package contract_test

import (
	"fmt"
	"testing"
)

// ===================================
// Security Review Findings — A/B Tests
//
// Convention: each finding has a RED test (fails against the buggy code)
// and a GREEN test (passes against the fixed code).
//
// TOKEN-02 and TOKEN-01 are latent defects whose faulty path is NOT
// reachable through the embedded-wasm harness under the repo's pinned
// toolchain (TinyGo 0.39.0 lowers `string(uint8Var)` to a single byte; the
// host's getEnvKey("contract.owner") returns "" rather than a nil pointer).
// Their fix was therefore extracted into the pure, host-independent package
// `magi_token/tokencore`, which IS exercised by the standard Go toolchain.
// The behavioral RED->GREEN for those two lives in
// tokencore/tokencore_test.go (run `go test ./tokencore/`):
//
//   * TestEncodeDecimals_AllValuesRoundTrip / _IsByteCastNotRuneCast
//     — RED with the rune cast (decimals 128..255 corrupt: e.g.
//       EncodeDecimals(128) -> 2 bytes "c280"), GREEN with the byte cast.
//   * TestCheckInitAuth_NilOwnerDoesNotPanic / _Matrix
//     — RED with the unguarded deref (real nil-pointer SIGSEGV), GREEN
//       with the symmetric owner nil-check.
//
// The tests below are the wasm-level behavioral / regression coverage that
// IS reachable through the harness: TOKEN-05 is a full behavioral
// RED->GREEN; the TOKEN-02/01 wasm tests are end-to-end regression guards
// proving the refactor did not change observable contract behavior.
// ===================================

// --- TOKEN-02 (HIGH): decimals encoding ---

// Documents the defect class with the standard gc compiler: the rune
// conversion the contract used to do corrupts the reader byte for d>=128,
// while the fixed byte form round-trips. (Constant documentation; the
// flipping RED/GREEN gate is in tokencore/tokencore_test.go.)
func TestTOKEN02_RuneConversionIsCorrupting(t *testing.T) {
	readerByte0 := func(s string) uint8 { return s[0] } // internal.go (*d)[0] equivalent
	for _, d := range []uint8{128, 200, 255} {
		if got := readerByte0(string(rune(d))); got == d {
			t.Errorf("expected string(rune(%d)) to corrupt the reader byte", d)
		}
		if fixed := readerByte0(string([]byte{d})); fixed != d {
			t.Errorf("fixed form must round-trip decimals=%d, got %d", d, fixed)
		}
	}
}

// End-to-end regression guard: with the fix in place the contract
// round-trips decimals byte-exact through the embedded wasm.
func TestTOKEN02_WasmRoundTripGuard(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "getInfo", []byte(""), nil, "hive:anyone", true, uint(100_000_000),
		`{"name":"Magi Token","symbol":"MAGI","decimals":"3","maxSupply":"1000000000"}`)
}

// --- TOKEN-05 (MEDIUM): Init does not validate Decimals ---
//
// Full behavioral RED -> GREEN through the embedded wasm:
// RED on buggy source (Init accepts decimals=200), GREEN on fixed (aborts).

func TestTOKEN05_InitFailsDecimalsTooLarge(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(fmt.Sprintf(`{"name":"Magi Token","symbol":"MAGI","decimals":%d,"maxSupply":"1000000000"}`, 200))
	CallContract(t, ct, "init", payload, nil, ownerAddress, false, uint(100_000_000), "")
}

// Boundary: decimals == 18 (max valid) must still succeed.
func TestTOKEN05_InitAcceptsDecimals18(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Magi Token","symbol":"MAGI","decimals":18,"maxSupply":"1000000000"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "getInfo", []byte(""), nil, "hive:anyone", true, uint(100_000_000),
		`{"name":"Magi Token","symbol":"MAGI","decimals":"18","maxSupply":"1000000000"}`)
}

// --- TOKEN-01 (LOW): unguarded *owner deref in Init ---

// Documents WHY the missing guard is unsafe: dereferencing a nil *string
// panics. (Constant documentation; the flipping RED/GREEN gate is
// tokencore.TestCheckInitAuth_NilOwnerDoesNotPanic.)
func TestTOKEN01_NilOwnerDerefPanics(t *testing.T) {
	caller := "hive:tibfox"
	func() {
		defer func() {
			if recover() == nil {
				t.Errorf("expected nil-owner deref to panic without the guard")
			}
		}()
		var owner *string // nil, as a missing contract.owner would be
		_ = caller != *owner
	}()
}

// End-to-end regression guard: with a valid owner the guard does not change
// behavior — legitimate owner init succeeds, non-owner is rejected.
func TestTOKEN01_GuardTransparentWasm(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, "hive:notowner", false, uint(100_000_000), "")
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
}
