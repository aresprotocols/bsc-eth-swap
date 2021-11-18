package executor

import (
	common "github.com/binance-chain/bsc-eth-swap/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethcmm "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"math/big"

	"github.com/binance-chain/bsc-eth-swap/model"
)

type Executor interface {
	GetBlockAndTxEvents(height int64) (*common.BlockAndEventLogs, error)
	GetChainName() string
	GetExplorerUrl() string
	GetTokenBalance(contract ethcmm.Address, who ethcmm.Address) (*big.Int, error)
}

// ===================  SwapStarted =============
var (
	SwapStartedEventName        = "SwapStarted"
	ETH2BSCSwapStartedEventHash = ethcmm.HexToHash("0xf60309f865a6aa297da5fac6188136a02e5acfdf6e8f6d35257a9f4e9653170f")
	BSC2ETHSwapStartedEventHash = ethcmm.HexToHash("0x49c08ff11118922c1e8298915531eff9ef6f8b39b44b3e9952b75d47e1d0cdd0")
)

type ETH2BSCSwapStartedEvent struct {
	ERC20Addr ethcmm.Address
	FromAddr  ethcmm.Address
	Amount    *big.Int
	FeeAmount *big.Int
}

func (ev *ETH2BSCSwapStartedEvent) ToSwapStartTxLog(log *types.Log) *model.SwapStartTxLog {
	pack := &model.SwapStartTxLog{
		TokenAddr:   ev.ERC20Addr.String(),
		FromAddress: ev.FromAddr.String(),
		Amount:      ev.Amount.String(),

		FeeAmount: ev.FeeAmount.String(),
		BlockHash: log.BlockHash.Hex(),
		TxHash:    log.TxHash.String(),
		Height:    int64(log.BlockNumber),
	}
	return pack
}

func ParseETH2BSCSwapStartEvent(abi *abi.ABI, log *types.Log) (*ETH2BSCSwapStartedEvent, error) {
	var ev ETH2BSCSwapStartedEvent

	err := abi.UnpackIntoInterface(&ev, SwapStartedEventName, log.Data)
	if err != nil {
		return nil, err
	}

	ev.ERC20Addr = ethcmm.BytesToAddress(log.Topics[1].Bytes())
	ev.FromAddr = ethcmm.BytesToAddress(log.Topics[2].Bytes())

	return &ev, nil
}

type BSC2ETHSwapStartedEvent struct {
	BEP20Addr ethcmm.Address
	ERC20Addr ethcmm.Address
	FromAddr  ethcmm.Address
	Amount    *big.Int
	FeeAmount *big.Int
}

func (ev *BSC2ETHSwapStartedEvent) ToSwapStartTxLog(log *types.Log) *model.SwapStartTxLog {
	pack := &model.SwapStartTxLog{
		TokenAddr:   ev.BEP20Addr.String(),
		FromAddress: ev.FromAddr.String(),
		Amount:      ev.Amount.String(),

		FeeAmount: ev.FeeAmount.String(),
		BlockHash: log.BlockHash.Hex(),
		TxHash:    log.TxHash.String(),
		Height:    int64(log.BlockNumber),
	}
	return pack
}

func ParseBSC2ETHSwapStartEvent(abi *abi.ABI, log *types.Log) (*BSC2ETHSwapStartedEvent, error) {
	var ev BSC2ETHSwapStartedEvent

	err := abi.UnpackIntoInterface(&ev, SwapStartedEventName, log.Data)
	if err != nil {
		return nil, err
	}

	ev.BEP20Addr = ethcmm.BytesToAddress(log.Topics[1].Bytes())
	ev.ERC20Addr = ethcmm.BytesToAddress(log.Topics[2].Bytes())
	ev.FromAddr = ethcmm.BytesToAddress(log.Topics[3].Bytes())

	return &ev, nil
}

// =================  SwapPairRegister ===================
var (
	SwapPairRegisterEventName = "SwapPairRegister"
	SwapPairRegisterEventHash = ethcmm.HexToHash("0x2bf34a4c27e048d089386ccfb7b7d28131d7107d7b168300bf4f48a66b3b0b14")
)

type SwapPairRegisterEvent struct {
	Sponsor            ethcmm.Address
	ContractAddr       ethcmm.Address
	TargetContractAddr ethcmm.Address
	Name               string
	Symbol             string
	Decimals           uint8
}

func (ev *SwapPairRegisterEvent) ToSwapPairRegisterLog(log *types.Log) *model.SwapPairRegisterTxLog {
	pack := &model.SwapPairRegisterTxLog{
		ERC20Addr: ev.ContractAddr.String(),
		BEP20Addr: ev.TargetContractAddr.String(),
		Sponsor:   ev.Sponsor.String(),
		Symbol:    ev.Symbol,
		Name:      ev.Name,
		Decimals:  int(ev.Decimals),

		BlockHash: log.BlockHash.Hex(),
		TxHash:    log.TxHash.String(),
		Height:    int64(log.BlockNumber),
	}
	return pack
}

func ParseSwapPairRegisterEvent(abi *abi.ABI, log *types.Log) (*SwapPairRegisterEvent, error) {
	var ev SwapPairRegisterEvent

	err := abi.UnpackIntoInterface(&ev, SwapPairRegisterEventName, log.Data)
	if err != nil {
		return nil, err
	}
	ev.Sponsor = ethcmm.BytesToAddress(log.Topics[1].Bytes())
	ev.ContractAddr = ethcmm.BytesToAddress(log.Topics[2].Bytes())
	ev.TargetContractAddr = ethcmm.BytesToAddress(log.Topics[3].Bytes())
	return &ev, nil
}
