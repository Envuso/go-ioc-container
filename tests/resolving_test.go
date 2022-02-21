package tests

import (
	"testing"

	Container "github.com/Envuso/go-ioc-container"
)

//
// RESOLVING SERVICES FROM CONTAINER
//

func TestResolvingWithTypeCast(t *testing.T) {
	container := Container.CreateContainer()
	if !container.Bind(new(serviceAbstract), &serviceConcrete{}) {
		t.Fatal("Could not bind serviceAbstract to serviceConcrete")
	}

	serviceResolved := container.Make(new(serviceAbstract))

	if serviceResolved == nil {
		t.Fatal("Could not resolve serviceConcrete via serviceAbstract")
	}

	service, ok := serviceResolved.(serviceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of serviceConcrete is not serviceAbstract")
	}
	expected := "Hello World!"
	received := service.Message()
	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingWithVar(t *testing.T) {
	container := Container.CreateContainer()
	if !container.Bind((*serviceAbstract)(nil), serviceConcrete{}) {
		t.Fatal("Could not bind serviceAbstract to serviceConcrete")
	}

	var service serviceAbstract
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve serviceConcrete via serviceAbstract")
	}

	_, ok := service.(serviceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of serviceConcrete is not serviceAbstract")
	}

	expected := "Hello World!"
	received := service.Message()
	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingWithProvidedArgs(t *testing.T) {
	container := Container.CreateContainer()
	if !container.Bind(newServiceConcreteWithMessageArg) {
		t.Fatal("Could not bind via newServiceConcreteWithMessageArg func")
	}
	expected := "A custom message"

	var service serviceAbstract
	container.MakeTo(&service, expected)

	if service == nil {
		t.Fatal("Could not resolve serviceConcrete via serviceAbstract")
	}

	_, ok := service.(serviceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of serviceConcrete is not serviceAbstract")
	}

	received := service.Message()
	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingWithSingleProvidedArgAndRestFromContainer(t *testing.T) {
	container := Container.CreateContainer()
	if !container.Bind(newServiceConcreteWithMessageArgAndService) {
		t.Fatal("Could not bind via newServiceConcreteWithMessageArg func")
	}
	if !container.Bind(newAnotherService) {
		t.Fatal("Could not bind newAnotherService")
	}
	expectedMessage := "A custom message"

	var service serviceAbstract
	container.MakeTo(&service, expectedMessage, (*anotherServiceAbstract)(nil))

	if service == nil {
		t.Fatal("Could not resolve serviceConcrete via serviceAbstract")
	}

	_, ok := service.(serviceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of serviceConcrete is not serviceAbstract")
	}

	received := service.Message()
	if received != expectedMessage {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expectedMessage)
	}
}
