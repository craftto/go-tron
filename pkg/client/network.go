package client

import (
	"fmt"

	"github.com/craftto/go-tron/pkg/proto/api"
	"github.com/craftto/go-tron/pkg/proto/core"
	"go.uber.org/zap"
)

// Broadcast broadcast TX
func (g *GrpcClient) Broadcast(tx *core.Transaction) (*api.Return, error) {
	ctx, cancel := g.GetContext()
	defer cancel()

	result, err := g.Client.BroadcastTransaction(ctx, tx)
	if err != nil {
		return nil, err
	}

	if !result.GetResult() {
		return result, fmt.Errorf("result error: %s", result.GetMessage())
	}

	if result.GetCode() != api.Return_SUCCESS {
		return result, fmt.Errorf("result error(%s): %s", result.GetCode(), result.GetMessage())
	}

	return result, nil
}

// ListNodes provides list of network nodes
func (g *GrpcClient) ListNodes() (*api.NodeList, error) {
	ctx, cancel := g.GetContext()
	defer cancel()

	nodeList, err := g.Client.ListNodes(ctx, new(api.EmptyMessage))
	if err != nil {
		zap.L().Error("List nodes", zap.Error(err))
	}

	return nodeList, nil
}

// GetNextMaintenanceTime get next epoch timestamp
func (g *GrpcClient) GetNextMaintenanceTime() (*api.NumberMessage, error) {
	ctx, cancel := g.GetContext()
	defer cancel()

	return g.Client.GetNextMaintenanceTime(ctx, new(api.EmptyMessage))
}
