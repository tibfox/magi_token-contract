package tokencore

import (
	"testing"
)

// --- TOKEN-02: decimals encoding ---

// TestEncodeDecimals_AllValuesRoundTrip is the behavioral RED/GREEN test for
// TOKEN-02. With the fix (byte cast) every uint8 0..255 round-trips. With the
// pre-fix rune conversion (`string(d)`), values >= 128 become a 2-byte UTF-8
// sequence and the first byte no longer equals the original — this test fails
// for all 128 of those values.
func TestEncodeDecimals_AllValuesRoundTrip(t *testing.T) {
	for v := 0; v <= 255; v++ {
		d := uint8(v)
		got := DecodeDecimals(EncodeDecimals(d))
		if got != d {
			t.Fatalf("round-trip corrupted: EncodeDecimals(%d) -> DecodeDecimals = %d (want %d)", d, got, d)
		}
	}
}

// TestEncodeDecimals_IsByteCastNotRuneCast pins the exact corruption mechanism:
// the encoded form of a >=128 value must be a single byte equal to the value,
// NOT the 2-byte UTF-8 a rune conversion would produce.
func TestEncodeDecimals_IsByteCastNotRuneCast(t *testing.T) {
	for _, d := range []uint8{128, 200, 255} {
		enc := EncodeDecimals(d)
		if len(enc) != 1 {
			t.Fatalf("EncodeDecimals(%d): len=%d (%x) — looks like a rune cast; want 1 byte", d, len(enc), enc)
		}
		if enc[0] != d {
			t.Fatalf("EncodeDecimals(%d): first byte=%d; want %d", d, enc[0], d)
		}
		// Demonstrate the bug the fix prevents: the rune conversion really
		// does corrupt these values.
		buggy := string(d) //nolint:gosimple // intentional demonstration of the TOKEN-02 defect
		if buggy[0] == d {
			t.Fatalf("expected rune-cast string(%d) to corrupt the first byte, but it did not (toolchain-dependent)", d)
		}
	}
}

func TestDecodeDecimals_EmptyIsZero(t *testing.T) {
	if got := DecodeDecimals(""); got != 0 {
		t.Fatalf("DecodeDecimals(\"\") = %d, want 0", got)
	}
}

// --- TOKEN-01: Init auth precondition ---

func strptr(s string) *string { return &s }

// TestCheckInitAuth_NilOwnerDoesNotPanic is the behavioral RED/GREEN test for
// TOKEN-01. The pre-fix inline code dereferenced `owner` (`*caller != *owner`)
// after only nil-checking `caller`, so a nil owner pointer panicked. The fixed
// CheckInitAuth must return a clean rejection instead of panicking.
func TestCheckInitAuth_NilOwnerDoesNotPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("CheckInitAuth panicked on nil owner (TOKEN-01 regression): %v", r)
		}
	}()
	msg, ok := CheckInitAuth(nil, strptr("hive:alice"))
	if ok {
		t.Fatal("expected rejection when owner is nil")
	}
	if msg != MsgOwnerRequired {
		t.Fatalf("nil owner: msg=%q, want %q", msg, MsgOwnerRequired)
	}
}

func TestCheckInitAuth_Matrix(t *testing.T) {
	cases := []struct {
		name        string
		owner       *string
		caller      *string
		wantMsg     string
		wantOK      bool
	}{
		{"nil caller", strptr("hive:owner"), nil, MsgCallerRequired, false},
		{"nil owner", nil, strptr("hive:owner"), MsgOwnerRequired, false},
		{"both nil", nil, nil, MsgCallerRequired, false}, // caller checked first
		{"mismatch", strptr("hive:owner"), strptr("hive:attacker"), MsgNotOwner, false},
		{"match", strptr("hive:owner"), strptr("hive:owner"), "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			msg, ok := CheckInitAuth(c.owner, c.caller)
			if ok != c.wantOK || msg != c.wantMsg {
				t.Fatalf("CheckInitAuth = (%q,%v), want (%q,%v)", msg, ok, c.wantMsg, c.wantOK)
			}
		})
	}
}
