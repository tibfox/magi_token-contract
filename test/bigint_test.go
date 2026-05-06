package contract_test

import (
	"testing"
)

// ===================================
// Big Integer / String Amount Tests
// ===================================
// These tests validate that the contract correctly handles amounts as
// string-encoded integers, including values beyond uint64 range.

// --- Large Number Init & Mint ---

func TestInitWithLargeMaxSupply(t *testing.T) {
	ct := SetupContractTest()
	// Max supply beyond uint64 (2^64 = 18446744073709551616)
	payload := []byte(`{"name":"Big Token","symbol":"BIG","decimals":18,"maxSupply":"100000000000000000000000000000"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(100_000_000), "")
	// Verify info returns the large maxSupply as string
	CallContract(t, ct, "getInfo", []byte(""), nil, "hive:anyone", true, uint(100_000_000),
		`{"name":"Big Token","symbol":"BIG","decimals":"18","maxSupply":"100000000000000000000000000000"}`)
}

func TestMintAndTransferLargeAmounts(t *testing.T) {
	ct := SetupContractTest()
	// Init with huge max supply
	payload := []byte(`{"name":"Big Token","symbol":"BIG","decimals":18,"maxSupply":"100000000000000000000000000000"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(100_000_000), "")

	// Mint a large amount (beyond uint64)
	CallContract(t, ct, "mint", []byte(`{"amount":"50000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000),
		`{"totalSupply":"50000000000000000000000000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"50000000000000000000000000000"}`)

	// Transfer a large amount
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:alice","amount":"20000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"30000000000000000000000000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:alice"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"20000000000000000000000000000"}`)
}

func TestBurnLargeAmount(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Big Token","symbol":"BIG","decimals":18,"maxSupply":"100000000000000000000000000000"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"50000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")

	// Burn a large amount
	CallContract(t, ct, "burn", []byte(`{"amount":"10000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000),
		`{"totalSupply":"40000000000000000000000000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"40000000000000000000000000000"}`)
}

// --- Large Allowances ---

func TestAllowanceLargeAmounts(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Big Token","symbol":"BIG","decimals":18,"maxSupply":"100000000000000000000000000000"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"50000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")

	// Approve large allowance
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:dex","amount":"25000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"allowance":"25000000000000000000000000000"}`)

	// TransferFrom with large amount
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:tibfox","to":"hive:bob","amount":"5000000000000000000000000000"}`), nil, "hive:dex", true, uint(100_000_000), "")
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:bob"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"5000000000000000000000000000"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"allowance":"20000000000000000000000000000"}`)
}

func TestIncreaseDecreaseAllowanceLarge(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Big Token","symbol":"BIG","decimals":18,"maxSupply":"100000000000000000000000000000"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(100_000_000), "")

	// Start with large allowance
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:dex","amount":"10000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Increase by large amount
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:dex","amount":"5000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"allowance":"15000000000000000000000000000"}`)

	// Decrease by large amount
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:dex","amount":"3000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:tibfox","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"allowance":"12000000000000000000000000000"}`)
}

// --- Mint Exceeds Max Supply With Large Numbers ---

func TestMintExceedsLargeMaxSupply(t *testing.T) {
	ct := SetupContractTest()
	payload := []byte(`{"name":"Big Token","symbol":"BIG","decimals":18,"maxSupply":"100000000000000000000000000000"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(100_000_000), "")

	// Mint exactly max supply
	CallContract(t, ct, "mint", []byte(`{"amount":"100000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Cannot mint 1 more
	CallContract(t, ct, "mint", []byte(`{"amount":"1"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Invalid String Amount Inputs ---

func TestMintNegativeAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Negative amount should fail
	CallContract(t, ct, "mint", []byte(`{"amount":"-100"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestTransferNegativeAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "mint", []byte(`{"amount":"1000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	// Negative amount should fail
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"-50"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestMintNonNumericAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Non-numeric string should fail
	CallContract(t, ct, "mint", []byte(`{"amount":"abc"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

func TestMintDecimalAmount(t *testing.T) {
	ct := SetupContractTest()
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")
	// Decimal amount should fail (must be integer)
	CallContract(t, ct, "mint", []byte(`{"amount":"100.5"}`), nil, ownerAddress, false, uint(100_000_000), "")
}

// --- Full Flow With Large Numbers ---

func TestFullFlowLargeNumbers(t *testing.T) {
	ct := SetupContractTest()
	// Init with 256-bit scale max supply (requires more gas due to large number parsing)
	payload := []byte(`{"name":"Big Token","symbol":"BIG","decimals":18,"maxSupply":"115792089237316195423570985008687907853269984665640564039457584007913129639935"}`)
	CallContract(t, ct, "init", payload, nil, ownerAddress, true, uint(200_000_000), "")

	// Mint a very large amount
	CallContract(t, ct, "mint", []byte(`{"amount":"1000000000000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")

	// Transfer half
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:alice","amount":"500000000000000000000000000000000000"}`), nil, ownerAddress, true, uint(100_000_000), "")

	// Verify balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"500000000000000000000000000000000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:alice"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"500000000000000000000000000000000000"}`)

	// Alice burns some
	CallContract(t, ct, "burn", []byte(`{"amount":"100000000000000000000000000000000000"}`), nil, "hive:alice", true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000),
		`{"totalSupply":"900000000000000000000000000000000000"}`)

	// Approve and transferFrom with large numbers
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:dex","amount":"200000000000000000000000000000000000"}`), nil, "hive:alice", true, uint(100_000_000), "")
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:alice","to":"hive:bob","amount":"150000000000000000000000000000000000"}`), nil, "hive:dex", true, uint(100_000_000), "")

	// Verify final state
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:alice"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"250000000000000000000000000000000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:bob"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"balance":"150000000000000000000000000000000000"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:alice","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000),
		`{"allowance":"50000000000000000000000000000000000"}`)
}
