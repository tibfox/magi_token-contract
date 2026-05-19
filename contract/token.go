package main

import (
	"magi_token/sdk"
	"magi_token/tokencore"
	"math/big"
	"strconv"

	"github.com/CosmWasm/tinyjson/jlexer"
)

// ===================================
// MAGI Token - Exported WASM Functions
// ===================================

// ===================================
// Initialization
// ===================================

// Init initializes the token contract.
// Can only be called once by the contract owner (deployment account).
// Payload: {"name": "Token Name", "symbol": "TKN", "decimals": 3, "maxSupply": "1000000000"}
//
//go:wasmexport init
func Init(payload *string) *string {
	if isInit() {
		sdk.Abort("Already initialized")
	}

	// Only contract owner can initialize
	owner := sdk.GetEnvKey("contract.owner")
	caller := sdk.GetEnvKey("msg.caller")
	// TOKEN-01: nil-check both owner and caller before any dereference.
	if msg, ok := tokencore.CheckInitAuth(owner, caller); !ok {
		sdk.Abort(msg)
	}

	// Parse payload
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}
	var p InitPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	// Validate payload
	if p.Name == "" {
		sdk.Abort("Name required")
	}
	if len(p.Name) > maxNameLen {
		sdk.Abort("Name exceeds maximum length")
	}
	if p.Symbol == "" {
		sdk.Abort("Symbol required")
	}
	if len(p.Symbol) > maxSymbolLen {
		sdk.Abort("Symbol exceeds maximum length")
	}
	if p.MaxSupply == nil || p.MaxSupply.Sign() == 0 {
		sdk.Abort("MaxSupply must be greater than 0")
	}
	if p.Decimals > 18 {
		sdk.Abort("decimals must be <= 18")
	}

	// Store token properties
	sdk.StateSetObject("token_name", p.Name)
	sdk.StateSetObject("token_symbol", p.Symbol)
	sdk.StateSetObject("token_decimals", tokencore.EncodeDecimals(p.Decimals))
	sdk.StateSetObject("token_max_supply", string(bigIntToBytes(p.MaxSupply)))

	// Initialize contract state
	sdk.StateSetObject("isInit", "1")
	sdk.StateSetObject("supply", string(bigIntToBytes(new(big.Int))))
	sdk.StateSetObject("owner", *owner)

	emitInit(*owner, p.Name, p.Symbol, int(p.Decimals), p.MaxSupply)
	return jsonResponse(SuccessResponse{Success: true})
}

// ===================================
// Token Supply Actions
// ===================================

// Mint creates new tokens and assigns them to the owner.
// Payload: {"amount": "1000"}
// Only the owner can mint tokens.
//
//go:wasmexport mint
func Mint(payload *string) *string {
	assertInit()
	assertNotPaused()
	owner, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Must be owner to mint")
	}
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p MintPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Amount == nil || p.Amount.Sign() == 0 {
		sdk.Abort("Amount must be greater than 0")
	}
	supply := getSupply()
	newSupply := safeAdd(supply, p.Amount)
	if newSupply.Cmp(getMaxSupply()) > 0 {
		sdk.Abort("Exceeded max supply")
	}
	setSupply(newSupply)
	incBalance(owner, p.Amount)
	emitTransfer("", owner, p.Amount) // ERC-20: mint is transfer from zero address
	return jsonResponse(SuccessResponse{Success: true})
}

// Burn destroys tokens from the caller's balance.
// Payload: {"amount": "500"}
//
//go:wasmexport burn
func Burn(payload *string) *string {
	assertInit()
	assertNotPaused()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p BurnPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Amount == nil || p.Amount.Sign() == 0 {
		sdk.Abort("Amount must be greater than 0")
	}
	caller := sdk.GetEnvKey("msg.caller")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	decBalance(*caller, p.Amount)
	supply := getSupply()
	newSupply := safeSub(supply, p.Amount)
	setSupply(newSupply)
	emitTransfer(*caller, "", p.Amount) // ERC-20: burn is transfer to zero address
	return jsonResponse(SuccessResponse{Success: true})
}

// ===================================
// Token Transfer Actions
// ===================================

// Transfer moves tokens from caller to recipient.
// Payload: {"to": "hive:recipient", "amount": "100"}
//
//go:wasmexport transfer
func Transfer(payload *string) *string {
	assertInit()
	assertNotPaused()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p TransferPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.To == "" {
		sdk.Abort("Recipient required")
	}
	validateAddress(p.To)
	if p.Amount == nil || p.Amount.Sign() == 0 {
		sdk.Abort("Amount must be greater than 0")
	}

	caller := sdk.GetEnvKey("msg.caller")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	from := *caller

	if from == p.To {
		sdk.Abort("Cannot transfer to self")
	}

	decBalance(from, p.Amount)
	incBalance(p.To, p.Amount)
	emitTransfer(from, p.To, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// TransferFrom moves tokens from one address to another using allowance.
// Payload: {"from": "hive:owner", "to": "hive:recipient", "amount": "100"}
// Caller must have sufficient allowance from 'from' address.
//
//go:wasmexport transferFrom
func TransferFrom(payload *string) *string {
	assertInit()
	assertNotPaused()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p TransferFromPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.From == "" {
		sdk.Abort("From address required")
	}
	validateAddress(p.From)
	if p.To == "" {
		sdk.Abort("To address required")
	}
	validateAddress(p.To)
	if p.Amount == nil || p.Amount.Sign() == 0 {
		sdk.Abort("Amount must be greater than 0")
	}

	caller := sdk.GetEnvKey("msg.caller")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	spender := *caller

	if p.From == p.To {
		sdk.Abort("Cannot transfer to same address")
	}

	// Check and deduct allowance
	allowance := getAllowanceInternal(p.From, spender)
	if allowance.Cmp(p.Amount) < 0 {
		sdk.Abort("Insufficient allowance")
	}
	newAllowance := safeSub(allowance, p.Amount)
	setAllowanceInternal(p.From, spender, newAllowance)
	emitApproval(p.From, spender, newAllowance)

	// Transfer tokens
	decBalance(p.From, p.Amount)
	incBalance(p.To, p.Amount)
	emitTransfer(p.From, p.To, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// ===================================
// Allowance Actions
// ===================================

// Approve sets the allowance for a spender to spend caller's tokens.
// Payload: {"spender": "hive:spender", "amount": "100"}
//
//go:wasmexport approve
func Approve(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p ApprovePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Spender == "" {
		sdk.Abort("Spender required")
	}
	validateAddress(p.Spender)

	caller := sdk.GetEnvKey("msg.caller")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	owner := *caller

	if owner == p.Spender {
		sdk.Abort("Cannot approve self")
	}

	if p.Amount == nil {
		p.Amount = new(big.Int)
	}

	setAllowanceInternal(owner, p.Spender, p.Amount)
	emitApproval(owner, p.Spender, p.Amount)
	return jsonResponse(SuccessResponse{Success: true})
}

// IncreaseAllowance atomically increases the allowance for a spender.
// This is to prevent race conditions.
// Payload: {"spender": "hive:spender", "amount": "100"}
//
//go:wasmexport increaseAllowance
func IncreaseAllowance(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p AllowancePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Spender == "" {
		sdk.Abort("Spender required")
	}
	validateAddress(p.Spender)

	caller := sdk.GetEnvKey("msg.caller")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	owner := *caller

	if owner == p.Spender {
		sdk.Abort("Cannot approve self")
	}

	if p.Amount == nil {
		p.Amount = new(big.Int)
	}

	currentAllowance := getAllowanceInternal(owner, p.Spender)
	newAllowance := safeAdd(currentAllowance, p.Amount)
	setAllowanceInternal(owner, p.Spender, newAllowance)
	emitApproval(owner, p.Spender, newAllowance)
	return jsonResponse(SuccessResponse{Success: true})
}

// DecreaseAllowance atomically decreases the allowance for a spender.
// This is to prevent race conditions.
// Payload: {"spender": "hive:spender", "amount": "100"}
//
//go:wasmexport decreaseAllowance
func DecreaseAllowance(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p AllowancePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Spender == "" {
		sdk.Abort("Spender required")
	}
	validateAddress(p.Spender)

	caller := sdk.GetEnvKey("msg.caller")
	if caller == nil {
		sdk.Abort("Caller required")
	}
	owner := *caller

	if owner == p.Spender {
		sdk.Abort("Cannot approve self")
	}

	if p.Amount == nil {
		p.Amount = new(big.Int)
	}

	currentAllowance := getAllowanceInternal(owner, p.Spender)
	if currentAllowance.Cmp(p.Amount) < 0 {
		sdk.Abort("Decreased allowance below zero")
	}
	newAllowance := safeSub(currentAllowance, p.Amount)
	setAllowanceInternal(owner, p.Spender, newAllowance)
	emitApproval(owner, p.Spender, newAllowance)
	return jsonResponse(SuccessResponse{Success: true})
}

// ===================================
// Contract Management Actions
// ===================================

// ChangeOwner transfers contract ownership to a new address.
// Payload: {"newOwner": "hive:newowner"}
//
//go:wasmexport changeOwner
func ChangeOwner(payload *string) *string {
	assertInit()
	previousOwner, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Not owner")
	}
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p ChangeOwnerPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.NewOwner == "" {
		sdk.Abort("New owner required")
	}
	validateAddress(p.NewOwner)

	sdk.StateSetObject("owner", p.NewOwner)
	emitOwnerChange(previousOwner, p.NewOwner)
	return jsonResponse(SuccessResponse{Success: true})
}

// Pause pauses all token transfers. Only owner can pause.
//
//go:wasmexport pause
func Pause(_ *string) *string {
	assertInit()
	owner, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Not owner")
	}
	if isPaused() {
		sdk.Abort("Already paused")
	}
	sdk.StateSetObject("paused", "1")
	emitPaused(owner)
	return jsonResponse(SuccessResponse{Success: true})
}

// Unpause unpauses all token transfers. Only owner can unpause.
//
//go:wasmexport unpause
func Unpause(_ *string) *string {
	assertInit()
	owner, isOwner := getOwner()
	if !isOwner {
		sdk.Abort("Not owner")
	}
	if !isPaused() {
		sdk.Abort("Not paused")
	}
	sdk.StateSetObject("paused", "0")
	emitUnpaused(owner)
	return jsonResponse(SuccessResponse{Success: true})
}

// ===================================
// Read-Only Queries
// ===================================

// BalanceOf returns the token balance of an address.
// Payload: {"account": "hive:user"}
//
//go:wasmexport balanceOf
func BalanceOf(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p BalancePayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Account == "" {
		sdk.Abort("Account required")
	}

	bal := getBalanceInternal(p.Account)
	return jsonResponse(BalanceResponse{Balance: bal})
}

// TotalSupply returns the current total supply.
//
//go:wasmexport totalSupply
func TotalSupply(_ *string) *string {
	assertInit()
	supply := getSupply()
	return jsonResponse(SupplyResponse{TotalSupply: supply})
}

// GetOwnerExport returns the current contract owner.
//
//go:wasmexport getOwner
func GetOwnerExport(_ *string) *string {
	assertInit()
	owner, _ := getOwner()
	return jsonResponse(OwnerResponse{Owner: owner})
}

// Allowance returns the allowance for a spender on owner's tokens.
// Payload: {"owner": "hive:owner", "spender": "hive:spender"}
//
//go:wasmexport allowance
func Allowance(payload *string) *string {
	assertInit()
	if payload == nil || *payload == "" {
		sdk.Abort("Payload required")
	}

	var p AllowanceQueryPayload
	r := jlexer.Lexer{Data: []byte(*payload)}
	p.UnmarshalTinyJSON(&r)
	if r.Error() != nil {
		sdk.Abort("Invalid payload")
	}

	if p.Owner == "" {
		sdk.Abort("Owner required")
	}
	if p.Spender == "" {
		sdk.Abort("Spender required")
	}

	alw := getAllowanceInternal(p.Owner, p.Spender)
	return jsonResponse(AllowanceResponse{Allowance: alw})
}

// GetInfo returns token metadata.
//
//go:wasmexport getInfo
func GetInfo(_ *string) *string {
	assertInit()
	return jsonResponse(InfoResponse{
		Name:      getTokenName(),
		Symbol:    getTokenSymbol(),
		Decimals:  strconv.Itoa(int(getTokenDecimals())),
		MaxSupply: getMaxSupply(),
	})
}

// IsPausedExport returns whether the contract is paused.
//
//go:wasmexport isPaused
func IsPausedExport(_ *string) *string {
	assertInit()
	return jsonResponse(PausedResponse{Paused: isPaused()})
}
