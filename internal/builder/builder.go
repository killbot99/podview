package builder

type ResourceBuilder interface {
	SetName() *ResourceBuilder
	SetNamespace() *ResourceBuilder
}
