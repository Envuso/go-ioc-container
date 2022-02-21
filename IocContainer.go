package Container

import (
	"reflect"
	"unsafe"

	"golang.org/x/exp/maps"
)

// IocContainerConfig - Holds configuration values... soon I will add some more, make them work fully
// Right now this is a placeholder
type IocContainerConfig struct {
	OnlyInjectStructFieldsWithInjectTag bool
}

type IocContainer struct {
	Config *IocContainerConfig

	// Store our singleton instances
	// instances map[reflect.Type]*IocContainerBinding

	// Our resolved singleton instances
	resolved map[reflect.Type]any

	// Store our abstract -> concrete bindings
	// If a type doesn't have an abstract type
	// We'll store concrete -> concrete
	bindings map[reflect.Type]*IocContainerBinding

	// Store aliases of Concrete -> Abstract, so we can resolve from concrete
	// when we only bound Abstract -> Concrete
	concretes map[reflect.Type]reflect.Type

	// When we register a tagged type, we'll store the tag string and then an array
	// of types for this tag, we can then use these types to resolve the bindings
	tagged map[string][]reflect.Type

	// If our container is a child container, we'll have a pointer to our parent
	parent *IocContainer
}

// CreateContainer - Create a new container instance
func CreateContainer() *IocContainer {
	c := &IocContainer{
		Config: &IocContainerConfig{OnlyInjectStructFieldsWithInjectTag: false},

		resolved:  make(map[reflect.Type]any),
		bindings:  make(map[reflect.Type]*IocContainerBinding),
		concretes: make(map[reflect.Type]reflect.Type),
		tagged:    make(map[string][]reflect.Type),
	}

	containerInstances = append(containerInstances, c.pointer())

	return c
}

var Container = CreateContainer()
var containerInstances = []unsafe.Pointer{}

// CreateChildContainer - Returns a new container, any failed look-ups of our
// child container, will then be looked up in the parent, or returned nil
func (container *IocContainer) CreateChildContainer() *IocContainer {
	c := &IocContainer{
		Config:    &IocContainerConfig{OnlyInjectStructFieldsWithInjectTag: false},
		resolved:  make(map[reflect.Type]any),
		bindings:  make(map[reflect.Type]*IocContainerBinding),
		concretes: make(map[reflect.Type]reflect.Type),
		tagged:    make(map[string][]reflect.Type),
	}

	c.parent = container

	containerInstances = append(containerInstances, c.pointer())

	return c
}

// ClearInstances - This will just remove any singleton instances from the container
// When they are next resolved via Make/MakeTo, they will be instantiated again
func (container *IocContainer) ClearInstances() {
	maps.Clear(container.resolved)
}

// Reset - Reset will empty all bindings in this container, you will have to register
// any bindings again before you can resolve them.
func (container *IocContainer) Reset() {
	maps.Clear(container.resolved)
	maps.Clear(container.bindings)
	maps.Clear(container.concretes)
}

// ParentContainer - Returns the parent container, if one exists
func (container *IocContainer) ParentContainer() *IocContainer {
	return container.parent
}
