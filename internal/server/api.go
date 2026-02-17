package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/eggfriedrice24/tnnl/internal/protocol"

	"github.com/google/uuid"
)

type API struct {
	store      *Store
	networkKey string
}

func NewAPI(store *Store, networkKey string) *API {
	api := &API{
		store:      store,
		networkKey: networkKey,
	}

	return api
}

// middleware to check Authorization: Bearer <network-key>
func (a *API) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get header: "Authorization: Bearer mysecretkey"
		auth := r.Header.Get("Authorization")
		expected := "Bearer " + a.networkKey

		if auth != expected {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r) // call the actual handler
	}
}

// POST /api/register - Register new device
func (a *API) handleRegister(w http.ResponseWriter, r *http.Request) {
	// 1. Decode request body
	var req protocol.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// 2. Check if device already exists (by public key)
	device, err := a.store.GetDeviceByPublicKey(req.PublicKey)

	if err != nil {
		// 3. New device - create it
		device = &protocol.Device{
			ID:        uuid.New().String(),
			Name:      req.Name,
			PublicKey: req.PublicKey,
			Endpoints: req.Endpoints,
			VirtualIP: a.store.AssignVirtualIP(),
			LastSeen:  time.Now(),
			Online:    true,
		}
		a.store.AddDevice(device)
	} else {
		// 4. Existing device - update it
		device.Endpoints = req.Endpoints
		device.LastSeen = time.Now()
		device.Online = true
	}

	// 5. Build response with peer list (exclude self)
	allDevices := a.store.ListDevices()
	peers := make([]protocol.Device, 0)
	for _, d := range allDevices {
		if d.ID != device.ID {
			peers = append(peers, *d)
		}
	}

	resp := protocol.RegisterResponse{
		DeviceID:  device.ID,
		VirtualIP: device.VirtualIP,
		Peers:     peers,
	}

	// 6. Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// GET /api/peers - List all devices
func (a *API) handlePeers(w http.ResponseWriter, r *http.Request) {
	devices := a.store.ListDevices()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}

// POST /api/heartbeat - Update presence
func (a *API) handleHeartbeat(w http.ResponseWriter, r *http.Request) {
	var req protocol.HeartbeatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := a.store.UpdateLastSeen(req.DeviceID); err != nil {
		http.Error(w, "Device not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Setup routes and return handler
func (a *API) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/register", a.authMiddleware(a.handleRegister))
	mux.HandleFunc("/api/peers", a.authMiddleware(a.handlePeers))
	mux.HandleFunc("/api/heartbeat", a.authMiddleware(a.handleHeartbeat))

	return mux
}
