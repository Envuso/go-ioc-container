package container

import (
	"log"
	"reflect"

	"github.com/modern-go/reflect2"
)

// Bind - Add a binding to the container, we can do this in a few different ways...
//
// ------
//
// Function binding:
//
// For example:
//  ContainerInstance.Bind(func() AbstractServiceInterface { return NewService() })
// AbstractServiceInterface will be our resolver, NewService() is the bound type to resolve.
//
// We can also just provide a concrete binding
// For example...
//  ContainerInstance.Bind(NewService)
//
// ------
//
// Abstract Interface -> Concrete Implementation binding:
//
// For example:
//  ContainerInstance.Bind((*AbstractServiceInterface)(nil), NewService)
//
// Assuming the NewService func would return something like "*Service"
// This doesn't look great, but it's the only way I know of, to pass an interface as value
//
// If "NewService" returns AbstractServiceInterface, you can just do
//  ContainerInstance.Bind(NewService)
//
// ------
//
// Concrete binding:
//
// For example:
//  ContainerInstance.Bind(SomeConcreteService{})
// or
//  ContainerInstance.Bind(new(SomeConcreteService))
//
func (container *ContainerInstance) Bind(bindingDef ...any) bool {

	definition := reflect.TypeOf(bindingDef[0])

	// Handle Function/Concrete binding
	if len(bindingDef) == 1 {
		if definition.Kind() == reflect.Func {
			container.addFunctionBinding(definition, bindingDef[0])
			return true
		}

		container.addConcreteBinding(definition, bindingDef[0])
		return true
	}

	// Handle Abstract -> Concrete binding

	abstractType := getAbstractReturnType(definition)
	if abstractType == nil {
		log.Printf("Failed to get type of abstract: %s", definition.String())
		return false
	}

	concreteBindingType := reflect.TypeOf(bindingDef[1])
	concreteType := getConcreteReturnType(concreteBindingType)
	if concreteType == nil {
		log.Printf("Failed to get type of concrete: %s", concreteBindingType.String())
		return false
	}

	container.addBinding(abstractType, &Binding{
		bindingType:      "Abstract",
		abstractType:     abstractType,
		concreteType:     concreteType,
		resolverFunction: bindingDef[1],
		invocable:        CreateInvocable(concreteType),
	})

	return true
}

// Singleton - Bind a "class" that should only be instantiated once when resolved
// in the future, the initial instantiation of this type will be returned
func (container *ContainerInstance) Singleton(singleton any, concreteResolverFunc ...any) bool {
	singletonType := reflect.TypeOf(singleton)

	// We can provide a function to singleton
	if singletonType.Kind() == reflect.Func && concreteResolverFunc == nil {
		if singletonType.NumOut() == 0 {
			log.Printf("Please make sure your singleton function provider(%s) has at-least one return type.", singletonType.String())
			log.Printf("Without it, we don't have a type to register this singleton under.")
			log.Printf("Your type will not be registered in the container.")
			return false
		}

		container.addSingletonBinding(getConcreteReturnType(singletonType.Out(0)), &Binding{
			bindingType: "Singleton",

			resolverFunction:   singleton,
			isFunctionResolver: true,

			abstractType: singletonType,
			concreteType: singletonType,

			invocable: CreateInvocableFunction(singleton),
		})

		return true
	}

	// We can provide a type instance directly to singleton
	singletonConcrete := getConcreteReturnType(singletonType)
	if singletonConcrete == nil {
		log.Printf("Failed to get type of singleton: %s", singletonType.String())
		log.Printf("Your singleton will not be registered in the container.")
		return false
	}

	// If we don't have a resolver func, we're just defining the singleton type...
	if concreteResolverFunc == nil {
		container.addSingletonBinding(singletonConcrete, &Binding{
			bindingType: "Singleton",

			isFunctionResolver: false,

			abstractType: singletonConcrete,
			concreteType: singletonConcrete,
			invocable:    CreateInvocable(singletonConcrete),
		})
		return true
	}

	// We can provide a type instance to singleton but use
	// concreteResolverFunc to resolve the initial singleton instance

	resolverFunc := concreteResolverFunc[0]
	resolverFuncType := reflect.TypeOf(resolverFunc)

	if resolverFuncType.Kind() != reflect.Func {
		log.Printf("Trying to register singleton(%s) -> resolver binding but resolver is not a function", singletonType.String())
		log.Printf("Your singleton will not be registered in the container.")
		return false
	}

	container.addSingletonBinding(singletonConcrete, &Binding{
		bindingType: "Singleton",

		isFunctionResolver: true,
		resolverFunction:   resolverFunc,

		abstractType: singletonConcrete,
		concreteType: singletonConcrete,
		invocable:    CreateInvocableFunction(resolverFunc),
	})

	return true
}

// Instance - This is similar to Singleton, except with Singleton we provide a type to instantiate
// With instance, we provide an already instantiated value to the container
func (container *ContainerInstance) Instance(instance any) bool {
	instanceType := reflect.TypeOf(instance)

	singletonConcrete := getConcreteReturnType(instanceType)

	if singletonConcrete == nil {
		log.Printf("Failed to get type of instance singleton: %s", instanceType.String())
		log.Printf("Your instance will not be registered in the container.")
		return false
	}

	container.addSingletonBinding(singletonConcrete, &Binding{
		bindingType: "Singleton",

		isFunctionResolver: false,

		abstractType: singletonConcrete,
		concreteType: singletonConcrete,
		invocable:    CreateInvocable(singletonConcrete),
	})

	// Our instance is already instantiated, we'll pass it straight to resolved
	container.resolved[singletonConcrete] = instance

	return true
}

// IsBound - Check if the provided value type exists in our container
func (container *ContainerInstance) IsBound(binding any) bool {
	return container.getBindingType(binding) != nil
}

// Make - Try to make a new instance of the provided value and return it
// This requires a type cast to work nicely...
// For example:
//  service := ContainerInstance.Make((*ServiceAbstract)(nil))
func (container *ContainerInstance) Make(abstract any, parameters ...any) any {
	binding := container.getBindingType(abstract)

	if binding == nil {
		log.Printf("Failed to resolve binding for abstract type %s", reflect.TypeOf(abstract).String())
		return nil
	}

	return container.makeFromBinding(binding, parameters...)
}

// MakeTo - Try to make a new instance of the provided value and assign it to your arg
// For example:
//  var service ServiceAbstract
//  ContainerInstance.MakeTo(&service)
func (container *ContainerInstance) MakeTo(makeTo any, parameters ...any) {
	makeToVal := reflect.ValueOf(makeTo)

	if makeToVal.Kind() != reflect.Pointer {
		log.Printf("Call to ContainerInstance.MakeTo(), the makeTo arg must be a pointer to your receiving var. Ex; var service ServiceAbstract; ContainerInstance.MakeTo(&service)...")
		return
	}

	makeToElem := makeToVal.Elem()
	makeToType := makeToElem.Type()
	if !makeToElem.CanSet() {
		log.Printf("Call to ContainerInstance.MakeTo(), the makeTo arg cannot be set?")
		return
	}

	resolved := container.Make(makeToType, parameters...)
	if resolved == nil {
		return
	}

	resolvedValue := reflect.ValueOf(resolved)

	if makeToVal.Kind() == reflect.Ptr && resolvedValue.Kind() != reflect.Ptr {
		reflect2.TypeOf(makeTo).UnsafeSet(
			makeToVal.UnsafePointer(),
			reflect2.PtrOf(resolved),
		)
		return
	}

	// ptr := reflect.NewAt(makeToValIndirect.Type(), unsafe.Pointer(makeToValIndirect.UnsafeAddr())).Elem()
	// ptr.Set(resolvedValue.Addr())

	makeToElem.Set(resolvedValue)
}
