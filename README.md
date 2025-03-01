# Docker Registry UI

Docker Registry UI is a web-based user interface for managing Docker images and repositories. It provides an easy-to-use interface for browsing, searching, and managing Docker images in a Docker registry.

## Project Structure

asd
.env
.gitignore
bun.lockb
cmd/
    server/
        main.go
go.mod
go.sum
internal/
    api/
        client.go
        images.go
        tags.go
        types.go
    config/
        config.go
    handlers/
        handlers.go
    templates/
        components/
            footer_templ.go
            footer.templ
            grid_background_templ.go
            grid_background.templ
            header_templ.go
            header.templ
            image_list_templ.go
            image_list.templ
            modal_templ.go
            ...
        layouts/
            ...
        pages/
            ...
package.json
static/
    css/
        main.css
        output.css
    js/
        tailwindcss.4.0.9.js
asd

## Prerequisites

- Go 1.16 or later
- Node.js and npm (for building static assets)
- Docker (optional, for running the Docker registry)

## Building from Source

1. Clone the repository:

asdsh
git clone https://github.com/yourusername/docker-registry-ui.git
cd docker-registry-ui
asd

2. Install Go dependencies:

asdsh
go mod download
asd

3. Install Node.js dependencies and build static assets:

asdsh
npm install
npm run build
asd

4. Build the Go server:

asdsh
go build -o bin/server cmd/server/main.go
asd

## Running the Server

1. Create a `.env` file with the following content:

asd
DOCKER_REGISTRY_URL=https://your-docker-registry.com
DOCKER_REGISTRY_USERNAME=your-username
DOCKER_REGISTRY_PASSWORD=your-password
asd

2. Run the server:

asdsh
./bin/server
asd

The server will start on `http://localhost:8080`.

## Self-Hosting

To self-host the Docker Registry UI, you can use Docker to run the server in a container.

1. Build the Docker image:

asdsh
docker build -t docker-registry-ui .
asd

2. Run the Docker container:

asdsh
docker run -d -p 8080:8080 --env-file .env docker-registry-ui
asd

The server will be accessible at `http://localhost:8080`.

## License

This project is licensed under the MIT License. See the LICENSE file for details.