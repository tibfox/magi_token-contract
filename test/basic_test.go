package contract_test

import (
	"testing"
)

// ===================================
// Core Token Tests (Init, Mint, Transfer, Burn, ChangeOwner)
// ===================================

// --- Init Tests ---

func TestInitSuccess(t *testing.T) {
	ct := SetupContractTest()
	// Init should succeed when called by contract owner
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Verify owner is set
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:tibfox"}`)
	// Verify supply is 0
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"0"}`)
}

func TestInitFailsNonCreator(t *testing.T) {
	ct := SetupContractTest()
	// Init should fail when called by non-owner
	CallContract(t, ct, "init", DefaultInitPayload, nil, "hive:someone", false, uint(100_000_000), "")
}

func TestInitFailsDouble(t *testing.T) {
	ct := SetupContractTest()
	// First init succeeds
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Second init fails
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Mint Tests ---

func TestMintSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Mint 1000 tokens
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check supply
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"1000"}`)
	// Check balance
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"1000"}`)
}

func TestMintFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Non-owner cannot mint
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, "hive:someone", false, uint(100_000_000), "")
}

func TestMintFailsNotInitialized(t *testing.T) {
	ct := SetupContractTest()
	// Cannot mint before init
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestMintFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot mint 0 tokens
	CallContract(t, ct, "mint", []byte(`{"amount":"0"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestMintFailsExceedsMaxSupply(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot mint more than max supply (1 billion)
	CallContract(t, ct, "mint", []byte(`{"amount":"1000000001"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Transfer Tests ---

func TestTransferSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer 500 to someone
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"500"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"500"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:someone"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"500"}`)
}

func TestTransferFailsInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot transfer more than balance
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"500"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFailsToSelf(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot transfer to self
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:tibfox","amount":"500"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot transfer 0 tokens
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"0"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Burn Tests ---

func TestBurnSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Burn 300 tokens
	CallContract(t, ct, "burn", []byte(`{"amount":"300"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Check supply and balance
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"700"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"700"}`)
}

func TestBurnFailsInsufficientBalance(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"100"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot burn more than balance
	CallContract(t, ct, "burn", []byte(`{"amount":"500"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestBurnFailsZeroAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot burn 0 tokens
	CallContract(t, ct, "burn", []byte(`{"amount":"0"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- ChangeOwner Tests ---

func TestChangeOwnerSuccess(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Change owner
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:newowner"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Verify new owner
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:newowner"}`)
	// Old owner cannot mint
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, false, uint(100_000_000), "")
	// New owner can mint
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, "hive:newowner", true, uint(100_000_000), "")
}

func TestChangeOwnerFailsNonOwner(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Non-owner cannot change owner
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:newowner"}`), nil, "hive:someone", false, uint(100_000_000), "")
}

// --- GetInfo Test ---

func TestGetInfo(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Get token info
	CallContract(t, ct, "getInfo", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"name":"Magi Token","symbol":"MAGI","decimals":"3","maxSupply":"1000000000"}`)
}

// --- Integration Test - Full Flow ---

func TestFullFlow(t *testing.T) {
	ct := SetupContractTest()
	// Initialize
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Mint tokens
	CallContract(t, ct, "mint", []byte(`{"amount":"10000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer to user A
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:userA","amount":"3000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Transfer to user B
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:userB","amount":"2000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// User A transfers to user B
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:userB","amount":"1000"}`), nil, "hive:userA", true, uint(100_000_000), "")
	// User B burns some
	CallContract(t, ct, "burn", []byte(`{"amount":"500"}`), nil, "hive:userB", true, uint(100_000_000), "")
	// Check final balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"5000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:userA"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"2000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:userB"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"2500"}`)
	// Check final supply
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"9500"}`)
}
