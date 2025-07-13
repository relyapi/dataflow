package biz

import (
	"context"
	"data-flow/api/v1/flow"

	"github.com/go-kratos/kratos/v2/log"
)

// ItemRepo is a Greater repo.
type ItemRepo interface {
	QueryDataSinkById(context.Context, string) (*flow.Sink, error)
	GetAllSink(context.Context) ([]*flow.Sink, error)
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

func (s ItemUseCase) GetAllSink(ctx context.Context) ([]*flow.Sink, error) {
	return s.repo.GetAllSink(ctx)
}
