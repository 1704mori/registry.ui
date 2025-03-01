package api

import (
	"fmt"
	"time"
)

// ImageInfo contains information about a Docker image repository
type ImageInfo struct {
	Name      string    `json:"name"`
	Tags      []TagInfo `json:"tags"`
	LatestTag TagInfo   `json:"latest_tag"`
}

// TagInfo contains information about a specific tag
type TagInfo struct {
	Name          string    `json:"name"`
	Repository    string    `json:"repository"`
	Digest        string    `json:"digest"`
	Size          int64     `json:"size"`
	ConfigDigest  string    `json:"config_digest"`
	CreatedAt     time.Time `json:"created_at"`
	LayersCount   int       `json:"layers_count"`
	SchemaVersion int       `json:"schema_version"`
}

// ManifestV2 represents a Docker image manifest v2
type ManifestV2 struct {
	SchemaVersion int    `json:"schemaVersion"`
	MediaType     string `json:"mediaType"`
	Config        struct {
		MediaType string `json:"mediaType"`
		Size      int64  `json:"size"`
		Digest    string `json:"digest"`
	} `json:"config"`
	Layers []struct {
		MediaType string `json:"mediaType"`
		Size      int64  `json:"size"`
		Digest    string `json:"digest"`
	} `json:"layers"`
}

// FormatSize returns a human-readable string for the size
func (t TagInfo) FormatSize() string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	size := float64(t.Size)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", size/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", size/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", size/KB)
	default:
		return fmt.Sprintf("%d B", t.Size)
	}
}
