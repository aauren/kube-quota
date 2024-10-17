/*
Copyright © 2024 Aaron U'Ren
*/
package cmd

import (
	"context"
	"fmt"
	"log"

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
	addQuota  bool
	ns        string
	quotaName string
)

func init() {
	rootCmd.AddCommand(workloadCmd)

	workloadCmd.Flags().StringVarP(&ns, "namespace", "n", "", "namespace to search within")
	workloadCmd.Flags().BoolVarP(&addQuota, "add-quota", "a", false, "add quota to bottom of results")
	workloadCmd.Flags().StringVarP(&quotaName, "quota-name", "q", "", "specific name of the quota you want to search for (by default "+
		"it will show a single quota within the requested namespace if there is only one found)")
}

func workloadRun(_ *cobra.Command, _ []string) {
	err := workloadValidateInput()
	if err != nil {
		log.Fatalf("Encountered error while parsing input: %v", err)
	}

	ctx := context.Background()
	pl, err := workloads.GetPodsByNamespace(ctx, ns)
	if err != nil {
		klog.Fatalf("could not get pods by namespace: %v", err)
	}

	tq := quota.QuotaForPodList(pl)
	tbl := cli.TabularizeTotalQuota(tq)

	tbl.Print()
}

func workloadValidateInput() error {
	if ns == "" {
		return fmt.Errorf("workload command currently requires a namespace be specified")
	}

	return nil
}
