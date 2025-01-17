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

package anchore

import (
	"context"
	"strconv"

	"emperror.dev/errors"

	"github.com/banzaicloud/pipeline/.gen/pipeline/pipeline"
	"github.com/banzaicloud/pipeline/internal/anchore"
	"github.com/banzaicloud/pipeline/internal/common"
)

// PolicyService policy related operations
type PolicyService interface {
	ListPolicies(ctx context.Context, orgID uint, clusterID uint) (interface{}, error)
	GetPolicy(ctx context.Context, orgID uint, clusterID uint, policyID string) (interface{}, error)
	CreatePolicy(ctx context.Context, orgID uint, clusterID uint, policy pipeline.PolicyBundleRecord) (interface{}, error)
	DeletePolicy(ctx context.Context, orgID uint, clusterID uint, policyID string) error
	UpdatePolicy(ctx context.Context, orgID uint, clusterID uint, policyID string, policyActivate pipeline.PolicyBundleActivate) error
}

type policyService struct {
	configProvider anchore.ConfigProvider

	logger common.Logger
}

func NewPolicyService(configProvider anchore.ConfigProvider, logger common.Logger) PolicyService {
	return policyService{
		configProvider: configProvider,
		logger:         logger,
	}
}

func (p policyService) ListPolicies(ctx context.Context, orgID uint, clusterID uint) (interface{}, error) {
	fnCtx := map[string]interface{}{"orgID": orgID, "clusterID": clusterID}
	p.logger.Info("retrieving policies ...", fnCtx)

	anchoreClient, err := p.getAnchoreClient(ctx, orgID, clusterID)
	if err != nil {
		p.logger.Debug("failed to get anchore client", fnCtx)

		return nil, err
	}

	policyList, err := anchoreClient.ListPolicies(ctx)
	if err != nil {
		p.logger.Debug("failure while retrieving policies", fnCtx)

		return nil, errors.WrapIf(err, "failed to retrieve policies")
	}

	p.logger.Info("policies successfully retrieved", fnCtx)
	return policyList, nil
}

func (p policyService) GetPolicy(ctx context.Context, orgID uint, clusterID uint, policyID string) (interface{}, error) {
	fnCtx := map[string]interface{}{"orgID": orgID, "clusterID": clusterID, "policyID": policyID}
	p.logger.Info("retrieving policy ...", fnCtx)

	anchoreClient, err := p.getAnchoreClient(ctx, orgID, clusterID)
	if err != nil {
		p.logger.Debug("failed to get anchore client", fnCtx)

		return nil, err
	}

	policyItem, err := anchoreClient.GetPolicy(ctx, policyID)
	if err != nil {
		p.logger.Debug("failure while retrieving policy", fnCtx)

		return nil, errors.WrapIf(err, "failed to retrieve policy")
	}

	p.logger.Info("policies successfully retrieved", fnCtx)
	return policyItem, nil
}

func (p policyService) CreatePolicy(ctx context.Context, orgID uint, clusterID uint, policy pipeline.PolicyBundleRecord) (interface{}, error) {
	fnCtx := map[string]interface{}{"orgID": orgID, "clusterID": clusterID, "policy": policy}
	p.logger.Info("creating policy ...", fnCtx)

	anchoreClient, err := p.getAnchoreClient(ctx, orgID, clusterID)
	if err != nil {
		p.logger.Debug("failed to get anchore client", fnCtx)

		return nil, err
	}

	policyItem, err := anchoreClient.CreatePolicy(ctx, policy)
	if err != nil {
		p.logger.Debug("failure while creating policy", fnCtx)

		return nil, errors.WrapIf(err, "failed to create policy")
	}

	p.logger.Info("policies successfully created", fnCtx)
	return policyItem, nil
}

func (p policyService) DeletePolicy(ctx context.Context, orgID uint, clusterID uint, policyID string) error {
	fnCtx := map[string]interface{}{"orgID": orgID, "clusterID": clusterID, "policyID": policyID}
	p.logger.Info("deleting policy ...", fnCtx)

	anchoreClient, err := p.getAnchoreClient(ctx, orgID, clusterID)
	if err != nil {
		p.logger.Debug("failed to get anchore client", fnCtx)

		return errors.WrapIf(err, "failed to get anchore client")
	}

	if err := anchoreClient.DeletePolicy(ctx, policyID); err != nil {
		p.logger.Debug("failure while deleting policy", fnCtx)

		return errors.WrapIf(err, "failed to delete policy")
	}

	p.logger.Info("policy successfully deleted", fnCtx)
	return nil
}

func (p policyService) UpdatePolicy(ctx context.Context, orgID uint, clusterID uint, policyID string, policyActivate pipeline.PolicyBundleActivate) error {
	fnCtx := map[string]interface{}{"orgID": orgID, "clusterID": clusterID, "policyID": policyID}
	p.logger.Info("updating policy ...", fnCtx)

	anchoreClient, err := p.getAnchoreClient(ctx, orgID, clusterID)
	if err != nil {
		p.logger.Debug("failed to get anchore client", fnCtx)

		return errors.WrapIf(err, "failed to get anchore client")
	}

	activate, err := strconv.ParseBool(policyActivate.Params.Active)
	if err != nil {
		p.logger.Debug("failed to parse activate param", fnCtx)

		return errors.WrapIf(err, "failed to parse activate param")
	}

	if err := anchoreClient.UpdatePolicy(ctx, policyID, activate); err != nil {
		p.logger.Debug("failure while updating policy", fnCtx)

		return errors.WrapIf(err, "failed to update policy")
	}

	p.logger.Info("policy successfully updated", fnCtx)
	return nil
}

// getAnchoreClient returns p rest client wrapper instance with the proper configuration
// todo this method may be extracted to p common place to be reused by other services
func (p policyService) getAnchoreClient(ctx context.Context, orgID uint, clusterID uint) (AnchoreClient, error) {
	config, err := p.configProvider.GetConfiguration(ctx, clusterID)
	if err != nil {
		return nil, err
	}

	return NewAnchoreClient(config.User, config.Password, config.Endpoint, p.logger), nil
}
