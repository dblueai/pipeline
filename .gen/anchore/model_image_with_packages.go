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

// ImageWithPackages An image record that contains packages
type ImageWithPackages struct {
	Image    ImageReference     `json:"image,omitempty"`
	Packages []PackageReference `json:"packages,omitempty"`
}
