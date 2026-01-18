package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"efrc/internal/protocol"
)

type APIClient struct {
	serverURL  string
	networkKey string
	httpClient *http.Client
}

func NewAPIClient(serverURL, networkKey string) *APIClient {
	return &APIClient{
		serverURL:  serverURL,
		networkKey: networkKey,
		httpClient: &http.Client{},
	}
}

// Register sends registration request to server
func (c *APIClient) Register(req *protocol.RegisterRequest) (*protocol.RegisterResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequest("POST", c.serverURL+"/api/register", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.networkKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	var result protocol.RegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetPeers fetches all peers from server
func (c *APIClient) GetPeers() ([]protocol.Device, error) {
	httpReq, err := http.NewRequest("GET", c.serverURL+"/api/peers", nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+c.networkKey)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("server error: %s", string(body))
	}

	var peers []protocol.Device
	if err := json.NewDecoder(resp.Body).Decode(&peers); err != nil {
		return nil, err
	}

	return peers, nil
}
