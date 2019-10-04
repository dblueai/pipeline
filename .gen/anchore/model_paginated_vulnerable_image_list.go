/*
 * Anchore Engine API Server
 *
 * This is the Anchore Engine API. Provides the primary external API for users of the service.
 *
 * API version: 0.1.12
 * Contact: nurmi@anchore.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package anchore

// PaginatedVulnerableImageList Pagination wrapped list of images with vulnerabilties that match some filter
type PaginatedVulnerableImageList struct {
	// The page number returned (should match the requested page query string param)
	Page string `json:"page,omitempty"`
	// True if additional pages exist (page + 1) or False if this is the last page
	NextPage string `json:"next_page,omitempty"`
	// The number of items sent in this response
	ReturnedCount int32             `json:"returned_count,omitempty"`
	Images        []VulnerableImage `json:"images,omitempty"`
}
