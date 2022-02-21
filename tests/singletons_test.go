package tests

import (
	"testing"

	Container "github.com/Envuso/go-ioc-container"
)

//
// BINDING SINGLETONS TO CONTAINER
//

func TestBindingSingletonInstanceToContainer(t *testing.T) {
	container := Container.CreateContainer()

	didBind := container.Instance(createSingletonServiceOne())
	if didBind == false {
		t.Fatalf("Failed to bind singleton `new(serviceConcrete)` to the container.")
	}

	if !container.IsBound(new(serviceConcrete)) {
		t.Fatalf("Singleton `new(serviceConcrete)` is not bound to the container.")
	}
}

func TestBindingSingletonTypeAndResolverFuncToContainer(t *testing.T) {
	container := Container.CreateContainer()

	didBind := container.Singleton(new(serviceConcrete), createSingletonServiceTwo)

	if didBind == false {
		t.Fatalf("Failed to bind singleton `new(serviceConcrete)` with resolver func to the container.")
	}

	if !container.IsBound(new(serviceConcrete)) {
		t.Fatalf("Singleton `new(serviceConcrete)` with resolver func is not bound to the container.")
	}
}

func TestBindingSingletonFuncToContainer(t *testing.T) {
	container := Container.CreateContainer()

	didBind := container.Singleton(createSingletonServiceThree)

	if didBind == false {
		t.Fatalf("Failed to bind singleton resolver func for serviceConcrete to the container.")
	}

	if !container.IsBound(new(serviceConcrete)) {
		t.Fatalf("Singleton resolver func serviceConcrete is not bound to the container.")
	}
}

//
// RESOLVING SINGLETONS FROM CONTAINER
//

func TestResolvingSingletonInstance(t *testing.T) {
	container := Container.CreateContainer()
	singleton := createSingletonServiceOne()
	container.Instance(singleton)

	var service *serviceConcrete
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve singleton serviceConcrete")
	}

	expected := singleton.message
	received := service.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}

	service.message = "Im now modified?"
	var newService *serviceConcrete
	container.MakeTo(&newService)

	expected = "Im now modified?"
	received = newService.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingSingletonTypeAndResolverFunc(t *testing.T) {
	container := Container.CreateContainer()
	singleton := createSingletonServiceTwo()
	container.Singleton(new(serviceConcrete), createSingletonServiceTwo)

	var service *serviceConcrete
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve singleton serviceConcrete")
	}

	expected := singleton.message
	received := service.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingSingletonFuncFunc(t *testing.T) {
	container := Container.CreateContainer()
	singleton := createSingletonServiceThree()
	container.Singleton(createSingletonServiceThree)

	var service *serviceConcrete
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve singleton serviceConcrete")
	}

	expected := singleton.message
	received := service.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}
