// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workflow

import (
	"time"

	"go.uber.org/cadence/workflow"

	"github.com/banzaicloud/pipeline/cluster"
	"github.com/banzaicloud/pipeline/internal/cluster/clustersetup"

	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
)

const CreateClusterWorkflowName = "eks-create-cluster"

// CreateClusterWorkflowInput holds data needed by the create cluster workflow
type CreateClusterWorkflowInput struct {
	ClusterID        uint
	ClusterName      string
	ClusterUID       string
	OrganizationID   uint
	OrganizationName string
	SecretID         string
	Region           string
	// 64 chars length unique unique identifier that consists the base of the create CloudFormation request
	// this is provided to workflow to assure idempotency in case of a workflow rerun, to submit unique id's for
	// each stack request we add the stack name to this token.
	AWSClientRequestToken string

	VpcCidr      string
	VpcID        string
	RouteTableID string

	VpcCloudFormationTemplate string

	PostHooks pkgCluster.PostHooks
}

// CreateClusterWorkflow executes the Cadence workflow responsible for creating and configuring an EKS cluster
func CreateClusterWorkflow(ctx workflow.Context, input CreateClusterWorkflowInput) error {
	cwo := workflow.ChildWorkflowOptions{
		ExecutionStartToCloseTimeout: 30 * time.Minute,
		TaskStartToCloseTimeout:      40 * time.Minute,
	}
	ctx = workflow.WithChildOptions(ctx, cwo)

	// create infra child workflow
	infraInput := CreateInfrastructureWorkflowInput{
		Region:                    input.Region,
		OrganizationID:            input.OrganizationID,
		SecretID:                  input.SecretID,
		ClusterName:               input.ClusterName,
		VpcCidr:                   input.VpcCidr,
		VpcID:                     input.VpcID,
		RouteTableID:              input.RouteTableID,
		VpcCloudFormationTemplate: input.VpcCloudFormationTemplate,
		AWSClientRequestToken:     input.AWSClientRequestToken,
	}

	err := workflow.ExecuteChildWorkflow(ctx, CreateInfraWorkflowName, infraInput).Get(ctx, nil)
	if err != nil {
		return err
	}

	// get k8s config
	var configSecretID string

	// cluster setup child workflow
	{
		workflowInput := clustersetup.WorkflowInput{
			ConfigSecretID: configSecretID,
			Cluster: clustersetup.Cluster{
				ID:   input.ClusterID,
				UID:  input.ClusterUID,
				Name: input.ClusterName,
			},
			Organization: clustersetup.Organization{
				ID:   input.OrganizationID,
				Name: input.OrganizationName,
			},
		}

		future := workflow.ExecuteChildWorkflow(ctx, clustersetup.WorkflowName, workflowInput)
		if err := future.Get(ctx, nil); err != nil {
			return err
		}
	}

	// run posthooks child workflow
	postHookWorkflowInput := cluster.RunPostHooksWorkflowInput{
		ClusterID: input.ClusterID,
		PostHooks: cluster.BuildWorkflowPostHookFunctions(input.PostHooks, true),
	}

	err = workflow.ExecuteChildWorkflow(ctx, cluster.RunPostHooksWorkflowName, postHookWorkflowInput).Get(ctx, nil)
	if err != nil {
		return err
	}

	return nil

}
