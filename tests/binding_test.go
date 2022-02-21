package tests

import (
	"testing"

	Container "github.com/Envuso/go-ioc-container"
)

//
// BINDING SERVICES TO CONTAINER
//

func TestBindingAbstractInterfaceToConcreteImplementation(t *testing.T) {
	container := Container.CreateContainer()
	if !container.Bind((*serviceAbstract)(nil), serviceConcrete{}) {
		t.Fatal("Could not bind serviceAbstract to serviceConcrete")
	}

	if !container.IsBound((*serviceAbstract)(nil)) {
		t.Fatal("Could not find binding for serviceAbstract")
	}
	if !container.IsBound(serviceConcrete{}) {
		t.Fatal("Could not find binding for serviceConcrete")
	}
}

func TestBindingAbstractInterfaceToConcreteImplementationUsingPtr(t *testing.T) {
	container := Container.CreateContainer()
	if !container.Bind((*serviceAbstract)(nil), &serviceConcrete{}) {
		t.Fatal("Could not bind serviceAbstract to serviceConcrete")
	}

	if !container.IsBound((*serviceAbstract)(nil)) {
		t.Fatal("Could not find binding for serviceAbstract")
	}
	if !container.IsBound(serviceConcrete{}) {
		t.Fatal("Could not find binding for serviceConcrete")
	}
}

func TestBindingAbstractToConcreteViaFunction(t *testing.T) {
	container := Container.CreateContainer()
	bound := container.Bind(func() serviceAbstract {
		return &serviceConcrete{}
	})

	if !bound {
		t.Fatal("Could not bind serviceAbstract to serviceConcrete via function")
	}
	if !container.IsBound((*serviceAbstract)(nil)) {
		t.Fatal("Could not find binding for serviceAbstract")
	}
}

func TestBindingAbstractToConcreteViaFunctionWithAbstractToFunctionArgs(t *testing.T) {
	container := Container.CreateContainer()
	bound := container.Bind((*serviceAbstract)(nil), func() *serviceConcrete {
		return &serviceConcrete{}
	})

	if !bound {
		t.Fatal("Could not bind serviceAbstract to serviceConcrete via function")
	}
	if !container.IsBound((*serviceAbstract)(nil)) {
		t.Fatal("Could not find binding for serviceAbstract")
	}
	if !container.IsBound(serviceConcrete{}) {
		t.Fatal("Could not find binding for serviceConcrete")
	}
}
