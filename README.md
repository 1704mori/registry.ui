# Registry.UI

Registry.UI is a web-based user interface for managing Docker images and repositories. It provides an easy-to-use interface for browsing, searching, and managing Docker images in a Docker registry.

## Prerequisites

- Go 1.24 or later
- Docker (optional, for running the Docker registry)

## Building from Source

1. Clone the repository:

```sh
git clone https://github.com/1704mori/registry.ui.git
cd registry-ui
```

2. Install Go dependencies:

```sh
go mod download
make install-tools
```

3. Run the Application
```sh
air .
```

## Run the Application

1. Run as Service

```sh
curl -fsSL https://raw.githubusercontent.com/1704mori/registry.ui/refs/heads/master/install.sh -o install.sh
chmod +x install.sh
sudo ./install.sh
```

2. Run as Docker

```sh
curl -fsSL https://raw.githubusercontent.com/1704mori/registry.ui/refs/heads/master/install.sh -o install.sh
chmod +x install.sh
sudo ./install.sh docker
```

The server will start on `http://localhost:8080`.

## License

This project is licensed under the GPL-3.0 License - see the [LICENSE](LICENSE) file for details.