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
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/prometheus/common/log"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	varConfig   string
	varLimit    string
	varTags     []string
	varLimits   []string
	varLimitMap map[string]int
)

// Ex is the default command executor
var EX *exec.Executor

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pigsty",
	Short: "Pigsty Command-Line Interface v0.8",
	Long: `
NAME
    pigsty -- Pigsty Command-Line Interface v0.8 

SYNOPSIS               

    meta               setup meta nodes            init|fetch|repo|cache|ansible
    node               setup database nodes        init|tune|dcs|remove|ping|bash|ssh|admin 
    pgsql              setup postgres clusters     init|node|dcs|postgres|template|business|monitor|service|monly|remove
    infra              setup infrastructure        init|ca|dns|prometheus|grafana|loki|haproxy|target
    clean              clean pgsql clusters        all|service|monitor|postgres|dcs
    config             mange pigsty config file    init|edit|info|dump|path
    serve              run pigsty API server       init|start|stop|restart|reload|status
    demo               setup local demo            init|up|new|clean|start|dns
    log                watch system log            query|postgres|patroni|pgbouncer|message
    pg                 pg operational tasks        user|db|svc|hba|log|psql|deploy|backup|restore|vacuum|repack


EXAMPLES

    1. infra summary
        pigsty infra

    2. pgsql clusters summary
        pigsty pgsql

    3. pigsty nodes summary
        pigsty node

    4. create pgsql cluster 'pg-test'
        pigsty pgsql init -l pg-test

    5. add new instance 10.10.10.13 of cluster 'pg-test'
        pigsty pgsql init -l 10.10.10.13

    6. remove cluster 'pg-test'
        pigsty clean -l pg-test

    7. create user dbuser_vonng on cluster 'pg-test'
        pigsty pg user dbuser_vonng -l pg-test

    8. create database test2 on cluster 'pg-test'
        pigsty pg db test -l pg-test


`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&varConfig, "inventory", "i", "./pigsty.yml", "inventory file")
	rootCmd.PersistentFlags().StringVarP(&varLimit, "limit", "l", "", "limit execution hosts")
	rootCmd.PersistentFlags().StringSliceVarP(&varTags, "tags", "t", []string{}, "limit execution tasks")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	log.Debugf("args: --config=%s --limit=%s --tags=%s", varConfig, varLimit, varTags)

	// parse limit string into list and maps
	varLimits = strings.Split(varLimit, ",")
	varLimitMap = make(map[string]int)
	for i, l := range varLimits {
		varLimitMap[l] = i
	}

	// use PIGSTY_CONFIG env instead of default args
	if envConfigPath := os.Getenv("PIGSTY_CONFIG"); varConfig == `./pigsty.yml` && envConfigPath != "" {
		varConfig = envConfigPath
		log.Debugf("get config path from env PIGSTY_CONFIG: %s", varConfig)
	}

	// TODO: load config via postgres if PGURL is given instead of inventory filepath
	if strings.HasPrefix(varConfig, "postgres://") {
		log.Info("pgsql inventory not implemented yet")
		os.Exit(1)
	}

	// build command executor from config path
	EX = exec.NewExecutor(varConfig)
	if EX == nil {
		log.Fatal("fail to create playbook executor")
		os.Exit(1)
	}
}
