package coordination

import (
	"context"
	"sync"
	"time"

	"go.uber.org/zap"
)

// ConnectionPool manages a pool of data node client connections
type ConnectionPool struct {
	clients     map[string]*DataNodeClient // nodeID -> client
	mu          sync.RWMutex
	logger      *zap.Logger
	preConnect  bool // Whether to pre-connect during registration
	healthCheck bool // Whether to run periodic health checks
	stopChan    chan struct{}
	wg          sync.WaitGroup
}

// ConnectionPoolConfig configures the connection pool
type ConnectionPoolConfig struct {
	PreConnect          bool          // Connect eagerly during registration (default: true)
	HealthCheckInterval time.Duration // Interval for health checks (default: 30s)
	EnableHealthCheck   bool          // Enable periodic health checks (default: true)
	ReconnectOnFailure  bool          // Automatically reconnect on health check failure (default: true)
}

// DefaultConnectionPoolConfig returns default configuration
func DefaultConnectionPoolConfig() *ConnectionPoolConfig {
	return &ConnectionPoolConfig{
		PreConnect:          true,
		HealthCheckInterval: 30 * time.Second,
		EnableHealthCheck:   true,
		ReconnectOnFailure:  true,
	}
}

// NewConnectionPool creates a new connection pool
func NewConnectionPool(logger *zap.Logger, config *ConnectionPoolConfig) *ConnectionPool {
	if config == nil {
		config = DefaultConnectionPoolConfig()
	}

	pool := &ConnectionPool{
		clients:     make(map[string]*DataNodeClient),
		logger:      logger,
		preConnect:  config.PreConnect,
		healthCheck: config.EnableHealthCheck,
		stopChan:    make(chan struct{}),
	}

	// Start health check goroutine if enabled
	if config.EnableHealthCheck {
		pool.wg.Add(1)
		go pool.runHealthChecks(config.HealthCheckInterval, config.ReconnectOnFailure)
	}

	return pool
}

// RegisterClient registers a data node client and optionally pre-connects
func (cp *ConnectionPool) RegisterClient(ctx context.Context, client *DataNodeClient) error {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	nodeID := client.NodeID()

	// Check if already registered
	if _, exists := cp.clients[nodeID]; exists {
		cp.logger.Debug("Data node client already registered", zap.String("node_id", nodeID))
		return nil
	}

	// Pre-connect if enabled
	if cp.preConnect && !client.IsConnected() {
		cp.logger.Info("Pre-connecting to data node", zap.String("node_id", nodeID))

		// Use a short timeout for pre-connection to avoid blocking
		connectCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()

		if err := client.Connect(connectCtx); err != nil {
			cp.logger.Warn("Failed to pre-connect to data node, will retry later",
				zap.String("node_id", nodeID),
				zap.Error(err))
			// Don't fail registration - connection will be established on first use
		} else {
			cp.logger.Info("Successfully pre-connected to data node", zap.String("node_id", nodeID))
		}
	}

	cp.clients[nodeID] = client
	cp.logger.Info("Registered data node client in connection pool",
		zap.String("node_id", nodeID),
		zap.Bool("pre_connected", client.IsConnected()))

	return nil
}

// UnregisterClient removes a client from the pool
func (cp *ConnectionPool) UnregisterClient(nodeID string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	if client, exists := cp.clients[nodeID]; exists {
		// Disconnect gracefully
		if client.IsConnected() {
			if err := client.Disconnect(); err != nil {
				cp.logger.Warn("Error disconnecting client during unregistration",
					zap.String("node_id", nodeID),
					zap.Error(err))
			}
		}
		delete(cp.clients, nodeID)
		cp.logger.Info("Unregistered data node client from connection pool", zap.String("node_id", nodeID))
	}
}

// GetClient retrieves a client from the pool
func (cp *ConnectionPool) GetClient(nodeID string) (*DataNodeClient, bool) {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	client, exists := cp.clients[nodeID]
	return client, exists
}

// GetAllClients returns all registered clients
func (cp *ConnectionPool) GetAllClients() map[string]*DataNodeClient {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	// Return copy to avoid concurrent map access
	clients := make(map[string]*DataNodeClient, len(cp.clients))
	for nodeID, client := range cp.clients {
		clients[nodeID] = client
	}
	return clients
}

// Close closes all connections and stops health checks
func (cp *ConnectionPool) Close() error {
	// Stop health check goroutine
	if cp.healthCheck {
		close(cp.stopChan)
		cp.wg.Wait()
	}

	// Disconnect all clients
	cp.mu.Lock()
	defer cp.mu.Unlock()

	for nodeID, client := range cp.clients {
		if client.IsConnected() {
			if err := client.Disconnect(); err != nil {
				cp.logger.Warn("Error disconnecting client during pool close",
					zap.String("node_id", nodeID),
					zap.Error(err))
			}
		}
	}

	cp.logger.Info("Connection pool closed")
	return nil
}

// runHealthChecks periodically checks connection health and reconnects if needed
func (cp *ConnectionPool) runHealthChecks(interval time.Duration, reconnect bool) {
	defer cp.wg.Done()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	cp.logger.Info("Started connection pool health checks",
		zap.Duration("interval", interval),
		zap.Bool("auto_reconnect", reconnect))

	for {
		select {
		case <-ticker.C:
			cp.performHealthCheck(reconnect)
		case <-cp.stopChan:
			cp.logger.Info("Stopping connection pool health checks")
			return
		}
	}
}

// performHealthCheck checks all connections and reconnects if needed
func (cp *ConnectionPool) performHealthCheck(reconnect bool) {
	cp.mu.RLock()
	clients := make([]*DataNodeClient, 0, len(cp.clients))
	for _, client := range cp.clients {
		clients = append(clients, client)
	}
	cp.mu.RUnlock()

	for _, client := range clients {
		nodeID := client.NodeID()

		if !client.IsConnected() {
			cp.logger.Warn("Data node connection is not connected",
				zap.String("node_id", nodeID))

			if reconnect {
				cp.logger.Info("Attempting to reconnect to data node", zap.String("node_id", nodeID))

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				if err := client.Connect(ctx); err != nil {
					cp.logger.Error("Failed to reconnect to data node",
						zap.String("node_id", nodeID),
						zap.Error(err))
				} else {
					cp.logger.Info("Successfully reconnected to data node", zap.String("node_id", nodeID))
				}
				cancel()
			}
		}
	}
}

// Stats returns connection pool statistics
func (cp *ConnectionPool) Stats() map[string]interface{} {
	cp.mu.RLock()
	defer cp.mu.RUnlock()

	connected := 0
	disconnected := 0

	for _, client := range cp.clients {
		if client.IsConnected() {
			connected++
		} else {
			disconnected++
		}
	}

	return map[string]interface{}{
		"total_clients":        len(cp.clients),
		"connected":            connected,
		"disconnected":         disconnected,
		"health_check_enabled": cp.healthCheck,
		"pre_connect_enabled":  cp.preConnect,
	}
}
