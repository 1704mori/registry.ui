package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	var ociIndex OCIImageIndex
	var totalSize int64
	var layersCount int

	// Determine if the response is an OCI image index
	if err := json.Unmarshal(body, &ociIndex); err == nil && ociIndex.SchemaVersion == 2 && len(ociIndex.Manifests) > 0 {
		for _, manifest := range ociIndex.Manifests {
			// Fetch each manifest to calculate the total size
			manifestBody, err := c.fetchManifest(repository, manifest.Digest)
			if err != nil {
				return nil, err
			}

			var subManifest ManifestV2
			if err := json.Unmarshal(manifestBody, &subManifest); err != nil {
				return nil, err
			}

			for _, layer := range subManifest.Layers {
				totalSize += layer.Size
			}
		}
		layersCount = len(ociIndex.Manifests)
	} else {
		// Otherwise, assume it's a Docker image manifest v2
		if err := json.Unmarshal(body, &manifest); err != nil {
			return nil, err
		}

		// Calculate total size
		for _, layer := range manifest.Layers {
			totalSize += layer.Size
		}
		layersCount = len(manifest.Layers)
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
		LayersCount:   layersCount,
		SchemaVersion: manifest.SchemaVersion,
	}

	return tagInfo, nil
}

// fetchManifest fetches the manifest for a given digest
func (c *Client) fetchManifest(repository, digest string) ([]byte, error) {
	path := fmt.Sprintf("/%s/manifests/%s", repository, digest)
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

	return io.ReadAll(resp.Body)
}

// DeleteTag deletes a specific tag from the registry
func (c *Client) DeleteTag(name, tag string) error {
	path := fmt.Sprintf("/%s/tags/%s", name, tag)
	log.Printf("Deleting tag: %s", path)
	req, err := c.newRequest("DELETE", path)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusAccepted {
		return nil
	} else if resp.StatusCode == http.StatusMethodNotAllowed {
		return fmt.Errorf("tag deletion not supported by this registry")
	}
	return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}
