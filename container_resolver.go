package container

import (
	"errors"
	"log"
	"reflect"
	"unsafe"
)

// resolve - This works in a couple of different ways:
//
// Singletons:
// - If type exists in container.resolved
//   - Return the value
// - If it doesn't and:
//   - It has been bound by a function, we'll call the function
//   - It has been bound by a type, we'll instantiate the type
// - We'll then add the result of the above into container.resolved
//
// Function bindings:
// - Call the function and inject the args, return it
//
// Type bindings:
// - Instantiate the type, return it
func (container *ContainerInstance) resolve(binding *Binding, parameters ...any) any {
	if binding.isSingleton {
		return container.resolveSingleton(binding, parameters...)
	}

	if binding.isFunctionResolver {
		return container.resolveFromFunctionResolver(binding, parameters...)
	}

	return binding.invocable.InstantiateWith(container)
}

// resolveStructFields - Attempt to resolve all the fields from the container, for the specified struct
func (container *ContainerInstance) resolveStructFields(instanceType reflect.Type, instance reflect.Value) reflect.Value {
	if instanceType == nil {
		panic(errors.New("container: invalid structure"))
	}

	structType := indirectType(instanceType)
	if structType.Kind() != reflect.Struct {
		panic(errors.New("container: invalid structure"))
	}

	structValue := instance
	if instance.Kind() == reflect.Ptr {
		structValue = instance.Elem()
	}

	for i := 0; i < structValue.NumField(); i++ {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		if container.Config.OnlyInjectStructFieldsWithInjectTag {
			if tag, ok := fieldType.Tag.Lookup("inject"); ok {
				print("Inject tag is : " + tag)
			}
			continue
		}

		fieldBinding := container.getBindingType(field.Type())
		if fieldBinding != nil {
			resolved := container.makeFromBinding(fieldBinding)
			if resolved != nil {
				ptr := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
				ptr.Set(reflect.ValueOf(resolved))
			}
		}
	}

	return instance
}

// resolveFunctionArgs - Resolves the args of our function we bound to the container
// parameters is an array of values that we wish to provide which is optional
// parameters will first be assigned starting at index 0 of the functions args
// Then we'll look at the function args, and if we assigned a value from the parameters already
// it will use that, otherwise we'll look the type up in the container and resolve it
func (container *ContainerInstance) resolveFunctionArgs(function reflect.Value, parameters ...any) []reflect.Value {
	inArgCount := 0

	if !function.IsValid() || function.IsZero() {
		return []reflect.Value{}
	}

	functionType := function.Type()
	if functionType.Kind() == reflect.Ptr {
		inArgCount = functionType.Elem().NumIn()
	} else {
		inArgCount = functionType.NumIn()
	}

	// We'll put the types of all in args into this array,
	// so we don't have to keep running .In()
	// Theoretically more performant?
	inArgTypes := make([]reflect.Type, inArgCount)
	for i := 0; i < inArgCount; i++ {
		inArgTypes[i] = functionType.In(i)
	}

	// Keep track of the args we've assigned, since we assign from
	// parameters and then assign remaining types from container
	assignedCount := 0
	assignedArgs := make([]bool, inArgCount)

	// The args for the function
	args := make([]reflect.Value, inArgCount)

	// We'll call this to assign an arg and mark it as assigned
	assignArg := func(i int, arg reflect.Value) {
		args[i] = arg
		assignedArgs[i] = true
		assignedCount++
	}

	// Assign parameter values from the provided parameters list first
	if len(parameters) > 0 {
		for i := 0; i < len(parameters); i++ {
			paramVal := reflect.ValueOf(parameters[i])

			// We'll only assign the param from the provided list, if the type matches?
			inArg := inArgTypes[i]
			if inArg == paramVal.Type() {
				assignArg(i, paramVal)
			}
		}

		// If our provided parameters fulfils all the function args, let's just early return
		if assignedCount >= inArgCount {
			return args
		}
	}

	// Now we'll try to resolve any other types from the container
	for i := 0; i < inArgCount; i++ {
		// We already assigned it, let's skip...
		if assignedArgs[i] {
			continue
		}

		// Now we'll attempt to resolve in inArg from the container...
		// If it can be resolved/exists, we'll provide the value
		// Otherwise, we'll create a new zero type of the arg
		resolved, didResolve := container.resolveFunctionArg(inArgTypes[i])
		if !didResolve {
			log.Printf("Assigning empty arg for arg(%d) on resolving type %s", i, functionType.Out(0).String())
		}

		assignArg(i, resolved)
	}

	return args
}

// resolveFunctionArg - Used in resolveFunctionArgs, we pass an arg type and attempt to
// resolve it from the container, if the type doesn't exist in the container
// we'll return a zero value version of the type
func (container *ContainerInstance) resolveFunctionArg(arg reflect.Type) (reflect.Value, bool) {
	argBinding := container.getBindingType(arg)
	if argBinding == nil {
		return reflect.New(arg), false
	}

	resolved := container.Make(argBinding)
	if resolved == nil {
		return reflect.New(arg), false
	}

	return reflect.ValueOf(resolved), true
}

// resolveFromFunctionResolver - Call the bound concrete function and provide any args,
// from parameters & the container. If our bound function returns an error for the second
// return value, and there is an error, our code will panic.
func (container *ContainerInstance) resolveFromFunctionResolver(binding *Binding, parameters ...any) any {

	instanceReturnValues := binding.invocable.CallMethodWith(container, parameters...)

	// If we have two return values... it's possible arg 1 is our implementation, arg 2 is an error?
	// If this is the case, we'll panic, idk what to do here.
	if len(instanceReturnValues) >= 2 {
		instance := instanceReturnValues[0]

		if err, ok := instanceReturnValues[1].Interface().(error); ok {
			panic(err.Error())
		}

		return instance.Interface()
	}

	return instanceReturnValues[0].Interface()
}

// resolveSingleton - Works similarly to resolve, except we're doing the function/type binding parts
// If our instance already exists in container.resolved, we'll return it from there
func (container *ContainerInstance) resolveSingleton(binding *Binding, parameters ...any) any {
	if instance, ok := container.resolved[binding.concreteType]; ok {
		return instance
	}

	var resolvedInstance any

	if binding.isFunctionResolver {
		resolvedInstance = container.resolveFromFunctionResolver(binding, parameters...)
		if resolvedInstance == nil {
			return nil
		}
	} else {
		resolvedInstance = binding.invocable.InstantiateWith(container)
		if resolvedInstance == nil {
			return nil
		}
	}

	container.resolved[binding.concreteType] = resolvedInstance

	return resolvedInstance
}
