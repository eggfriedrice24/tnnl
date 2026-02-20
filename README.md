# tnnl

A self-hosted mesh VPN that connects your devices into a private network, allowing them to communicate as if on the same LAN — from anywhere in the world. Think Tailscale, but self-hosted.

```
┌─────────────────────────────────────────────────────────────┐
│                    CONTROL SERVER                           │
│                   (VPS with public IP)                      │
│                                                             │
│  • Device registry (public keys, endpoints, IPs)            │
│  • Authenticates devices via network key                    │
│  • Assigns virtual IPs (10.100.0.x)                         │
│  • Coordinates peer discovery                               │
└─────────────────────────────────────────────────────────────┘
                            │
            ┌───────────────┼───────────────┐
            ▼               ▼               ▼
     ┌──────────┐    ┌──────────┐    ┌──────────┐
     │  Laptop  │◄══►│  Phone   │◄══►│  Server  │
     └──────────┘    └──────────┘    └──────────┘
            ▲               ▲               ▲
            └═══════════════╧═══════════════┘
                    WireGuard P2P tunnels
```

## Features

- **WireGuard-based** — fast, modern, and secure (ChaCha20-Poly1305)
- **Mesh topology** — devices connect directly to each other (P2P)
- **Central coordination** — lightweight server handles discovery, not traffic
- **Simple auth** — shared network key (like a WiFi password)
- **Virtual IPs** — each device gets a stable `10.100.0.x` address
- **Web dashboard** — React UI for monitoring your network

## Quick Start

### Prerequisites

- Go 1.23+
- Node.js 18+ and pnpm (for the dashboard)

### 1. Start the control server

```bash
go run ./cmd/tnnl-server --network-key "your-secret-key" --port 2424
```

### 2. Set up a device

```bash
# Initialize device (generates WireGuard keypair)
go run ./cmd/tnnl init --name "laptop"

# Connect to your server
go run ./cmd/tnnl login your-secret-key --server http://your-server:2424

# Register with the network
go run ./cmd/tnnl register
```

### 3. Check status

```bash
# View device info
go run ./cmd/tnnl status

# List all devices in the network
go run ./cmd/tnnl peers
```

Repeat step 2 on each device you want to add to your network.

### 4. Web dashboard (optional)

```bash
cd web && pnpm install && pnpm dev
```

## Building

```bash
go build -o tnnl ./cmd/tnnl
go build -o tnnl-server ./cmd/tnnl-server
go build -o tnnl-client ./cmd/tnnl-client
```

## Architecture

| Component | Binary | Purpose |
|-----------|--------|---------|
| Control Server | `tnnl-server` | Central coordination, device registry, IP assignment |
| CLI | `tnnl` | User commands (init, login, register, status, peers) |
| Daemon | `tnnl-client` | Background service maintaining the VPN (WIP) |
| Dashboard | `web/` | React UI for network monitoring |

### How it works

1. Each device generates a WireGuard keypair on `tnnl init`
2. Devices authenticate with the control server using a shared network key
3. The server assigns each device a virtual IP and tracks its public key + endpoints
4. Devices discover each other through the server and establish direct WireGuard tunnels
5. All traffic between devices is encrypted end-to-end by WireGuard

### API

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `/api/register` | POST | Register a device |
| `/api/peers` | GET | List all devices |
| `/api/heartbeat` | POST | Update device presence |

All endpoints require `Authorization: Bearer <network-key>`.

## Project Structure

```
tnnl/
├── cmd/
│   ├── tnnl/              # CLI tool
│   ├── tnnl-server/       # Control server
│   └── tnnl-client/       # Daemon
├── internal/
│   ├── protocol/          # Shared types (Device, RegisterRequest, etc.)
│   ├── server/            # HTTP handlers, in-memory store
│   ├── client/            # Config management, API client
│   ├── wg/                # WireGuard key generation
│   └── netutil/           # Network utilities (planned)
├── web/                   # React dashboard (Vite + TanStack Router + shadcn)
└── plan/                  # Design docs and task tracking
```

## Configuration

Device config is stored at `~/.tnnl/config.json`:

```json
{
  "server_url": "http://localhost:2424",
  "network_key": "your-secret-key",
  "device_id": "uuid...",
  "device_name": "laptop",
  "private_key": "base64...",
  "public_key": "base64...",
  "virtual_ip": "10.100.0.1"
}
```

## Roadmap

- [x] Control server with REST API
- [x] CLI (init, login, register, status, peers)
- [x] WireGuard key generation
- [x] Client-server communication
- [x] Web dashboard scaffolding
- [ ] Endpoint discovery (local IPs + STUN)
- [ ] WireGuard interface management
- [ ] Daemon with heartbeat and peer sync
- [ ] WebSocket real-time peer updates
- [ ] DERP relay for NAT fallback
- [ ] NAT hole punching
- [ ] Cross-platform support (Linux, macOS, Windows)

## Security

- **Network key** — shared secret to join the network
- **WireGuard keys** — per-device Curve25519 keypairs; private keys never leave the device
- **Encryption** — all mesh traffic encrypted with WireGuard (ChaCha20-Poly1305)
- **No root required** — uses wireguard-go userspace implementation

## License

MIT
