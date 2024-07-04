package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/envoy"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"istio.io/istio/pkg/util/protomarshal"
)

const envoyRouterCheckTool string = "router_check_tool"

type RootCommand struct {
	Gateway          string
	Verbosity        int
	GenerateOnly     bool
	RouterCheckFlags RouterCheckFlags
}

type RouterCheckFlags struct {
	ConfigDir               string
	TestDir                 string
	Details                 bool
	DisableDeprecationCheck bool
	OnlyShowFailures        bool
	FailUnder               float64
	CoverageAll             bool
	OutputDir               string
	DetailedCoverage        bool
}

func NewCmdRoot() (*cobra.Command, error) {
	ctx := context.Background()
	rootCmd := &RootCommand{}
	cmd := &cobra.Command{
		Use: "istio-router-check",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if _, err := exec.LookPath(envoyRouterCheckTool); err != nil {
				return fmt.Errorf("missing route table check tool: %w", err)
			}
			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: slog.Level(rootCmd.Verbosity),
			}))
			ctx = logr.NewContextWithSlogLogger(ctx, logger)
			cmd.SetContext(ctx)
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			return rootCmd.Run(cmd.Context())
		},
		SilenceUsage: true,
	}
	// Router check flags
	cmd.Flags().StringVarP(&rootCmd.RouterCheckFlags.ConfigDir, "config-dir", "c", "", "directory containing virtualservices")
	cmd.Flags().StringVarP(&rootCmd.RouterCheckFlags.TestDir, "test-dir", "t", "", "directory containing tests")
	cmd.Flags().BoolVarP(&rootCmd.RouterCheckFlags.Details, "details", "d", true, "print detailed information about the test results")
	cmd.Flags().BoolVarP(&rootCmd.RouterCheckFlags.DisableDeprecationCheck, "disable-deprecation-check", "", true, "disable deprecation check")
	cmd.Flags().BoolVarP(&rootCmd.RouterCheckFlags.OnlyShowFailures, "only-show-failures", "", false, "only show failures")
	cmd.Flags().Float64VarP(&rootCmd.RouterCheckFlags.FailUnder, "fail-under", "f", 0.0, "threshold for failure")
	cmd.Flags().BoolVarP(&rootCmd.RouterCheckFlags.CoverageAll, "covall", "", false, "measure coverage by checking all route fields")
	cmd.Flags().StringVarP(&rootCmd.RouterCheckFlags.OutputDir, "output-dir", "o", "", "output directory for coverage information")
	cmd.Flags().BoolVarP(&rootCmd.RouterCheckFlags.DetailedCoverage, "detailed-coverage", "", false, "print detailed coverage information")

	// Root flags
	cmd.Flags().IntVarP(&rootCmd.Verbosity, "", "v", 1, "log verbosity level")
	cmd.Flags().StringVarP(&rootCmd.Gateway, "gateway", "", "", "Only consider VirtualServices bound to this gateway")
	cmd.Flags().BoolVarP(&rootCmd.GenerateOnly, "generate-only", "", false, "Only generate the tests and routes, do not run the router check tool")

	requiredFlags := []string{"config-dir", "test-dir"}
	for _, flag := range requiredFlags {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			return nil, fmt.Errorf("failed to mark flag %q required: %w", flag, err)
		}
	}

	return cmd, nil
}

func (c *RootCommand) Run(ctx context.Context) error {
	log := logr.FromContextOrDiscard(ctx)
	tempDir, err := os.MkdirTemp("", ".router-check-tool-")
	if err != nil {
		return fmt.Errorf("failed to create temp dir: %w", err)
	}

	testsFile := filepath.Join(tempDir, "tests.json")
	routeFile := filepath.Join(tempDir, "routes.json")

	if err := c.prepareTests(ctx, testsFile); err != nil {
		return fmt.Errorf("failed to prepare tests: %w", err)
	}

	if err := c.prepareRoutes(ctx, routeFile); err != nil {
		return fmt.Errorf("failed to prepare routes: %w", err)
	}

	if c.GenerateOnly {
		log.Info("skipping router check tool")
		return nil
	}

	args := c.routerCheckFlags(routeFile, testsFile)
	routerCheck := exec.Command(envoyRouterCheckTool, args...)

	log.V(3).Info("running command", "command", routerCheck.String())
	out, err := routerCheck.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", out)
	}

	fmt.Println(string(out))
	return nil
}

func (c *RootCommand) prepareTests(ctx context.Context, outputFile string) error {
	log := logr.FromContextOrDiscard(ctx)

	log.V(3).Info("reading tests", "dir", c.RouterCheckFlags.TestDir)
	tests, err := envoy.ReadTests(c.RouterCheckFlags.TestDir)
	if err != nil {
		return fmt.Errorf("failed to read test files: %w", err)
	}

	rawTests, err := json.Marshal(tests)
	if err != nil {
		return fmt.Errorf("failed to marshal tests: %w", err)
	}
	log.V(3).Info("writing tests", "file", outputFile)
	err = os.WriteFile(outputFile, rawTests, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write tests: %w", err)
	}
	return nil
}

func (c *RootCommand) prepareRoutes(ctx context.Context, outputFile string) error {
	log := logr.FromContextOrDiscard(ctx)

	log.V(3).Info("reading virtualservices")
	cfg, err := envoy.ReadCRDs(c.RouterCheckFlags.ConfigDir)
	if err != nil {
		return fmt.Errorf("failed to read config files: %w", err)
	}

	routeGen := envoy.NewRouteGenerator(
		envoy.WithConfigs(cfg),
		envoy.WithGateway(c.Gateway),
	)
	routes, err := routeGen.Routes()
	if err != nil {
		return fmt.Errorf("failed to generate routes: %w", err)
	}
	if len(routes) != 1 {
		return fmt.Errorf("expected exactly one route, got %d. Parsed %d configs", len(routes), len(cfg))
	}
	route := routes[0]
	raw, err := protomarshal.ToJSON(route)
	if err != nil {
		return fmt.Errorf("failed to marshal route: %w", err)
	}

	log.V(3).Info("writing route", "route", route.Name, "file", outputFile)
	err = os.WriteFile(outputFile, []byte(raw), os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write route: %w", err)
	}
	return nil
}

func (c *RootCommand) routerCheckFlags(configFile, testFile string) []string {
	args := []string{
		"--config-path", configFile,
		"--test-path", testFile,
	}
	if c.RouterCheckFlags.Details {
		args = append(args, "--details")
	}
	if c.RouterCheckFlags.DisableDeprecationCheck {
		args = append(args, "--disable-deprecation-check")
	}
	if c.RouterCheckFlags.OnlyShowFailures {
		args = append(args, "--only-show-failures")
	}
	if c.RouterCheckFlags.FailUnder != 0.0 {
		args = append(args, "--fail-under", fmt.Sprintf("%f", c.RouterCheckFlags.FailUnder))
	}
	if c.RouterCheckFlags.CoverageAll {
		args = append(args, "--covall")
	}
	if c.RouterCheckFlags.OutputDir != "" {
		args = append(args, "--output-dir", c.RouterCheckFlags.OutputDir)
	}
	if c.RouterCheckFlags.DetailedCoverage {
		args = append(args, "--detailed-coverage")
	}

	return args
}
