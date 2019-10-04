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

// SystemStatusResponse System status response
type SystemStatusResponse struct {
	// A list of service objects
	ServiceStates []Service `json:"service_states,omitempty"`
}
