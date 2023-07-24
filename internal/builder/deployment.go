package builder

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

type DeploymentBuilder interface {
	SetName(string) *ResourceBuilder
	// SetNamespace(string) *ResourceBuilder
	// SetReplicas(int) *ResourceBuilder
	// SetImageRegistry(string) *ResourceBuilder
	// SetImageTag(string) *ResourceBuilder
	// SetCPULimit(resource.Quantity) *ResourceBuilder
	// SetCPURequest(resource.Quantity) *ResourceBuilder
	// SetMemoryLimit(resource.Quantity) *ResourceBuilder
	SetMemoryRequest(resource.Quantity) *ResourceBuilder
	Done() *appsv1.Deployment
}

type DeploymentBuild struct {
	Deployment *appsv1.Deployment
}

func NewDeployment() *DeploymentBuild {
	return &DeploymentBuild{}
}

func (d *DeploymentBuild) SetName(s string) *DeploymentBuild {
	return d
}

func (d *DeploymentBuild) Done() *appsv1.Deployment {

	return d.Deployment
}
func (d *DeploymentBuild) SetMemoryRequest(resource.Quantity) *DeploymentBuild {

	return d
}
