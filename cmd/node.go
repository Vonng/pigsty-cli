/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"context"
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "node",
	Short: "setup database node",
	Long: `node -- setup pigsty database nodes

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, cls := range EX.Config.Clusters {
			if cls.Name == "meta" {
				continue
			}
			for _, ins := range cls.Instances {
				if varLimit != "" && !ins.MatchNames(varLimits) {
					continue
				}
				if EX.Config.IsMetaNode(ins.IP) {
					fmt.Printf("%-15s\t*%-31s\t%d.%s.%s\n", ins.IP, ins.Name, ins.Seq, ins.Role, ins.Cluster.Name)
				} else {
					fmt.Printf("%-15s\t%-32s\t%d.%s.%s\n", ins.IP, ins.Name, ins.Seq, ins.Role, ins.Cluster.Name)
				}
			}
		}
		return nil
	},
}

var nodeInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init database node",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("node.yml"),
			exec.WithName("node init"),
			exec.WithLimit(varLimit),
			exec.WithTags(varTags...),
		)
		if varForce {
			job.Opts.ExtraVars["dcs_exists_action"] = "clean"
		}
		return job.Run(context.TODO())
	},
}

var nodeRemoveCmd = &cobra.Command{
	Use:   "tune",
	Short: "remove node",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("node-remove.yml"),
			exec.WithName("node remove"),
			exec.WithLimit(varLimit),
		)
		if varForce {
			job.Opts.ExtraVars["yum_remove"] = true
		}
		return job.Run(context.TODO())
	},
}

var nodeDcsCmd = &cobra.Command{
	Use:   "init",
	Short: "init database node consul",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("node.yml"),
			exec.WithName("node dcs init"),
			exec.WithLimit(varLimit),
			exec.WithTags("dcs"),
		)
		if varForce {
			job.Opts.ExtraVars["dcs_exists_action"] = "clean"
		}
		return job.Run(context.TODO())
	},
}

var nodeTuneCmd = &cobra.Command{
	Use:   "tune",
	Short: "tune database node",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("node.yml"),
			exec.WithName("node tune"),
			exec.WithLimit(varLimit),
			exec.WithTags("node_tuned"),
		)
		if varMode != "" {
			varMode = strings.ToLower(varMode)
			if varMode == "oltp" || varMode == "olap" || varMode == "crit" || varMode == "tiny" {
				logrus.Warnf("unknown profile %s specified", varMode)
			}
			job.Opts.ExtraVars["node_tune"] = varMode
		}
		return job.Run(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(nodesCmd)

	// node init
	nodeInitCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")
	nodesCmd.AddCommand(nodeInitCmd)

	// node remove
	nodeRemoveCmd.Flags().BoolVarP(&varForce, "force", "f", false, "uninstall packages")
	nodesCmd.AddCommand(nodeRemoveCmd)

	// node dcs
	nodeDcsCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")
	nodesCmd.AddCommand(nodeDcsCmd)

	// node tune
	nodeTuneCmd.Flags().StringVarP(&varMode, "mode", "m", "", "pgsql config template: oltp|olap|crit|tiny|other...")
	nodesCmd.AddCommand(nodeTuneCmd)
}
