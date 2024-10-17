/*
Copyright © 2024 Aaron U'Ren

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"os"

	goflags "flag"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "kube-quota",
	Short: "A small Kubernetes CLI that looks at Kubernetes resources usage and quotas",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	fs := goflags.NewFlagSet("", goflags.PanicOnError)
	klog.InitFlags(fs)
	rootCmd.Flags().AddGoFlagSet(fs)
}

func getFlagString(cmd *cobra.Command, flagName string) string {
	val, err := cmd.Flags().GetString(flagName)
	if err != nil {
		klog.Fatalf("Could not get string flag: %s - %v", flagName, err)
	}
	return val
}

func getFlagBool(cmd *cobra.Command, flagName string) bool {
	val, err := cmd.Flags().GetBool(flagName)
	if err != nil {
		klog.Fatalf("Could not get boolean flag: %s - %v", flagName, err)
	}
	return val
}
