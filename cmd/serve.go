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
	"github.com/Vonng/pigsty-cli/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	varServeListenAddress string
	varServePublicDir     string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch pigsty API server",
	Long: `Usage:
    pigsty serve -L|--listen-addr listen_address
                 -D|--public-dir  resource dir

EXAMPLE:

    # run server
        pigsty serve -L :9633 -D ./public -i ~/pigsty/pigsty.yml

    # get config
        curl http://localhost:9633/api/v1/config
   
    # post config
        curl -X POST http://localhost:9633/api/v1/config -d@<config.yml>
   
    # get cluster info
        curl -X GET http://localhost:9633/api/v1/pgsql/:cluster/info
   
    # create cluster (SSE: type=:cluster )
        curl -X GET http://localhost:9633/api/v1/pgsql/:cluster/init
        args: force=true will force init (remove existing cluster)
        args: tags=only execute partial of the playbook
   
    # remove cluster (SSE: type=:cluster )
        curl -X GET http://localhost:9633/api/v1/pgsql/:cluster/remove
    
    # e.g: remove pg-test cluster
        curl http://localhost:9633/api/v1/pgsql/pg-test/remove

    # remove pg-test monitor and service
        curl http://localhost:9633/api/v1/pgsql/pg-test/remove?tags=monitor,service

`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("pigsty server @ %s , use config %s, public %s", varServeListenAddress, varConfig, varServeListenAddress)
		server.InitDefaultServer(varConfig, varServePublicDir, varServeListenAddress)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&varServeListenAddress, "listen-addr", "L", ":9633", "listen address")
	serveCmd.Flags().StringVarP(&varServePublicDir, "public-dir", "D", "./public", "public resource dir")
}
