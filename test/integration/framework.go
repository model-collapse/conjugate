package integration

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/conjugate/conjugate/pkg/common/config"
	"github.com/conjugate/conjugate/pkg/coordination"
	"github.com/conjugate/conjugate/pkg/data"
	"github.com/conjugate/conjugate/pkg/master"
	"go.uber.org/zap"
)

// TestCluster represents a multi-node test cluster
type TestCluster struct {
	t              testing.TB
	logger         *zap.Logger
	baseDir        string
	masterNodes    []*MasterNodeWrapper
	coordNodes     []*CoordNodeWrapper
	dataNodes      []*DataNodeWrapper
	shutdownCh     chan struct{}
	wg             sync.WaitGroup
	startTime      time.Time
	mu             sync.RWMutex
}

// MasterNodeWrapper wraps a master node with test utilities
type MasterNodeWrapper struct {
	Node   *master.MasterNode
	Config *config.MasterConfig
	DataDir string
	Started bool
	mu     sync.RWMutex
}

// CoordNodeWrapper wraps a coordination node with test utilities
type CoordNodeWrapper struct {
	Node    *coordination.CoordinationNode
	Config  *config.CoordinationConfig
	Started bool
	mu      sync.RWMutex
}

// DataNodeWrapper wraps a data node with test utilities
type DataNodeWrapper struct {
	Node    *data.DataNode
	Config  *config.DataNodeConfig
	DataDir string
	Started bool
	mu      sync.RWMutex
}

// ClusterConfig defines the cluster topology
type ClusterConfig struct {
	NumMasters      int
	NumCoordination int
	NumData         int
	StartPorts      PortRange
}

// PortRange defines port allocation
type PortRange struct {
	MasterRaftBase int
	MasterGRPCBase int
	CoordRESTBase  int
	DataGRPCBase   int
}

// getFreePort asks the kernel for a free open port that is ready to use.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

// allocatePortRange allocates a contiguous range of free ports
func allocatePortRange(count int) (int, error) {
	// Get a free port as base
	basePort, err := getFreePort()
	if err != nil {
		return 0, err
	}

	// Try to verify the next ports are also available
	// If any are taken, try again with a new base
	for attempt := 0; attempt < 10; attempt++ {
		allFree := true
		for i := 0; i < count; i++ {
			addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("localhost:%d", basePort+i))
			if err != nil {
				allFree = false
				break
			}
			l, err := net.ListenTCP("tcp", addr)
			if err != nil {
				allFree = false
				break
			}
			l.Close()
		}

		if allFree {
			return basePort, nil
		}

		// Try again with new base
		basePort, err = getFreePort()
		if err != nil {
			return 0, err
		}
	}

	return 0, fmt.Errorf("failed to allocate %d contiguous free ports after 10 attempts", count)
}

// DynamicClusterConfig returns a cluster configuration with dynamically allocated ports.
// This is safe for parallel test execution.
func DynamicClusterConfig() (*ClusterConfig, error) {
	return DynamicClusterConfigWithSize(3, 1, 2)
}

// DynamicClusterConfigWithSize returns a cluster configuration with custom node counts
// and dynamically allocated ports. This is safe for parallel test execution.
func DynamicClusterConfigWithSize(numMasters, numCoord, numData int) (*ClusterConfig, error) {
	// Calculate total ports needed
	// Masters need: raft port + grpc port = 2 per master
	// Coord needs: rest port = 1 per coord
	// Data needs: grpc port = 1 per data
	totalPorts := numMasters*2 + numCoord + numData + 3 // +3 for safety margin

	basePort, err := allocatePortRange(totalPorts)
	if err != nil {
		return nil, err
	}

	return &ClusterConfig{
		NumMasters:      numMasters,
		NumCoordination: numCoord,
		NumData:         numData,
		StartPorts: PortRange{
			MasterRaftBase: basePort,
			MasterGRPCBase: basePort + numMasters,
			CoordRESTBase:  basePort + numMasters*2,
			DataGRPCBase:   basePort + numMasters*2 + numCoord,
		},
	}, nil
}

// DefaultClusterConfig returns a sensible default configuration with fixed ports.
// WARNING: This uses hardcoded ports and is NOT safe for parallel test execution.
// Use DynamicClusterConfig() instead for parallel tests.
func DefaultClusterConfig() *ClusterConfig {
	return &ClusterConfig{
		NumMasters:      3,
		NumCoordination: 1,
		NumData:         2,
		StartPorts: PortRange{
			MasterRaftBase: 19300,
			MasterGRPCBase: 19400,
			CoordRESTBase:  19500,
			DataGRPCBase:   19600,
		},
	}
}

// NewTestClusterWithDynamicPorts creates a test cluster with dynamically allocated ports.
// This is the recommended way to create test clusters for parallel execution.
func NewTestClusterWithDynamicPorts(t testing.TB) (*TestCluster, error) {
	cfg, err := DynamicClusterConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to allocate dynamic ports: %w", err)
	}
	return NewTestCluster(t, cfg)
}

// NewTestCluster creates a new test cluster
func NewTestCluster(t testing.TB, cfg *ClusterConfig) (*TestCluster, error) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}

	baseDir := t.TempDir()

	cluster := &TestCluster{
		t:          t,
		logger:     logger,
		baseDir:    baseDir,
		shutdownCh: make(chan struct{}),
		startTime:  time.Now(),
	}

	// Create master nodes
	if err := cluster.createMasterNodes(cfg); err != nil {
		return nil, fmt.Errorf("failed to create master nodes: %w", err)
	}

	// Create coordination nodes
	if err := cluster.createCoordinationNodes(cfg); err != nil {
		return nil, fmt.Errorf("failed to create coordination nodes: %w", err)
	}

	// Create data nodes
	if err := cluster.createDataNodes(cfg); err != nil {
		return nil, fmt.Errorf("failed to create data nodes: %w", err)
	}

	return cluster, nil
}

// createMasterNodes creates master node instances
func (tc *TestCluster) createMasterNodes(cfg *ClusterConfig) error {
	// Build peer list
	var peers []string
	for i := 0; i < cfg.NumMasters; i++ {
		raftPort := cfg.StartPorts.MasterRaftBase + i
		peers = append(peers, fmt.Sprintf("127.0.0.1:%d", raftPort))
	}

	for i := 0; i < cfg.NumMasters; i++ {
		nodeID := fmt.Sprintf("master-%d", i)
		dataDir := filepath.Join(tc.baseDir, nodeID)

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("failed to create data dir: %w", err)
		}

		masterCfg := &config.MasterConfig{
			NodeID:      nodeID,
			BindAddr:    "127.0.0.1",
			RaftPort:    cfg.StartPorts.MasterRaftBase + i,
			GRPCPort:    cfg.StartPorts.MasterGRPCBase + i,
			DataDir:     dataDir,
			Peers:       peers,
			LogLevel:    "debug",
			MetricsPort: 19700 + i,
		}

		node, err := master.NewMasterNode(masterCfg, tc.logger)
		if err != nil {
			return fmt.Errorf("failed to create master node %s: %w", nodeID, err)
		}

		wrapper := &MasterNodeWrapper{
			Node:    node,
			Config:  masterCfg,
			DataDir: dataDir,
			Started: false,
		}

		tc.masterNodes = append(tc.masterNodes, wrapper)
		tc.logger.Info("Created master node", zap.String("node_id", nodeID))
	}

	return nil
}

// createCoordinationNodes creates coordination node instances
func (tc *TestCluster) createCoordinationNodes(cfg *ClusterConfig) error {
	// Use first master as default master address
	masterAddr := fmt.Sprintf("127.0.0.1:%d", cfg.StartPorts.MasterGRPCBase)

	for i := 0; i < cfg.NumCoordination; i++ {
		nodeID := fmt.Sprintf("coord-%d", i)

		coordCfg := &config.CoordinationConfig{
			NodeID:     nodeID,
			BindAddr:   "127.0.0.1",
			RESTPort:   cfg.StartPorts.CoordRESTBase + i,
			MasterAddr: masterAddr,
			LogLevel:   "debug",
		}

		// Use a custom prometheus registry for each node to avoid duplicate registration in tests
		registry := prometheus.NewRegistry()
		node, err := coordination.NewCoordinationNodeWithRegistry(coordCfg, tc.logger, registry)
		if err != nil {
			return fmt.Errorf("failed to create coordination node %s: %w", nodeID, err)
		}

		wrapper := &CoordNodeWrapper{
			Node:    node,
			Config:  coordCfg,
			Started: false,
		}

		tc.coordNodes = append(tc.coordNodes, wrapper)
		tc.logger.Info("Created coordination node", zap.String("node_id", nodeID))
	}

	return nil
}

// createDataNodes creates data node instances
func (tc *TestCluster) createDataNodes(cfg *ClusterConfig) error {
	masterAddr := fmt.Sprintf("127.0.0.1:%d", cfg.StartPorts.MasterGRPCBase)

	for i := 0; i < cfg.NumData; i++ {
		nodeID := fmt.Sprintf("data-%d", i)
		dataDir := filepath.Join(tc.baseDir, nodeID)

		if err := os.MkdirAll(dataDir, 0755); err != nil {
			return fmt.Errorf("failed to create data dir: %w", err)
		}

		dataCfg := &config.DataNodeConfig{
			NodeID:      nodeID,
			BindAddr:    "127.0.0.1",
			GRPCPort:    cfg.StartPorts.DataGRPCBase + i,
			DataDir:     dataDir,
			MasterAddr:  masterAddr,
			StorageTier: "hot",
			MaxShards:   100,
			LogLevel:    "debug",
		}

		node, err := data.NewDataNode(dataCfg, tc.logger)
		if err != nil {
			return fmt.Errorf("failed to create data node %s: %w", nodeID, err)
		}

		wrapper := &DataNodeWrapper{
			Node:    node,
			Config:  dataCfg,
			DataDir: dataDir,
			Started: false,
		}

		tc.dataNodes = append(tc.dataNodes, wrapper)
		tc.logger.Info("Created data node", zap.String("node_id", nodeID))
	}

	return nil
}

// Start starts all nodes in the cluster
func (tc *TestCluster) Start(ctx context.Context) error {
	tc.logger.Info("Starting test cluster")

	// Start master nodes first
	for _, wrapper := range tc.masterNodes {
		if err := tc.startMasterNode(ctx, wrapper); err != nil {
			return fmt.Errorf("failed to start master node: %w", err)
		}
	}

	// Wait for leader election
	tc.logger.Info("Waiting for leader election")
	if err := tc.WaitForLeader(10 * time.Second); err != nil {
		return fmt.Errorf("leader election failed: %w", err)
	}

	// Start data nodes BEFORE coordination nodes
	// This ensures data nodes register with master before coordination tries to discover them
	for _, wrapper := range tc.dataNodes {
		if err := tc.startDataNode(ctx, wrapper); err != nil {
			return fmt.Errorf("failed to start data node: %w", err)
		}
	}

	// Give data nodes a moment to register with master
	time.Sleep(500 * time.Millisecond)

	// Start coordination nodes
	// They will discover data nodes during their Start() method
	for _, wrapper := range tc.coordNodes {
		if err := tc.startCoordNode(ctx, wrapper); err != nil {
			return fmt.Errorf("failed to start coordination node: %w", err)
		}
	}

	tc.logger.Info("Test cluster started successfully",
		zap.Int("masters", len(tc.masterNodes)),
		zap.Int("coordination", len(tc.coordNodes)),
		zap.Int("data", len(tc.dataNodes)))

	return nil
}

// startMasterNode starts a single master node
func (tc *TestCluster) startMasterNode(ctx context.Context, wrapper *MasterNodeWrapper) error {
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	if wrapper.Started {
		return nil
	}

	if err := wrapper.Node.Start(ctx); err != nil {
		return fmt.Errorf("failed to start master node %s: %w", wrapper.Config.NodeID, err)
	}

	wrapper.Started = true
	tc.logger.Info("Started master node", zap.String("node_id", wrapper.Config.NodeID))
	return nil
}

// startCoordNode starts a single coordination node
func (tc *TestCluster) startCoordNode(ctx context.Context, wrapper *CoordNodeWrapper) error {
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	if wrapper.Started {
		return nil
	}

	// Coordination node may fail to connect to master initially, that's ok
	_ = wrapper.Node.Start(ctx)

	wrapper.Started = true
	tc.logger.Info("Started coordination node", zap.String("node_id", wrapper.Config.NodeID))
	return nil
}

// startDataNode starts a single data node
func (tc *TestCluster) startDataNode(ctx context.Context, wrapper *DataNodeWrapper) error {
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	if wrapper.Started {
		return nil
	}

	if err := wrapper.Node.Start(ctx); err != nil {
		return fmt.Errorf("failed to start data node %s: %w", wrapper.Config.NodeID, err)
	}

	wrapper.Started = true
	tc.logger.Info("Started data node", zap.String("node_id", wrapper.Config.NodeID))
	return nil
}

// Stop stops all nodes in the cluster
func (tc *TestCluster) Stop() error {
	tc.logger.Info("Stopping test cluster")

	ctx := context.Background()

	// Stop data nodes first
	for _, wrapper := range tc.dataNodes {
		if err := tc.stopDataNode(ctx, wrapper); err != nil {
			tc.logger.Error("Failed to stop data node", zap.Error(err))
		}
	}

	// Stop coordination nodes
	for _, wrapper := range tc.coordNodes {
		if err := tc.stopCoordNode(ctx, wrapper); err != nil {
			tc.logger.Error("Failed to stop coordination node", zap.Error(err))
		}
	}

	// Stop master nodes last
	for _, wrapper := range tc.masterNodes {
		if err := tc.stopMasterNode(ctx, wrapper); err != nil {
			tc.logger.Error("Failed to stop master node", zap.Error(err))
		}
	}

	close(tc.shutdownCh)
	tc.wg.Wait()

	tc.logger.Info("Test cluster stopped")
	return nil
}

// stopMasterNode stops a single master node
func (tc *TestCluster) stopMasterNode(ctx context.Context, wrapper *MasterNodeWrapper) error {
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	if !wrapper.Started {
		return nil
	}

	if err := wrapper.Node.Stop(ctx); err != nil {
		return err
	}

	wrapper.Started = false
	return nil
}

// stopCoordNode stops a single coordination node
func (tc *TestCluster) stopCoordNode(ctx context.Context, wrapper *CoordNodeWrapper) error {
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	if !wrapper.Started {
		return nil
	}

	if err := wrapper.Node.Stop(ctx); err != nil {
		return err
	}

	wrapper.Started = false
	return nil
}

// stopDataNode stops a single data node
func (tc *TestCluster) stopDataNode(ctx context.Context, wrapper *DataNodeWrapper) error {
	wrapper.mu.Lock()
	defer wrapper.mu.Unlock()

	if !wrapper.Started {
		return nil
	}

	if err := wrapper.Node.Stop(ctx); err != nil {
		return err
	}

	wrapper.Started = false
	return nil
}

// WaitForLeader waits for a leader to be elected
func (tc *TestCluster) WaitForLeader(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {
		for _, wrapper := range tc.masterNodes {
			if wrapper.Node.IsLeader() {
				tc.logger.Info("Leader elected", zap.String("leader", wrapper.Config.NodeID))
				return nil
			}
		}
		time.Sleep(100 * time.Millisecond)
	}

	return fmt.Errorf("no leader elected within %v", timeout)
}

// GetLeader returns the current leader node
func (tc *TestCluster) GetLeader() *MasterNodeWrapper {
	for _, wrapper := range tc.masterNodes {
		if wrapper.Node.IsLeader() {
			return wrapper
		}
	}
	return nil
}

// GetMasterNode returns a master node by index
func (tc *TestCluster) GetMasterNode(index int) *MasterNodeWrapper {
	if index < 0 || index >= len(tc.masterNodes) {
		return nil
	}
	return tc.masterNodes[index]
}

// GetCoordNode returns a coordination node by index
func (tc *TestCluster) GetCoordNode(index int) *CoordNodeWrapper {
	if index < 0 || index >= len(tc.coordNodes) {
		return nil
	}
	return tc.coordNodes[index]
}

// GetDataNode returns a data node by index
func (tc *TestCluster) GetDataNode(index int) *DataNodeWrapper {
	if index < 0 || index >= len(tc.dataNodes) {
		return nil
	}
	return tc.dataNodes[index]
}

// NumMasterNodes returns the number of master nodes
func (tc *TestCluster) NumMasterNodes() int {
	return len(tc.masterNodes)
}

// NumCoordNodes returns the number of coordination nodes
func (tc *TestCluster) NumCoordNodes() int {
	return len(tc.coordNodes)
}

// NumDataNodes returns the number of data nodes
func (tc *TestCluster) NumDataNodes() int {
	return len(tc.dataNodes)
}

// Uptime returns the cluster uptime
func (tc *TestCluster) Uptime() time.Duration {
	return time.Since(tc.startTime)
}

// WaitForClusterReady waits for the cluster to be fully operational
func (tc *TestCluster) WaitForClusterReady(timeout time.Duration) error {
	// Wait for leader
	if err := tc.WaitForLeader(timeout); err != nil {
		return err
	}

	// Wait for data nodes to register with master and coordination to discover them
	// Data nodes send heartbeats immediately after startup and register
	// Coordination node discovers data nodes during its Start() method
	// We wait a bit to ensure this process completes
	expectedDataNodes := len(tc.dataNodes)
	tc.logger.Info("Waiting for data node registration and discovery",
		zap.Int("expected_nodes", expectedDataNodes))

	// Wait 2 seconds for nodes to register and be discovered
	// This is more than enough for the registration and discovery process
	time.Sleep(2 * time.Second)

	tc.logger.Info("Cluster should be ready")
	return nil
}
