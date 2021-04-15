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
	"fmt"
	"github.com/spf13/cobra"
)

// nodesCmd represents the nodes command
var nodesCmd = &cobra.Command{
	Use:   "node",
	Short: "setup database nodes",
	Long: `node -- setup pigsty database nodes

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

func init() {
	rootCmd.AddCommand(nodesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// nodesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// nodesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
