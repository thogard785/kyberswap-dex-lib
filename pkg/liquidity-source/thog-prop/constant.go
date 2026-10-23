package thogprop

import (
	"errors"

	"github.com/KyberNetwork/kyberswap-dex-lib/pkg/valueobject"
)

const (
	DexType           = valueobject.ExchangeThogProp
	methodGetPoolData = "getPoolData"
	maxTokens         = 32
	wmon              = "0x3bd359c1119da7da1d913d1c4d2b7c461115433a"
	// Measured using Foundry 1.8.6 with Monad page-storage pricing and cold calls.
	// Same-category swaps skip the covariance-risk work; see testdata/gas-v7.json.gz.
	sameCategoryGas  = 350000
	crossCategoryGas = 550000
)

var (
	ErrInvalidState = errors.New("thog-prop: invalid schema-7 snapshot")
	ErrInvalidToken = errors.New("thog-prop: unknown or identical token")
	ErrUnavailable  = errors.New("thog-prop: quote unavailable")
	ErrOverflow     = errors.New("thog-prop: arithmetic overflow or division by zero")
	ErrBalance      = errors.New("thog-prop: inventory capacity exceeded")
)
