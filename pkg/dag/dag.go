package dag

import (
	"fmt"
	"reflect"
	"strings"
)

// AcyclicGraph implements the data structure of the DAG.
type AcyclicGraph[K comparable, T any] struct {
	vertices         map[K]T
	inboundEdges     map[K]Set[K]
	outboundEdges    map[K]Set[K]
	ancestorsCache   map[K]Set[K]
	descendantsCache map[K]Set[K]
}

// NewAcyclicGraph creates / initializes a new AcyclicGraph.
func NewAcyclicGraph[K comparable, T any]() *AcyclicGraph[K, T] {
	return &AcyclicGraph[K, T]{
		vertices:         make(map[K]T),
		inboundEdges:     make(map[K]Set[K]),
		outboundEdges:    make(map[K]Set[K]),
		ancestorsCache:   make(map[K]Set[K]),
		descendantsCache: make(map[K]Set[K]),
	}
}

// AddVertex adds the vertex v to the AcyclicGraph.
//
// Returns:
// 	- ErrVertexIDEmpty if id is empty.
// 	- ErrVertexDuplicate if id is already part of the graph.
func (g *AcyclicGraph[K, T]) AddVertex(id K, v T) error {
	if reflect.ValueOf(id).IsZero() {
		return ErrVertexIDEmpty
	}

	if _, exists := g.vertices[id]; exists {
		return ErrVertexDuplicate
	}

	g.vertices[id] = v

	return nil
}

// GetVertex returns a vertex by its id.
//
// Returns:
// 	- ErrVertexIDEmpty if id is empty.
// 	- ErrVertexNotFound if id is unknown.
func (g *AcyclicGraph[K, T]) GetVertex(id K) (any, error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, err
	}

	return g.vertices[id], nil
}

// DeleteVertex deletes the vertex with the given id.
// DeleteVertex also deletes all attached edges (inbound and outbound).
// DeleteVertex returns an error if id is empty or unknown.
func (g *AcyclicGraph[K, T]) DeleteVertex(id K) error {
	if err := g.checkVertexID(id); err != nil {
		return err
	}

	// Get descendents and ancestors as they are now.
	descendants, _ := g.getDescendantsIDs(id)
	ancestors, _ := g.getAncestorsIDs(id)

	// Delete id in outbound edges of parents.
	if _, exists := g.inboundEdges[id]; exists {
		for parent := range g.inboundEdges[id] {
			delete(g.outboundEdges[parent], id)
		}
	}

	// Delete id in inbound edges of children.
	if _, exists := g.outboundEdges[id]; exists {
		for child := range g.outboundEdges[id] {
			delete(g.inboundEdges[child], id)
		}
	}

	// Delete in- and outbound of id itself.
	delete(g.inboundEdges, id)
	delete(g.outboundEdges, id)

	// For id and all its descendants, delete cached ancestors.
	for _, descendant := range descendants {
		delete(g.ancestorsCache, descendant)
	}
	delete(g.ancestorsCache, id)

	// For id and all its ancestors, delete cached descendants.
	for _, ancestor := range ancestors {
		delete(g.descendantsCache, ancestor)
	}
	delete(g.descendantsCache, id)

	// Delete id itself.
	delete(g.vertices, id)

	return nil
}

// AddEdge adds an edge between sourceID and targetID.
// AddEdge returns an error if sourceID or targetID are empty or unknown,
// if the edge already exists, or if the new edge would create a loop.
func (g *AcyclicGraph[K, T]) AddEdge(sourceID, targetID K) error {
	if err := g.checkEdgeIDs(sourceID, targetID); err != nil {
		return err
	}

	if isEdge, _ := g.IsEdge(sourceID, targetID); isEdge {
		return ErrEdgeDuplicate
	}

	// Get descendents and ancestors as they are now.
	descendants, _ := g.getDescendantsIDs(targetID)
	ancestors, _ := g.getAncestorsIDs(sourceID)

	// Check if we're creating a loop.
	if descendants.Includes(sourceID) {
		return ErrEdgeLoop
	}

	// Target ID is a child of source ID.
	if _, exists := g.outboundEdges[sourceID]; !exists {
		g.outboundEdges[sourceID] = make(Set[K])
	}
	g.outboundEdges[sourceID].Add(targetID)

	// Source ID is a parent of target ID.
	if _, exists := g.inboundEdges[targetID]; !exists {
		g.inboundEdges[targetID] = make(Set[K])
	}
	g.inboundEdges[targetID].Add(sourceID)

	// For target and all its descendants, delete cached ancestors.
	for descendant := range descendants {
		delete(g.ancestorsCache, descendant)
	}
	delete(g.ancestorsCache, targetID)

	// For source and all its ancestors, delete cached descendants.
	for ancestor := range ancestors {
		delete(g.descendantsCache, ancestor)
	}
	delete(g.descendantsCache, sourceID)

	return nil
}

// IsEdge returns true if there exists an edge between sourceID and targetID.
// IsEdge returns false if there is no such edge.
// IsEdge returns an error if sourceID or targetID are empty, unknown, or the same.
func (g *AcyclicGraph[K, T]) IsEdge(sourceID, targetID K) (bool, error) {
	if err := g.checkEdgeIDs(sourceID, targetID); err != nil {
		return false, err
	}

	if _, exists := g.outboundEdges[sourceID]; !exists {
		return false, nil
	}

	if !g.outboundEdges[sourceID].Includes(targetID) {
		return false, nil
	}

	return true, nil
}

// DeleteEdge deletes the edge between sourceID and targetID.
// DeleteEdge returns an error if sourceID or targetID are empty or unknown,
// or if there is no edge between sourceID and targetID.
func (g *AcyclicGraph[K, T]) DeleteEdge(sourceID, targetID K) error {
	if err := g.checkEdgeIDs(sourceID, targetID); err != nil {
		return err
	}

	if isEdge, _ := g.IsEdge(sourceID, targetID); !isEdge {
		return ErrEdgeNotFound
	}

	// Get descendents and ancestors as they are now.
	descendants, _ := g.getDescendantsIDs(sourceID)
	ancestors, _ := g.getAncestorsIDs(targetID)

	// Delete outbound and inbound.
	g.outboundEdges[sourceID].Delete(targetID)
	g.inboundEdges[targetID].Delete(sourceID)

	// For sourceID and all its descendants, delete cached ancestors.
	for descendant := range descendants {
		delete(g.ancestorsCache, descendant)
	}
	delete(g.ancestorsCache, sourceID)

	// For targetID and all its ancestors, delete cached descendants.
	for ancestor := range ancestors {
		delete(g.descendantsCache, ancestor)
	}
	delete(g.descendantsCache, targetID)

	return nil
}

// GetOrder returns the number of vertices in the graph.
func (g *AcyclicGraph[K, T]) GetOrder() int {
	return len(g.vertices)
}

// GetSize returns the number of edges in the graph.
func (g *AcyclicGraph[K, T]) GetSize() int {
	size := 0
	for _, v := range g.outboundEdges {
		size += len(v)
	}

	return size
}

// GetLeaves returns all vertices without children.
func (g *AcyclicGraph[K, T]) GetLeaves() map[K]T {
	leaves := make(map[K]T)

	for k, v := range g.vertices {
		if targetIDs, ok := g.outboundEdges[k]; !ok || len(targetIDs) == 0 {
			leaves[k] = v
		}
	}

	return leaves
}

// GetRoots returns all vertices without parents.
func (g *AcyclicGraph[K, T]) GetRoots() map[K]T {
	roots := make(map[K]T)

	for k, v := range g.vertices {
		if sourceIDs, ok := g.inboundEdges[k]; !ok || len(sourceIDs) == 0 {
			roots[k] = v
		}
	}

	return roots
}

// GetVertices returns all vertices.
func (g *AcyclicGraph[K, T]) GetVertices() map[K]T {
	return copyMap(g.vertices)
}

// GetParents returns all the immediate parents (inbound vertices) of
// the vertex with the specified id.
// GetParents returns an error if id is empty or unknown.
func (g *AcyclicGraph[K, T]) GetParents(id K) (map[K]T, error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, err
	}

	parents := make(map[K]T)

	for pid := range g.inboundEdges[id] {
		parents[pid] = g.vertices[pid]
	}

	return parents, nil
}

// GetChildren returns all the immediate children (outbound vertices) of
// the vertex with the specified id.
// GetChildren returns an error if id is empty or unknown.
func (g *AcyclicGraph[K, T]) GetChildren(id K) (map[K]T, error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, err
	}

	children := make(map[K]T)

	for cid := range g.outboundEdges[id] {
		children[cid] = g.vertices[cid]
	}

	return children, nil
}

// GetAncestors return all ancestors of the vertex with the specified id.
// GetAncestors returns an error if id is empty or unknown.
//
// Note: in order to get the ancestors, GetAncestors populates
// the ancestor-cache as needed.
// Depending on order and size of the sub-graph of the vertex with
// the specified id, this may take a long time and consume a lot of memory.
func (g *AcyclicGraph[K, T]) GetAncestors(id K) (map[K]T, error) {
	ids, err := g.getAncestorsIDs(id)
	if err != nil {
		return nil, err
	}

	return g.loadFromCache(ids), nil
}

func (g *AcyclicGraph[K, T]) getAncestorsIDs(id K) (Set[K], error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, err
	}

	cache, exists := g.ancestorsCache[id]
	if exists {
		return cache, nil
	}

	cache = make(Set[K])

	if parents, ok := g.inboundEdges[id]; ok {
		for parent := range parents {
			parentAncestors, _ := g.getAncestorsIDs(parent)
			for _, ancestor := range parentAncestors {
				cache.Add(ancestor)
			}
			cache.Add(parent)
		}
	}

	g.ancestorsCache[id] = cache

	return cache, nil
}

// GetOrderedAncestors returns all ancestors of the vertex with the specified id
// in a breath-first order.
// Only the first occurrence of each vertex is returned.
// GetOrderedAncestors returns an error if id is empty or unknown.
//
// Note: there is no order between sibling vertices.
// Two consecutive runs of GetOrderedAncestors may return different results.
func (g *AcyclicGraph[K, T]) GetOrderedAncestors(id K) ([]K, error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, err
	}

	ids, _, _ := g.AncestorsWalker(id)

	var ancestors []K
	for aid := range ids {
		ancestors = append(ancestors, aid)
	}

	return ancestors, nil
}

// GetDescendants return all descendants of the vertex with the specified id.
// GetDescendants returns an error if id is empty or unknown.
//
// Note: in order to get the descendants, GetDescendants populates the
// descendants-cache as needed.
// Depending on order and size of the sub-graph of the vertex
// with the specified id this may take a long time and consume a lot of memory.
func (g *AcyclicGraph[K, T]) GetDescendants(id K) (map[K]T, error) {
	ids, err := g.getDescendantsIDs(id)
	if err != nil {
		return nil, err
	}

	return g.loadFromCache(ids), nil
}

func (g *AcyclicGraph[K, T]) getDescendantsIDs(id K) (Set[K], error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, err
	}

	cache, exists := g.descendantsCache[id]
	if exists {
		return cache, nil
	}

	cache = make(Set[K])

	if children, ok := g.outboundEdges[id]; ok {
		for child := range children {
			childDescendants, _ := g.getDescendantsIDs(child)

			for _, descendant := range childDescendants {
				cache.Add(descendant)
			}
			cache.Add(child)
		}
	}

	g.descendantsCache[id] = cache

	return cache, nil
}

// GetOrderedDescendants returns all descendants of the vertex with
// the specified id in a breath-first order.
// Only the first occurrence of each vertex is returned.
// GetOrderedDescendants returns an error if id is empty or unknown.
//
// Note: there is no order between sibling vertices.
// Two consecutive runs of GetOrderedDescendants may return different results.
func (g *AcyclicGraph[K, T]) GetOrderedDescendants(id K) ([]K, error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, err
	}

	ids, _, _ := g.DescendantsWalker(id)

	var descendants []K
	for did := range ids {
		descendants = append(descendants, did)
	}

	return descendants, nil
}

// AncestorsWalker returns a channel and subsequently returns / walks all
// ancestors of the vertex with the specified id in a breath first order.
// The second channel returned may be used to stop further walking.
// AncestorsWalker returns an error, if id is empty or unknown.
//
// Note: there is no order between sibling vertices.
// Two consecutive runs of AncestorsWalker may return different results.
func (g *AcyclicGraph[K, T]) AncestorsWalker(id K) (chan K, chan bool, error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, nil, err
	}

	ids := make(chan K)
	signal := make(chan bool, 1)

	go func() {
		g.walkAncestors(id, ids, signal)
		close(ids)
		close(signal)
	}()

	return ids, signal, nil
}

func (g *AcyclicGraph[K, T]) walkAncestors(id K, ids chan K, signal chan bool) {
	var (
		fifo    []K
		visited = make(Set[K])
	)

	for parent := range g.inboundEdges[id] {
		visited.Add(parent)
		fifo = append(fifo, parent)
	}

	for {
		if len(fifo) == 0 {
			return
		}

		top := fifo[0]
		fifo = fifo[1:]

		for parent := range g.inboundEdges[top] {
			if !visited.Includes(parent) {
				visited.Add(parent)
				fifo = append(fifo, parent)
			}
		}

		select {
		case <-signal:
			return
		default:
			ids <- top
		}
	}
}

// DescendantsWalker returns a channel and subsequently returns / walks all
// descendants of the vertex with the specified id in a breath first order.
// The second channel returned may be used to stop further walking.
// DescendantsWalker returns an error, if id is empty or unknown.
//
// Note: there is no order between sibling vertices.
// Two consecutive runs of DescendantsWalker may return different results.
func (g *AcyclicGraph[K, T]) DescendantsWalker(id K) (chan K, chan bool, error) {
	if err := g.checkVertexID(id); err != nil {
		return nil, nil, err
	}

	ids := make(chan K)
	signal := make(chan bool, 1)

	go func() {
		g.walkDescendants(id, ids, signal)
		close(ids)
		close(signal)
	}()

	return ids, signal, nil
}

func (g *AcyclicGraph[K, T]) walkDescendants(id K, ids chan K, signal chan bool) {
	var (
		fifo    []K
		visited = make(Set[K])
	)

	for child := range g.outboundEdges[id] {
		visited.Add(child)
		fifo = append(fifo, child)
	}

	for {
		if len(fifo) == 0 {
			return
		}

		top := fifo[0]
		fifo = fifo[1:]

		for child := range g.outboundEdges[top] {
			if !visited.Includes(child) {
				visited.Add(child)
				fifo = append(fifo, child)
			}
		}

		select {
		case <-signal:
			return
		default:
			ids <- top
		}
	}
}

// ReduceTransitively transitively reduce the graph.
//
// Note: in order to do the reduction, the descendant-cache of all vertices is
// populated (i.e. the transitive closure).
// Depending on order and size of the AcyclicGraph, this may take a long time
// and consume a lot of memory.
func (g *AcyclicGraph[K, T]) ReduceTransitively() {
	graphChanged := false

	// Populate the descendents cache for all roots (i.e. the whole graph).
	for root := range g.GetRoots() {
		_, _ = g.getDescendantsIDs(root)
	}

	for id := range g.vertices {
		childrenDescendants := make(Set[K])

		for childID := range g.outboundEdges[id] {
			// Collect child descendants.
			for descendant := range g.descendantsCache[childID] {
				childrenDescendants.Add(descendant)
			}
		}

		for childID := range g.outboundEdges[id] {
			// Remove the edge between v and child,
			// only if child is a descendant of the children of v.
			if childrenDescendants.Includes(childID) {
				g.outboundEdges[id].Delete(childID)
				g.inboundEdges[childID].Delete(id)

				graphChanged = true
			}
		}
	}

	// Flush the descendants- and ancestor cache if the graph has changed.
	if graphChanged {
		g.flushCaches()
	}
}

func (g *AcyclicGraph[K, T]) String() string {
	var sb strings.Builder

	_, _ = fmt.Fprintf(&sb, "DAG Vertices: %d - Edges: %d\n", g.GetOrder(), g.GetSize())

	sb.WriteString("Vertices:\n")
	for k := range g.vertices {
		_, _ = fmt.Fprintf(&sb, "  %v\n", k)
	}

	sb.WriteString("Edges:\n")
	for v, children := range g.outboundEdges {
		for child := range children {
			_, _ = fmt.Fprintf(&sb, "  %v -> %v\n", v, child)
		}
	}

	return sb.String()
}

// checkVertexID checks that the specified vertex ID is valid for this graph.
func (g *AcyclicGraph[K, T]) checkVertexID(id K) error {
	if reflect.ValueOf(id).IsZero() {
		return ErrVertexIDEmpty
	}

	if _, exists := g.vertices[id]; !exists {
		return ErrVertexNotFound
	}

	return nil
}

// checkVertexID checks that the specified vertex IDs is a valid edge
// in this graph.
func (g *AcyclicGraph[K, T]) checkEdgeIDs(sourceID, targetID K) error {
	if reflect.ValueOf(sourceID).IsZero() {
		return ErrEdgeSourceIDEmpty
	}
	if reflect.ValueOf(targetID).IsZero() {
		return ErrEdgeTargetIDEmpty
	}
	if sourceID == targetID {
		return ErrEdgeSourceTargetIdentical
	}
	if _, exists := g.vertices[sourceID]; !exists {
		return ErrEdgeSourceIDNotFound
	}
	if _, exists := g.vertices[targetID]; !exists {
		return ErrEdgeTargetIDNotFound
	}

	return nil
}

func (g *AcyclicGraph[K, T]) flushCaches() {
	g.ancestorsCache = make(map[K]Set[K])
	g.descendantsCache = make(map[K]Set[K])
}

func (g *AcyclicGraph[K, T]) loadFromCache(cache Set[K]) map[K]T {
	res := make(map[K]T, len(cache))
	for k := range cache {
		res[k] = g.vertices[k]
	}

	return res
}

func copyMap[K comparable, T any](in map[K]T) map[K]T {
	out := make(map[K]T)
	for id, value := range in {
		out[id] = value
	}

	return out
}
