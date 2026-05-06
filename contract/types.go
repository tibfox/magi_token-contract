package main

import "math/big"

// ===================================
// MAGI Token - JSON Types (tinyjson)
// ===================================

// ===================================
// Payload Types (Input)
// ===================================

// InitPayload for init action
type InitPayload struct {
	Name      string   `json:"name"`
	Symbol    string   `json:"symbol"`
	Decimals  uint8    `json:"decimals"`
	MaxSupply *big.Int `json:"-"`
}

// TransferPayload for transfer action
type TransferPayload struct {
	To     string   `json:"to"`
	Amount *big.Int `json:"-"`
}

// TransferFromPayload for transferFrom action
type TransferFromPayload struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"-"`
}

// ApprovePayload for approve action
type ApprovePayload struct {
	Spender string   `json:"spender"`
	Amount  *big.Int `json:"-"`
}

// AllowancePayload for increaseAllowance/decreaseAllowance actions
type AllowancePayload struct {
	Spender string   `json:"spender"`
	Amount  *big.Int `json:"-"`
}

// AllowanceQueryPayload for getAllowance query
type AllowanceQueryPayload struct {
	Owner   string `json:"owner"`
	Spender string `json:"spender"`
}

// BalancePayload for getBalance query
type BalancePayload struct {
	Account string `json:"account"`
}

// MintPayload for mint action
type MintPayload struct {
	Amount *big.Int `json:"-"`
}

// BurnPayload for burn action
type BurnPayload struct {
	Amount *big.Int `json:"-"`
}

// ChangeOwnerPayload for changeOwner action
type ChangeOwnerPayload struct {
	NewOwner string `json:"newOwner"`
}

// ===================================
// Response Types (Output)
// ===================================

// BalanceResponse for balance queries
type BalanceResponse struct {
	Balance *big.Int `json:"-"`
}

// SupplyResponse for supply queries
type SupplyResponse struct {
	TotalSupply *big.Int `json:"-"`
}

// AllowanceResponse for allowance queries
type AllowanceResponse struct {
	Allowance *big.Int `json:"-"`
}

// OwnerResponse for owner queries
type OwnerResponse struct {
	Owner string `json:"owner"`
}

// InfoResponse for token info queries.
// Decimals is encoded as a decimal string so consumers (e.g. the DEX router's
// register_token) that expect string-typed numerics can parse it directly.
type InfoResponse struct {
	Name      string   `json:"name"`
	Symbol    string   `json:"symbol"`
	Decimals  string   `json:"decimals"`
	MaxSupply *big.Int `json:"-"`
}

// PausedResponse for isPaused queries
type PausedResponse struct {
	Paused bool `json:"paused"`
}

// SuccessResponse for mutation operations
type SuccessResponse struct {
	Success bool `json:"success"`
}

// ===================================
// Event Types
// ===================================

// InitEvent for contract initialization
type InitEvent struct {
	Type       string         `json:"type"`
	Attributes InitAttributes `json:"attributes"`
}

type InitAttributes struct {
	Owner     string   `json:"owner"`
	Name      string   `json:"name"`
	Symbol    string   `json:"symbol"`
	Decimals  int      `json:"decimals"`
	MaxSupply *big.Int `json:"-"`
}

// TransferEvent for token transfers
type TransferEvent struct {
	Type       string             `json:"type"`
	Attributes TransferAttributes `json:"attributes"`
}

type TransferAttributes struct {
	From   string   `json:"from"`
	To     string   `json:"to"`
	Amount *big.Int `json:"-"`
}

// ApprovalEvent for allowance approvals
type ApprovalEvent struct {
	Type       string             `json:"type"`
	Attributes ApprovalAttributes `json:"attributes"`
}

type ApprovalAttributes struct {
	Owner   string   `json:"owner"`
	Spender string   `json:"spender"`
	Amount  *big.Int `json:"-"`
}

// OwnerChangeEvent for ownership transfers
type OwnerChangeEvent struct {
	Type       string                `json:"type"`
	Attributes OwnerChangeAttributes `json:"attributes"`
}

type OwnerChangeAttributes struct {
	PreviousOwner string `json:"previousOwner"`
	NewOwner      string `json:"newOwner"`
}

// PausedEvent for pause action
type PausedEvent struct {
	Type       string           `json:"type"`
	Attributes PausedAttributes `json:"attributes"`
}

type PausedAttributes struct {
	By string `json:"by"`
}

// UnpausedEvent for unpause action
type UnpausedEvent struct {
	Type       string             `json:"type"`
	Attributes UnpausedAttributes `json:"attributes"`
}

type UnpausedAttributes struct {
	By string `json:"by"`
}
