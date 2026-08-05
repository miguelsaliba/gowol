# gowol

A basic web server and cli to that sends magic wol packets to other devices on the network.

## Installation

### Docker compose (Recommended for self hosting)

- Download the `compose.yaml` file in this repo.
- Modify it to your needs
- Run `docker compose up -d`
- It should now be running on `http://127.0.0.1:8080/`

### Docker

```bash
docker run --network host ghcr.io/miguelsaliba/gowol:latest
```

### Build from source

- Make sure Go is installed (minimum version 1.26)
- Clone the repo `git clone https://github.com/miguelsaliba/gowol/`
- Build using `make build` or `go build -o bin/ ./cmd/...`
- Run with `./bin/gowol` or `./bin/gowol-server`

### Releases (Planned)

No releases yet

## Usage

### CLI

```bash
Usage:
        gowol <mac_address> [arguments]

Arguments:
  -b, --broadcast string   Broadcast IP address to send the WOL packet to (default "255.255.255.255")
  -h, --help               Show usage information
  -p, --port uint          Port to send the WOL packet to (default 9)
  -v, --verbose            Enable verbose logging
```

### Server

Simply run with `gowol-server` or follow the docker installation and it should be running on http://127.0.0.1:8080/

## References

- https://wiki.wireshark.org/WakeOnLAN
- https://www.amd.com/content/dam/amd/en/documents/archived-tech-docs/white-papers/20213.pdf
