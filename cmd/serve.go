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
	varServeDataDir       string
	varServePublicDir     string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch pigsty API server",
	Long: `Usage:
    pigsty serve -i|--inventory   inventory file      (pigsty.yml by default)
                 -L|--listen-addr listen_address      (:9633 by default)
                 -P|--public-dir  public resource dir (embed by default)
                 -D|--data-dir     log dir            (/tmp/pigsty by default)
                  (will create <public_dir>/log for logging purpose)

EXAMPLE:

    # run server
        pigsty serve -i ~/pigsty/pigsty.yml

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
		logrus.Debugf("pigsty server @ %s , use config %s, data dir %s, public dir %s", varServeListenAddress, varConfig, varServeDataDir, varServePublicDir)
		server.InitDefaultServer(varServeListenAddress, varConfig, varServeDataDir, varServePublicDir)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&varServeListenAddress, "listen-addr", "L", ":9633", "listen address")
	serveCmd.Flags().StringVarP(&varServeDataDir, "data-dir", "D", "/tmp/pigsty", "temporary resource dir")
	serveCmd.Flags().StringVarP(&varServePublicDir, "public-dir", "P", "embed", "public resource dir")
}
