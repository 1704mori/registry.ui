package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ListRepositories returns a list of repositories in the registry
func (c *Client) ListRepositories() ([]string, error) {
	req, err := c.newRequest("GET", "/_catalog")
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repositories: status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Repositories []string `json:"repositories"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	return result.Repositories, nil
}

// GetImageInfo returns details about a specific image repository
func (c *Client) GetImageInfo(name string) (*ImageInfo, error) {
	// First, get the tags for this repository
	tags, err := c.ListTags(name)
	if err != nil {
		return nil, err
	}

	// For each tag, get the manifest
	info := &ImageInfo{
		Name: name,
		Tags: make([]TagInfo, 0, len(tags)),
	}

	for _, tag := range tags {
		tagInfo, err := c.GetTagInfo(name, tag)
		if err != nil {
			// Continue to the next tag if there's an error
			continue
		}
		info.Tags = append(info.Tags, *tagInfo)
	}

	// Set latest tag and total size
	for _, tag := range info.Tags {
		if tag.Name == "latest" {
			info.LatestTag = tag
			break
		}
	}

	if info.LatestTag.Name == "" && len(info.Tags) > 0 {
		info.LatestTag = info.Tags[0]
	}

	return info, nil
}
