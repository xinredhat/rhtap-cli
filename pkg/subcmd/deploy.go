package subcmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/redhat-appstudio/tssc/pkg/chartfs"
	"github.com/redhat-appstudio/tssc/pkg/config"
	"github.com/redhat-appstudio/tssc/pkg/flags"
	"github.com/redhat-appstudio/tssc/pkg/installer"
	"github.com/redhat-appstudio/tssc/pkg/k8s"

	"github.com/spf13/cobra"
)

// Deploy is the deploy subcommand.
type Deploy struct {
	cmd    *cobra.Command   // cobra command
	logger *slog.Logger     // application logger
	flags  *flags.Flags     // global flags
	cfg    *config.Config   // installer configuration
	cfs    *chartfs.ChartFS // embedded filesystem
	kube   *k8s.Kube        // kubernetes client

	chartPath          string // path of the chart when deploying a single chart
	valuesTemplatePath string // path to the values template file
}

var _ Interface = &Deploy{}

const deployDesc = `
Deploys the TSSC platform components. The installer looks the the informed
configuration to identify the products to be installed, and the dependencies to be
resolved.

The deployment configuration file describes the sequence of Helm charts to be
applied, on the attribute 'tssc.dependencies[]'.

The platform configuration is rendered from the values template file
(--values-template), this configuration payload is given to all Helm charts.

The installer resources are embedded in the executable, these resources are
employed by default, to use local files just point the "config.yaml" file to
find the dependencies in the local filesystem.

A single chart can be deployed by specifying its path. E.g.:
	tssc deploy charts/tssc-openshift
`

// Cmd exposes the cobra instance.
func (d *Deploy) Cmd() *cobra.Command {
	return d.cmd
}

// log logger with contextual information.
func (d *Deploy) log() *slog.Logger {
	return d.flags.LoggerWith(
		d.logger.With(flags.ValuesTemplateFlag, d.valuesTemplatePath))
}

// Complete verifies the object is complete.
func (d *Deploy) Complete(args []string) error {
	var err error
	if d.cfg, err = bootstrapConfig(d.cmd.Context(), d.kube); err != nil {
		return err
	}
	if len(args) == 1 {
		d.chartPath = args[0]
	}
	return nil
}

// Validate asserts the requirements to start the deployment are in place.
func (d *Deploy) Validate() error {
	return k8s.EnsureOpenShiftProject(
		d.cmd.Context(),
		d.log(),
		d.kube,
		d.cfg.Installer.Namespace,
	)
}

// Run deploys the enabled dependencies listed on the configuration.
func (d *Deploy) Run() error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	cfs, err := chartfs.NewChartFS(cwd)
	if err != nil {
		return err
	}

	d.log().Debug("Reading values template file")
	valuesTmpl, err := cfs.ReadFile(d.valuesTemplatePath)
	if err != nil {
		return fmt.Errorf("failed to read values template file: %w", err)
	}

	d.log().Debug("Installing dependencies...")
	var deps []config.Dependency
	if d.chartPath == "" {
		// Installing each Helm Chart dependency from the configuration, only
		// selecting the Helm Charts that are enabled.
		deps = d.cfg.GetEnabledDependencies(d.log())
	} else {
		// Installing a single Chart dependency
		dep, err := d.cfg.GetDependency(d.log(), d.chartPath)
		if err != nil {
			return err
		}
		deps = append(deps, *dep)
	}
	for ix, dep := range deps {
		fmt.Printf("\n\n############################################################\n")
		fmt.Printf("# [%d/%d] Deploying '%s' in '%s'.\n", ix+1, len(deps), dep.Chart, dep.Namespace)
		fmt.Printf("############################################################\n")

		i := installer.NewInstaller(d.log(), d.flags, d.kube, cfs, &dep)

		err := i.SetValues(d.cmd.Context(), &d.cfg.Installer, string(valuesTmpl))
		if err != nil {
			return err
		}
		if d.flags.Debug {
			i.PrintRawValues()
		}

		if err := i.RenderValues(); err != nil {
			return err
		}
		if d.flags.Debug {
			i.PrintValues()
		}

		err = i.Install(d.cmd.Context())
		// Delete temporary resources
		if err := k8s.RetryDeleteResources(d.cmd.Context(), d.kube, d.cfg.Installer.Namespace); err != nil {
			d.log().Debug(err.Error())
		}
		if err != nil {
			return err
		}
		fmt.Printf("############################################################\n\n")
	}

	fmt.Printf("Deployment complete!\n")
	return nil
}

// NewDeploy instantiates the deploy subcommand.
func NewDeploy(
	logger *slog.Logger,
	f *flags.Flags,
	cfs *chartfs.ChartFS,
	kube *k8s.Kube,
) Interface {
	d := &Deploy{
		cmd: &cobra.Command{
			Use:          "deploy [chart]",
			Short:        "Rollout TSSC platform components",
			Long:         deployDesc,
			SilenceUsage: true,
		},
		logger:    logger.WithGroup("deploy"),
		flags:     f,
		cfs:       cfs,
		kube:      kube,
		chartPath: "",
	}
	flags.SetValuesTmplFlag(d.cmd.PersistentFlags(), &d.valuesTemplatePath)
	return d
}
