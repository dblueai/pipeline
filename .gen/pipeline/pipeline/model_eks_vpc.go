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

type EksVpc struct {

	// The identifier of existing VPC to be used for creating the EKS cluster. If not provided a new VPC is created for the cluster.
	VpcId string `json:"vpcId,omitempty"`

	// The CIDR range for the VPC in case new VPC is created.
	Cidr string `json:"cidr,omitempty"`
}
