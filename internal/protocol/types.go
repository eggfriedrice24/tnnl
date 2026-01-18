package protocol

import "time"

// represents a registered device in the network
type Device struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	PublicKey string    `json:"public_key"`
	Endpoints []string  `json:"endpoints"`  // IP:port combos
	VirtualIP string    `json:"virtual_ip"` // assigned mesh IP (10.100.0.x)
	LastSeen  time.Time `json:"last_seen"`
	Online    bool      `json:"online"`
}

// RegisterRequest is sent when a device joins network
type RegisterRequest struct {
	Name      string   `json:"name"`
	PublicKey string   `json:"public_key"`
	Endpoints []string `json:"endpoints"`
}

// RegisterResponse is returned after successful registration
type RegisterResponse struct {
	DeviceID  string   `json:"device_id"`
	VirtualIP string   `json:"virtual_ip"`
	Peers     []Device `json:"peers"` // current peers in the network
}

// HeartbeatRequest is sent periodically to update presence
type HeartbeatRequest struct {
	DeviceID  string   `json:"device_id"`
	Endpoints []string `json:"endpoints"` // may change due to NAT
}
