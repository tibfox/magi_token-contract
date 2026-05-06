package contract_test

import (
	"testing"
)

// ===================================
// Token Lifecycle Test
// ===================================
// This test simulates a realistic token lifecycle from creation through
// various usage patterns including minting, transfers, allowances, burns,
// and supply updates. It covers common DeFi patterns like DEX swaps,
// treasury management, and community distributions.

func TestTokenLifecycle(t *testing.T) {
	ct := SetupContractTest()

	// =========================================
	// PHASE 1: Contract Creation & Initial Setup
	// =========================================

	// Initialize contract - contract owner becomes token owner
	CallContract(t, ct, "init", DefaultInitPayload, nil, ownerAddress, true, uint(100_000_000), "")

	// Verify initial state
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:tibfox"}`)
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"0"}`)
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":false}`)
	CallContract(t, ct, "getInfo", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"name":"Magi Token","symbol":"MAGI","decimals":"3","maxSupply":"1000000000"}`)

	// =========================================
	// PHASE 2: Initial Token Distribution
	// =========================================

	// Owner mints initial supply for treasury (100,000 tokens)
	CallContract(t, ct, "mint", []byte(`{"amount":"100000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"100000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"100000000"}`)

	// Distribute to team members
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:alice","amount":"10000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:bob","amount":"10000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:charlie","amount":"5000000"}`), nil, ownerAddress, true, uint(100_000_000), "")

	// Verify distribution
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"75000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:alice"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"10000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:bob"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"10000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:charlie"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"5000000"}`)

	// =========================================
	// PHASE 3: DEX Integration (Allowance Pattern)
	// =========================================

	// Alice approves DEX to spend her tokens for trading
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:dex","amount":"5000000"}`), nil, "hive:alice", true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:alice","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"5000000"}`)

	// Bob also approves DEX
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:dex","amount":"3000000"}`), nil, "hive:bob", true, uint(100_000_000), "")

	// DEX executes swap: Alice sells 2M tokens to Dave
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:alice","to":"hive:dave","amount":"2000000"}`), nil, "hive:dex", true, uint(100_000_000), "")

	// Verify post-swap state
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:alice"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"8000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:dave"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"2000000"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:alice","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"3000000"}`)

	// DEX executes another swap: Bob sells 1.5M tokens to Eve
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:bob","to":"hive:eve","amount":"1500000"}`), nil, "hive:dex", true, uint(100_000_000), "")

	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:bob"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"8500000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:eve"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"1500000"}`)

	// =========================================
	// PHASE 4: Peer-to-Peer Transfers
	// =========================================

	// Dave sends tokens to multiple recipients
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:frank","amount":"500000"}`), nil, "hive:dave", true, uint(100_000_000), "")
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:grace","amount":"300000"}`), nil, "hive:dave", true, uint(100_000_000), "")

	// Charlie sends to Eve
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:eve","amount":"1000000"}`), nil, "hive:charlie", true, uint(100_000_000), "")

	// Verify balances
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:dave"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"1200000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:frank"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"500000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:grace"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"300000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:eve"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"2500000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:charlie"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"4000000"}`)

	// =========================================
	// PHASE 5: Allowance Management (Increase/Decrease)
	// =========================================

	// Alice increases her DEX allowance for more trading
	CallContract(t, ct, "increaseAllowance", []byte(`{"spender":"hive:dex","amount":"2000000"}`), nil, "hive:alice", true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:alice","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"5000000"}`)

	// Bob decides to reduce his DEX allowance
	CallContract(t, ct, "decreaseAllowance", []byte(`{"spender":"hive:dex","amount":"1000000"}`), nil, "hive:bob", true, uint(100_000_000), "")
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:bob","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"500000"}`)

	// Eve approves a payment processor
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:payments","amount":"1000000"}`), nil, "hive:eve", true, uint(100_000_000), "")

	// Payment processor charges Eve for a service
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:eve","to":"hive:merchant","amount":"250000"}`), nil, "hive:payments", true, uint(100_000_000), "")

	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:eve"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"2250000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:merchant"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"250000"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:eve","spender":"hive:payments"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"750000"}`)

	// =========================================
	// PHASE 6: Token Burns (Deflationary Mechanism)
	// =========================================

	// Treasury burns some tokens to reduce supply
	CallContract(t, ct, "burn", []byte(`{"amount":"5000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"95000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"70000000"}`)

	// Alice burns some of her tokens
	CallContract(t, ct, "burn", []byte(`{"amount":"1000000"}`), nil, "hive:alice", true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"94000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:alice"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"7000000"}`)

	// =========================================
	// PHASE 7: Additional Minting (Supply Expansion)
	// =========================================

	// Owner mints additional tokens for a new partnership
	CallContract(t, ct, "mint", []byte(`{"amount":"10000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"104000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"80000000"}`)

	// Distribute new tokens to partner
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:partner","amount":"10000000"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:partner"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"10000000"}`)

	// =========================================
	// PHASE 8: Emergency Pause & Recovery
	// =========================================

	// Owner pauses contract due to security concern
	CallContract(t, ct, "pause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":true}`)

	// Transfers should fail while paused
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:someone","amount":"1000"}`), nil, "hive:alice", false, uint(100_000_000), "")
	CallContract(t, ct, "transferFrom", []byte(`{"from":"hive:alice","to":"hive:bob","amount":"1000"}`), nil, "hive:dex", false, uint(100_000_000), "")

	// But allowance management still works
	CallContract(t, ct, "approve", []byte(`{"spender":"hive:newdex","amount":"1000000"}`), nil, "hive:alice", true, uint(100_000_000), "")

	// Owner unpauses after security review
	CallContract(t, ct, "unpause", []byte(""), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":false}`)

	// Transfers work again
	CallContract(t, ct, "transfer", []byte(`{"to":"hive:bob","amount":"500000"}`), nil, "hive:alice", true, uint(100_000_000), "")

	// =========================================
	// PHASE 9: Ownership Transfer
	// =========================================

	// Transfer ownership to a multisig/DAO
	CallContract(t, ct, "changeOwner", []byte(`{"newOwner":"hive:dao"}`), nil, ownerAddress, true, uint(100_000_000), "")
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:dao"}`)

	// Old owner can no longer mint
	CallContract(t, ct, "mint", []byte(`{"amount":"1000000"}`), nil, ownerAddress, false, uint(100_000_000), "")

	// New owner (DAO) can mint
	CallContract(t, ct, "mint", []byte(`{"amount":"5000000"}`), nil, "hive:dao", true, uint(100_000_000), "")
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"109000000"}`)

	// =========================================
	// PHASE 10: Final State Verification
	// =========================================

	// Verify final balances of all participants
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:tibfox"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"70000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:dao"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"5000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:alice"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"6500000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:bob"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"9000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:charlie"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"4000000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:dave"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"1200000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:eve"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"2250000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:frank"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"500000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:grace"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"300000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:merchant"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"250000"}`)
	CallContract(t, ct, "balanceOf", []byte(`{"account":"hive:partner"}`), nil, "hive:anyone", true, uint(100_000_000), `{"balance":"10000000"}`)

	// Verify final supply matches sum of all balances
	// Total: 70M + 5M + 6.5M + 9M + 4M + 1.2M + 2.25M + 0.5M + 0.3M + 0.25M + 10M = 109M
	CallContract(t, ct, "totalSupply", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"totalSupply":"109000000"}`)

	// Verify remaining allowances
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:alice","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"5000000"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:alice","spender":"hive:newdex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"1000000"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:bob","spender":"hive:dex"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"500000"}`)
	CallContract(t, ct, "allowance", []byte(`{"owner":"hive:eve","spender":"hive:payments"}`), nil, "hive:anyone", true, uint(100_000_000), `{"allowance":"750000"}`)

	// Verify contract state
	CallContract(t, ct, "getOwner", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"owner":"hive:dao"}`)
	CallContract(t, ct, "isPaused", []byte(""), nil, "hive:anyone", true, uint(100_000_000), `{"paused":false}`)
}
