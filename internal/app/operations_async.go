package app

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// WaitTask retrieves an image/video task until success, failure, or cancellation.
// Callers own the deadline. No hidden polling deadline is imposed.
func (c *Client) WaitTask(ctx context.Context, id string, interval time.Duration) (json.RawMessage, error) {
	return c.waitOperation(ctx, OperationRequest{Operation: "task_get", PathParams: map[string]string{"id": id}}, interval)
}

// WaitAgent polls a beta agent result without replaying its generation request.
func (c *Client) WaitAgent(ctx context.Context, agentID, asyncID string, interval time.Duration) (json.RawMessage, error) {
	data, err := json.Marshal(map[string]string{"agent_id": agentID, "async_id": asyncID})
	if err != nil {
		return nil, err
	}
	return c.waitOperation(ctx, OperationRequest{Operation: "agent_result", JSON: data}, interval)
}

func (c *Client) waitOperation(ctx context.Context, input OperationRequest, interval time.Duration) (json.RawMessage, error) {
	if interval <= 0 {
		interval = time.Second
	}
	for {
		result, err := c.ExecuteOperation(ctx, input)
		if err != nil {
			return result, err
		}
		status, err := asyncStatus(result)
		if err != nil {
			return result, err
		}
		switch strings.ToUpper(status) {
		case "SUCCESS", "SUCCEEDED", "COMPLETED", "DONE":
			return result, nil
		case "FAIL", "FAILED", "FAILURE", "ERROR", "CANCELLED", "CANCELED":
			return result, fmt.Errorf("asynchronous task %s: %s", status, string(result))
		case "PROCESSING", "PENDING", "RUNNING", "QUEUED", "IN_PROGRESS", "SUBMITTED", "INIT":
		default:
			return result, fmt.Errorf("unknown asynchronous task status %q", status)
		}
		timer := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return result, ctx.Err()
		case <-timer.C:
		}
	}
}

func asyncStatus(result json.RawMessage) (string, error) {
	var value map[string]json.RawMessage
	if err := json.Unmarshal(result, &value); err != nil {
		return "", err
	}
	for _, key := range []string{"task_status", "status"} {
		var status string
		if json.Unmarshal(value[key], &status) == nil && status != "" {
			return status, nil
		}
	}
	for _, key := range []string{"data", "result"} {
		if child := value[key]; len(child) > 0 && string(child) != "null" {
			if status, err := asyncStatus(child); err == nil {
				return status, nil
			}
		}
	}
	return "", fmt.Errorf("asynchronous response has no task status")
}
