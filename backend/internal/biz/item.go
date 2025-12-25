package biz

import (
	"context"

	"github.com/tomeai/dataflow/api/v1/flow"

	"github.com/go-kratos/kratos/v2/log"
)

// ItemRepo is a Greater repo.
type ItemRepo interface {
	QueryDataSinkById(context.Context, string) (*flow.Sink, error)
	GetSinks(context.Context) ([]*flow.Sink, error)
}

type ItemUseCase struct {
	repo ItemRepo
	log  *log.Helper
}

func NewSourceUseCase(repo ItemRepo, logger log.Logger) *ItemUseCase {
	return &ItemUseCase{repo: repo, log: log.NewHelper(logger)}
}

func (s ItemUseCase) QueryDataSinkById(ctx context.Context, id string) (*flow.Sink, error) {
	return s.repo.QueryDataSinkById(ctx, id)
}

func (s ItemUseCase) GetSinks(ctx context.Context) ([]*flow.Sink, error) {
	return s.repo.GetSinks(ctx)
}
