package main

import (
	"magi_token/sdk"
	"magi_token/tokencore"
	"math/big"

	"github.com/CosmWasm/tinyjson/jwriter"
)

// ===================================
// Internal Helper Functions
// ===================================

// ===================================
// Input Validation
// ===================================

// Maximum allowed lengths for user-controlled input fields.
const (
	maxAddressLen = 256
	maxNameLen    = 64
	maxSymbolLen  = 16
)

// validateAddress checks that an address is within length bounds and contains
// no pipe characters, which are used as state key delimiters.
func validateAddress(account string) {
	if len(account) > maxAddressLen {
		sdk.Abort("Address exceeds maximum length")
	}
	for i := 0; i < len(account); i++ {
		if account[i] == '|' {
			sdk.Abort("Invalid character in address")
		}
	}
}

// ===================================
// Safe Math Utilities
// ===================================

// safeAdd performs a + b addition with big.Int.
func safeAdd(a, b *big.Int) *big.Int {
	return new(big.Int).Add(a, b)
}

// safeSub performs a - b subtraction. Aborts execution if b > a.
func safeSub(a, b *big.Int) *big.Int {
	if b.Cmp(a) > 0 {
		sdk.Abort("safeSub underflow")
	}
	return new(big.Int).Sub(a, b)
}

// ===================================
// big.Int Parsing
// ===================================

// parseBigInt parses a decimal string into a *big.Int. Aborts on invalid input.
func parseBigInt(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		sdk.Abort("Invalid amount: not a valid integer")
	}
	if v.Sign() < 0 {
		sdk.Abort("Invalid amount: negative value")
	}
	return v
}

// ===================================
// Balance Management
// ===================================

// balanceKey returns the state key for an account's balance.
func balanceKey(account string) string {
	return "bal|" + account
}

// incBalance increments token balance of an address.
func incBalance(account string, amount *big.Int) {
	oldBal := getBalanceInternal(account)
	newBal := safeAdd(oldBal, amount)
	sdk.StateSetObject(balanceKey(account), string(bigIntToBytes(newBal)))
}

// decBalance decrements token balance of an address. Aborts if insufficient balance.
func decBalance(account string, amount *big.Int) {
	oldBal := getBalanceInternal(account)
	if oldBal.Cmp(amount) < 0 {
		sdk.Abort("Insufficient balance")
	}
	newBal := safeSub(oldBal, amount)
	sdk.StateSetObject(balanceKey(account), string(bigIntToBytes(newBal)))
}

// getBalanceInternal retrieves token balance of an address.
func getBalanceInternal(account string) *big.Int {
	bal := sdk.StateGetObject(balanceKey(account))
	if bal == nil || *bal == "" {
		return new(big.Int)
	}
	return bytesToBigInt([]byte(*bal))
}

// ===================================
// Allowance Management
// ===================================

// allowanceKey returns the state key for an allowance (owner approves spender).
func allowanceKey(owner, spender string) string {
	return "alw|" + owner + "|" + spender
}

// getAllowanceInternal retrieves the allowance for a spender on owner's tokens.
func getAllowanceInternal(owner, spender string) *big.Int {
	alw := sdk.StateGetObject(allowanceKey(owner, spender))
	if alw == nil || *alw == "" {
		return new(big.Int)
	}
	return bytesToBigInt([]byte(*alw))
}

// setAllowanceInternal sets the allowance for a spender on owner's tokens.
func setAllowanceInternal(owner, spender string, amount *big.Int) {
	sdk.StateSetObject(allowanceKey(owner, spender), string(bigIntToBytes(amount)))
}

// ===================================
// Token Properties (from state)
// ===================================

// getTokenName retrieves the token name from state.
func getTokenName() string {
	n := sdk.StateGetObject("token_name")
	if n == nil {
		return ""
	}
	return *n
}

// getTokenSymbol retrieves the token symbol from state.
func getTokenSymbol() string {
	s := sdk.StateGetObject("token_symbol")
	if s == nil {
		return ""
	}
	return *s
}

// getTokenDecimals retrieves the token decimals from state.
func getTokenDecimals() uint8 {
	d := sdk.StateGetObject("token_decimals")
	if d == nil {
		return 0
	}
	return tokencore.DecodeDecimals(*d)
}

// getMaxSupply retrieves the max supply from state.
func getMaxSupply() *big.Int {
	m := sdk.StateGetObject("token_max_supply")
	if m == nil || *m == "" {
		return new(big.Int)
	}
	return bytesToBigInt([]byte(*m))
}

// ===================================
// Supply Management
// ===================================

// getSupply retrieves the current total supply.
func getSupply() *big.Int {
	s := sdk.StateGetObject("supply")
	if s == nil || *s == "" {
		return new(big.Int)
	}
	return bytesToBigInt([]byte(*s))
}

// setSupply sets the total supply.
func setSupply(amount *big.Int) {
	sdk.StateSetObject("supply", string(bigIntToBytes(amount)))
}

// ===================================
// big.Int <-> []byte Helper
// ===================================

// bigIntToBytes serializes a *big.Int to bytes (big-endian unsigned).
// Returns [0] for zero values.
func bigIntToBytes(val *big.Int) []byte {
	b := val.Bytes()
	if len(b) == 0 {
		return []byte{0}
	}
	return b
}

// bytesToBigInt deserializes bytes (big-endian unsigned) to a *big.Int.
func bytesToBigInt(b []byte) *big.Int {
	return new(big.Int).SetBytes(b)
}

// ===================================
// JSON Response Helper
// ===================================

func jsonResponse(marshaler interface{ MarshalTinyJSON(*jwriter.Writer) }) *string {
	w := jwriter.Writer{}
	marshaler.MarshalTinyJSON(&w)
	if w.Error != nil {
		sdk.Abort("JSON marshal error")
	}
	result := string(w.Buffer.BuildBytes())
	return &result
}
