package protoconf_loader

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
	"time"

	"github.com/protoconf/client-go/mockagent"
	"github.com/protoconf/protoconf/examples/protoconf/src/crawler"
	"google.golang.org/protobuf/proto"
)

func ExampleConfiguration() {
	// Initialize the configuration with the default configuration.
	crawlerConfig := &crawler.CrawlerService{LogLevel: 3}

	// Create a new mock agent.
	mock := mockagent.NewProtoconfAgentMock(crawlerConfig)

	// Create a new configuration with the crawler config and the mock agent.
	config, err := NewConfiguration(crawlerConfig, "crawler/config", WithAgentStub(mock.Stub))
	if err != nil {
		log.Fatalf("Error creating configuration: %v", err)
	}

	// Start the configuration watcher.
	err = config.WatchConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Mock a configuration update.
	mock.SendUpdate(&crawler.CrawlerService{LogLevel: 17})
	// Wait for the configuration to be updated.
	time.Sleep(1 * time.Millisecond)

	fmt.Printf("Config changed: %v", crawlerConfig)

	// Output: Config changed: log_level:17
}

func ExampleConfiguration_OnConfigChange() {
	// Initialize the configuration with the default configuration.
	crawlerConfig := &crawler.CrawlerService{LogLevel: 3}

	// Create a new mock agent.
	mock := mockagent.NewProtoconfAgentMock(crawlerConfig)

	// Create a new configuration with the crawler config and the mock agent.
	config, err := NewConfiguration(crawlerConfig, "crawler/config", WithAgentStub(mock.Stub))
	if err != nil {
		log.Fatalf("Error creating configuration: %v", err)
	}

	// Start the configuration watcher.
	err = config.WatchConfig(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// Create a channel to receive configuration updates.
	configUpdates := make(chan *crawler.CrawlerService)
	config.OnConfigChange(func(c proto.Message) {
		configUpdates <- c.(*crawler.CrawlerService)
	})

	// Mock a configuration update.
	mock.SendUpdate(&crawler.CrawlerService{LogLevel: 17})
	// Wait for the configuration to be updated.
	time.Sleep(1 * time.Millisecond)

	// Receive the configuration update.
	newConfig := <-configUpdates
	fmt.Printf("Config changed: %v", newConfig)

	// Output: Config changed: log_level:17
}

func ExampleConfiguration_concurrencySafety() {
	// Initialize the configuration with the default configuration.
	config := func() *atomic.Pointer[crawler.CrawlerService] {
		crawlerConfig := &crawler.CrawlerService{LogLevel: 3}
		configHolder := &atomic.Pointer[crawler.CrawlerService]{}
		configHolder.Store(crawlerConfig)

		// Create a new mock agent.
		mock := mockagent.NewProtoconfAgentMock(crawlerConfig)

		// Create a new configuration with the crawler config and the mock agent.
		config, err := NewConfiguration(crawlerConfig, "crawler/config", WithAgentStub(mock.Stub))
		if err != nil {
			log.Fatalf("Error creating configuration: %v", err)
		}

		// Start the configuration watcher.
		err = config.WatchConfig(context.Background())
		if err != nil {
			log.Fatal(err)
		}

		// Update the configuration holder when the configuration changes.
		config.OnConfigChange(func(c proto.Message) {
			configHolder.Store(c.(*crawler.CrawlerService))
		})
		mock.SendUpdate(&crawler.CrawlerService{LogLevel: 18})
		return configHolder
	}()

	// Mock a configuration update.
	// Wait for the configuration to be updated.
	time.Sleep(1 * time.Millisecond)

	// Receive the configuration update.
	fmt.Printf("Config changed: %v", config.Load())

	// Output: Config changed: log_level:18
}

func ExampleConfiguration_watchConfigOnDisk() {
	crawlerConfig := &crawler.CrawlerService{LogLevel: 3}
	tempDir := os.TempDir()
	configFilePath := filepath.Join(tempDir, "crawler_config.json")
	os.WriteFile(configFilePath, []byte("{\"logLevel\": 3}"), 0644)

	config, err := NewConfiguration(crawlerConfig, "crawler/config")
	if err != nil {
		log.Fatalf("Error creating configuration: %v", err)
	}
	config.LoadConfig(tempDir, "crawler_config.json")
	fmt.Println("Config loaded from disk:", crawlerConfig)

	config.WatchConfig(context.Background())

	os.WriteFile(configFilePath, []byte("{\"logLevel\": 4}"), 0644)
	time.Sleep(1 * time.Millisecond)

	fmt.Println("Config updated on disk:", crawlerConfig)
	// Output:
	// Config loaded from disk: log_level:3
	// Config updated on disk: log_level:4
}
