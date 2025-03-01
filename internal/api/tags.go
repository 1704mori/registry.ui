package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ListTags returns a list of tags for a specific repository
func (c *Client) ListTags(repository string) ([]string, error) {
	path := fmt.Sprintf("/%s/tags/list", repository)
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list tags: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Name string   `json:"name"`
		Tags []string `json:"tags"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Tags, nil
}

// GetTagInfo returns detailed information about a specific tag
func (c *Client) GetTagInfo(repository, tag string) (*TagInfo, error) {
	// First, get the manifest
	path := fmt.Sprintf("/%s/manifests/%s", repository, tag)
	req, err := c.newRequest("GET", path)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get manifest: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var manifest ManifestV2
	if err := json.Unmarshal(body, &manifest); err != nil {
		return nil, err
	}

	// Calculate total size
	var totalSize int64
	for _, layer := range manifest.Layers {
		totalSize += layer.Size
	}

	// Get digest from header
	digest := resp.Header.Get("Docker-Content-Digest")

	// Build tag info
	tagInfo := &TagInfo{
		Name:          tag,
		Repository:    repository,
		Digest:        digest,
		Size:          totalSize,
		ConfigDigest:  manifest.Config.Digest,
		CreatedAt:     time.Now(), // Registry API doesn't provide creation time
		LayersCount:   len(manifest.Layers),
		SchemaVersion: manifest.SchemaVersion,
	}

	return tagInfo, nil
}

// DeleteTag deletes a specific tag
func (c *Client) DeleteTag(repository, digest string) error {
	// The Docker Registry API deletes by digest, not by tag name
	path := fmt.Sprintf("/%s/manifests/%s", repository, digest)
	req, err := c.newRequest("DELETE", path)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return fmt.Errorf("failed to delete tag: status code %d", resp.StatusCode)
	}

	return nil
}
