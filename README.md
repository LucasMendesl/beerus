<p align="center">
    <img src="./assets/hakai.gif" align="center" width="50%">
</p>

<h1 align="center">Beerus</h1>

<p align="center">
	<em><code>❯ Clean up your Docker workspace by automatically removing unused resources based on customizable rules.</code></em>
</p>

## 🔗 Table of Contents

- [📍 Overview](#-overview)
- [👾 Features](#-features)
- [🚀 Getting Started](#-getting-started)
  - [☑️ Prerequisites](#-prerequisites)
  - [🤖 Usage](#🤖-usage)
  - [🧪 Testing](#🧪-testing)
- [🔰 Contributing](#-contributing)
- [🎗 License](#-license)
- [🙌 Acknowledgments](#-acknowledgments)

## 📍 Overview

Just as Beerus, the God of Destruction in Dragon Ball Super, maintains balance in Universe 7 by destroying what needs to be eliminated, this tool helps maintain balance in your Docker environment by efficiently cleaning up unused images and containers. Like the powerful deity who wakes up periodically to perform his duties, Beerus (the tool) periodically scans and removes Docker artifacts based on configurable policies, automatically removing:

- ✅ Unused containers based on restart policies and exit status.

- ✅ Dangling and expired images based on customizable age thresholds.

- ✅ Resources that are no longer needed but taking up space.

## 👾 Features

- 🔄 **Automatic Container Cleanup**
  - Removes exited containers based on restart policy
  - Configurable thresholds for containers with "always" restart policy
  - Monitors container exit events for immediate cleanup

- 🗑️ **Smart Image Management**
  - Removes dangling images
  - Age-based cleanup with configurable lifetime threshold
  - Handles untagged image events

- ⚡ **High Performance**
  - Concurrent processing of cleanup operations
  - Configurable concurrency levels
  - Event-driven architecture for real-time cleanup

- 🔧 **Highly Configurable**
  - YAML-based configuration
  - Environment variable support
  - CLI flags for all options
  - Adjustable cleanup thresholds
  - Flexible logging options

## 🚀 Getting Started

### ☑️ Prerequisites

Before getting started with __beerus__, ensure your runtime environment meets the following requirements:

- ✅ **Go 1.23 or higher**
- ✅ **Docker daemon running**
- ✅ **Access to Docker socket**

### 🤖 Usage

**Using `go`** &nbsp; [<img align="center" src="https://img.shields.io/badge/golang-00ADD8?&style={badge_style}&logo=go&logoColor=white" />](https://go.dev/)

```sh
❯ go run beerus.go hakai {args}
```

**Using `docker`** &nbsp; [<img align="center" src="https://img.shields.io/badge/Docker-2CA5E0.svg?style={badge_style}&logo=docker&logoColor=white" />](https://www.docker.com/)

```sh
# running using environment variables
❯ docker run \
-e BEERUS_IMAGES_LIFETIME_THRESHOLD=1 \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
lucasmendesl/beerus:0.1.0 hakai

#running using cli flags
❯ docker run \
-v /var/run/docker.sock:/var/run/docker.sock:ro \
ghcr.io/lucasmendesl/beerus:0.1.0 hakai --lifetime-threshold=100
```

#### 📦 Avaliable Registries

1. **ghcr.io/lucasmendesl/beerus (ghcr.io)**
2. **lucasmendesl/beerus:0.1.0  (dockerhub)**


#### 🛠 Configuration

This project features a **highly flexible and adaptable configuration system**, enabling users to define settings in the way that best suits their environment and workflow. Configuration can be managed through **YAML files, command-line flags, and environment variables**.

**Configuration Reference:**

| Option | Description | Default | Environment Variable | CLI Flag | YAML Path |
|--------|-------------|---------|---------------------|----------|-----------|
| Concurrency Level | Number of concurrent workers | 5 | `BEERUS_CONCURRENCY_LEVEL` | `--concurrency-level` | `beerus.concurrencyLevel` |
| Poll Check Interval | Resource check interval (hours) | 1 | `BEERUS_EXPIRING_POLL_CHECK_INTERVAL` | `--expiring-poll-check-interval` | `beerus.expiringPollCheckInterval` |
| Log Level | Logging verbosity | "info" | `BEERUS_LOG_LEVEL` | `--log-level` | `beerus.logging.level` |
| Log Format | Log output format | "text" | `BEERUS_LOG_FORMAT` | `--log-format` | `beerus.logging.format` |
| Image Lifetime | Age threshold for cleanup (days) | 100 | `BEERUS_IMAGES_LIFETIME_THRESHOLD` | `--lifetime-threshold` | `beerus.images.lifetimeThreshold` |
| Image Ignore Labels | Skip cleanup for these labels | [] | `BEERUS_IMAGES_IGNORE_LABELS` | `--image-ignore-labels` | `beerus.images.ignoreLabels` |
| Force Removal On Conflict | Allow to remove repository images that have more than one tag | false | `BEERUS_IMAGES_FORCE_REMOVAL_ON_CONFLICT` | `--force-removal-on-conflict` | `beerus.images.forceRemovalOnConflict` |
| Container Max Restarts | Max "always" policy restarts | 0 | `BEERUS_CONTAINERS_MAX_ALWAYS_RESTART_POLICY_COUNT` | `--max-always-restart-policy-count` | `beerus.containers.maxAlwaysRestartPolicyCount` |
| Container Ignore Labels | Skip cleanup for these labels | [] | `BEERUS_CONTAINERS_IGNORE_LABELS` | `--container-ignore-labels` | `beerus.containers.ignoreLabels` |
| Force Volume Cleanup | Remove associated volumes | false | `BEERUS_CONTAINERS_FORCE_VOLUME_CLEANUP` | `--force-volume-cleanup` | `beerus.containers.forceVolumeCleanup` |
| Force Link Cleanup | Remove associated links | false | `BEERUS_CONTAINERS_FORCE_LINK_CLEANUP` | `--force-link-cleanup` | `beerus.containers.forceLinkCleanup` |

**YAML Configuration File**

```yaml
version: "1.0"
beerus:
  # Number of concurrent workers for processing containers/images
  concurrencyLevel: 5

  # How often to check for expired resources (in hours)
  expiringPollCheckInterval: 1

  logging:
    # Log level: debug, info, warn, error
    level: "info"
    # Log format: json, text
    format: "text"

  images:
    # Remove images older than N days
    lifetimeThreshold: 100
    # Skip cleanup for images with these labels
    ignoreLabels:
      - "beerus.service.critical"
    # Force remove repository images that have more that one tag
    forceRemovalOnConflict: false

  containers:
    # Maximum restart count for containers with "always" policy
    # 0 means no limit
    maxAlwaysRestartPolicyCount: 5
    # Skip cleanup for containers with these labels
    ignoreLabels:
      - "beerus.service.critical"
    # Remove associated volumes on container cleanup
    forceVolumeCleanup: false
    # Remove associated links on container cleanup
    forceLinkCleanup: false
```

**Command-Line Flags**

```sh
beerus hakai \
  --concurrency-level=10 \
  --expiring-poll-check-interval=2 \
  --log-level=info \
  --log-format=json \
  --lifetime-threshold=60 \
  --image-ignore-labels="beerus.service.env.prod" \
  --container-ignore-labels="beerus.service.critical" \
  --max-always-restart-policy-count 10 \
  --force-volume-cleanup \
  --force-link-cleanup \
  --force-removal-on-conflict
```

**Environment Variables**

```sh
export BEERUS_LOG_LEVEL=debug
export BEERUS_LOG_FORMAT=text
export BEERUS_CONCURRENCY_LEVEL=10
export BEERUS_EXPIRING_POLL_CHECK_INTERVAL=24
export BEERUS_IMAGES_LIFETIME_THRESHOLD=30
export BEERUS_IMAGES_IGNORE_LABELS="beerus.env.prod,beerus.keep.image"

# you can avoid a usage of this, using the on-failure policy
export BEERUS_CONTAINERS_MAX_ALWAYS_RESTART=10
export BEERUS_CONTAINERS_IGNORE_LABELS="beerus.critical.service"
export BEERUS_CONTAINERS_FORCE_VOLUME_CLEANUP=true
```

### 🧪 Testing

This project follows Go’s standard testing framework and supports a **flexible and efficient testing workflow**. Unit tests can be executed using **Go’s built-in test command**

- **Running All Tests**
  - Execute all unit tests in the project:
    ```sh
    go test ./...
    ```

- **Running Tests with Code Coverage**
  - Generate a test coverage report:
    ```sh
    go test -cover ./...
    ```
  - Output coverage details to a file:
    ```sh
    go test -coverprofile=coverage.out ./...
    ```
  - View the coverage report in a browser:
    ```sh
    go tool cover -html=coverage.out
    ```

## 🔰 Contributing

- **💬 [Join the Discussions](https://github.com/lucasmendesl/beerus/discussions)**: Share your insights, provide feedback, or ask questions.
- **🐛 [Report Issues](https://github.com/lucasmendesl/beerus/issues)**: Submit bugs found or log feature requests for the `beerus` project.
- **💡 [Submit Pull Requests](https://github.com/lucasmendesl/beerus/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.

## 🎗 License

This project is protected under the [MIT License](https://choosealicense.com/licenses/mit/) License.

## 🙌 Acknowledgments

- Inspired by Dragon Ball Super's God of Destruction
- Docker community
- All contributors

---

<div align="center"> Made with ❤️ by Lucas Mendes Loureiro </div>
