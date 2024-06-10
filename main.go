package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pulumi/pulumi-gcp/sdk/v7/go/gcp/container"
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

	pulumi.Run(func(ctx *pulumi.Context) error {
		engineVersions, err := container.GetEngineVersions(ctx, &container.GetEngineVersionsArgs{})
		if err != nil {
			return err
		}
		masterVersion := engineVersions.LatestMasterVersion

		cluster, err := container.NewCluster(ctx, clusterName, &container.ClusterArgs{
			DeletionProtection: pulumi.Bool(false),
			Network:            pulumi.String(network),
			Subnetwork:         pulumi.String(subnetwork),
			InitialNodeCount:   pulumi.Int(nodeCount),
			MinMasterVersion:   pulumi.String(masterVersion),
			NodeVersion:        pulumi.String(masterVersion),
			NodeConfig: &container.ClusterNodeConfigArgs{
				MachineType: pulumi.String(machineType),
				OauthScopes: pulumi.StringArray{
					pulumi.String("https://www.googleapis.com/auth/compute"),
					pulumi.String("https://www.googleapis.com/auth/devstorage.read_only"),
					pulumi.String("https://www.googleapis.com/auth/logging.write"),
					pulumi.String("https://www.googleapis.com/auth/monitoring"),
				},
			},
		})
		if err != nil {
			return err
		}

		ctx.Export("kubeconfig", generateKubeconfig(cluster.Endpoint, cluster.Name, cluster.MasterAuth))

		return nil
	})
}

func generateKubeconfig(clusterEndpoint pulumi.StringOutput, clusterName pulumi.StringOutput,
	clusterMasterAuth container.ClusterMasterAuthOutput) pulumi.StringOutput {
	context := pulumi.Sprintf("demo_%s", clusterName)

	return pulumi.Sprintf(`apiVersion: v1
clusters:
- cluster:
    certificate-authority-data: %s
    server: https://%s
  name: %s
contexts:
- context:
    cluster: %s
    user: %s
  name: %s
current-context: %s
kind: Config
preferences: {}
users:
- name: %s
  user:
    exec:
      apiVersion: client.authentication.k8s.io/v1beta1
      command: gke-gcloud-auth-plugin
      installHint: Install gke-gcloud-auth-plugin for use with kubectl by following
        https://cloud.google.com/blog/products/containers-kubernetes/kubectl-auth-changes-in-gke
      provideClusterInfo: true
`,
		clusterMasterAuth.ClusterCaCertificate().Elem(),
		clusterEndpoint, context, context, context, context, context, context)
}
