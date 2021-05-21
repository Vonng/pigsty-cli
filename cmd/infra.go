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
	"context"
	"fmt"
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/spf13/cobra"
)

// infraCmd represents the infra command
var infraCmd = &cobra.Command{
	Use:   "infra",
	Short: "setup infrastructure",
	Long: `infra -- setup pigsty infrastructure on meta node

    init           complete infra init on meta node
    repo           setup local yum repo 
    ca             setup local ca 
    dns            setup dnsmasq nameserver
    prometheus     setup prometheus & alertmanager
    grafana        setup grafana service
    loki           setup loki logging collector
    haproxy        refresh haproxy admin page index      
    target         refresh prometheus static targets

`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(EX.Config.InfraInfo())
	},
}

var infraInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init pigsty infra on meta nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		job := EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra init"),
			exec.WithLimit(varLimit),
		)
		return job.Run(context.TODO())
	},
}

var infraRepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "setup pigsty repo on meta nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra repo init"),
			exec.WithTags("repo"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraNodeCmd = &cobra.Command{
	Use:   "node",
	Short: "setup node infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra node init"),
			exec.WithTags("node"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraCaCmd = &cobra.Command{
	Use:   "ca",
	Short: "setup ca on meta node",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra ca init"),
			exec.WithTags("ca"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraDnsCmd = &cobra.Command{
	Use:   "dns",
	Short: "setup dns infrastructure",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra dns init"),
			exec.WithTags("nameserver"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraPrometheusCmd = &cobra.Command{
	Use:   "prometheus",
	Short: "setup pigsty prometheus on meta nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra prometheus init"),
			exec.WithTags("prometheus"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraPrometheusReloadCmd = &cobra.Command{
	Use:   "reload",
	Short: "recreate prometheus targets and reload",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra prometheus reload"),
			exec.WithTags("prometheus_targets", "prometheus_reload"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraGrafanaCmd = &cobra.Command{
	Use:   "grafana",
	Short: "setup pigsty grafana on meta nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra grafana init"),
			exec.WithTags("grafana"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraLokiCmd = &cobra.Command{
	Use:   "loki",
	Short: "setup loki on meta nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra-loki.yml"),
			exec.WithName("infra loki init"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraPgsqlCmd = &cobra.Command{
	Use:   "pgsql",
	Short: "setup pgsql on meta nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra pgsql init"),
			exec.WithTags("pgsql"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraHaproxyCmd = &cobra.Command{
	Use:   "haproxy",
	Short: "update haproxy index page",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra haproxy index update"),
			exec.WithTags("nginx_haproxy", "nginx_restart"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

var infraTargetCmd = &cobra.Command{
	Use:   "target",
	Short: "update prometheus filesd targets",
	RunE: func(cmd *cobra.Command, args []string) error {
		return EX.NewJob(
			exec.WithPlaybook("infra.yml"),
			exec.WithName("infra filesd target"),
			exec.WithTags("prometheus_targets", "prometheus_reload"),
			exec.WithLimit(varLimit),
		).Run(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(infraCmd)
	infraCmd.AddCommand(infraInitCmd)
	infraCmd.AddCommand(infraRepoCmd)
	infraCmd.AddCommand(infraNodeCmd)
	infraCmd.AddCommand(infraCaCmd)
	infraCmd.AddCommand(infraDnsCmd)
	infraCmd.AddCommand(infraPrometheusCmd)
	infraCmd.AddCommand(infraGrafanaCmd)
	infraCmd.AddCommand(infraLokiCmd)
	infraCmd.AddCommand(infraPgsqlCmd)
	infraCmd.AddCommand(infraHaproxyCmd)
	infraCmd.AddCommand(infraTargetCmd)
	infraPrometheusCmd.AddCommand(infraPrometheusReloadCmd)
}
