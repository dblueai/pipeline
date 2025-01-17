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

package adapter

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"

	"emperror.dev/emperror"
	"emperror.dev/errors"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"

	"github.com/banzaicloud/pipeline/internal/cluster"
	"github.com/banzaicloud/pipeline/internal/common"
	sqljson "github.com/banzaicloud/pipeline/internal/database/sql/json"
	intPKE "github.com/banzaicloud/pipeline/internal/pke"
	"github.com/banzaicloud/pipeline/internal/providers/azure/pke"
	"github.com/banzaicloud/pipeline/model"
	pkgCluster "github.com/banzaicloud/pipeline/pkg/cluster"
)

const (
	GORMAzurePKEClustersTableName  = "azure_pke_clusters"
	GORMAzurePKENodePoolsTableName = "azure_pke_node_pools"
)

type gormAzurePKEClusterStore struct {
	db  *gorm.DB
	log common.Logger
}

func NewGORMAzurePKEClusterStore(db *gorm.DB, logger common.Logger) pke.AzurePKEClusterStore {
	return gormAzurePKEClusterStore{
		db:  db,
		log: logger,
	}
}

type gormAzurePKENodePoolModel struct {
	gorm.Model

	Autoscaling  bool
	ClusterID    uint `gorm:"unique_index:idx_azure_pke_np_cluster_id_name"`
	CreatedBy    uint
	DesiredCount uint
	InstanceType string
	Max          uint
	Min          uint
	Name         string `gorm:"unique_index:idx_azure_pke_np_cluster_id_name"`
	Roles        string
	SubnetName   string
	Zones        string
}

func (gormAzurePKENodePoolModel) TableName() string {
	return GORMAzurePKENodePoolsTableName
}

type gormAzurePKEClusterModel struct {
	ID                     uint `gorm:"primary_key"`
	ClusterID              uint `gorm:"unique_index:idx_azure_pke_cluster_id"`
	ResourceGroupName      string
	VirtualNetworkLocation string
	VirtualNetworkName     string

	ActiveWorkflowID  string
	KubernetesVersion string

	HTTPProxy httpProxyModel `gorm:"type:json"`

	Cluster   cluster.ClusterModel        `gorm:"foreignkey:ClusterID"`
	NodePools []gormAzurePKENodePoolModel `gorm:"foreignkey:ClusterID;association_foreignkey:ClusterID"`

	AccessPoints          accessPointsModel          `gorm:"type:json"`
	ApiServerAccessPoints apiServerAccessPointsModel `gorm:"type:json"`
}

func (gormAzurePKEClusterModel) TableName() string {
	return GORMAzurePKEClustersTableName
}

type httpProxyModel struct {
	HTTP       httpProxyOptionsModel `json:"http,omitempty"`
	HTTPS      httpProxyOptionsModel `json:"https,omitempty"`
	Exceptions []string              `json:"exceptions,omitempty"`
}

func (m *httpProxyModel) Scan(v interface{}) error {
	return sqljson.Scan(v, m)
}

func (m httpProxyModel) Value() (driver.Value, error) {
	return sqljson.Value(m)
}

func (m *httpProxyModel) fromEntity(e intPKE.HTTPProxy) {
	m.HTTP.fromEntity(e.HTTP)
	m.HTTPS.fromEntity(e.HTTPS)
	m.Exceptions = e.Exceptions
}

func (m httpProxyModel) toEntity() intPKE.HTTPProxy {
	return intPKE.HTTPProxy{
		HTTP:       m.HTTP.toEntity(),
		HTTPS:      m.HTTPS.toEntity(),
		Exceptions: m.Exceptions,
	}
}

type httpProxyOptionsModel struct {
	Host     string `json:"host"`
	Port     uint16 `json:"port,omitempty"`
	SecretID string `json:"secretId,omitempty"`
}

func (m *httpProxyOptionsModel) fromEntity(e intPKE.HTTPProxyOptions) {
	m.Host = e.Host
	m.Port = e.Port
	m.SecretID = e.SecretID
}

func (m httpProxyOptionsModel) toEntity() intPKE.HTTPProxyOptions {
	return intPKE.HTTPProxyOptions{
		Host:     m.Host,
		Port:     m.Port,
		SecretID: m.SecretID,
	}
}

type accessPointModel struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func (m *accessPointModel) fromEntity(e pke.AccessPoint) {
	m.Name = e.Name
	m.Address = e.Address
}

func (m accessPointModel) toEntity() pke.AccessPoint {
	return pke.AccessPoint{
		Name:    m.Name,
		Address: m.Address,
	}
}

type accessPointsModel []accessPointModel

func (m *accessPointsModel) Scan(v interface{}) error {
	return sqljson.Scan(v, m)
}

func (m accessPointsModel) Value() (driver.Value, error) {
	return sqljson.Value(m)
}

func (m *accessPointsModel) fromEntity(e pke.AccessPoints) {
	*m = make(accessPointsModel, len(e))
	for i, ap := range e {
		(*m)[i].fromEntity(ap)
	}
}

func (m accessPointsModel) toEntity() pke.AccessPoints {
	aps := make(pke.AccessPoints, len(m))
	for i, apm := range m {
		aps[i] = apm.toEntity()
	}
	return aps
}

type apiServerAccessPointModel string

func (m *apiServerAccessPointModel) fromEntity(e pke.APIServerAccessPoint) {
	*m = apiServerAccessPointModel(e)
}
func (m apiServerAccessPointModel) toEntity() pke.APIServerAccessPoint {
	return pke.APIServerAccessPoint(m)
}

type apiServerAccessPointsModel []apiServerAccessPointModel

func (m *apiServerAccessPointsModel) Scan(v interface{}) error {
	return sqljson.Scan(v, m)
}

func (m apiServerAccessPointsModel) Value() (driver.Value, error) {
	return sqljson.Value(m)
}

func (m *apiServerAccessPointsModel) fromEntity(e pke.APIServerAccessPoints) {
	*m = make(apiServerAccessPointsModel, len(e))
	for i, asap := range e {
		(*m)[i].fromEntity(asap)
	}
}

func (m apiServerAccessPointsModel) toEntity() pke.APIServerAccessPoints {
	asaps := make(pke.APIServerAccessPoints, len(m))
	for i, asapm := range m {
		asaps[i] = asapm.toEntity()
	}
	return asaps
}

func fillClusterFromClusterModel(cl *pke.PKEOnAzureCluster, model cluster.ClusterModel) {
	cl.CreatedBy = model.CreatedBy
	cl.CreationTime = model.CreatedAt
	cl.ID = model.ID
	cl.K8sSecretID = model.ConfigSecretID
	cl.Name = model.Name
	cl.OrganizationID = model.OrganizationID
	cl.SecretID = model.SecretID
	cl.SSHSecretID = model.SSHSecretID
	cl.Status = model.Status
	cl.StatusMessage = model.StatusMessage
	cl.UID = model.UID

	cl.ScaleOptions.DesiredCpu = model.ScaleOptions.DesiredCpu
	cl.ScaleOptions.DesiredGpu = model.ScaleOptions.DesiredGpu
	cl.ScaleOptions.DesiredMem = model.ScaleOptions.DesiredMem
	cl.ScaleOptions.Enabled = model.ScaleOptions.Enabled
	cl.ScaleOptions.Excludes = unmarshalStringSlice(model.ScaleOptions.Excludes)
	cl.ScaleOptions.KeepDesiredCapacity = model.ScaleOptions.KeepDesiredCapacity
	cl.ScaleOptions.OnDemandPct = model.ScaleOptions.OnDemandPct

	cl.Kubernetes.RBAC = model.RbacEnabled
	cl.Kubernetes.OIDC.Enabled = model.OidcEnabled
	cl.Monitoring = model.Monitoring
	cl.Logging = model.Logging
	cl.SecurityScan = model.SecurityScan
	cl.TtlMinutes = model.TtlMinutes
}

func marshalStringSlice(s []string) string {
	data, err := json.Marshal(s)
	emperror.Panic(errors.WrapIf(err, "failed to marshal string slice"))
	return string(data)
}

func unmarshalStringSlice(s string) (result []string) {
	if s == "" {
		// empty list in legacy format
		return nil
	}
	err := errors.WrapIf(json.Unmarshal([]byte(s), &result), "failed to unmarshal string slice")
	if err != nil {
		// try to parse legacy format
		result = strings.Split(s, ",")
	}
	return
}

func fillClusterFromAzurePKEClusterModel(cluster *pke.PKEOnAzureCluster, model gormAzurePKEClusterModel) error {
	fillClusterFromClusterModel(cluster, model.Cluster)

	cluster.ResourceGroup.Name = model.ResourceGroupName
	cluster.Location = model.VirtualNetworkLocation

	cluster.NodePools = make([]pke.NodePool, len(model.NodePools))
	for i, np := range model.NodePools {
		fillNodePoolFromModel(&cluster.NodePools[i], np)
	}

	cluster.VirtualNetwork.Name = model.VirtualNetworkName
	cluster.VirtualNetwork.Location = model.VirtualNetworkLocation

	cluster.Kubernetes.Version = model.KubernetesVersion
	cluster.ActiveWorkflowID = model.ActiveWorkflowID

	cluster.HTTPProxy = model.HTTPProxy.toEntity()
	cluster.AccessPoints = model.AccessPoints.toEntity()
	cluster.APIServerAccessPoints = model.ApiServerAccessPoints.toEntity()

	return nil
}

func fillNodePoolFromModel(nodePool *pke.NodePool, model gormAzurePKENodePoolModel) {
	nodePool.Autoscaling = model.Autoscaling
	nodePool.CreatedBy = model.CreatedBy
	nodePool.DesiredCount = model.DesiredCount
	nodePool.InstanceType = model.InstanceType
	nodePool.Max = model.Max
	nodePool.Min = model.Min
	nodePool.Name = model.Name
	nodePool.Roles = unmarshalStringSlice(model.Roles)
	nodePool.Subnet.Name = model.SubnetName
	nodePool.Zones = unmarshalStringSlice(model.Zones)
}

func fillModelFromNodePool(model *gormAzurePKENodePoolModel, nodePool pke.NodePool) {
	model.Autoscaling = nodePool.Autoscaling
	model.CreatedBy = nodePool.CreatedBy
	model.DesiredCount = nodePool.DesiredCount
	model.InstanceType = nodePool.InstanceType
	model.Max = nodePool.Max
	model.Min = nodePool.Min
	model.Name = nodePool.Name
	model.Roles = marshalStringSlice(nodePool.Roles)
	model.SubnetName = nodePool.Subnet.Name
	model.Zones = marshalStringSlice(nodePool.Zones)
}

func (s gormAzurePKEClusterStore) nodePools() *gorm.DB {
	return s.db.Table(GORMAzurePKENodePoolsTableName)
}

func (s gormAzurePKEClusterStore) clusterDetails() *gorm.DB {
	return s.db.Table(GORMAzurePKEClustersTableName)
}

func (s gormAzurePKEClusterStore) CreateNodePool(clusterID uint, nodePool pke.NodePool) error {
	var np gormAzurePKENodePoolModel
	fillModelFromNodePool(&np, nodePool)
	np.ClusterID = clusterID
	return getError(s.db.Create(&np), "failed to create node pool model")
}

func (s gormAzurePKEClusterStore) Create(params pke.CreateParams) (c pke.PKEOnAzureCluster, err error) {
	nodePools := make([]gormAzurePKENodePoolModel, len(params.NodePools))
	for i, np := range params.NodePools {
		fillModelFromNodePool(&nodePools[i], np)
	}

	model := gormAzurePKEClusterModel{
		Cluster: cluster.ClusterModel{
			CreatedBy:      params.CreatedBy,
			Name:           params.Name,
			Location:       params.Location,
			Cloud:          pkgCluster.Azure,
			Distribution:   pkgCluster.PKE,
			OrganizationID: params.OrganizationID,
			SecretID:       params.SecretID,
			SSHSecretID:    params.SSHSecretID,
			Status:         pkgCluster.Creating,
			StatusMessage:  pkgCluster.CreatingMessage,
			RbacEnabled:    params.RBAC,
			OidcEnabled:    params.OIDC,
			ScaleOptions: model.ScaleOptions{
				Enabled:             params.ScaleOptions.Enabled,
				DesiredCpu:          params.ScaleOptions.DesiredCpu,
				DesiredMem:          params.ScaleOptions.DesiredMem,
				DesiredGpu:          params.ScaleOptions.DesiredGpu,
				OnDemandPct:         params.ScaleOptions.OnDemandPct,
				Excludes:            marshalStringSlice(params.ScaleOptions.Excludes),
				KeepDesiredCapacity: params.ScaleOptions.KeepDesiredCapacity,
			},
		},
		KubernetesVersion:      params.KubernetesVersion,
		ResourceGroupName:      params.ResourceGroupName,
		VirtualNetworkLocation: params.Location,
		VirtualNetworkName:     params.VirtualNetworkName,
		NodePools:              nodePools,
	}

	model.HTTPProxy.fromEntity(params.HTTPProxy)
	model.AccessPoints.fromEntity(params.AccessPoints)
	model.ApiServerAccessPoints.fromEntity(params.APIServerAccessPoints)

	{
		// Adapting to legacy format. TODO: Please remove this as soon as possible.
		for _, f := range params.Features {
			switch f.Kind {
			case "InstallLogging":
				model.Cluster.Logging = true
			case "InstallMonitoring":
				model.Cluster.Monitoring = true
			case "InstallAnchoreImageValidator":
				model.Cluster.SecurityScan = true
			}
		}
	}

	if err = getError(s.db.Preload("Cluster").Preload("NodePools").Create(&model), "failed to create cluster model"); err != nil {
		return
	}

	if err := fillClusterFromAzurePKEClusterModel(&c, model); err != nil {
		return c, errors.WrapIf(err, "failed to fill cluster from model")
	}

	return
}

func (s gormAzurePKEClusterStore) DeleteNodePool(clusterID uint, nodePoolName string) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}
	if nodePoolName == "" {
		return errors.New("empty node pool name")
	}

	model := gormAzurePKENodePoolModel{
		ClusterID: clusterID,
		Name:      nodePoolName,
	}
	if err := getError(s.db.Where(model).First(&model), "failed to load model from database"); err != nil {
		return err
	}

	return getError(s.db.Delete(model), "failed to delete model from database")
}

func (s gormAzurePKEClusterStore) Delete(clusterID uint) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := cluster.ClusterModel{
		ID: clusterID,
	}
	if err := getError(s.db.Where(model).First(&model), "failed to load model from database"); err != nil {
		return err
	}

	return getError(s.db.Delete(model), "failed to soft-delete model from database")
}

func (s gormAzurePKEClusterStore) GetByID(clusterID uint) (cluster pke.PKEOnAzureCluster, _ error) {
	if err := validateClusterID(clusterID); err != nil {
		return cluster, errors.WrapIf(err, "invalid cluster ID")
	}

	model := gormAzurePKEClusterModel{
		ClusterID: clusterID,
	}
	if err := getError(s.db.Preload("Cluster").Preload("NodePools").Where(&model).First(&model), "failed to load model from database"); err != nil {
		return cluster, err
	}
	if err := fillClusterFromAzurePKEClusterModel(&cluster, model); err != nil {
		return cluster, errors.WrapIf(err, "failed to fill cluster from model")
	}
	return
}

func (s gormAzurePKEClusterStore) SetStatus(clusterID uint, status, message string) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := cluster.ClusterModel{
		ID: clusterID,
	}
	if err := getError(s.db.Where(&model).First(&model), "failed to load cluster model"); err != nil {
		return err
	}

	if status != model.Status || message != model.StatusMessage {
		fields := map[string]interface{}{
			"status":        status,
			"statusMessage": message,
		}

		statusHistory := cluster.StatusHistoryModel{
			ClusterID:   model.ID,
			ClusterName: model.Name,

			FromStatus:        model.Status,
			FromStatusMessage: model.StatusMessage,
			ToStatus:          status,
			ToStatusMessage:   message,
		}
		if err := getError(s.db.Save(&statusHistory), "failed to save status history"); err != nil {
			return err
		}

		return getError(s.db.Model(&model).Updates(fields), "failed to update cluster model")
	}

	return nil
}

func (s gormAzurePKEClusterStore) UpdateClusterAccessPoints(clusterID uint, accessPoints pke.AccessPoints) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := gormAzurePKEClusterModel{
		ClusterID: clusterID,
	}
	if err := getError(s.db.Where(&model).First(&model), "failed to load cluster model"); err != nil {
		return err
	}

	s.log.Debug("access points from db", map[string]interface{}{"accesspoints": model.AccessPoints})

	for i := range model.AccessPoints {
		for _, update := range accessPoints {
			if model.AccessPoints[i].Name == update.Name {
				model.AccessPoints[i].fromEntity(update)
			}
		}
	}

	s.log.Debug("updated access points from db", map[string]interface{}{"accesspoints": model.AccessPoints})

	updates := gormAzurePKEClusterModel{AccessPoints: model.AccessPoints}
	return getError(s.db.Model(&model).Updates(updates), "failed to update PKE-on-Azure cluster access points model")
}

func (s gormAzurePKEClusterStore) SetActiveWorkflowID(clusterID uint, workflowID string) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := gormAzurePKEClusterModel{
		ClusterID: clusterID,
	}

	return getError(s.db.Model(&model).Where("cluster_id = ?", clusterID).Update("ActiveWorkflowID", workflowID), "failed to update PKE-on-Azure cluster model")
}

func (s gormAzurePKEClusterStore) SetConfigSecretID(clusterID uint, secretID string) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := cluster.ClusterModel{
		ID: clusterID,
	}

	fields := map[string]interface{}{
		"ConfigSecretID": secretID,
	}

	return getError(s.db.Model(&model).Updates(fields), "failed to update cluster model")
}

func (s gormAzurePKEClusterStore) SetSSHSecretID(clusterID uint, secretID string) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := cluster.ClusterModel{
		ID: clusterID,
	}

	fields := map[string]interface{}{
		"SSHSecretID": secretID,
	}

	return getError(s.db.Model(&model).Updates(fields), "failed to update cluster model")
}

func (s gormAzurePKEClusterStore) GetConfigSecretID(clusterID uint) (string, error) {
	if err := validateClusterID(clusterID); err != nil {
		return "", errors.WrapIf(err, "invalid cluster ID")
	}

	model := cluster.ClusterModel{
		ID: clusterID,
	}
	if err := getError(s.db.Where(&model).First(&model), "failed to load cluster model"); err != nil {
		return "", err
	}
	return model.ConfigSecretID, nil
}

func (s gormAzurePKEClusterStore) SetFeature(clusterID uint, feature string, state bool) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := cluster.ClusterModel{
		ID: clusterID,
	}

	features := map[string]bool{
		"SecurityScan": true,
		"Logging":      true,
		"Monitoring":   true,
	}

	if !features[feature] {
		return fmt.Errorf("unknown feature: %q", feature)
	}

	fields := map[string]interface{}{
		feature: state,
	}

	return getError(s.db.Model(&model).Updates(fields), "failed to update %q feature state", feature)
}

func (s gormAzurePKEClusterStore) SetNodePoolSizes(clusterID uint, nodePoolName string, min, max, desiredCount uint, autoscaling bool) error {
	if err := validateClusterID(clusterID); err != nil {
		return errors.WrapIf(err, "invalid cluster ID")
	}

	model := gormAzurePKENodePoolModel{
		ClusterID: clusterID,
		Name:      nodePoolName,
	}

	fields := map[string]interface{}{
		"DesiredCount": desiredCount,
		"Min":          min,
		"Max":          max,
		"Autoscaling":  autoscaling,
	}

	return getError(s.db.Model(&model).Where("cluster_id = ? AND name = ?", clusterID, nodePoolName).Updates(fields), "failed to update nodepool model")
}

// Migrate executes the table migrations for the provider.
func Migrate(db *gorm.DB, logger logrus.FieldLogger) error {
	tables := []interface{}{
		&gormAzurePKENodePoolModel{},
		&gormAzurePKEClusterModel{},
	}

	var tableNames string
	for _, table := range tables {
		tableNames += fmt.Sprintf(" %s", db.NewScope(table).TableName())
	}

	logger.WithFields(logrus.Fields{
		"provider":    pke.PKEOnAzure,
		"table_names": strings.TrimSpace(tableNames),
	}).Info("migrating provider tables")

	return db.AutoMigrate(tables...).Error
}

func validateClusterID(clusterID uint) error {
	if clusterID == 0 {
		return errors.New("cluster ID cannot be 0")
	}
	return nil
}

func getError(db *gorm.DB, message string, args ...interface{}) error {
	err := db.Error
	if gorm.IsRecordNotFoundError(err) {
		err = recordNotFoundError{}
	}
	if len(args) == 0 {
		err = errors.WrapIf(err, message)
	} else {
		err = errors.WrapIff(err, message, args...)
	}
	return err
}

type recordNotFoundError struct{}

func (recordNotFoundError) Error() string {
	return "record was not found"
}

func (recordNotFoundError) NotFound() bool {
	return true
}
