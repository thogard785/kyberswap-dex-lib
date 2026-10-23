# ThogAMM schema-7 adapter

ThogAMM is a portfolio-inventory AMM on Monad (chain ID 143). Integrate its proxy,
`0x80c74517bcc2d67ffe02d3ed886796272f647210`; the verified implementation is
`0x709ba37580a7faf034522163353bd1547d2fb0ba`. Pool ID:
`0xce389e78282dedac7b18ba7f775b7602d2ab3ab171bbd6711eb0239be6ef4dcc`.

[Integration guide and exact integer specification](https://github.com/thogard785/thog-amm-integration-guide/blob/main/THOGAMM_AGGREGATOR_INTEGRATION.md)
· [Proxy on Monadscan](https://monadscan.com/address/0x80c74517bcc2d67ffe02d3ed886796272f647210)
· [Reference Rust model](https://github.com/thogard785/thog-amm-poll-sdk/tree/v0.3.0/model)

Discovery reads `getPoolIds` and `getTokens`. Each refresh makes one atomic
`getPoolData(0,32)` call, receiving all 126 words, native MON, all ERC-20 balances
and dynamic token/category metadata. A refresh can replace or reorder the
registry. Schema 7 / wire 7 / model 1 is required; legacy `makerSnapshot` state
cannot price the deployed model.

The tracker requests another refresh after one second even without pool events;
refresh on each new head when possible. `entity.Pool.Timestamp` records receipt
time and `BlockNumber` identifies the actual quoted state. Tracker errors must
remain visible to the pool service; an old serialized state remains a historical
quote, not evidence of a live feed. The simulator uses the contract's exact
block-age rule and does not convert elapsed wall time into invented blocks.

Prices use 56-bit mantissas. The curve charges separate bid/ask/category/pair
widths, positive covariance-inventory risk, an independent quadratic input
volume charge and integrated wrapper costs. Sliding output refunds participate
in an exact integer quadratic solve, including the full discriminant and
adjacent-integer check. Arithmetic uses checked `uint256.Int` / `int256.Int`;
there is no floating point or `big.Int` arithmetic in the quote model.

`CalcAmountOut` matches the read-only `makerQuoteExactInput` at the captured block,
including its strict 90% receiving-token/category caps. It does not assume a
future block or a warm AuctionHandler. For an explicitly known execution
context, `State.QuoteExecution` and `State.QuoteExactOutput` reproduce the
contract's single 10 bp reduction when gas price exceeds twice base fee or the
AuctionHandler is warm; settlement checks the full caps. `State.MarginalPrice`
and `State.Limits` reproduce the ERC-7815 read-only surfaces. Router execution
still needs the same complete state and context, input approval/funds, deadline
and minimum-output protection.

`CalcAmountOut` returns a `SwapInfo` containing the exact next balances.
`UpdateBalance` applies those values without recomputing; subsequent quotes
rebuild risk and wrapper inventory from the updated complete portfolio.
Clones copy writable balance state, and quote calls do not mutate it. Gas estimates
are 350,000 within a risk category and 550,000 across categories. The deployed
cold-call measurements span 222,072–323,820 and 343,043–521,306 respectively,
across all directed pairs and $1/$100/$5,000 sizes at block 115786825. These use
Foundry 1.8.6 with current Monad page-storage pricing; the earlier 660,000
estimate used an older simulator. Estimate the complete router transaction
before submission and remeasure after catalog, policy or token-code changes.

## Verification

The same independently captured corpus tests this Go model and both Rust SDKs:

| Corpus | Cases | Exact successes | Matching rejections |
| --- | ---: | ---: | ---: |
| Three mainnet snapshots, every directed pair | 2,016 | 1,308 | 708 |
| Deployed-bytecode Monad-fork executions, exact input/output and gas/warmth context | 1,680 | 1,261 | 419 |
| Fork boundary states: volume, fixed/sliding wrappers, caps, price age | 5,544 | 2,793 | 2,751 |

Tests compare integer amounts, raw marginal fractions, limits and rejection
behavior. Tracker tests verify the single RPC call and serialized-state round
trip; simulator tests verify purity, clone isolation and exact balance updates.
The fixture directory has a nested module so downstream module downloads omit
it. [Capture tools and provenance](https://github.com/thogard785/thog-amm-integration-guide#verification)
are public; all fork changes are local and no transactions are broadcast.

```sh
GOTOOLCHAIN=go1.27.2 go test ./pkg/liquidity-source/thog-prop -race
GOTOOLCHAIN=go1.27.2 go vet ./pkg/liquidity-source/thog-prop
GOTOOLCHAIN=go1.27.2 go test ./pkg/liquidity-source/thog-prop -run '^$' -bench . -benchmem
```
