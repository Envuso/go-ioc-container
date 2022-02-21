package tests

import (
	"testing"

	Container "github.com/Envuso/go-ioc-container"
)

//
// CHILD CONTAINERS
//

func TestChildBindingDoesNotAffectParent(t *testing.T) {
	container := Container.CreateContainer()

	childContainer := container.CreateChildContainer()
	childContainer.Singleton(createSingletonServiceOne)

	if !childContainer.IsBound(new(serviceConcrete)) {
		t.Fatal("Could not find binding in child container for serviceConcrete")
	}

	if container.IsBound(new(serviceConcrete)) {
		t.Fatal("Service Concrete is bound to parent container and shouldn't be")
	}

}

func TestResettingContainer(t *testing.T) {
	container := Container.CreateContainer()
	container.Singleton(createSingletonServiceOne)

	if !container.IsBound(new(serviceConcrete)) {
		t.Fatal("Could not find binding in container")
	}

	container.Reset()

	if container.IsBound(new(serviceConcrete)) {
		t.Fatal("Binding is still set in container after reset?")
	}
}

func TestClearingContainerInstances(t *testing.T) {
	container := Container.CreateContainer()
	container.Singleton(createSingletonServiceOne)

	if !container.IsBound(new(serviceConcrete)) {
		t.Fatal("Could not find binding in container")
	}

	var service *serviceConcrete
	container.MakeTo(&service)
	service.message = "before clear"

	container.ClearInstances()

	if !container.IsBound(new(serviceConcrete)) {
		t.Fatal("Could not find binding in container")
	}

	var nextService *serviceConcrete
	container.MakeTo(&nextService)

	if nextService.Message() == "before clear" {
		t.Fatal("serviceConcrete resolved after ClearInstances() call, still has value from previous singleton")
	}

}
