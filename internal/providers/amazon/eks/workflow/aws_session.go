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
	"emperror.dev/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/banzaicloud/pipeline/pkg/providers/amazon"
	"github.com/banzaicloud/pipeline/secret"
	"github.com/banzaicloud/pipeline/secret/verify"
)

type AWSSessionFactory struct {
}

func NewAWSSessionFactory() *AWSSessionFactory {
	return &AWSSessionFactory{}
}

// New creates a new AWS session.
func (f *AWSSessionFactory) New(organizationID uint, secretID string, region string) (*session.Session, error) {
	keyvals := []interface{}{
		"orgID", organizationID,
		"secretID", secretID,
	}

	sir, err := secret.Store.Get(organizationID, secretID)
	if err != nil {
		return nil, errors.WrapIfWithDetails(err, "failed to get AWS secret", keyvals...)
	}

	err = sir.ValidateSecretType(amazon.Provider)

	if err != nil {
		return nil, errors.WithDetails(err, keyvals...)
	}

	awsCred := verify.CreateAWSCredentials(sir.Values)

	return session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: awsCred,
	})
}
