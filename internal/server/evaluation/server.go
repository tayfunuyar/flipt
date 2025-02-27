package evaluation

import (
	"context"

	"go.flipt.io/flipt/internal/storage"
	"go.flipt.io/flipt/rpc/flipt"
	"go.flipt.io/flipt/rpc/flipt/evaluation"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Storer is the minimal abstraction for interacting with the storage layer for evaluation.
type Storer interface {
	GetFlag(ctx context.Context, namespaceKey, key string) (*flipt.Flag, error)
	GetEvaluationRules(ctx context.Context, namespaceKey string, flagKey string) ([]*storage.EvaluationRule, error)
	GetEvaluationDistributions(ctx context.Context, ruleID string) ([]*storage.EvaluationDistribution, error)
	GetEvaluationRollouts(ctx context.Context, namespaceKey, flagKey string) ([]*storage.EvaluationRollout, error)
}

// Server serves the Flipt evaluate v2 gRPC Server.
type Server struct {
	logger    *zap.Logger
	store     Storer
	evaluator *Evaluator
	evaluation.UnimplementedEvaluationServiceServer
}

// New is constructs a new Server.
func New(logger *zap.Logger, store Storer) *Server {
	return &Server{
		logger:    logger,
		store:     store,
		evaluator: NewEvaluator(logger, store),
	}
}

// RegisterGRPC registers the EvaluateServer onto the provided gRPC Server.
func (s *Server) RegisterGRPC(server *grpc.Server) {
	evaluation.RegisterEvaluationServiceServer(server, s)
}
