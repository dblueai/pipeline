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

	"github.com/gofrs/uuid"
	"go.uber.org/cadence"
	"go.uber.org/cadence/workflow"
)

const CreateInfraWorkflowName = "eks-create-infra"

// CreateInfrastructureWorkflowInput holds data needed by the create EKS cluster infrastructure workflow
type CreateInfrastructureWorkflowInput struct {
	Region         string
	OrganizationID uint
	SecretID       string

	ClusterName               string
	VpcID                     string
	RouteTableID              string
	VpcCidr                   string
	VpcCloudFormationTemplate string
}

// CreateInfrastructureWorkflow executes the Cadence workflow responsible for creating EKS
// cluster infrastructure such as VPC, subnets, EKS master nodes, worker nodes, etc
func CreateInfrastructureWorkflow(ctx workflow.Context, input CreateInfrastructureWorkflowInput) error {
	ao := workflow.ActivityOptions{
		ScheduleToStartTimeout: 5 * time.Minute,
		StartToCloseTimeout:    10 * time.Minute,
		WaitForCancellation:    true,
		RetryPolicy: &cadence.RetryPolicy{
			InitialInterval:    2 * time.Second,
			BackoffCoefficient: 1.5,
			MaximumInterval:    30 * time.Second,
			MaximumAttempts:    30,
		},
	}

	ctx = workflow.WithActivityOptions(ctx, ao)

	// create VPC activity
	{
		var activityOutput CreateVpcActivityOutput

		activityInput := &CreateVpcActivityInput{
			OrganizationID:         input.OrganizationID,
			SecretID:               input.SecretID,
			Region:                 input.Region,
			ClusterName:            input.ClusterName,
			VpcID:                  input.VpcID,
			RouteTableID:           input.RouteTableID,
			VpcCidr:                input.VpcCidr,
			CloudFormationTemplate: input.VpcCloudFormationTemplate,
			StackName:              generateStackNameForCluster(input.ClusterName),
			AWSClientRequestToken:  uuid.Must(uuid.NewV4()).String(),
		}

		if err := workflow.ExecuteActivity(ctx, CreateVpcActivityName, activityInput).Get(ctx, &activityOutput); err != nil {
			return err
		}
	}
	return nil
}
