package client

import (
	"bytes"
	"context"
	"io"

	"github.com/bnb-chain/greenfield-storage-provider/pkg/log"
	"github.com/bnb-chain/greenfield-storage-provider/store/piecestore/piece"
	"github.com/bnb-chain/greenfield-storage-provider/store/piecestore/storage"
)

// PieceStoreAPI provides an interface to enable mocking the
// StoreClient's API operation. This makes unit test to test your code easier.
//
//go:generate mockgen -source=./piece_store_client.go -destination=./mock/piece_store_mock.go -package=mock
type PieceStoreAPI interface {
	GetPiece(ctx context.Context, key string, offset, limit int64) ([]byte, error)
	PutPiece(key string, value []byte) error
}

type StoreClient struct {
	ps *piece.PieceStore
}

func NewStoreClient(pieceConfig *storage.PieceStoreConfig) (*StoreClient, error) {
	ps, err := piece.NewPieceStore(pieceConfig)
	if err != nil {
		return nil, err
	}
	return &StoreClient{ps: ps}, nil
}

func (client *StoreClient) GetPiece(ctx context.Context, key string, offset, limit int64) ([]byte, error) {
	rc, err := client.ps.Get(ctx, key, offset, limit)
	if err != nil {
		log.Errorw("get piece data from piece store failed", "error", err)
		return nil, err
	}
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, rc)
	if err != nil {
		log.Errorw("copy data failed", "error", err)
		return nil, err
	}
	return buf.Bytes(), nil
}

func (client *StoreClient) PutPiece(key string, value []byte) error {
	return client.ps.Put(context.Background(), key, bytes.NewReader(value))
}
