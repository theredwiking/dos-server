# dos-server

## Dependencies
- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Fresh](https://github.com/gravityblast/fresh)
- [Make](https://www.gnu.org/software/make/)

## Getting Started
1. Clone the repository
2. Make sure you have dos-server-sdk.json (Part of report submission)
3. run ```bash
    go mod download && go mod verify
```

## Development
1. Run the server
```bash
    GOOGLE_APPLICATION_CREDENTIALS={path to dos-server-sdk.json} fresh
```

## Deployment
### Docker
1. Build the docker image
```bash
    docker build -t dos-server .
```

### Executable
1. Build the executable
```bash
    make build
```
