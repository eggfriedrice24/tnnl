package server

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/eggfriedrice24/tnnl/internal/protocol"
)

type Store struct {
	devices map[string]*protocol.Device
	mu      sync.RWMutex
	nextIP  int
}

func NewStore() *Store {
	store := &Store{
		devices: make(map[string]*protocol.Device),
		nextIP:  1,
	}

	return store
}

func (s *Store) AddDevice(d *protocol.Device) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	s.devices[d.ID] = d

	return nil
}

func (s *Store) GetDevice(id string) (*protocol.Device, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	device, ok := s.devices[id]

	if !ok {
		return nil, errors.New("device not found")
	}

	return device, nil
}

func (s *Store) GetDeviceByPublicKey(pubKey string) (*protocol.Device, error) {
	s.mu.RLock()

	defer s.mu.RUnlock()

	for _, device := range s.devices {
		if device.PublicKey == pubKey {
			return device, nil
		}
	}

	return nil, errors.New("device not found")
}

func (s *Store) ListDevices() []*protocol.Device {
	s.mu.RLock()

	defer s.mu.RUnlock()

	devices := make([]*protocol.Device, 0, len(s.devices))

	for _, device := range s.devices {
		devices = append(devices, device)
	}

	return devices
}

func (s *Store) UpdateLastSeen(id string) error {
	s.mu.Lock()

	defer s.mu.Unlock()

	device, ok := s.devices[id]
	if !ok {
		return errors.New("device not found")
	}

	device.LastSeen = time.Now()
	return nil
}

func (s *Store) AssignVirtualIP() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	ip := fmt.Sprintf("10.100.0.%d", s.nextIP)
	s.nextIP++
	return ip
}
