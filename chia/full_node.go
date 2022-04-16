package chia

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	defaultHeightToleranceInBlocks = 5000
	defaultRunEveryMinutes         = 0
	statusNotSynced                = "Not Synced"
	statusSyncing                  = "Syncing"
)

type FullNode struct {
	Id     string
	Height int64
	Synced bool
}

func getConnectedNodes() ([]FullNode, error) {

	out, err := exec.Command("chia", "show", "-c").CombinedOutput()
	if err != nil {
		return nil, err
	}

	nodes := make([]FullNode, 0, 50)
	scanLineForNodeHeight := false
	for _, line := range strings.Split(string(out), "\n") {

		if strings.HasPrefix(line, "FULL_NODE") {
			lineFields := strings.Fields(line)

			nodes = append(nodes, FullNode{
				Id: strings.TrimSuffix(lineFields[3], "..."),
			})
			scanLineForNodeHeight = true
			continue
		}

		if scanLineForNodeHeight {
			lineFields := strings.Fields(line)
			fullNode := nodes[len(nodes)-1]

			height, err := strconv.ParseInt(lineFields[1], 10, 0)
			if err != nil {
				return nil, fmt.Errorf("Error when converting height of full node %s: %s\n", fullNode.Id, err)
			}

			fullNode.Height = height
			nodes[len(nodes)-1] = fullNode

			scanLineForNodeHeight = !scanLineForNodeHeight
		}
	}

	return nodes, nil
}

func getOwnNodeStatus() (FullNode, error) {

	out, err := exec.Command("chia", "show", "-s").CombinedOutput()
	if err != nil {
		log.Fatalln(err)
	}

	fullNode := FullNode{}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "Node ID") {
			lineFields := strings.Fields(line)
			fullNode.Id = lineFields[2]
		}

		if strings.Contains(line, "Height") && strings.Contains(line, "Time") {
			lineFields := strings.Fields(line)
			height, err := strconv.ParseInt(lineFields[8], 10, 0)
			if err != nil {
				log.Fatalf("Error when converting height of full node %s: %s\n", fullNode.Id, err)
			}
			fullNode.Height = height
		}

		if strings.Contains(line, "Current Blockchain Status:") {
			fullNode.Synced = !strings.Contains(line, statusSyncing) && !strings.Contains(line, statusNotSynced)
		}
	}

	if len(fullNode.Id) == 0 {
		return FullNode{}, fmt.Errorf("Full node is not up running.")
	}

	return fullNode, nil
}

func filterNodesWhichAreFarBehind(
	connectedNodes []FullNode, ownNode FullNode, heightTolerance int64,
) []FullNode {
	nodesToRemove := make([]FullNode, 0, 50)
	for _, node := range connectedNodes {
		if node.Height < ownNode.Height-heightTolerance {
			nodesToRemove = append(nodesToRemove, node)
		}
	}
	return nodesToRemove
}

func getErrorLoggingFn(runIndefinitely bool) func(format string, a ...interface{}) {
	if runIndefinitely {
		return log.Fatalf
	}

	return func(format string, a ...interface{}) {
		fmt.Printf(format, a...)
	}
}

func RunFullNodeCheck(runEveryMins, heightTolerance int64) {

	runIndefinitely := runEveryMins == 0

	// we go fatal if we don't run indefinitely and trouble arises.
	logErrorFn := getErrorLoggingFn(runIndefinitely)

	for iRun := uint64(0); runIndefinitely && iRun < 1; iRun++ {
		connectedNodes, err := getConnectedNodes()
		if err != nil {
			logErrorFn("Error when fetching connected nodes: %s\n", err)
		}
		fmt.Printf("Found %d connected nodes\n", len(connectedNodes))

		ownFullNode, err := getOwnNodeStatus()
		if err != nil {
			logErrorFn("Error when fetching own node status: %s\n", err)
		}
		fmt.Printf("Own node status %+v\n", ownFullNode)

		nodesToRemove := filterNodesWhichAreFarBehind(
			connectedNodes, ownFullNode, heightTolerance)
		fmt.Printf("Removing %d nodes\n", len(nodesToRemove))

		for _, node := range nodesToRemove {
			fmt.Printf("Removing node %s with height %d\n", node.Id, node.Height)
			_, err := exec.Command("chia", "show", "-r", node.Id).CombinedOutput()
			if err != nil {
				logErrorFn("Error removing node '%s': %s", node.Id, err)
			}
		}

		time.Sleep(time.Duration(runEveryMins) * time.Minute)
	}

}
