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
	"github.com/Vonng/pigsty-cli/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	varServerListenAddress string
	varServerDataDir       string
	varServerPublicDir     string
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Launch pigsty API server",
	Long: `Usage:
    pigsty server -i|--inventory   inventory file      (pigsty.yml by default)
                 -L|--listen-addr listen_address      (:9633 by default)
                 -P|--public-dir  public resource dir (embed by default)
                 -D|--data-dir     log dir            (/tmp/pigsty by default)
                  (will create <public_dir>/log for logging purpose)

EXAMPLE:

    # run server
        pigsty server -i ~/pigsty/pigsty.yml

    # get config
        curl http://localhost:9633/api/v1/config
   
    # post config (YAML config as body)
        curl -X POST http://localhost:9633/api/v1/config -d@<pigsty.yml>
   
    # list jobs
        curl -X GET http://localhost:9633/api/v1/jobs

    # get current job
        curl -X GET http://localhost:9633/api/v1/job

    # create new job ( pgsql init @ pg-test )
        curl -X POST http://localhost:9633/api/v1/job?playbook=pgsql&cluster=pg-test

    # create new job ( pgsql remove @ pg-test2 )
        curl -X POST http://localhost:9633/api/v1/job?playbook=pgsql-remove&cluster=pg-test2

    # cancel job
        curl -X DELETE http://localhost:9633/api/v1/job

    # list logs
        curl -X GET http://localhost:9633/api/v1/logs

    # get latest job log
        curl -X GET http://localhost:9633/api/v1/log/latest

    # get job log by job id
        curl -X GET http://localhost:9633/api/v1/log/:jobid

`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Debugf("pigsty server run @ %s , use config %s, data dir %s, public dir %s", varServerListenAddress, varConfig, varServerDataDir, varServerPublicDir)
		server.InitDefaultServer(varServerListenAddress, varConfig, varServerDataDir, varServerPublicDir)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&varServerListenAddress, "listen-addr", "L", ":9633", "listen address")
	serverCmd.Flags().StringVarP(&varServerDataDir, "data-dir", "D", "/tmp/pigsty", "temporary resource dir")
	serverCmd.Flags().StringVarP(&varServerPublicDir, "public-dir", "P", "embed", "public resource dir")
}
