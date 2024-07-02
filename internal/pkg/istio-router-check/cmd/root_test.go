package cmd_test

import (
	"os"
	"testing"

	"github.com/getyourguide/istio-config-validator/internal/pkg/istio-router-check/cmd"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	cmdRoot, err := cmd.NewCmdRoot()
	require.NoError(t, err)
	require.NotNil(t, cmdRoot)

	t.Run("it should fail with missing required flags", func(t *testing.T) {
		// TODO(cainelli): add router_check_tool to CI
		if os.Getenv("CI") == "true" {
			t.Skip("skip as it requires router_check_tool binary not yet in CI")
		}
		err = cmdRoot.Execute()
		require.ErrorContains(t, err, "required flag(s)")
		require.ErrorContains(t, err, "config-dir")
		require.ErrorContains(t, err, "test-dir")
	})

	t.Run("it should run the test", func(t *testing.T) {
		// TODO(cainelli): add router_check_tool to CI
		if os.Getenv("CI") == "true" {
			t.Skip("skip as it requires router_check_tool binary not yet in CI")
		}
		cmdRoot.SetArgs([]string{"--config-dir", "testdata/virtualservice.yml", "--test-dir", "testdata/test.yml"})
		err = cmdRoot.Execute()
		require.NoError(t, err)
	})
}
