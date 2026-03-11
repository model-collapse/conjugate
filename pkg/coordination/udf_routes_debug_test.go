package coordination

import (
	"context"
	"strings"
	"testing"

	"github.com/conjugate/conjugate/pkg/common/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

// TestUDFRoutesRegistration tests that UDF routes are properly registered
func TestUDFRoutesRegistration(t *testing.T) {
	logger, _ := zap.NewDevelopment()

	cfg := &config.CoordinationConfig{
		NodeID:     "coord-debug",
		BindAddr:   "127.0.0.1",
		RESTPort:   9299,
		MasterAddr: "127.0.0.1:8000",
	}

	// Use a dedicated prometheus registry to avoid duplicate registration panics
	node, err := NewCoordinationNodeWithRegistry(cfg, logger, prometheus.NewRegistry())
	require.NoError(t, err)
	defer node.Stop(context.Background())

	t.Logf("UDF Registry initialized: %v", node.udfRegistry != nil)
	t.Logf("UDF Runtime initialized: %v", node.udfRuntime != nil)

	// Get all routes from Gin
	routes := node.ginRouter.Routes()

	t.Logf("Total routes registered: %d", len(routes))

	udfRouteCount := 0
	for _, route := range routes {
		t.Logf("Route: %s %s", route.Method, route.Path)
		if strings.HasPrefix(route.Path, "/api/v1/udfs") {
			udfRouteCount++
			t.Logf("  ^^ UDF route found!")
		}
	}

	if node.udfRegistry != nil {
		require.Greater(t, udfRouteCount, 0, "UDF routes should be registered when UDF registry is available")
	}
}
