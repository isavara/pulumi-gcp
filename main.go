package main

import (
	"fmt"
	"os"
	"strconv"

	"gcp-go-gke/pulumiutil"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	// Get environment variables
	gcpProject := os.Getenv("GCP_PROJECT")
	network := os.Getenv("NETWORK")
	subnetwork := os.Getenv("SUB_NETWORK")
	machineType := os.Getenv("MACHINE_TYPE")
	nodeCountStr := os.Getenv("NODE_COUNT")
	clusterName := os.Getenv("GKE_CLUSTER_NAME")

	// Print environment variables for debugging
	fmt.Println("GCP_PROJECT:", gcpProject)
	fmt.Println("NETWORK:", network)
	fmt.Println("SUB_NETWORK:", subnetwork)
	fmt.Println("MACHINE_TYPE:", machineType)
	fmt.Println("NODE_COUNT:", nodeCountStr)
	fmt.Println("GKE_CLUSTER_NAME:", clusterName)

	// Convert nodeCountStr to an integer
	nodeCount, err := strconv.Atoi(nodeCountStr)
	if err != nil {
		fmt.Println("Error converting NODE_COUNT to integer:", err)
		return
	}

	// Create GKE configuration
	config := &pulumiutil.GKEConfig{
		GCPProject:  gcpProject,
		Network:     network,
		Subnetwork:  subnetwork,
		MachineType: machineType,
		NodeCount:   nodeCount,
		ClusterName: clusterName,
	}

	// Run Pulumi to create the GKE cluster
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Create the GKE cluster
		cluster, err := pulumiutil.CreateGKECluster(ctx, config)
		if err != nil {
			return err
		}

		// Export the kubeconfig for the created cluster
		ctx.Export("kubeconfig", pulumiutil.GenerateKubeconfig(cluster.Endpoint, cluster.Name, cluster.MasterAuth))

		return nil
	})
}
