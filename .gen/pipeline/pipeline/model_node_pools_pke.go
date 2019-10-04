/*
 * Pipeline API
 *
 * Pipeline is a feature rich application platform, built for containers on top of Kubernetes to automate the DevOps experience, continuous application development and the lifecycle of deployments.
 *
 * API version: latest
 * Contact: info@banzaicloud.com
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package pipeline

type NodePoolsPke struct {
	Name string `json:"name"`

	Roles []string `json:"roles"`

	// user provided custom node labels to be placed onto the nodes of the node pool
	Labels map[string]string `json:"labels,omitempty"`

	// Enables/disables autoscaling of this node pool through Kubernetes cluster autoscaler.
	Autoscaling bool `json:"autoscaling"`

	Provider string `json:"provider"`

	ProviderConfig map[string]interface{} `json:"providerConfig"`

	Hosts []PkeHosts `json:"hosts,omitempty"`
}
