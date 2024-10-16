/*
Copyright © 2024 Aaron U'Ren
*/
package cmd

import (
	"context"
	"log"

	"github.com/aauren/kube-quota/pkg/cli"
	kubequota "github.com/aauren/kube-quota/pkg/kubernetes/quota"
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

func init() {
	rootCmd.AddCommand(workloadCmd)

	workloadCmd.Flags().StringP("namespace", "n", "", "namespace to search within")
	workloadCmd.Flags().BoolP("add-quota", "a", false, "add quota to bottom of results")
	workloadCmd.Flags().StringP("quota-name", "q", "", "specific name of the quota you want to search for (by default it will show a "+
		"single quota within the requested namespace if there is only one found)")
}

func workloadRun(cmd *cobra.Command, _ []string) {
	err := workloadValidateInput(cmd)
	if err != nil {
		log.Fatalf("Encountered error while parsing input: %v", err)
	}

	ctx := context.Background()
	ns := getFlagString(cmd, "namespace")
	aq := getFlagBool(cmd, "add-quota")
	qn := getFlagString(cmd, "quota-name")
	pl, err := workloads.GetPodsByNamespace(ctx, ns)
	if err != nil {
		klog.Fatalf("could not get pods by namespace: %v", err)
	}

	tq := quota.QuotaForPodList(pl)

	tbl := cli.CreateTableWriter()
	cli.TabularizeTotalQuota(tbl, tq)

	if aq {
		kq, err := kubequota.FindByNSAndName(ctx, ns, qn)
		if err != nil {
			klog.Fatalf("could not get pods by namespace: %v", err)
		}
		q := quota.ForKubeQuota(kq)
		tbl.ResetHeaders()
		cli.AddNewSection(tbl, "quota")
		cli.TabularizeKubeQuota(tbl, q)
	}

	tbl.Render()
}

func workloadValidateInput(cmd *cobra.Command) error {
	err := cmd.MarkFlagRequired("namespace")
	if err != nil {
		return err
	}

	return nil
}
