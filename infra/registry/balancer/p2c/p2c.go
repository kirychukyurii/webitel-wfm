package p2c

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/webitel/webitel-wfm/infra/registry"
	"github.com/webitel/webitel-wfm/infra/registry/node/ewma"
)

const (
	forcePick = time.Second * 3

	// Name is p2c(Pick of 2 choices) balancer name
	Name = "p2c"
)

var _ registry.Balancer = (*Balancer)(nil)

// Option is p2c builder option.
type Option func(o *options)

// options is p2c builder options
type options struct{}

// New creates a p2c selector.
func New(opts ...Option) registry.Selector {
	return NewBuilder(opts...).Build()
}

// Balancer is p2c selector.
type Balancer struct {
	mu     sync.Mutex
	r      *rand.Rand
	picked int64
}

// choose two distinct nodes.
func (s *Balancer) prePick(nodes []registry.WeightedNode) (nodeA registry.WeightedNode, nodeB registry.WeightedNode) {
	s.mu.Lock()
	a := s.r.Intn(len(nodes))
	b := s.r.Intn(len(nodes) - 1)
	s.mu.Unlock()
	if b >= a {
		b = b + 1
	}

	nodeA, nodeB = nodes[a], nodes[b]

	return
}

// Pick pick a node.
func (s *Balancer) Pick(_ context.Context, nodes []registry.WeightedNode) (registry.WeightedNode, registry.DoneFunc, error) {
	if len(nodes) == 0 {
		return nil, nil, registry.ErrNoAvailable
	}

	if len(nodes) == 1 {
		done := nodes[0].Pick()
		return nodes[0], done, nil
	}

	var pc, upc registry.WeightedNode
	nodeA, nodeB := s.prePick(nodes)

	// meta.Weight is the weight set by the service publisher in discovery
	if nodeB.Weight() > nodeA.Weight() {
		pc, upc = nodeB, nodeA
	} else {
		pc, upc = nodeA, nodeB
	}

	// If the failed node has never been selected once during forceGap, it is forced to be selected once
	// Take advantage of forced opportunities to trigger updates of success rate and delay
	if upc.PickElapsed() > forcePick && atomic.CompareAndSwapInt64(&s.picked, 0, 1) {
		pc = upc
		atomic.StoreInt64(&s.picked, 0)
	}

	done := pc.Pick()

	return pc, done, nil
}

// NewBuilder returns a selector builder with p2c balancer
func NewBuilder(opts ...Option) registry.Builder {
	var option options
	for _, opt := range opts {
		opt(&option)
	}

	return &registry.DefaultBuilder{
		Balancer: &Builder{},
		Node:     &ewma.Builder{},
	}
}

// Builder is p2c builder
type Builder struct{}

// Build creates Balancer
func (b *Builder) Build() registry.Balancer {
	return &Balancer{r: rand.New(rand.NewSource(time.Now().UnixNano()))}
}
