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
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudformation"

	internalAmazon "github.com/banzaicloud/pipeline/internal/providers/amazon"
)

// getStackTags returns the tags that are placed onto CF template stacks.
// These tags  are propagated onto the resources created by the CF template.
func getStackTags(clusterName, stackType string) []*cloudformation.Tag {
	return append([]*cloudformation.Tag{
		{Key: aws.String("banzaicloud-pipeline-cluster-name"), Value: aws.String(clusterName)},
		{Key: aws.String("banzaicloud-pipeline-stack-type"), Value: aws.String(stackType)},
	}, internalAmazon.PipelineTags()...)
}

func generateStackNameForCluster(clusterName string) string {
	return "pipeline-eks-" + clusterName
}

// EKSActivityInput holds common input data for all activities
type EKSActivityInput struct {
	OrganizationID uint
	SecretID       string

	Region string

	ClusterName string

	// 64 chars length unique unique identifier that identifies the create CloudFormation
	AWSClientRequestToken string
}
