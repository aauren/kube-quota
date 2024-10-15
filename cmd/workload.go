/*
Copyright © 2024 Aaron U'Ren
*/
package cmd

import (
	"context"

	"github.com/aauren/kube-quota/pkg/cli"
	"github.com/aauren/kube-quota/pkg/kubernetes/workloads"
	"github.com/aauren/kube-quota/pkg/quota"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// workloadCmd represents the quotaUsed command
var workloadCmd = &cobra.Command{
	Use:   "workload",
	Short: "See the quota usage by workload",
	Run:   workloadRun,
}

var (
	ns string
)

func init() {
	rootCmd.AddCommand(workloadCmd)

	workloadCmd.Flags().StringVarP(&ns, "namespace", "n", "", "namespace to search within")
}

func workloadRun(_ *cobra.Command, _ []string) {
	ctx := context.Background()
	pl, err := workloads.GetPodsByNamespace(ctx, ns)
	if err != nil {
		klog.Fatalf("could not get pods by namespace: %v", err)
	}

	tq := quota.QuotaForPodList(pl)
	tbl := cli.TabularizeTotalQuota(tq)

	tbl.Print()
}
