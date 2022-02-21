package container

import (
	"log"
	"reflect"
	"unsafe"
)

// hasBinding - Look up a Type in the container and return whether it exists
func (container *ContainerInstance) hasBinding(binding reflect.Type) bool {
	_, ok := container.bindings[binding]

	return ok
}

// getBindingType - Try to get a binding type from the binding arg in a few different ways
// We'll first assume we're checking for an abstract type binding...
// If we didn't get it from the abstract, we'll then check for the concrete...
// Now as a last ditch effort, we'll look the bindingType up in container.concretes
// If we can't find anything, we return nil
func (container *ContainerInstance) getBindingType(binding any) reflect.Type {
	bindingType := getType(binding)

	// First, we'll check if we have this type as a singleton binding
	testType := getConcreteReturnType(bindingType)
	if container.hasBinding(testType) {
		return testType
	}

	// We'll first assume we're checking for an abstract type binding...
	testType = getAbstractReturnType(bindingType)
	if container.hasBinding(testType) {
		return testType
	}

	// If we didn't get it from the abstract, we'll then check for the concrete...
	testType = getConcreteReturnType(bindingType)
	if container.hasBinding(testType) {
		return testType
	}

	// Now as a last ditch effort, we'll look the bindingType up in container.concretes
	if potentialAbstract, ok := container.concretes[bindingType]; ok {
		if container.hasBinding(potentialAbstract) {
			return potentialAbstract
		}
	}

	return nil
}

// makeFromBinding - Once we've obtained our binding type from
// Make, we'll then check the containers bindings
// If it doesn't exist, and we have a parent container we'll then call makeFromBinding on the
// parent container. Which will either recurse until a resolve is made, or return nil
func (container *ContainerInstance) makeFromBinding(binding reflect.Type, parameters ...any) any {
	containerBinding, ok := container.bindings[binding]
	if !ok {
		if container.parent != nil {
			return container.parent.makeFromBinding(binding, parameters...)
		}
		log.Printf("Failed to resolve container binding for abstract type %s", binding.String())
		return nil
	}

	return container.resolve(containerBinding, parameters...)
}

func (container *ContainerInstance) pointer() unsafe.Pointer {
	return reflect.ValueOf(container).UnsafePointer()
}

func removeAllChildContainerInstances() {
	// mainContainerPointer := Container.pointer()
	//
	// newContainerInstances := []unsafe.Pointer{}
	//
	// for _, containerPtr := range containerInstances {
	// 	var container = *((*ContainerInstance)(containerPtr))
	//
	// 	if container.parent == nil && container.pointer() != mainContainerPointer {
	//
	// 	}
	// }
}
