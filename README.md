<p align="center">
  <img src="[https://camo.githubusercontent.com/cd8f010b18b0a33d15b48fe7e5c81d0ddc94f7412509c7c99162b6e12755024a/68747470733a2f2f70726f746f636f6e662e6769746875622e696f2f70726f746f636f6e662f6173736574732f70726f746f636f6e665f6e65772e737667](https://camo.githubusercontent.com/4b6b61abd932fea50206c69287ba6485373497eefc5a416b4874060879363224/68747470733a2f2f70726f746f636f6e662e6769746875622e696f2f70726f746f636f6e662f6173736574732f70726f746f636f6e665f6e65772e737667)" width="100" alt="project-logo">
</p>
<p align="center">
    <h1 align="center"><a href="https://github.com/protoconf/protoconf">Protoconf</a> GOLANG client</h1>
</p>

<p align="center">
    <em>codify configuration, instant delivery</em>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/protoconf/client-go?style=default&logo=opensourceinitiative&logoColor=white&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/protoconf/client-go?style=default&logo=git&logoColor=white&color=0080ff" alt="last-commit">
    <img src="https://codecov.io/gh/protoconf/client-go/graph/badge.svg?token=OSiGAMlJMN)](https://codecov.io/gh/protoconf/client-go">
<p align="center">
	<!-- default option, no dependency badges. -->
</p>

<br><!-- TABLE OF CONTENTS -->

<details>
  <summary>Table of Contents</summary><br>

- [📍 Overview](#-overview)
- [🧩 Features](#-features)
- [⚙️ Installation](#️-installation)
- [🤖 Usage](#-usage)
- [🤝 Contributing](#-contributing)
- [🎗 License](#-license)
</details>
<hr>

## 📍 Overview

`protoconf_loader` is a Go package that provides a flexible and robust configuration management solution.
<br>
It supports loading configuration from files, watching for file changes, and subscribing to configuration updates from a Protoconf server.

---

## 🧩 Features

- Load configuration from files
- Watch configuration files for changes
- Subscribe to configuration updates from a Protoconf server
- Hot Reload of configuration changes
- Thread-safe access to configuration values
- Support for various data types (int32, int64, uint32, uint64, float32, float64, bool, string)
- Custom error handling and logging

## ⚙️ Installation

To install the package, use the following command:

```bash
go get github.com/protoconf/client-go
```

## 🤖 Usage

Here's a basic example of how to use the protoconf_loader:

```go
package main

import (
    "context"
    "log"

    "github.com/protoconf/client-go/"
    "github.com/protoconf/protoconf/examples/protoconf/src/crawler"
    "google.golang.org/protobuf/proto"
)

func main() {
    // Initialize the configuration with the default configuration.
    crawlerConfig := &crawler.CrawlerService{LogLevel: 3}
    // Create a new configuration instance
    cfg, err := protoconf_loader.NewConfiguration(crawlerConfig, "path/to/config")
    if err != nil {
        log.Fatalf("Failed to create configuration: %v", err)
    }

    // Load the configuration
    err = cfg.LoadConfig("path/to/config", "config_file_name")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    config.OnConfigChange(func(c proto.Message) {
		log("got new config",c)
	})
    // Start watching for changes
    err = config.WatchConfig(context.Background())
    if err != nil {
        log.Fatalf("Failed to start watching configuration: %v", err)
    }


    // Stop watching when done
    cfg.StopWatching()
}
```

## 🤝 Contributing

Contributions are welcome! Here are several ways you can contribute:

- **[Report Issues](https://github.com/protoconf/client-go/issues)**: Submit bugs found or log feature requests for the `client-go` project.
- **[Submit Pull Requests](https://github.com/protoconf/client-go/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.
- **[Join the Discussions](https://discord.protoconf.sh/)**: Share your insights, provide feedback, or ask questions.

## 🎗 License

This project is protected under the [MIT](https://choosealicense.com/licenses/mit/) License.<br>
For more details, refer to the [LICENSE](https://github.com/protoconf/client-go/blob/main/LICENSE) file.
