package dag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildTestGraph(t *testing.T, vertices map[string]int, inbound map[string]Set[string], outbound map[string]Set[string]) *AcyclicGraph[string, int] {
	t.Helper()

	g := NewAcyclicGraph[string, int]()
	if vertices != nil {
		g.vertices = vertices
	}
	if inbound != nil {
		g.inboundEdges = inbound
	}
	if outbound != nil {
		g.outboundEdges = outbound
	}

	return g
}

func TestNewAcyclicGraph(t *testing.T) {
	dag := NewAcyclicGraph[string, int]()

	assert.Zero(t, dag.GetOrder())
	assert.Zero(t, dag.GetSize())
}

func TestAcyclicGraph_GetOrder(t *testing.T) {
	tests := []struct {
		name   string
		source *AcyclicGraph[string, int]
		want   int
	}{
		{
			name:   "Empty",
			source: buildTestGraph(t, map[string]int{}, nil, nil),
			want:   0,
		},
		{
			name: "Vertices",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			want: 2,
		},
		{
			name: "Edge",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			want: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.source.GetOrder())
		})
	}
}

func TestAcyclicGraph_GetSize(t *testing.T) {
	tests := []struct {
		name   string
		source *AcyclicGraph[string, int]
		want   int
	}{
		{
			name:   "Empty",
			source: buildTestGraph(t, map[string]int{}, nil, nil),
			want:   0,
		},
		{
			name: "Vertices - No Edge",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			want: 0,
		},
		{
			name: "One Edge",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.source.GetSize())
		})
	}
}

func TestAcyclicGraph_AddVertex(t *testing.T) {
	type args struct {
		id string
		v  int
	}
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		args    args
		want    *AcyclicGraph[string, int]
		wantErr error
	}{
		{
			name:   "Empty ID",
			source: NewAcyclicGraph[string, int](),
			args: args{
				id: "",
				v:  1,
			},
			want:    nil,
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Duplicate ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
			}, nil, nil),
			args: args{
				id: "foo",
				v:  2,
			},
			wantErr: ErrVertexDuplicate,
		},
		{
			name:   "Valid add",
			source: NewAcyclicGraph[string, int](),
			args: args{
				id: "foo",
				v:  1,
			},
			want: &AcyclicGraph[string, int]{
				vertices: map[string]int{
					"foo": 1,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.AddVertex(tt.args.id, tt.args.v)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want.vertices, tt.source.vertices)
			assert.Equal(t, tt.want.GetOrder(), tt.source.GetOrder())
			assert.Equal(t, tt.want.GetSize(), tt.source.GetSize())
		})
	}
}

func TestAcyclicGraph_GetVertex(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    any
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Ok",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			id:   "foo",
			want: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.source.GetVertex(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)
			assert.NoError(t, err)
		})
	}
}

func TestAcyclicGraph_DeleteVertex(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    *AcyclicGraph[string, int]
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Ok - no edges",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			id: "foo",
			want: buildTestGraph(t, map[string]int{
				"bar": 2,
			}, nil, nil),
		},
		{
			name: "Ok - inbound edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			id: "foo",
			want: buildTestGraph(t,
				map[string]int{
					"bar": 2,
				},
				map[string]Set[string]{},
				map[string]Set[string]{
					"bar": {},
				},
			),
		},
		{
			name: "Ok - outbound edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			id: "foo",
			want: buildTestGraph(t,
				map[string]int{
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {},
				},
				map[string]Set[string]{},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.DeleteVertex(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want.vertices, tt.source.vertices)
			assert.Equal(t, tt.want.inboundEdges, tt.source.inboundEdges)
			assert.Equal(t, tt.want.outboundEdges, tt.source.outboundEdges)
			assert.Equal(t, tt.want.GetOrder(), tt.source.GetOrder())
			assert.Equal(t, tt.want.GetSize(), tt.source.GetSize())
		})
	}
}

func TestAcyclicGraph_AddEdge(t *testing.T) {
	tests := []struct {
		name     string
		source   *AcyclicGraph[string, int]
		sourceID string
		targetID string
		want     *AcyclicGraph[string, int]
		wantErr  error
	}{
		{
			name: "Empty Source ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "",
			targetID: "bar",
			wantErr:  ErrEdgeSourceIDEmpty,
		},
		{
			name: "Empty Target ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "",
			wantErr:  ErrEdgeTargetIDEmpty,
		},
		{
			name: "Missing Source ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "baz",
			targetID: "bar",
			wantErr:  ErrEdgeSourceIDNotFound,
		},
		{
			name: "Missing Target ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "baz",
			wantErr:  ErrEdgeTargetIDNotFound,
		},
		{
			name: "Identical IDs",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "foo",
			wantErr:  ErrEdgeSourceTargetIdentical,
		},
		{
			name: "Duplicate",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			wantErr:  ErrEdgeDuplicate,
		},
		{
			name: "Loop",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			wantErr:  ErrEdgeLoop,
		},
		{
			name: "Ok",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			sourceID: "foo",
			targetID: "bar",
			want: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
		},
		{
			name: "Ok - Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"baz": "baz"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
		},
		{
			name: "Ok - Ancestors",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"baz": "baz"},
				},
				map[string]Set[string]{
					"baz": {"foo": "foo"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.AddEdge(tt.sourceID, tt.targetID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want.vertices, tt.source.vertices)
			assert.Equal(t, tt.want.inboundEdges, tt.source.inboundEdges)
			assert.Equal(t, tt.want.outboundEdges, tt.source.outboundEdges)
			assert.Equal(t, tt.want.GetOrder(), tt.source.GetOrder())
			assert.Equal(t, tt.want.GetSize(), tt.source.GetSize())
		})
	}
}

func TestAcyclicGraph_IsEdge(t *testing.T) {
	tests := []struct {
		name     string
		source   *AcyclicGraph[string, int]
		sourceID string
		targetID string
		want     bool
		wantErr  error
	}{
		{
			name: "Empty Source ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "",
			targetID: "bar",
			wantErr:  ErrEdgeSourceIDEmpty,
		},
		{
			name: "Empty Target ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "",
			wantErr:  ErrEdgeTargetIDEmpty,
		},
		{
			name: "Missing Source ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "baz",
			targetID: "bar",
			wantErr:  ErrEdgeSourceIDNotFound,
		},
		{
			name: "Missing Target ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "baz",
			wantErr:  ErrEdgeTargetIDNotFound,
		},
		{
			name: "Identical IDs",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "foo",
			wantErr:  ErrEdgeSourceTargetIdentical,
		},
		{
			name: "Edge",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want:     true,
		},
		{
			name: "No Edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{},
				map[string]Set[string]{},
			),
			sourceID: "bar",
			targetID: "foo",
			want:     false,
		},
		{
			name: "No Edge on source",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 2,
				},
				map[string]Set[string]{
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"baz": {"bar": "bar"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want:     false,
		},
		{
			name: "No Edge on target",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 2,
				},
				map[string]Set[string]{
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"baz": "baz"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want:     false,
		},
		{
			name: "Wrong Edge on source ",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 2,
				},
				map[string]Set[string]{
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"baz": "baz"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want:     false,
		},
		{
			name: "Wrong Edge on source and target",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 2,
				},
				map[string]Set[string]{
					"baz": {
						"foo": "foo",
						"bar": "bar",
					},
				},
				map[string]Set[string]{
					"foo": {"baz": "baz"},
					"bar": {"baz": "baz"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := tt.source.IsEdge(tt.sourceID, tt.targetID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, res)
		})
	}
}

func TestAcyclicGraph_DeleteEdge(t *testing.T) {
	tests := []struct {
		name     string
		source   *AcyclicGraph[string, int]
		sourceID string
		targetID string
		want     *AcyclicGraph[string, int]
		wantErr  error
	}{
		{
			name: "Empty Source ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "",
			targetID: "bar",
			wantErr:  ErrEdgeSourceIDEmpty,
		},
		{
			name: "Empty Target ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "",
			wantErr:  ErrEdgeTargetIDEmpty,
		},
		{
			name: "Missing Source ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "baz",
			targetID: "bar",
			wantErr:  ErrEdgeSourceIDNotFound,
		},
		{
			name: "Missing Target ID",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "baz",
			wantErr:  ErrEdgeTargetIDNotFound,
		},
		{
			name: "Identical IDs",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "foo",
			wantErr:  ErrEdgeSourceTargetIdentical,
		},
		{
			name: "Missing Edge",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			sourceID: "foo",
			targetID: "bar",
			wantErr:  ErrEdgeNotFound,
		},
		{
			name: "Ok",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"bar": {},
				},
				map[string]Set[string]{
					"foo": {},
				},
			),
		},
		{
			name: "Ok - Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {},
					"bar": {"baz": "baz"},
				},
			),
		},
		{
			name: "Ok - Ancestors",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
			),
			sourceID: "foo",
			targetID: "bar",
			want: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {},
					"foo": {"baz": "baz"},
				},
				map[string]Set[string]{
					"foo": {},
					"baz": {"foo": "foo"},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.source.DeleteEdge(tt.sourceID, tt.targetID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want.vertices, tt.source.vertices)
			assert.Equal(t, tt.want.inboundEdges, tt.source.inboundEdges)
			assert.Equal(t, tt.want.outboundEdges, tt.source.outboundEdges)
			assert.Equal(t, tt.want.GetOrder(), tt.source.GetOrder())
			assert.Equal(t, tt.want.GetSize(), tt.source.GetSize())
		})
	}
}

func TestAcyclicGraph_GetLeaves(t *testing.T) {
	tests := []struct {
		name   string
		source *AcyclicGraph[string, int]
		want   map[string]int
	}{
		{
			name: "No Vertices",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: map[string]int{},
		},
		{
			name: "No Edges (all leaves)",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			want: map[string]int{
				"foo": 1,
				"bar": 2,
			},
		},
		{
			name: "Two nodes - one edge",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			want: map[string]int{
				"foo": 1,
			},
		},
		{
			name: "Three nodes - one edge",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"baz": "baz"},
				},
			),
			want: map[string]int{
				"foo": 1,
				"baz": 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaves := tt.source.GetLeaves()

			assert.Equal(t, tt.want, leaves)
		})
	}
}

func TestAcyclicGraph_GetRoots(t *testing.T) {
	tests := []struct {
		name   string
		source *AcyclicGraph[string, int]
		want   map[string]int
	}{
		{
			name: "No Vertices",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: map[string]int{},
		},
		{
			name: "No Edges (all roots)",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			want: map[string]int{
				"foo": 1,
				"bar": 2,
			},
		},
		{
			name: "Two nodes - one edge",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "Three nodes - one edge",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"baz": "baz"},
				},
			),
			want: map[string]int{
				"foo": 1,
				"bar": 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			leaves := tt.source.GetRoots()

			assert.Equal(t, tt.want, leaves)
		})
	}
}

func TestAcyclicGraph_GetVertices(t *testing.T) {
	tests := []struct {
		name   string
		source *AcyclicGraph[string, int]
		want   map[string]int
	}{
		{
			name:   "Empty Graph",
			source: buildTestGraph(t, nil, nil, nil),
			want:   map[string]int{},
		},
		{
			name: "Nodes",
			source: buildTestGraph(t, map[string]int{
				"foo": 1,
				"bar": 2,
			}, nil, nil),
			want: map[string]int{
				"foo": 1,
				"bar": 2,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := tt.source.GetVertices()

			assert.Equal(t, tt.want, res)
			assert.NotSame(t, tt.source.vertices, res)
		})
	}

}

func TestAcyclicGraph_GetParents(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    map[string]int
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Empty graph",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: map[string]int{},
		},
		{
			name: "No Edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				nil,
				nil,
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Ancestors, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "One Ancestor, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "No Ancestors, Two Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "One Ancestor, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "No Ancestors, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Parents, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
		{
			name: "No Ancestors, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Parents, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"qux": {"foo": "foo"},
					"quz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"foo": {
						"qux": "qux",
						"quz": "quz",
					},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.source.GetParents(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				assert.Equal(t, tt.want, items)
			}
		})
	}
}

func TestAcyclicGraph_GetChildren(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    map[string]int
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Empty graph",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: map[string]int{},
		},
		{
			name: "No Edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				nil,
				nil,
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Ancestors, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "One Ancestor, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
			),
			id: "foo",
			want: map[string]int{
				"baz": 3,
			},
		},
		{
			name: "No Ancestors, Two Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "One Ancestor, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "No Ancestors, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "Two Parents, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "No Ancestors, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
		{
			name: "Two Parents, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"qux": {"foo": "foo"},
					"quz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"foo": {
						"qux": "qux",
						"quz": "quz",
					},
				},
			),
			id: "foo",
			want: map[string]int{
				"qux": 4,
				"quz": 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.source.GetChildren(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				assert.Equal(t, tt.want, items)
			}
		})
	}
}

func TestAcyclicGraph_GetAncestors(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    map[string]int
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Empty graph",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: map[string]int{},
		},
		{
			name: "No Edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				nil,
				nil,
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Ancestors, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
		{
			name: "One Ancestor, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "No Ancestors, Two Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "One Ancestor, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "No Ancestors, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Parents, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
		{
			name: "No Ancestors, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Parents, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"qux": {"foo": "foo"},
					"quz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"foo": {
						"qux": "qux",
						"quz": "quz",
					},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.source.GetAncestors(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				assert.Equal(t, tt.want, items)
			}
		})
	}
}

func TestAcyclicGraph_GetOrderedAncestors(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    [][]string
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Empty graph",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: [][]string{},
		},
		{
			name: "No Edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				nil,
				nil,
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "Two Ancestors, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar", "baz"},
			},
		},
		{
			name: "Three Ancestors on two levels, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"qux": "qux",
					},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
					"qux": {"foo": "foo"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar", "qux"},
				{"baz"},
			},
		},
		{
			name: "No Ancestors, Three Descendants on two levels",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"qux": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"bar": {"qux": "qux"},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "One Ancestor, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar"},
			},
		},
		{
			name: "No Ancestors, Two Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "One Ancestor, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar"},
			},
		},
		{
			name: "No Ancestors, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "Two Parents, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar", "baz"},
			},
		},
		{
			name: "No Ancestors, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "Two Parents, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"qux": {"foo": "foo"},
					"quz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"foo": {
						"qux": "qux",
						"quz": "quz",
					},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar", "baz"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.source.GetOrderedAncestors(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				for _, w := range tt.want {
					tst := items[:len(w)]
					assert.ElementsMatch(t, w, tst)
					items = items[len(w):]
				}
			}
		})
	}
}

func TestAcyclicGraph_GetDescendants(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    map[string]int
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Empty graph",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: map[string]int{},
		},
		{
			name: "No Edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				nil,
				nil,
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "Two Ancestors, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "One Ancestor, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
			),
			id: "foo",
			want: map[string]int{
				"baz": 3,
			},
		},
		{
			name: "No Ancestors, Two Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
		{
			name: "One Ancestor, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "No Ancestors, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
			},
		},
		{
			name: "Two Parents, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
			),
			id:   "foo",
			want: map[string]int{},
		},
		{
			name: "No Ancestors, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
			),
			id: "foo",
			want: map[string]int{
				"bar": 2,
				"baz": 3,
			},
		},
		{
			name: "Two Parents, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"qux": {"foo": "foo"},
					"quz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"foo": {
						"qux": "qux",
						"quz": "quz",
					},
				},
			),
			id: "foo",
			want: map[string]int{
				"qux": 4,
				"quz": 5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.source.GetDescendants(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				assert.Equal(t, tt.want, items)
			}
		})
	}
}

func TestAcyclicGraph_GetOrderedDescendants(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    [][]string
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "Empty graph",
			source: buildTestGraph(t,
				map[string]int{},
				nil,
				nil,
			),
			want: [][]string{},
		},
		{
			name: "No Edges",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				nil,
				nil,
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "Two Ancestors, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "Three Ancestors on two levels, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"qux": "qux",
					},
					"bar": {"baz": "baz"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
					"qux": {"foo": "foo"},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "No Ancestors, Three Descendants on two levels",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"qux": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"bar": {"qux": "qux"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar", "baz", "qux"},
			},
		},
		{
			name: "One Ancestor, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"foo": {"baz": "baz"},
				},
			),
			id: "foo",
			want: [][]string{
				{"baz"},
			},
		},
		{
			name: "No Ancestors, Two Descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"bar": "bar"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
					"bar": {"baz": "baz"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar", "baz"},
			},
		},
		{
			name: "One Ancestor, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "No Ancestors, One descendant",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {"bar": "bar"},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar"},
			},
		},
		{
			name: "Two Parents, No descendants",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
			),
			id:   "foo",
			want: [][]string{},
		},
		{
			name: "No Ancestors, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
				},
			),
			id: "foo",
			want: [][]string{
				{"bar", "baz"},
			},
		},
		{
			name: "Two Parents, Two Children",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
					},
					"qux": {"foo": "foo"},
					"quz": {"foo": "foo"},
				},
				map[string]Set[string]{
					"bar": {"foo": "foo"},
					"baz": {"foo": "foo"},
					"foo": {
						"qux": "qux",
						"quz": "quz",
					},
				},
			),
			id: "foo",
			want: [][]string{
				{"qux", "quz"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := tt.source.GetOrderedDescendants(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				for _, w := range tt.want {
					tst := items[:len(w)]
					assert.ElementsMatch(t, w, tst)
					items = items[len(w):]
				}
			}
		})
	}
}

func TestAcyclicGraph_AncestorsWalker(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    [][]string
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "",
			source: func() *AcyclicGraph[string, int] {
				g := NewAcyclicGraph[string, int]()
				_ = g.AddVertex("01", 1)
				_ = g.AddVertex("02", 1)
				_ = g.AddVertex("03", 1)
				_ = g.AddVertex("04", 1)
				_ = g.AddVertex("05", 1)
				_ = g.AddVertex("06", 1)
				_ = g.AddVertex("07", 1)
				_ = g.AddVertex("08", 1)
				_ = g.AddVertex("09", 1)
				_ = g.AddVertex("10", 1)

				_ = g.AddEdge("01", "02")
				_ = g.AddEdge("01", "03")
				_ = g.AddEdge("02", "04")
				_ = g.AddEdge("02", "05")
				_ = g.AddEdge("04", "06")
				_ = g.AddEdge("05", "06")
				_ = g.AddEdge("06", "07")
				_ = g.AddEdge("07", "08")
				_ = g.AddEdge("07", "09")
				_ = g.AddEdge("08", "10")
				_ = g.AddEdge("09", "10")

				return g
			}(),
			id: "10",
			want: [][]string{
				{"08", "09"},
				{"07"},
				{"06"},
				{"04", "05"},
				{"02"},
				{"01"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, _, err := tt.source.AncestorsWalker(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			var items []string
			for id := range ids {
				items = append(items, id)
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				for _, w := range tt.want {
					tst := items[:len(w)]
					assert.ElementsMatch(t, w, tst)
					items = items[len(w):]
				}
			}
		})
	}
}

func TestAcyclicGraph_AncestorsWalker_Signal(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		stopID  string
		want    []string
		wantErr error
	}{
		{
			name: "",
			source: func() *AcyclicGraph[string, int] {
				g := NewAcyclicGraph[string, int]()
				_ = g.AddVertex("foo", 1)
				_ = g.AddVertex("bar", 2)
				_ = g.AddVertex("baz", 3)
				_ = g.AddVertex("qux", 4)
				_ = g.AddVertex("quz", 5)

				_ = g.AddEdge("foo", "bar")
				_ = g.AddEdge("bar", "baz")
				_ = g.AddEdge("bar", "qux")
				_ = g.AddEdge("qux", "quz")

				return g
			}(),
			id:     "quz",
			stopID: "bar",
			want: []string{
				"qux",
				"bar",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, signal, err := tt.source.AncestorsWalker(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			var items []string
			for id := range ids {
				items = append(items, id)
				if id == tt.stopID {
					signal <- true
					break
				}
			}

			assert.Equal(t, tt.want, items)
		})
	}
}

func TestAcyclicGraph_DescendantsWalker(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		id      string
		want    [][]string
		wantErr error
	}{
		{
			name: "Empty ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "",
			wantErr: ErrVertexIDEmpty,
		},
		{
			name: "Missing ID",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
				},
				nil,
				nil,
			),
			id:      "baz",
			wantErr: ErrVertexNotFound,
		},
		{
			name: "",
			source: func() *AcyclicGraph[string, int] {
				g := NewAcyclicGraph[string, int]()
				_ = g.AddVertex("01", 1)
				_ = g.AddVertex("02", 1)
				_ = g.AddVertex("03", 1)
				_ = g.AddVertex("04", 1)
				_ = g.AddVertex("05", 1)
				_ = g.AddVertex("06", 1)
				_ = g.AddVertex("07", 1)
				_ = g.AddVertex("08", 1)
				_ = g.AddVertex("09", 1)
				_ = g.AddVertex("10", 1)

				_ = g.AddEdge("01", "02")
				_ = g.AddEdge("01", "03")
				_ = g.AddEdge("02", "04")
				_ = g.AddEdge("02", "05")
				_ = g.AddEdge("04", "06")
				_ = g.AddEdge("05", "06")
				_ = g.AddEdge("06", "07")
				_ = g.AddEdge("07", "08")
				_ = g.AddEdge("07", "09")
				_ = g.AddEdge("08", "10")
				_ = g.AddEdge("09", "10")

				return g
			}(),
			id: "01",
			want: [][]string{
				{"02", "03"},
				{"04", "05"},
				{"06"},
				{"07"},
				{"08", "09"},
				{"10"},
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, _, err := tt.source.DescendantsWalker(tt.id)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			var items []string
			for id := range ids {
				items = append(items, id)
			}

			if len(tt.want) == 0 {
				assert.Empty(t, items)
			} else {
				for _, w := range tt.want {
					tst := items[:len(w)]
					assert.ElementsMatch(t, w, tst)
					items = items[len(w):]
				}
			}
		})
	}
}

func TestAcyclicGraph_DescendantsWalker_Signal(t *testing.T) {
	tests := []struct {
		name    string
		source  *AcyclicGraph[string, int]
		startID string
		stopID  string
		want    []string
		wantErr error
	}{
		{
			name: "",
			source: func() *AcyclicGraph[string, int] {
				g := NewAcyclicGraph[string, int]()
				_ = g.AddVertex("foo", 1)
				_ = g.AddVertex("bar", 2)
				_ = g.AddVertex("baz", 3)
				_ = g.AddVertex("qux", 4)
				_ = g.AddVertex("quz", 5)

				_ = g.AddEdge("foo", "bar")
				_ = g.AddEdge("bar", "baz")
				_ = g.AddEdge("bar", "qux")
				_ = g.AddEdge("qux", "quz")

				return g
			}(),
			startID: "foo",
			stopID:  "qux",
			want: []string{
				"bar",
				"baz",
				"qux",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ids, signal, err := tt.source.DescendantsWalker(tt.startID)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)

			var items []string
			for id := range ids {
				items = append(items, id)
				if id == tt.stopID {
					signal <- true
					break
				}
			}

			assert.Equal(t, tt.want, items)
		})
	}
}

func TestAcyclicGraph_ReduceTransitively(t *testing.T) {
	tests := []struct {
		name   string
		source *AcyclicGraph[string, int]
		want   *AcyclicGraph[string, int]
	}{
		{
			name: "Redundant edge between foo and cor",
			source: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
					"cor": 6,
				},
				map[string]Set[string]{
					"foo": {},
					"bar": {
						"foo": "foo",
					},
					"baz": {
						"foo": "foo",
					},
					"qux": {
						"foo": "foo",
					},
					"quz": {
						"foo": "foo",
					},
					"cor": {
						"foo": "foo",
						"bar": "bar",
						"baz": "baz",
						"qux": "qux",
						"quz": "quz",
					},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
						"qux": "qux",
						"quz": "quz",
						"cor": "cor",
					},
					"bar": {
						"cor": "cor",
					},
					"baz": {
						"cor": "cor",
					},
					"qux": {
						"cor": "cor",
					},
					"quz": {
						"cor": "cor",
					},
				},
			),
			want: buildTestGraph(t,
				map[string]int{
					"foo": 1,
					"bar": 2,
					"baz": 3,
					"qux": 4,
					"quz": 5,
					"cor": 6,
				},
				map[string]Set[string]{
					"foo": {},
					"bar": {
						"foo": "foo",
					},
					"baz": {
						"foo": "foo",
					},
					"qux": {
						"foo": "foo",
					},
					"quz": {
						"foo": "foo",
					},
					"cor": {
						"bar": "bar",
						"baz": "baz",
						"qux": "qux",
						"quz": "quz",
					},
				},
				map[string]Set[string]{
					"foo": {
						"bar": "bar",
						"baz": "baz",
						"qux": "qux",
						"quz": "quz",
					},
					"bar": {
						"cor": "cor",
					},
					"baz": {
						"cor": "cor",
					},
					"qux": {
						"cor": "cor",
					},
					"quz": {
						"cor": "cor",
					},
				},
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.source.ReduceTransitively()

			assert.Equal(t, tt.want, tt.source)
		})
	}
}

func TestAcyclicGraph_String(t *testing.T) {
	g := NewAcyclicGraph[string, int]()

	_ = g.AddVertex("1", 1)
	_ = g.AddVertex("2", 2)
	_ = g.AddVertex("3", 3)
	_ = g.AddVertex("4", 4)

	_ = g.AddEdge("1", "2")
	_ = g.AddEdge("2", "3")
	_ = g.AddEdge("2", "4")

	expected := []string{
		"DAG Vertices: 4 - Edges: 3",
		`
DAG Vertices: 4 - Edges: 3
Vertices:
  1
  2
  3
  4
Edges:
  1 -> 2
  2 -> 3
  2 -> 4
`[1:],
	}

	res := g.String()
	assert.Equal(t, expected[0], res[:len(expected[0])])
	assert.Len(t, res, len(expected[1]))
}

func BenchmarkAcyclicGraph_AddVertices(b *testing.B) {
	dag := NewAcyclicGraph[int, int]()

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		_ = dag.AddVertex(n, n)
	}
}

func BenchmarkAcyclicGraph_AddEdges(b *testing.B) {
	dag := NewAcyclicGraph[int, int]()

	_ = dag.AddVertex(0, 0)
	for n := 1; n <= b.N; n++ {
		_ = dag.AddVertex(n, n)
	}

	b.ResetTimer()
	for n := 1; n <= b.N; n++ {
		_ = dag.AddEdge(0, n)
	}
}
