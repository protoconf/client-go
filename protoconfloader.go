package protoconf_loader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	pc "github.com/protoconf/protoconf/agent/api/proto/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/fsnotify/fsnotify"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	AgentDefaultAddress = ":4300"
)

type Configuration struct {
	msg              proto.Message
	configPath       string
	logger           *slog.Logger
	isLoaded         *atomic.Bool
	isWatchingFile   *atomic.Bool
	isWatchingAgent  *atomic.Bool
	configFile       string
	fsnotifyWatcher  *fsnotify.Watcher
	mu               sync.RWMutex
	UnmarshalOptions protojson.UnmarshalOptions
	quit             chan bool
	CancelAgent      context.CancelFunc
	Host             string
	Port             int
	agentStub        pc.ProtoconfServiceClient

	onConfigChange func(p proto.Message)
}

type Option func(*Configuration)

func WithAgentStub(stub pc.ProtoconfServiceClient) Option {
	return func(c *Configuration) {
		c.agentStub = stub
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(c *Configuration) {
		c.logger = logger
	}
}

// NewConfiguration creates a new Configuration instance with the given proto.Message,
// config path and optional options.
// It initializes the fsnotify watcher, sets the unmarshal options, and initializes other fields.
// If any error occurs during the watcher creation, it returns an error.
func NewConfiguration(p proto.Message, configPath string, opts ...Option) (*Configuration, error) {
	fsnotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	config := &Configuration{
		msg:             p,
		configPath:      configPath,
		logger:          slog.Default(),
		isLoaded:        &atomic.Bool{},
		isWatchingFile:  &atomic.Bool{},
		isWatchingAgent: &atomic.Bool{},
		mu:              sync.RWMutex{},

		fsnotifyWatcher: fsnotifyWatcher,
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
		quit:           make(chan bool),
		onConfigChange: nil,
	}
	for _, opt := range opts {
		opt(config)
	}
	return config, nil
}

// LoadConfig loads the configuration from the specified configPath and configName.
// If the configuration is already loaded, it returns nil without doing anything.
// It sets the configFile field to the joined path of configPath and configName,
// then calls the loadConfig method to actually load the configuration.
// Finally, it sets the isLoaded field to true and returns nil.
// If there is an error during loading the configuration, it returns the error.
func (c *Configuration) LoadConfig(configPath string, configName string) error {
	if c.isLoaded.Load() {
		return nil
	}
	c.configFile = filepath.Join(configPath, configName)
	err := c.loadConfig()
	if err != nil {
		return err
	}
	c.isLoaded.Store(true)
	return nil
}

// WatchConfig starts watching the configuration file and the agent for changes.
// It returns an error if the configuration is not loaded yet.
// It watches for file changes and agent updates using separate goroutines.
// The method logs the successful start of watching and returns nil upon successful completion.
func (c *Configuration) WatchConfig(ctx context.Context) error {
	if !c.isLoaded.Load() {
		return errors.New("config is not loaded yet")
	}

	// Watch config file changes
	if !c.isWatchingFile.Load() {
		if err := c.watchFileChanges(); err != nil {
			return err
		}
		c.isWatchingFile.Store(true)
	}

	// Watch agent changes
	if !c.isWatchingAgent.Load() {
		if err := c.listenToChanges(c.configPath, ctx); err != nil {
			return err
		}
		c.isWatchingAgent.Store(true)
	}

	c.logger.Info(
		"Successfully watching config",
		slog.String("watching_agent_path", c.configPath),
		slog.String("watching_file", c.configFile))

	return nil
}

// watchFileChanges is a method of the Configuration struct that starts watching the configuration file for changes.
// It adds the configuration file to the fsnotifyWatcher and starts a goroutine to handle file events.
// When a write event occurs, it calls the loadConfig method to reload the configuration.
// If there is an error while watching or loading the configuration file, it logs the error.
// The method returns an error if there is an error adding the file to the fsnotifyWatcher.
func (c *Configuration) watchFileChanges() error {
	err := c.fsnotifyWatcher.Add(c.configFile)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-c.fsnotifyWatcher.Events:
				if !ok {
					c.logger.Error("error while watching config file", slog.Any("event", event))
					c.isWatchingFile.Store(false)
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					if err := c.loadConfig(); err != nil {
						c.logger.Error("error while watching and loading config file", slog.Any("error", err))
					}
				}
			case err := <-c.fsnotifyWatcher.Errors:
				c.logger.Error("error while watching config file", slog.Any("error", err))
				return
			}
		}
	}()

	return nil
}

// StopWatching stops watching the configuration file and the agent for changes.
// It closes the fsnotifyWatcher and cancels the context for the agent connection.
func (c *Configuration) StopWatching() {
	if c.isWatchingFile.Load() {
		c.fsnotifyWatcher.Close()
	}
	if c.isWatchingAgent.Load() {
		c.CancelAgent()
		c.isWatchingAgent.Store(false)
	}
}

// LoadConfig loads the configuration from the specified configPath and configName.
// If the configuration is already loaded, it returns nil.
// It sets the configFile field to the joined path of configPath and configName,
// then calls the loadConfig method to actually load the configuration.
// Finally, it sets the isLoaded field to true and returns nil.
// If there is an error during loading the configuration, it returns the error.
func (c *Configuration) loadConfig() error {
	var (
		ErrReadConfigFile  = errors.New("error reading config file")
		ErrUnmarshalConfig = errors.New("error unmarshaling config")
	)
	configReader, err := os.ReadFile(c.configFile)
	if err != nil {
		c.logger.Error("error reading config file", slog.Any("error", err))
		return ErrReadConfigFile
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	err = c.UnmarshalOptions.Unmarshal(configReader, c.msg)
	if err != nil {
		c.logger.Error("error unmarshaling config file", slog.Any("error", err))
		return ErrUnmarshalConfig
	}

	if c.onConfigChange != nil {
		c.onConfigChange(c.msg)
	}
	return nil
}

// listenToChanges is a method of the Configuration struct that establishes a connection to a server using gRPC and subscribes to receive configuration updates.
// It takes a path string and a context.Context as parameters.
// It first gets the hostname to use for the connection by calling the getHostname method.
// Then it dials the server using the obtained address and insecure transport credentials.
// If there is an error while dialing, it logs the error and returns it.
// It creates a new ProtoconfServiceClient using the connection.
// It creates a new context with cancellation capability using the provided context.
// It subscribes for configuration updates by calling the SubscribeForConfig method of the ProtoconfServiceClient.
// If there is an error while subscribing, it logs the error and returns it.
// It starts a goroutine to handle the received configuration updates by calling the handleConfigUpdates method.
// Finally, it returns nil.
func (c *Configuration) listenToChanges(path string, ctx context.Context) error {
	psc := c.agentStub
	if psc == nil {
		address := c.getHostname()

		conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			c.logger.Error("Error connecting to server ", slog.String("address", address), slog.Any("error", err))
			return err
		}

		psc = pc.NewProtoconfServiceClient(conn)
	}
	var agentCtx context.Context
	agentCtx, c.CancelAgent = context.WithCancel(ctx)
	stream, err := psc.SubscribeForConfig(agentCtx, &pc.ConfigSubscriptionRequest{Path: path})
	if err != nil {
		c.logger.Error("Error subscribing for config", slog.String("path", path), slog.Any("error", err))
		return err
	}
	go c.handleConfigUpdates(stream, path)
	return nil
}

// handleConfigUpdates listens for changes to the configuration and invokes the OnConfigChange function.
func (c *Configuration) handleConfigUpdates(stream pc.ProtoconfService_SubscribeForConfigClient, path string) {
	for {
		select {
		case <-c.quit:
			c.logger.Info("Stopping listening to changes due to quit signal")
			return // Exit the goroutine gracefully
		default:
			// Read the next update from the stream
			update, err := stream.Recv()
			if err == io.EOF {
				c.logger.Error("Connection closed while streaming config path", slog.String("path", path))
				return
			}
			if err != nil {
				c.logger.Error("Error unmarshaling config", slog.String("path", path), slog.Any("error", err))
			}

			// Unmarshal the update into the configuration
			c.mu.Lock()
			err = update.GetValue().UnmarshalTo(c.msg)
			c.mu.Unlock()
			if err != nil {
				c.logger.Error("Error while streaming config path", slog.String("path", path), slog.Any("error", err))
				// Implement appropriate error handling here
				continue // Continue to the next iteration
			}

			// Invoke the OnConfigChange function with the updated configuration
			if c.onConfigChange != nil {
				c.onConfigChange(c.msg)
			}
		}
	}
}

// OnConfigChange sets the event handler that is called when a config file changes.
func (c *Configuration) OnConfigChange(run func(p proto.Message)) {
	c.onConfigChange = run
}

// Atomic executes the given function atomically.
func (c *Configuration) WithLock(f func() error) error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return f()
}

// Get returns the result of the given function.
func Get[T any](c *Configuration, f func() T) T {
	c.mu.Lock()
	defer c.mu.Unlock()
	return f()
}

func (c *Configuration) Int32(f func() int32) int32       { return Get(c, f) }
func (c *Configuration) Int64(f func() int64) int64       { return Get(c, f) }
func (c *Configuration) UInt32(f func() uint32) uint32    { return Get(c, f) }
func (c *Configuration) UInt64(f func() uint64) uint64    { return Get(c, f) }
func (c *Configuration) Float32(f func() float32) float32 { return Get(c, f) }
func (c *Configuration) Float64(f func() float64) float64 { return Get(c, f) }
func (c *Configuration) Bool(f func() bool) bool          { return Get(c, f) }
func (c *Configuration) String(f func() string) string    { return Get(c, f) }
func (c *Configuration) Any(f func() any) any             { return Get(c, f) }

// getHostname returns the hostname to use for the agent connection.
// It uses the Host and Port fields if they are set, otherwise it uses the default address.
func (c *Configuration) getHostname() string {
	address := fmt.Sprintf("%v:%v", c.Host, c.Port)
	// Use default if not supplied
	if address == ":0" {
		address = AgentDefaultAddress
	}

	return address
}
