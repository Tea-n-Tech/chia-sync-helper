package chia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	someNode = fullNode{
		Id:     "123",
		Height: 0,
		Synced: false,
	}
)

func getFullNodeSlice(nNodes int) []fullNode {
	nodes := make([]fullNode, 0, nNodes)
	for iNode := 0; iNode < nNodes; iNode++ {
		nodes = append(nodes, someNode)
	}
	return nodes
}

func TestDecideWhichNodesToRemove(t *testing.T) {

	testCases := []struct {
		name                  string
		nNodesTotal           int64
		nNodesBehind          int
		expectedNodesToRemove int
	}{
		{
			name:                  "No nodes at all",
			nNodesTotal:           0,
			nNodesBehind:          0,
			expectedNodesToRemove: 0,
		},
		{
			name:                  "One node not behind",
			nNodesTotal:           1,
			nNodesBehind:          0,
			expectedNodesToRemove: 0,
		},
		{
			name:                  "One node which is behind",
			nNodesTotal:           1,
			nNodesBehind:          1,
			expectedNodesToRemove: 1,
		},
		{
			name:                  "Too many nodes behind uneven case",
			nNodesTotal:           7,
			nNodesBehind:          5,
			expectedNodesToRemove: 2,
		},
		{
			name:                  "Too many nodes behind even case",
			nNodesTotal:           6,
			nNodesBehind:          5,
			expectedNodesToRemove: 2,
		},
		{
			name:                  "Even case",
			nNodesTotal:           6,
			nNodesBehind:          3,
			expectedNodesToRemove: 0,
		},
		{
			name:                  "More than half nodes okay",
			nNodesTotal:           6,
			nNodesBehind:          2,
			expectedNodesToRemove: 0,
		},
	}

	for _, tCase := range testCases {
		t.Logf("Running test: %s\n", tCase.name)
		nodesToRemove := decideWhichNodesToRemove(
			tCase.nNodesTotal, getFullNodeSlice(tCase.nNodesBehind))
		assert.Equal(t, tCase.expectedNodesToRemove, len(nodesToRemove))
	}
}
