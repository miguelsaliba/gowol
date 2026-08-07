# gowol

A tiny self-hostable web server and cli that sends magic wol packets to other devices on the network.

## Installation

#### Docker compose (Recommended for self hosting)

- Download the [`compose.yaml`](./compose.yaml) file in this repo.
- Modify it to your needs
- Run `docker compose up -d`
- It should now be running on http://127.0.0.1:2611/

You can now also run commands with the cli using `docker compose exec gowol gowol --help`

#### Docker

```bash
docker run --network host ghcr.io/miguelsaliba/gowol:latest
```

#### Build from source

- Make sure Go is installed (minimum version 1.26)
- Clone the repo `git clone https://github.com/miguelsaliba/gowol/`
- Build using `make build` (or `make cli` or `make server` if you want only the cli or server respectively)
- Run with `./bin/gowol` or `./bin/gowol-server`

#### Releases (Planned)

No releases yet

## Usage

#### CLI

```bash
Usage:
        gowol <mac_address> [arguments]

Arguments:
  -b, --broadcast string   Broadcast IP address to send the WOL packet to (default "255.255.255.255")
  -h, --help               Show usage information
  -p, --port uint          Port to send the WOL packet to (default 9)
  -v, --verbose            Enable verbose logging
```

#### Server

Simply run with `gowol-server` or follow the docker installation and it should be running on http://127.0.0.1:2611/

## Configuration

Server configuration is done through environment variables:

| Variable   | Description         |
|------------|---------------------|
| `GOWOL_PORT` | Port (default 2611) |

## References

- https://wiki.wireshark.org/WakeOnLAN
- https://www.amd.com/content/dam/amd/en/documents/archived-tech-docs/white-papers/20213.pdf
