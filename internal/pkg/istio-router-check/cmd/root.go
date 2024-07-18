package cmd

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/envoy"
	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/helpers"
	"github.com/getyourguide/istio-config-validator/internal/pkg/parser"
	"github.com/go-logr/logr"
	"github.com/spf13/cobra"
	"istio.io/api/networking/v1alpha3"
	"istio.io/istio/pkg/util/protomarshal"
)

type RootCommand struct {
	ConfigDir    string
	TestDir      string
	ConvertTests bool
	OutputDir    string
	Gateway      string
	Verbosity    int
}

const (
	LevelInfo  = 0
	LevelDebug = 9
)

func NewCmdRoot() (*cobra.Command, error) {
	ctx := context.Background()
	rootCmd := &RootCommand{}
	cmd := &cobra.Command{
		Use: "istio-router-check",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

	cmd.Flags().IntVarP(&rootCmd.Verbosity, "", "v", LevelInfo, "Log verbosity level")
	cmd.Flags().StringVarP(&rootCmd.Gateway, "gateway", "", "", "Only consider VirtualServices bound to this gateway (i.e: istio-system/istio-ingressgateway)")
	cmd.Flags().StringVarP(&rootCmd.ConfigDir, "config-dir", "c", "", "Directory with Istio VirtualService and Gateway files")
	cmd.Flags().StringVarP(&rootCmd.TestDir, "test-dir", "t", "", "Directory with Envoy test files")
	cmd.Flags().BoolVarP(&rootCmd.ConvertTests, "convert-tests", "", false, "Convert istio-config-validator tests into Envoy tests")
	cmd.Flags().StringVarP(&rootCmd.OutputDir, "output-dir", "o", "", "Directory to output Envoy routes and tests")

	for _, flag := range []string{"output-dir", "config-dir", "test-dir"} {
		if err := cmd.MarkFlagRequired(flag); err != nil {
			return nil, fmt.Errorf("failed to mark flag %q required: %w", flag, err)
		}
		if err := cmd.MarkFlagDirname(flag); err != nil {
			return nil, fmt.Errorf("failed to mark flag %q as dirname: %w", flag, err)
		}
	}

	if err := cmd.Flags().MarkHidden("convert-tests"); err != nil {
		return nil, fmt.Errorf("failed to mark flag hidden: %w", err)
	}

	return cmd, nil
}

func (c *RootCommand) Run(ctx context.Context) error {
	if err := os.MkdirAll(c.OutputDir, os.ModePerm); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := c.prepareRoutes(ctx); err != nil {
		return fmt.Errorf("failed to prepare routes: %w", err)
	}

	if c.ConvertTests {
		if err := c.prepareTests(ctx); err != nil {
			return fmt.Errorf("failed to prepare istio-config-validator tests: %w", err)
		}
		return nil
	}

	if err := c.prepareEnvoyTests(ctx); err != nil {
		return fmt.Errorf("failed to prepare envoy tests: %w", err)
	}

	return nil
}

func (c *RootCommand) prepareEnvoyTests(ctx context.Context) error {
	log := logr.FromContextOrDiscard(ctx)
	if c.TestDir == "" {
		log.V(LevelDebug).Info("no envoy test directory provided")
		return nil
	}

	log.Info("reading tests", "dir", c.TestDir)
	tests, err := envoy.ReadTests(c.TestDir)
	if err != nil {
		return fmt.Errorf("failed to read envoy test files: %w", err)
	}

	rawTests, err := json.Marshal(tests)
	if err != nil {
		return fmt.Errorf("failed to marshal tests: %w", err)
	}
	outputFile := filepath.Join(c.OutputDir, "tests.json")
	log.Info("writing tests", "file", outputFile)
	err = os.WriteFile(outputFile, rawTests, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write tests: %w", err)
	}
	return nil
}

func (c *RootCommand) prepareRoutes(ctx context.Context) error {
	log := logr.FromContextOrDiscard(ctx)

	log.Info("reading virtualservices")
	cfg, err := envoy.ReadCRDs(c.ConfigDir)
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
	if len(routes) <= 0 {
		return fmt.Errorf("expected at least one route, got %d. Parsed %d configs", len(routes), len(cfg))
	}
	for _, route := range routes {
		raw, err := protomarshal.ToJSON(route)
		if err != nil {
			return fmt.Errorf("failed to marshal route: %w", err)
		}
		routeName := fmt.Sprintf("route_%s_%s.json", cmp.Or(strings.ReplaceAll(c.Gateway, "/", "_"), "sidecar"), route.Name)
		routeFile := filepath.Join(c.OutputDir, routeName)
		log.Info("writing route", "route", route.Name, "file", routeFile)
		err = os.WriteFile(routeFile, []byte(raw), os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to write route: %w", err)
		}
	}
	return nil
}

func (c *RootCommand) prepareTests(ctx context.Context) error {
	if c.TestDir == "" {
		return nil
	}
	log := logr.FromContextOrDiscard(ctx)
	log.Info("reading tests", "dir", c.TestDir)

	oldFiles, err := helpers.WalkYAML(c.TestDir)
	if err != nil {
		return fmt.Errorf("could not read directory %s: %w", c.TestDir, err)
	}
	strict := true
	testCases, err := parser.ParseTestCases(oldFiles, strict)
	if err != nil {
		return fmt.Errorf("parsing testcases failed: %w", err)
	}
	var newTests envoy.Tests
	for _, tc := range testCases {
		inputs, err := tc.Request.Unfold()
		if err != nil {
			return fmt.Errorf("could not unfold request: %w", err)
		}
		if !tc.WantMatch {
			log.V(LevelDebug).Info("skipping negative test", "test", tc.Description, "reason", "router_check_tool does not support negative tests")
			continue
		}
		if tc.Rewrite != nil {
			log.V(LevelDebug).Info("skipping rewrite test", "test", tc.Description, "reason", "format assertion is different in envoy tests")
			continue
		}
		for _, req := range inputs {
			var reqHeaders []envoy.Header
			for key, value := range req.Headers {
				reqHeaders = append(reqHeaders, envoy.Header{Key: key, Value: value})
			}
			input := envoy.Input{
				SSL:                      true,
				Authority:                req.Authority,
				Method:                   cmp.Or(req.Method, http.MethodGet),
				Path:                     cmp.Or(req.URI, "/"),
				AdditionalRequestHeaders: reqHeaders,
			}
			validate, err := convertValidate(input, tc)
			if err != nil {
				return fmt.Errorf("could not convert test %q: %w", tc.Description, err)
			}
			newTests.Tests = append(newTests.Tests, envoy.Test{
				TestName: fmt.Sprintf("%s: method=%q authority=%q path=%q headers=%+v", tc.Description, input.Method, input.Authority, input.Path, input.AdditionalRequestHeaders),
				Input:    input,
				Validate: validate,
			})
		}
	}
	outputFile := filepath.Join(c.OutputDir, "tests.json")
	log.Info("writing tests", "file", outputFile)
	raw, err := json.Marshal(newTests)
	if err != nil {
		return fmt.Errorf("failed to marshal tests: %w", err)
	}
	err = os.WriteFile(outputFile, raw, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write tests: %w", err)
	}

	return nil
}

func convertValidate(input envoy.Input, tc *parser.TestCase) (envoy.Validate, error) {
	output := envoy.Validate{}
	if tc.Route != nil {
		var route *v1alpha3.HTTPRouteDestination
		for _, r := range tc.Route {
			if r.GetWeight() >= route.GetWeight() {
				route = r
			}
		}
		output.ClusterName = fmt.Sprintf("outbound|%d|%s|%s",
			cmp.Or(route.GetDestination().GetPort().GetNumber(), 80),
			route.GetDestination().GetSubset(),
			route.GetDestination().GetHost(),
		)
	}
	if tc.Redirect != nil {
		authority := cmp.Or(tc.Redirect.GetAuthority(), input.Authority)
		scheme := cmp.Or(tc.Redirect.GetScheme(), "https")
		output.PathRedirect = fmt.Sprintf("%s://%s%s", scheme, authority, tc.Redirect.GetUri())
	}
	return output, nil
}
