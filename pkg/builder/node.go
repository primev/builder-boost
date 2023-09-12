package builder

import (
	"context"
	"io"

	"golang.org/x/exp/slog"
)

// Options are the options for the builder node
type Options struct {
	Logger               *slog.Logger
	ServerAddr           string
	TracingEndpoint      string
	RollupKey            string
	RollupAddr           string
	RollupContract       string
	BuilderToken         string
	MetricsEnabled       bool
	DAContract           string
	DARPCAddr            string
	InclusionListEnabled bool
}

type Node struct {
	io.Closer
}

func NewNode(
	ctx context.Context,
	opts Options,
) (*Node, error) {
	return &Node{}, nil
}
