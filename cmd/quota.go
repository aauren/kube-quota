package cmd

import (
	"context"

	"github.com/aauren/kube-quota/pkg/cli"
	kubequota "github.com/aauren/kube-quota/pkg/kubernetes/quota"
	"github.com/aauren/kube-quota/pkg/quota"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// workloadCmd represents the quotaUsed command
var quotaCmd = &cobra.Command{
	Use:   "quota",
	Short: "See the quota",
	Args:  cobra.MaximumNArgs(1),
	Run:   quotaRun,
}

func init() {
	rootCmd.AddCommand(quotaCmd)

	quotaCmd.Flags().StringP("namespace", "n", "", "namespace to search within")
	quotaCmd.Flags().StringP("quota-name", "q", "", "specific name of the quota you want to search for (by default "+
		"it will show a single quota within the requested namespace if there is only one found)")
}

func quotaRun(cmd *cobra.Command, args []string) {
	err := quotaValidateInput(cmd)
	if err != nil {
		klog.Fatalf("Encountered error while parsing input: %v", err)
	}

	ctx := context.Background()
	ns := getFlagString(cmd, "namespace")
	quotaName := ""
	if len(args) > 0 {
		quotaName = args[0]
	}

	kq, err := kubequota.FindByNSAndName(ctx, ns, quotaName)
	if err != nil {
		klog.Fatalf("could not get pods by namespace: %v", err)
	}

	q := quota.ForKubeQuota(kq)

	tbl := cli.CreateTableWriter()
	cli.TabularizeKubeQuota(tbl, q)
	tbl.Render()
}

func quotaValidateInput(cmd *cobra.Command) error {
	err := cmd.MarkFlagRequired("namespace")
	if err != nil {
		return err
	}

	return nil
}
