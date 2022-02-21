package container

import (
	"log"
	"reflect"
)

// addFunctionBinding - Create a new container binding from the function
// This resolves the return type of the function as the Abstract
// And the functions return value is our Concrete
func (container *ContainerInstance) addFunctionBinding(definition reflect.Type, resolver any) {
	numOut := definition.NumOut()
	if numOut == 0 {
		log.Printf("Trying to register binding but it doesnt have a return type...")
		return
	}
	if numOut > 1 {
		log.Printf("Registering a function binding with > 1 return args. Only the first arg is handled.")
	}

	// if definition.NumIn() > 0 {
	// 	log.Printf("Function binding has args... if these args cannot be found in the container when resolving, your code will error.")
	// }

	resolverType := reflect.TypeOf(resolver)

	container.addBinding(indirectType(definition.Out(0)), &Binding{
		bindingType: "Function",

		resolverFunction:   resolver,
		isFunctionResolver: true,

		abstractType: definition,
		concreteType: resolverType,

		invocable: CreateInvocableFunction(resolver),
	})
}

// addConcreteBinding - Create a new container binding from the concrete value
// This will set our abstract type to the concrete type and the concrete type will be our concrete type..
// This just allows us to easily bind things to the container if we don't care about abstracts
func (container *ContainerInstance) addConcreteBinding(definition reflect.Type, concrete any) {
	concreteType := definition
	if definition.Kind() == reflect.Ptr {
		concreteType = definition.Elem()
	}

	// concreteWrapperFuncType := reflect.TypeOf(func() any {
	// 	return concrete
	// })
	//
	// concreteFunc := reflect.MakeFunc(concreteWrapperFuncType, func(args []reflect.Value) []reflect.Value {
	// 	return []reflect.Value{reflect.ValueOf(concrete)}
	// })

	container.addBinding(concreteType, &Binding{
		bindingType: "Concrete",

		isFunctionResolver: false,

		abstractType: concreteType,
		concreteType: definition,

		invocable: CreateInvocable(concreteType),
	})
}

// addBinding - Convenience function to add a Binding for the type &
// create a reverse lookup for Concrete -> Abstract
func (container *ContainerInstance) addBinding(abstractType reflect.Type, binding *Binding) {
	container.bindings[abstractType] = binding
	container.concretes[binding.concreteType] = abstractType
}

func (container *ContainerInstance) addSingletonBinding(singletonType reflect.Type, binding *Binding) {
	binding.isSingleton = true
	container.addBinding(singletonType, binding)
}
