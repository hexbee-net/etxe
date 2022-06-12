package dag

import (
	"errors"
)

var (
	ErrVertexIDEmpty             = errors.New("vertex ID is empty")
	ErrVertexDuplicate           = errors.New("vertex is already in the graph")
	ErrVertexNotFound            = errors.New("vertex is not in the graph")
	ErrEdgeSourceIDEmpty         = errors.New("edge source ID is empty")
	ErrEdgeTargetIDEmpty         = errors.New("edge target ID is empty")
	ErrEdgeSourceIDNotFound      = errors.New("edge source ID is not in the graph")
	ErrEdgeTargetIDNotFound      = errors.New("edge target ID is not in the graph")
	ErrEdgeSourceTargetIdentical = errors.New("edge source and target IDs are identical")
	ErrEdgeDuplicate             = errors.New("edge is already in the graph")
	ErrEdgeNotFound              = errors.New("edge is not in the graph")
	ErrEdgeLoop                  = errors.New("edge would create a loop")
)
