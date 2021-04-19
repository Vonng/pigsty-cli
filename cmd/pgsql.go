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

const ()

var (
	varFormatDetail bool
	varFormatYaml   bool
	varFormatJson   bool
	varMode         string // oltp olap crit tiny
	varForce        bool
)

// pgsqlCmd represents the pgsql command
var pgsqlCmd = &cobra.Command{
	Use:   "pgsql",
	Short: "setup pgsql clusters",
	Long: `SYNOPSIS:
    
    pgsql list                      show pgsql cluster definition
    pgsql init                      init new postgres clusters or instances
    pgsql node                      init pgsql node
    pgsql dcs                       init pgsql dcs (consul)
    pgsql postgres                  init postgres service (postgres|patroni|pgbouncer)
    pgsql monitor                   init monitor components
    pgsql service                   init services provider
    pgsql pgbouncer                 init pgbouncer service
    pgsql template                  init postgres template database
    pgsql business                  init postgres business users and databases
    pgsql promtail                  init promtail log collect agent
    pgsql config                    init patroni config template
    pgsql monly                     init monitor system in monitor-only mode
    pgsql hba                       init hba rule files
    pgsql remove                    remove postgres cluster or instances

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, c := range EX.Config.Clusters {
			if c.Name == "meta" || varLimit != "" && !c.MatchNames(varLimits) {
				continue // skip meta group, and skip unmatched clusters if limit is set
			}
			fmt.Println(c.Repr(parseOutputFormat()))
		}
		return nil
	},
}

var pgsqlListCmd = &cobra.Command{
	Use:   "list",
	Short: "list pgsql clusters",
	Long:  `list -- list pgsql clusters`,
	RunE: func(cmd *cobra.Command, args []string) error {
		for _, c := range EX.Config.Clusters {
			if c.Name == "meta" || varLimit != "" && !c.MatchNames(varLimits) {
				continue // skip meta group, and skip unmatched clusters if limit is set
			}
			fmt.Println(c.Repr(parseOutputFormat()))
		}
		return nil
	},
}

var pgsqlInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init pgsql on targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("pgsql init"),
			exec.WithLimit(varLimit),
			exec.WithTags(varTags...),
		)
		if varForce {
			job.Opts.ExtraVars["pg_exists_action"] = "clean"
			job.Opts.ExtraVars["dcs_exists_action"] = "clean"
		}
		return job.Run(context.TODO())
	},
}

var pgsqlNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "init pgsql node",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("node init"),
			exec.WithLimit(varLimit),
			exec.WithTags("node"),
		)
		if varMode != "" {
			if strings.HasSuffix(varMode, ".yml") {
				varMode = strings.TrimSuffix(varMode, ".yml")
			}
			varMode = strings.ToLower(varMode)
			job.Opts.ExtraVars["node_tune"] = varMode
		}
		return job.Run(context.TODO())
	},
}

var pgsqlDcsCmd = &cobra.Command{
	Use:   "dcs",
	Short: "init pgsql dcs service",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("dcs init"),
			exec.WithLimit(varLimit),
			exec.WithTags("dcs"),
		)
		if varForce {
			job.Opts.ExtraVars["dcs_exists_action"] = "clean"
		}
		return job.Run(context.TODO())
	},
}

var pgsqlPostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "init pgsql postgres service",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("postgres init"),
			exec.WithLimit(varLimit),
			exec.WithTags("postgres"),
		)
		if varForce {
			job.Opts.ExtraVars["pg_exists_action"] = "clean"
		}
		return job.Run(context.TODO())
	},
}

var pgsqlMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "init pgsql monitor",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("monitor init"),
			exec.WithLimit(varLimit),
			exec.WithTags("monitor"),
		)
		return job.Run(context.TODO())
	},
}

var pgsqlServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "init pgsql service",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("service init"),
			exec.WithLimit(varLimit),
			exec.WithTags("service"),
		)
		return job.Run(context.TODO())
	},
}

var pgsqlPromtailCmd = &cobra.Command{
	Use:   "promtail",
	Short: "init promtail",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-promtail.yml"),
			exec.WithName("init promtail"),
			exec.WithLimit(varLimit),
		)
		if varForce {
			job.Opts.ExtraVars["promtail_clean"] = true
		}
		return job.Run(context.TODO())
	},
}

/************************************************************************
*  special subtasks
*************************************************************************/
var pgsqlPgbouncerCmd = &cobra.Command{
	Use:   "pgbouncer",
	Short: "init pgsql pgbouncer service",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("pgbouncer init"),
			exec.WithLimit(varLimit),
			exec.WithTags("pgbouncer"),
		)
		return job.Run(context.TODO())
	},
}

var pgsqlTemplateCmd = &cobra.Command{
	Use:   "template",
	Short: "init pgsql template",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("template init"),
			exec.WithLimit(varLimit),
			exec.WithTags("pg_init"),
		)
		return job.Run(context.TODO())
	},
}

var pgsqlBusinessCmd = &cobra.Command{
	Use:   "business",
	Short: "init pgsql business user & db",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("userdb init"),
			exec.WithLimit(varLimit),
			exec.WithTags("pg_user", "pg_db"),
		)
		return job.Run(context.TODO())
	},
}

var pgsqlConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "config pgsql with template",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("pgsql config"),
			exec.WithTags("pg_config"),
			exec.WithLimit(varLimit),
		)
		if varMode != "" {
			if !strings.HasSuffix(varMode, ".yml") {
				varMode += ".yml"
			}
			varMode = strings.ToLower(varMode)
			job.Opts.ExtraVars["pg_conf"] = varMode
		}
		return job.Run(context.TODO())
	},
}

var pgsqlMonlyCmd = &cobra.Command{
	Use:   "monly",
	Short: "init pgsql monitor only",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-monitor.yml"),
			exec.WithName("monly init"),
			exec.WithLimit(varLimit),
			exec.WithTags("monitor"),
		)
		return job.Run(context.TODO())
	},
}

var pgsqlHbaCmd = &cobra.Command{
	Use:   "hba",
	Short: "init pgsql hba rules",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql.yml"),
			exec.WithName("hba init"),
			exec.WithLimit(varLimit),
			exec.WithTags("pg_hba"),
		)
		job.Opts.ExtraVars["pg_reload"] = true
		return job.Run(context.TODO())
	},
}

var pgsqlRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove pgsql from targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("pgsql init"),
			exec.WithLimit(varLimit),
		)
		if varForce {
			job.Opts.Forks = "10"
		}
		return job.Run(context.TODO())
	},
}

// parseOutputFormat will parse -d -j -y flags and turn into format string
func parseOutputFormat() string {
	if (varFormatYaml && varFormatJson) || (varFormatYaml && varFormatDetail) || (varFormatJson && varFormatDetail) {
		logrus.Errorf("format args -d -j -y can not be used together")
		return "default"
	}
	if varFormatYaml {
		return "yaml"
	}
	if varFormatJson {
		return "json"
	}
	if varFormatDetail {
		return "detail"
	}
	return "default"
}

func init() {
	rootCmd.AddCommand(pgsqlCmd)

	// pgsql list
	pgsqlCmd.AddCommand(pgsqlListCmd)
	pgsqlCmd.Flags().BoolVarP(&varFormatDetail, "detail", "d", false, "detail format")
	pgsqlCmd.Flags().BoolVarP(&varFormatYaml, "yaml", "y", false, "yaml output")
	pgsqlCmd.Flags().BoolVarP(&varFormatJson, "json", "j", false, "json output")

	// pgsql init
	pgsqlCmd.AddCommand(pgsqlInitCmd)
	pgsqlInitCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")

	// pgsql node
	pgsqlCmd.AddCommand(pgsqlNodeCmd)
	pgsqlNodeCmd.Flags().StringVarP(&varMode, "mode", "m", "", "pgsql node tune template: oltp|olap|crit|tiny|other...")

	// pgsql dcs
	pgsqlCmd.AddCommand(pgsqlDcsCmd)
	pgsqlDcsCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")

	// pgsql postgres
	pgsqlCmd.AddCommand(pgsqlPostgresCmd)
	pgsqlPostgresCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")

	// pgsql monitor
	pgsqlCmd.AddCommand(pgsqlMonitorCmd)

	// pgsql service
	pgsqlCmd.AddCommand(pgsqlServiceCmd)

	// pgsql promtail (beta)
	pgsqlCmd.AddCommand(pgsqlPromtailCmd)
	pgsqlPromtailCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")

	/**********************
	* special subtasks
	**********************/
	// pgsql pgbouncer
	pgsqlCmd.AddCommand(pgsqlPgbouncerCmd)

	// pgsql template
	pgsqlCmd.AddCommand(pgsqlTemplateCmd)

	// pgsql business
	pgsqlCmd.AddCommand(pgsqlBusinessCmd)

	// pgsql config
	pgsqlCmd.AddCommand(pgsqlConfigCmd)
	pgsqlConfigCmd.Flags().StringVarP(&varMode, "mode", "m", "", "pgsql config template: oltp|olap|crit|tiny|other...")

	// pgsql monly
	pgsqlCmd.AddCommand(pgsqlMonlyCmd)

	// pgsql hba
	pgsqlCmd.AddCommand(pgsqlHbaCmd)

	// pgsql remove
	pgsqlCmd.AddCommand(pgsqlRemoveCmd)
	pgsqlRemoveCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")

}
