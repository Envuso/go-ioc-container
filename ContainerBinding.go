package Container

import "reflect"

type IocContainerBinding struct {
	// Function, Concrete, Abstract
	// Function is a function we call to get a resolved value
	// Concrete is resolvable by passing a type of the concrete
	// Abstract is our Abstract -> Concrete resolver. We provide an interface, and get a service implementation.
	// Singleton
	bindingType string

	// The function to call when our bindingType is Function
	resolverFunction any

	// If we have a resolver function to call upon resolve
	// If this is false, we'll just create a new instance of concreteType
	isFunctionResolver bool

	// Set to true when we create this binding as a singleton
	isSingleton bool

	// Our abstract type, this is usually an interface
	// If we only bound a concrete implementation, this will also be our concrete
	abstractType reflect.Type
	// This is our actually resolvable concrete type
	concreteType reflect.Type
}
