/*
Copyright © 2021 Ruohang Feng <rh@vonng.com>

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

// metaCmd represents the meta command
var metaCmd = &cobra.Command{
	Use:   "meta",
	Short: "setup meta node",
	Long: `SYNOPSIS:
    
    init               setup meta nodes
    fetch              fetch pigsty resource from internet (pigsty.tgz)
    repo               unzip files/pkg.tgz to /www/pigsty and add file repo
    cache              cache local repo pkgs to files/pkg.tgz
    
EXAMPLES:

    1. download pigsty resource to ~/pigsty
        pigsty meta fetch -d ~/pigsty
    2. unzip pigsty offline installation pkgs
        pigsty meta repo
    3. make new pkg cache to files/pkg.tgz
        pigsty meta cache

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("meta called")
	},
}

func init() {
	rootCmd.AddCommand(metaCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// metaCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// metaCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
