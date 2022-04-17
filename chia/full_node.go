package chia

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	statusNotSynced = "Not Synced"
	statusSyncing   = "Syncing"
)

type fullNode struct {
	Id     string
	Height int64
	Synced bool
}

func getConnectedNodes() ([]fullNode, error) {

	out, err := exec.Command("chia", "show", "-c").CombinedOutput()
	if err != nil {
		return nil, err
	}

	nodes := make([]fullNode, 0, 50)
	scanLineForNodeHeight := false
	for _, line := range strings.Split(string(out), "\n") {

		if strings.HasPrefix(line, "FULL_NODE") {
			lineFields := strings.Fields(line)

			nodes = append(nodes, fullNode{
				Id: strings.TrimSuffix(lineFields[3], "..."),
			})
			scanLineForNodeHeight = true
			continue
		}

		if scanLineForNodeHeight {
			lineFields := strings.Fields(line)
			node := nodes[len(nodes)-1]

			height, err := strconv.ParseInt(lineFields[1], 10, 0)
			if err != nil {
				return nil, fmt.Errorf("Error when converting height of full node %s: %s\n", node.Id, err)
			}

			node.Height = height
			nodes[len(nodes)-1] = node

			scanLineForNodeHeight = !scanLineForNodeHeight
		}
	}

	return nodes, nil
}

func getOwnNodeStatus() (fullNode, error) {

	out, err := exec.Command("chia", "show", "-s").CombinedOutput()
	if err != nil {
		return fullNode{}, err
	}

	node := fullNode{}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "Node ID") {
			lineFields := strings.Fields(line)
			node.Id = lineFields[2]
		}

		if strings.Contains(line, "Height") && strings.Contains(line, "Time") {
			lineFields := strings.Fields(line)
			height, err := strconv.ParseInt(lineFields[8], 10, 0)
			if err != nil {
				return fullNode{}, fmt.Errorf("Error when converting height of full node %s: %s", node.Id, err)
			}
			node.Height = height
		}

		if strings.Contains(line, "Current Blockchain Status:") {
			node.Synced = !strings.Contains(line, statusSyncing) && !strings.Contains(line, statusNotSynced)
		}
	}

	if len(node.Id) == 0 {
		return fullNode{}, fmt.Errorf("Full node is not up running.")
	}

	return node, nil
}

func filterNodesWhichAreFarBehind(
	connectedNodes []fullNode, ownNode fullNode, heightTolerance int64,
) []fullNode {
	nodesToRemove := make([]fullNode, 0, 50)
	for _, node := range connectedNodes {
		if node.Height < ownNode.Height-heightTolerance {
			nodesToRemove = append(nodesToRemove, node)
		}
	}
	return nodesToRemove
}

func RunFullNodeCheck(runEverySeconds, heightTolerance int64) {

	runInfinitely := runEverySeconds != 0

	for iRun := uint64(0); runInfinitely || iRun < 1; iRun++ {
		connectedNodes, err := getConnectedNodes()
		if err != nil {
			fmt.Printf("Error when fetching connected nodes: %s\n", err)
			continue
		}
		fmt.Printf("Found %d connected nodes\n", len(connectedNodes))

		ownFullNode, err := getOwnNodeStatus()
		if err != nil {
			fmt.Printf("Error when fetchincliArgs.runEverySeconds g own node status: %s\n", err)
			continue
		}
		fmt.Printf("Own node status %+v\n", ownFullNode)

		nodesToRemove := filterNodesWhichAreFarBehind(
			connectedNodes, ownFullNode, heightTolerance)
		fmt.Printf("Removing %d nodes\n", len(nodesToRemove))

		for _, node := range nodesToRemove {
			fmt.Printf("Removing node %s with height %d\n", node.Id, node.Height)
			_, err := exec.Command("chia", "show", "-r", node.Id).CombinedOutput()
			if err != nil {
				fmt.Printf("Error removing node '%s': %s", node.Id, err)
				continue
			}
		}

		time.Sleep(time.Duration(runEverySeconds) * time.Second)
	}

}
