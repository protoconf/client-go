package protoconf_loader

import (
	"context"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	protoconfagent "github.com/protoconf/protoconf/agent/api/proto/v1"
	pb "github.com/protoconf/protoconf/examples/protoconf/src/crawler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
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

	mock := newProtoconfAgentMock(&pb.CrawlerService{LogLevel: 12})
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
	agent := newProtoconfAgentMock(&pb.CrawlerService{LogLevel: 21})
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

type mockResponse struct {
	wait     time.Duration
	response *protoconfagent.ConfigUpdate
}

type protoconfAgentMock struct {
	protoconfagent.UnimplementedProtoconfServiceServer
	responses []*mockResponse
	Chan      chan *protoconfagent.ConfigUpdate
	Stub      protoconfagent.ProtoconfServiceClient
	initial   proto.Message
}

func newProtoconfAgentMock(initial proto.Message, responses ...*mockResponse) *protoconfAgentMock {
	buffer := 1024 * 1024
	lis := bufconn.Listen(buffer)
	ctx := context.Background()
	rpcServer := grpc.NewServer()
	srv := &protoconfAgentMock{responses: responses, Chan: make(chan *protoconfagent.ConfigUpdate), initial: initial}
	protoconfagent.RegisterProtoconfServiceServer(rpcServer, srv)
	go func() {
		context.AfterFunc(ctx, func() {
			rpcServer.GracefulStop()
		})
		if err := rpcServer.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.DialContext(ctx, "",
		grpc.WithContextDialer(func(context.Context, string) (net.Conn, error) {
			return lis.Dial()
		}), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}
	srv.Stub = protoconfagent.NewProtoconfServiceClient(conn)
	return srv
}

func (p protoconfAgentMock) SubscribeForConfig(request *protoconfagent.ConfigSubscriptionRequest, srv protoconfagent.ProtoconfService_SubscribeForConfigServer) error {
	go func() {
		for _, r := range p.responses {
			time.Sleep(r.wait)
			srv.Send(r.response)
		}
	}()
	for {
		select {
		case <-srv.Context().Done():
			return srv.Context().Err()
		case update := <-p.Chan:
			srv.Send(update)
		}
	}
}

func (p protoconfAgentMock) SendUpdate(msg proto.Message) {
	v, _ := anypb.New(msg)
	update := &protoconfagent.ConfigUpdate{
		Value: v,
	}
	go func() {
		p.Chan <- update
	}()
}
