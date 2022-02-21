package tests

import (
	"testing"

	Container "github.com/Envuso/go-ioc-container"
)

//
// INVOCABLE
//

func TestCallingFunctionUsingDI(t *testing.T) {
	container := Container.CreateContainer()
	container.Bind(createSingletonServiceOne)

	container.Call(func(concrete *serviceConcrete) {
		if concrete == nil {
			t.Fatal("Concrete wasnt resolved in function call")
		}
	})
}

func TestCallingFunctionOnStructUsingDI(t *testing.T) {
	// container := Container.CreateContainer()
	// container.Bind(createSingletonServiceOne)
	//
	// invocable := CreateInvocable(reflect.TypeOf(serviceConcrete{}))
	// result := invocable.CallMethodByNameWith("Message", container)
	//
	// message := result[0].Interface().(string)
	// print(message)

}

func TestCreatingStructUsingDI(t *testing.T) {
	type TestStruct struct {
		Resolved anotherServiceAbstract
	}
	container := Container.CreateContainer()

	container.Bind(newAnotherService)
	if !container.IsBound(new(anotherServiceAbstract)) {
		t.Fatal("Could not find anotherServiceAbstract binding in container")
	}

	container.Bind(new(TestStruct))
	if !container.IsBound(new(TestStruct)) {
		t.Fatal("Could not find TestStruct binding in container")
	}

	var service *TestStruct
	container.MakeTo(&service)

	// if service == nil -
	// 	t.Fatal("Could not resolve TestStruct from container using MakeTo")
	// }

	if service.Resolved == nil {
		t.Fatal("Struct was not initialized with service injection of anotherServiceAbstract")
	}

}
