package notification

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"insider_task/internal/configs"
	"io"
	"net/http"
)

type Client struct {
	config *configs.NotificationClient
}

func NewClient(config *configs.NotificationClient) *Client {
	return &Client{
		config: config,
	}
}

func (c *Client) SendMessage(
	ctx context.Context,
	request *MessageRequest,
) (*MessageResponse, error) {
	jsonPayload, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.config.ServiceBaseURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result MessageResponse

	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
