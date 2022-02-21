package tests

import (
	"testing"

	Container "github.com/Envuso/go-ioc-container"
)

//
// TAGGED BINDINGS
//

func TestAddingTaggedBindings(t *testing.T) {
	container := Container.CreateContainer()
	container.Bind(newServiceConcrete)
	container.Bind(newServiceConcreteTwo)

	if !container.IsBound(new(serviceConcrete)) || !container.IsBound(new(serviceConcreteTwo)) {
		t.Fatal("Could not find bindings in container")
	}

	didTag := container.Tag("SingletonServices", new(serviceConcrete), new(serviceConcreteTwo))
	if !didTag {
		t.Fatal("Failed to tag services")
	}

	tagged := container.Tagged("SingletonServices")

	if len(tagged) == 0 {
		t.Fatal("Failed to get tagged services")
	}
	if len(tagged) != 2 {
		t.Fatalf("We should only have two tagged services, we have %d", len(tagged))
	}

	if _, ok := tagged[0].(*serviceConcrete); !ok {
		t.Fatal("First tagged service is not serviceConcrete")
	}

	if _, ok := tagged[1].(*serviceConcreteTwo); !ok {
		t.Fatal("First tagged service is not serviceConcreteTwo")
	}
}
