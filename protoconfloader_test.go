package protoconf_loader

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/protoconf/client-go/mockagent"
	pb "github.com/protoconf/protoconf/examples/protoconf/src/crawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

/**
 * Test_LoadConfig is a test function that tests the LoadConfig method of the Configuration struct.
 * It tests the successful loading of a config file, as well as the failure cases of file not found and unmarshal error.
 *
 * Args:
 *   - configPath (string): The path to the config file.
 *   - configName (string): The name of the config file.
 *   - serviceName (string): The name of the service.
 *
 * Returns:
 *   - None
 */

/**
 * TestConfigFileChanges is a test function that tests the behavior of the Configuration struct when the config file changes.
 * It initializes the Configuration struct, loads the config, watches the config, and verifies the changes in the config file.
 *
 * Args:
 *   - None
 *
 * Returns:
 *   - None
 */

func Test_LoadConfig(t *testing.T) {
	type args struct {
		configPath  string
		configName  string
		serviceName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test load config success",
			args: args{
				configPath:  "./test_data",
				configName:  "test_config.json",
				serviceName: "crawler/text_crawler",
			},
			wantErr: false,
		},
		{
			name: "test load config fail: file not found",
			args: args{
				configPath:  "./test_data",
				configName:  "not_exist_config.json",
				serviceName: "crawler/text_crawler",
			},
			wantErr: true,
		},
		{
			name: "test load config fail: unmarshal error",
			args: args{
				configPath:  "./test_data",
				configName:  "invalid_config.json",
				serviceName: "crawler/text_crawler",
			},
			wantErr: true,
		},
	}
	config := pb.CrawlerService{}
	handler := slog.NewJSONHandler(io.Discard, nil)
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			c, _ := NewConfiguration(&config, tt.args.serviceName, WithLogger(slog.New(handler)))
			if err := c.LoadConfig(tt.args.configPath, tt.args.configName); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

/**
 * TestConfigFileChanges is a test function that tests the behavior of the Configuration struct when the config file changes.
 * It initializes the Configuration struct with a protobuf message, a service name, and a logger.
 * Then, it loads the config file, writes a new config file, watches for changes, waits for the config to be reloaded, and stops watching.
 * Finally, it asserts that the config has been reloaded correctly.
 */
func TestConfigFileChanges(t *testing.T) {
	// Initialize the Configuration struct
	testDir := t.TempDir()

	mock := mockagent.NewProtoconfAgentMock(&pb.CrawlerService{LogLevel: 12})
	c := &pb.CrawlerService{}
	handler := slog.NewJSONHandler(io.Discard, nil)
	logger := slog.New(handler)
	config, _ := NewConfiguration(c, "test_config.json", WithLogger(logger), WithAgentStub(mock.Stub))

	// Load the config
	err := os.WriteFile(filepath.Join(testDir, "config.json"), []byte("{\"logLevel\":12}"), 0644)
	assert.NoError(t, err)
	assert.NotNil(t, config)
	err = config.LoadConfig(testDir, "config.json")

	assert.NoError(t, err)
	assert.Equal(t, int32(12), c.LogLevel)
	// Watch the config
	var x = atomic.Int32{}
	x.Store(1)
	config.OnConfigChange(func(p proto.Message) {
		x.Store(11)
	})
	assert.Equal(t, int32(1), x.Load())
	err = config.WatchConfig(context.Background())
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(testDir, "config.json"), []byte("{\"logLevel\":17}"), 0644)
	assert.NoError(t, err)

	// Wait for the config to be re loaded
	time.Sleep(1 * time.Second)
	config.WithLock(func() error {
		assert.Equal(t, int32(17), c.GetLogLevel())
		return nil
	})
	assert.Equal(t, int32(17), config.Int32(c.GetLogLevel))
	assert.Equal(t, int32(11), x.Load())
	config.StopWatching()

}

func TestConfigFileChangesWithAgent(t *testing.T) {
	c := &pb.CrawlerService{}
	agent := mockagent.NewProtoconfAgentMock(&pb.CrawlerService{LogLevel: 21})
	handler := slog.NewJSONHandler(io.Discard, nil)
	config, err := NewConfiguration(c, "test_config.json", WithLogger(slog.New(handler)), WithAgentStub(agent.Stub))
	assert.NoError(t, err)
	require.NotNil(t, config)

	testDir := t.TempDir()
	err = os.WriteFile(filepath.Join(testDir, "config.json"), []byte("{\"logLevel\":21}"), 0644)
	require.NoError(t, err)
	err = config.LoadConfig(testDir, "config.json")
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = config.WatchConfig(ctx)
	require.NoError(t, err)
	time.Sleep(1 * time.Second)
	assert.Equal(t, int32(21), c.LogLevel)
	agent.SendUpdate(&pb.CrawlerService{LogLevel: 22})
	time.Sleep(1 * time.Second)
	assert.Equal(t, int32(22), config.Int32(c.GetLogLevel))
}

func TestConfigFileChangesWithNoAgentConnection(t *testing.T) {
	c := &pb.CrawlerService{}
	handler := slog.NewJSONHandler(io.Discard, nil)
	config, err := NewConfiguration(c, "test_config.json", WithLogger(slog.New(handler)))
	config.Host = "none"
	assert.NoError(t, err)
	testDir := t.TempDir()
	err = os.WriteFile(filepath.Join(testDir, "config.json"), []byte("{\"logLevel\":21}"), 0644)
	require.NoError(t, err)
	err = config.LoadConfig(testDir, "config.json")
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = config.WatchConfig(ctx)
	require.Error(t, err)

}

func TestGetHostName(t *testing.T) {
	tests := []struct {
		name             string
		conf             Configuration
		expectedHostname string
	}{
		{
			name:             "default",
			conf:             Configuration{},
			expectedHostname: ":4300",
		},
		{
			name:             "user-provided",
			conf:             Configuration{Host: "localhost", Port: 4301},
			expectedHostname: "localhost:4301",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.conf.getHostname(), tt.expectedHostname)
		})
	}
}
