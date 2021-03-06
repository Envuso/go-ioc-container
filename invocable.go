package container

import (
	"log"
	"reflect"
)

// CreateInvocable - Binding should be a struct or function
func CreateInvocable(bindingType reflect.Type) *Invocable {
	isInvoc, invocableType := isInvocable(bindingType)
	if !isInvoc {
		log.Printf("type passed to CreateInvocable (%s) is not an invocable type(function or struct)", bindingType.String())
		return nil
	}

	return &Invocable{
		bindingType:   bindingType,
		typeOfBinding: invocableType,
	}
}

// CreateInvocableFunction - Pass a function reference through - skips the need to get/resolve the type etc
func CreateInvocableFunction(function any) *Invocable {
	bindingType := getType(function)

	isInvoc, invocableType := isInvocable(bindingType)
	if !isInvoc {
		log.Printf("type passed to CreateInvocableFunction is not an invocable type(function or struct)")
		return nil
	}

	return &Invocable{
		instance:       getVal(function),
		bindingType:    bindingType,
		typeOfBinding:  invocableType,
		isInstantiated: true,
	}
}

// CreateInvocableStruct - Pass a struct reference through - skips the need to get/resolve the type etc
func CreateInvocableStruct(structRef any) *Invocable {
	bindingType := getType(structRef)

	isInvoc, invocableType := isInvocable(bindingType)
	if !isInvoc {
		log.Printf("type passed to CreateInvocableStruct (%s) is not an invocable type(function or struct)", bindingType.String())
		return nil
	}

	return &Invocable{
		instance:       getVal(structRef),
		bindingType:    bindingType,
		typeOfBinding:  invocableType,
		isInstantiated: true,
	}
}

type Invocable struct {
	bindingType reflect.Type

	// struct or func
	typeOfBinding string

	// Instantiated Value
	instance       reflect.Value
	isInstantiated bool
}

func (invocable *Invocable) instantiate() {
	if invocable.typeOfBinding == "func" {
		invocable.instantiateFunction()
	}
	if invocable.typeOfBinding == "struct" {
		invocable.instantiateStruct()
	}
}

func (invocable *Invocable) instantiateFunction() {
	if invocable.typeOfBinding != "func" {
		log.Printf("Cannot Instantiate type of %s. instantiateFunction() can only instantiate functions.", invocable.bindingType.String())
		return
	}

	invocable.instance = reflect.New(invocable.bindingType)
	invocable.isInstantiated = true

}

func (invocable *Invocable) instantiateStruct() {
	if invocable.typeOfBinding != "struct" {
		log.Printf("Cannot Instantiate type of %s. instantiateStruct() can only instantiate structs.", invocable.bindingType.String())
		return
	}

	invocable.instance = reflect.New(invocable.bindingType)
	invocable.isInstantiated = true
}

func (invocable *Invocable) InstantiateStructAndFill(container *ContainerInstance) reflect.Value {
	if !invocable.isInstantiated {
		invocable.instantiate()
	}

	return container.resolveStructFields(invocable.bindingType, invocable.instance)
}

// InstantiateWith - Instantiate a struct and fill its fields with values from the container
func (invocable *Invocable) InstantiateWith(container *ContainerInstance) any {
	resolvedStruct := invocable.InstantiateStructAndFill(container)

	return resolvedStruct.Interface()
}

// CallMethodByNameWith - Call the method and assign its parameters from the passed parameters & container
func (invocable *Invocable) CallMethodByNameWith(methodName string, container *ContainerInstance, parameters ...any) []reflect.Value {
	if invocable.typeOfBinding != "struct" {
		panic("CallMethodByNameWith is only usable when the Invocable instance is created with a struct.")
	}
	if !invocable.isInstantiated {
		invocable.instantiate()
	}

	structInstance := invocable.InstantiateStructAndFill(container)
	method := structInstance.MethodByName(methodName)

	return method.Call(container.ResolveFunctionArgs(method, parameters...))
}

func (invocable *Invocable) CallMethodByNameWithArgInterceptor(methodName string, container *ContainerInstance, interceptor FuncArgResolverInterceptor, parameters ...any) []reflect.Value {
	if invocable.typeOfBinding != "struct" {
		panic("CallMethodByNameWith is only usable when the Invocable instance is created with a struct.")
	}
	if !invocable.isInstantiated {
		invocable.instantiate()
	}

	structInstance := invocable.InstantiateStructAndFill(container)
	method := structInstance.MethodByName(methodName)

	return method.Call(container.ResolveFunctionArgsWithInterceptor(method, interceptor, parameters...))
}

// CallMethodWith - Call the method and assign its parameters from the passed parameters & container
func (invocable *Invocable) CallMethodWith(container *ContainerInstance, parameters ...any) []reflect.Value {
	if !invocable.isInstantiated {
		invocable.instantiate()
	}

	return invocable.instance.Call(
		container.ResolveFunctionArgs(invocable.instance, parameters...),
	)
}
