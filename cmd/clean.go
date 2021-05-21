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
	"github.com/Vonng/pigsty-cli/exec"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "remove pgsql cluster/instance",
	Long: `SYNOPSIS:

    # you MUST limit clean targets, other wise all instances will be purged!
    clean <-l cluster>              remove pgsql cluster
    clean <-l instance|ip>          remove pgsql instance
    
    # clean specific component
    clean all                       remove pgsql and uninstall packages
    clean service                   remove pgsql service
    clean monitor                   remove pgsql monitor
    clean postgres                  remove pgsql postgres pgbouncer patroni
    clean dcs                       remove consul dcs agent 
    clean packages                  remove pgsql and dcs packages

`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("pgsql remove"),
			exec.WithLimit(varLimit),
			exec.WithTags(varTags...),
		)
		return job.Run(context.TODO())
	},
}

var cleanAllCmd = &cobra.Command{
	Use:   "all",
	Short: "clean all component and uninstall packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("pgsql remove all"),
			exec.WithLimit(varLimit),
			exec.WithTags("service", "monitor", "postgres", "dcs", "pkgs"),
		)
		return job.Run(context.TODO())
	},
}

var cleanServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "clean service only",
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("clean service"),
			exec.WithLimit(varLimit),
			exec.WithTags("service"),
		)
		return job.Run(context.TODO())
	},
}

var cleanMonitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "clean monitor only",
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("clean monitor"),
			exec.WithLimit(varLimit),
			exec.WithTags("monitor"),
		)
		return job.Run(context.TODO())
	},
}

var cleanPostgresCmd = &cobra.Command{
	Use:   "postgres",
	Short: "clean postgres service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("clean postgres"),
			exec.WithLimit(varLimit),
			exec.WithTags("postgres"),
		)
		return job.Run(context.TODO())
	},
}

var cleanDcsCmd = &cobra.Command{
	Use:   "dcs",
	Short: "clean dcs component",
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("clean dcs"),
			exec.WithLimit(varLimit),
			exec.WithTags("dcs"),
		)
		return job.Run(context.TODO())
	},
}

var cleanPackagesCmd = &cobra.Command{
	Use:   "packages",
	Short: "clean service only",
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-remove.yml"),
			exec.WithName("clean packages"),
			exec.WithLimit(varLimit),
			exec.WithTags("packages"),
		)
		return job.Run(context.TODO())
	},
}

var cleanPromtailCmd = &cobra.Command{
	Use:   "promtail",
	Short: "clean promtail service",
	RunE: func(cmd *cobra.Command, args []string) error {
		if varLimit == "" && !varForce {
			logrus.Fatalf("YOU MUST USE LIMIT(-l) WITH CLEAN! or use force(-f) to overwrite")
			os.Exit(1)
		}
		job := EX.NewJob(
			exec.WithPlaybook("pgsql-promtail.yml"),
			exec.WithName("clean promtail"),
			exec.WithLimit(varLimit),
			exec.WithTags("promtail_clean"),
		)
		job.Opts.ExtraVars["promtail_clean"] = true
		return job.Run(context.TODO())
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
	cleanCmd.Flags().BoolVarP(&varForce, "force", "f", false, "force execution")

	// clean sub command
	cleanCmd.AddCommand(cleanAllCmd)
	cleanCmd.AddCommand(cleanServiceCmd)
	cleanCmd.AddCommand(cleanMonitorCmd)
	cleanCmd.AddCommand(cleanPostgresCmd)
	cleanCmd.AddCommand(cleanDcsCmd)
	cleanCmd.AddCommand(cleanPackagesCmd)
	cleanCmd.AddCommand(cleanPromtailCmd)

}
