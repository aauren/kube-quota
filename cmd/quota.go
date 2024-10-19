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

	// Create our context and get any arguments the user may have set
	ctx := context.Background()
	ns := getFlagString(cmd, "namespace")
	quotaName := ""
	if len(args) > 0 {
		quotaName = args[0]
	}

	// Get all of our data and format it.
	kq, err := kubequota.FindByNSAndName(ctx, ns, quotaName)
	if err != nil {
		klog.Fatalf("could not get pods by namespace: %v", err)
	}
	q := quota.ForKubeQuota(kq)

	// Setup our table and add our header.
	tbl := cli.CreateTableWriter()
	cli.AddTableHeader(tbl, []string{"Name"}, q)

	// Add our data to the table.
	err = cli.AddRow(tbl, q, []string{"Quota"})
	if err != nil {
		klog.Fatalf("Could not add row to table: %v", err)
	}

	// Render our table
	tbl.Render()
}

func quotaValidateInput(cmd *cobra.Command) error {
	err := cmd.MarkFlagRequired("namespace")
	if err != nil {
		return err
	}

	return nil
}
