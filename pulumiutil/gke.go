package pulumiutil

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// GKEConfig holds the configuration for creating a GKE cluster
type GKEConfig struct {
	GCPProject  string
	Network     string
	Subnetwork  string
	MachineType string
	NodeCount   int
	ClusterName string
}

// NewGKEConfigFromEnv initializes a GKEConfig from environment variables
func NewGKEConfigFromEnv() (*GKEConfig, error) {
	nodeCount, err := strconv.Atoi(os.Getenv("NODE_COUNT"))
	if err != nil {
		return nil, fmt.Errorf("Error converting NODE_COUNT to integer: %v", err)
	}

	return &GKEConfig{
		GCPProject:  os.Getenv("GCP_PROJECT"),
		Network:     os.Getenv("NETWORK"),
		Subnetwork:  os.Getenv("SUB_NETWORK"),
		MachineType: os.Getenv("MACHINE_TYPE"),
		NodeCount:   nodeCount,
		ClusterName: os.Getenv("GKE_CLUSTER_NAME"),
	}, nil
}

// CreateGKECluster creates a GKE cluster with the specified configuration
func CreateGKECluster(ctx *pulumi.Context, config *GKEConfig) (*container.Cluster, error) {
	engineVersions, err := container.GetEngineVersions(ctx, &container.GetEngineVersionsArgs{})
	if err != nil {
		return nil, err
	}
	masterVersion := engineVersions.LatestMasterVersion

	cluster, err := container.NewCluster(ctx, config.ClusterName, &container.ClusterArgs{
		DeletionProtection: pulumi.Bool(false),
		Network:            pulumi.String(config.Network),
		Subnetwork:         pulumi.String(config.Subnetwork),
		InitialNodeCount:   pulumi.Int(config.NodeCount),
		MinMasterVersion:   pulumi.String(masterVersion),
		NodeVersion:        pulumi.String(masterVersion),
		NodeConfig: &container.ClusterNodeConfigArgs{
			MachineType: pulumi.String(config.MachineType),
			OauthScopes: pulumi.StringArray{
				pulumi.String("https://www.googleapis.com/auth/compute"),
				pulumi.String("https://www.googleapis.com/auth/devstorage.read_only"),
				pulumi.String("https://www.googleapis.com/auth/logging.write"),
				pulumi.String("https://www.googleapis.com/auth/monitoring"),
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return cluster, nil
}
