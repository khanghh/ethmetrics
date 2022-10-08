package core

import (
	"context"
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/metrics"
)

type Block struct {
	*types.Header
	Hash         common.Hash
	Size         common.StorageSize
	Transactions []common.Hash
}

func (s *Block) UnmarshalJSON(data []byte) error {
	s.Header = &types.Header{}
	var aux struct {
		Hash         common.Hash   `json:"hash"`
		Size         string        `json:"size"`
		Transactions []common.Hash `json:"transactions"`
	}
	if err := s.Header.UnmarshalJSON(data); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	s.Hash = aux.Hash
	size, err := hexutil.DecodeUint64(aux.Size)
	if err != nil {
		return err
	}
	s.Size = common.StorageSize(size)
	s.Transactions = aux.Transactions
	return nil
}

type MetricsCollector interface {
	Setup(ctx *Ctx) error
	Collect(ctx *Ctx)
}

type MetricsPublisher interface {
	PublishMetrics(ctx context.Context, registry metrics.Registry)
}
