package Container

import (
	"log"
	"reflect"
	"unsafe"

	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

type IocContainer struct {
	// Store our singleton instances
	// instances map[reflect.Type]*IocContainerBinding

	// Our resolved singleton instances
	resolved map[reflect.Type]any

	// Store our abstract -> concrete bindings
	// If a type doesn't have an abstract type
	// We'll store concrete -> concrete
	bindings map[reflect.Type]*IocContainerBinding

	// Store aliases of Concrete -> Abstract, so we can resolve from concrete
	// when we only bound Abstract -> Concrete
	concretes map[reflect.Type]reflect.Type

	// When we register a tagged type, we'll store the tag string and then an array
	// of types for this tag, we can then use these types to resolve the bindings
	tagged map[string][]reflect.Type

	// If our container is a child container, we'll have a pointer to our parent
	parent *IocContainer
}

// CreateContainer - Create a new container instance
func CreateContainer() *IocContainer {
	c := &IocContainer{
		resolved:  make(map[reflect.Type]any),
		bindings:  make(map[reflect.Type]*IocContainerBinding),
		concretes: make(map[reflect.Type]reflect.Type),
		tagged:    make(map[string][]reflect.Type),
	}

	containerInstances = append(containerInstances, c.pointer())

	return c
}

var Container = CreateContainer()
var containerInstances = []unsafe.Pointer{}

func removeAllChildContainerInstances() {
	// mainContainerPointer := Container.pointer()
	//
	// newContainerInstances := []unsafe.Pointer{}
	//
	// for _, containerPtr := range containerInstances {
	// 	var container = *((*IocContainer)(containerPtr))
	//
	// 	if container.parent == nil && container.pointer() != mainContainerPointer {
	//
	// 	}
	// }
}

// CreateChildContainer - Returns a new container, any failed look-ups of our
// child container, will then be looked up in the parent, or returned nil
func (container *IocContainer) CreateChildContainer() *IocContainer {
	c := &IocContainer{
		resolved:  make(map[reflect.Type]any),
		bindings:  make(map[reflect.Type]*IocContainerBinding),
		concretes: make(map[reflect.Type]reflect.Type),
		tagged:    make(map[string][]reflect.Type),
	}

	c.parent = container

	containerInstances = append(containerInstances, c.pointer())

	return c
}

func (container *IocContainer) pointer() unsafe.Pointer {
	return reflect.ValueOf(container).UnsafePointer()
}

// ClearInstances - This will just remove any singleton instances from the container
// When they are next resolved via Make/MakeTo, they will be instantiated again
func (container *IocContainer) ClearInstances() {
	maps.Clear(container.resolved)
}

// Reset - Reset will empty all bindings in this container, you will have to register
// any bindings again before you can resolve them.
func (container *IocContainer) Reset() {
	maps.Clear(container.resolved)
	maps.Clear(container.bindings)
	maps.Clear(container.concretes)
}

// ParentContainer - Returns the parent container, if one exists
func (container *IocContainer) ParentContainer() *IocContainer {
	return container.parent
}

// Bind - Add a binding to the container, we can do this in a few different ways...
//
// ------
//
// Function binding:
//
// For example:
//  IocContainer.Bind(func() AbstractServiceInterface { return NewService() })
// AbstractServiceInterface will be our resolver, NewService() is the bound type to resolve.
//
// We can also just provide a concrete binding
// For example...
//  IocContainer.Bind(NewService)
//
// ------
//
// Abstract Interface -> Concrete Implementation binding:
//
// For example:
//  IocContainer.Bind((*AbstractServiceInterface)(nil), NewService)
//
// Assuming the NewService func would return something like "*Service"
// This doesn't look great, but it's the only way I know of, to pass an interface as value
//
// If "NewService" returns AbstractServiceInterface, you can just do
//  IocContainer.Bind(NewService)
//
// ------
//
// Concrete binding:
//
// For example:
//  IocContainer.Bind(SomeConcreteService{})
// or
//  IocContainer.Bind(new(SomeConcreteService))
//
func (container *IocContainer) Bind(bindingDef ...any) bool {

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

	container.addBinding(abstractType, &IocContainerBinding{
		bindingType:      "Abstract",
		abstractType:     abstractType,
		concreteType:     concreteType,
		resolverFunction: bindingDef[1],
	})

	return true
}

// Singleton - Bind a "class" that should only be instantiated once when resolved
// in the future, the initial instantiation of this type will be returned
func (container *IocContainer) Singleton(singleton any, concreteResolverFunc ...any) bool {
	singletonType := reflect.TypeOf(singleton)

	// We can provide a function to singleton
	if singletonType.Kind() == reflect.Func && concreteResolverFunc == nil {
		if singletonType.NumOut() == 0 {
			log.Printf("Please make sure your singleton function provider(%s) has at-least one return type.", singletonType.String())
			log.Printf("Without it, we don't have a type to register this singleton under.")
			log.Printf("Your type will not be registered in the container.")
			return false
		}

		container.addSingletonBinding(getConcreteReturnType(singletonType.Out(0)), &IocContainerBinding{
			bindingType: "Singleton",

			resolverFunction:   singleton,
			isFunctionResolver: true,

			abstractType: singletonType,
			concreteType: singletonType,
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
		container.addSingletonBinding(singletonConcrete, &IocContainerBinding{
			bindingType: "Singleton",

			isFunctionResolver: false,

			abstractType: singletonConcrete,
			concreteType: singletonConcrete,
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

	container.addSingletonBinding(singletonConcrete, &IocContainerBinding{
		bindingType: "Singleton",

		isFunctionResolver: true,
		resolverFunction:   resolverFunc,

		abstractType: singletonConcrete,
		concreteType: singletonConcrete,
	})

	return true
}

// Instance - This is similar to Singleton, except with Singleton we provide a type to instantiate
// With instance, we provide an already instantiated value to the container
func (container *IocContainer) Instance(instance any) bool {
	instanceType := reflect.TypeOf(instance)

	singletonConcrete := getConcreteReturnType(instanceType)

	if singletonConcrete == nil {
		log.Printf("Failed to get type of instance singleton: %s", instanceType.String())
		log.Printf("Your instance will not be registered in the container.")
		return false
	}

	container.addSingletonBinding(singletonConcrete, &IocContainerBinding{
		bindingType: "Singleton",

		isFunctionResolver: false,

		abstractType: singletonConcrete,
		concreteType: singletonConcrete,
	})

	// Our instance is already instantiated, we'll pass it straight to resolved
	container.resolved[singletonConcrete] = instance

	return true
}

// addFunctionBinding - Create a new container binding from the function
// This resolves the return type of the function as the Abstract
// And the functions return value is our Concrete
func (container *IocContainer) addFunctionBinding(definition reflect.Type, resolver any) {
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

	container.addBinding(indirectType(definition.Out(0)), &IocContainerBinding{
		bindingType: "Function",

		resolverFunction:   resolver,
		isFunctionResolver: true,

		abstractType: definition,
		concreteType: reflect.TypeOf(resolver),
	})
}

// addConcreteBinding - Create a new container binding from the concrete value
// This will set our abstract type to the concrete type and the concrete type will be our concrete type..
// This just allows us to easily bind things to the container if we don't care about abstracts
func (container *IocContainer) addConcreteBinding(definition reflect.Type, concrete any) {
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

	container.addBinding(concreteType, &IocContainerBinding{
		bindingType: "Concrete",

		isFunctionResolver: false,

		abstractType: concreteType,
		concreteType: definition,
	})
}

// resolveFunctionArgs - Resolves the args of our function we bound to the container
// parameters is an array of values that we wish to provide which is optional
// parameters will first be assigned starting at index 0 of the functions args
// Then we'll look at the function args, and if we assigned a value from the parameters already
// it will use that, otherwise we'll look the type up in the container and resolve it
func (container *IocContainer) resolveFunctionArgs(function reflect.Value, parameters ...any) []reflect.Value {
	functionType := function.Type()
	inArgCount := functionType.NumIn()

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
func (container *IocContainer) resolveFunctionArg(arg reflect.Type) (reflect.Value, bool) {
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
func (container *IocContainer) resolveFromFunctionResolver(binding *IocContainerBinding, parameters ...any) any {
	function := reflect.ValueOf(binding.resolverFunction)

	instanceReturnValues := function.Call(
		container.resolveFunctionArgs(function, parameters...),
	)

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
func (container *IocContainer) resolveSingleton(binding *IocContainerBinding, parameters ...any) any {
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
		resolvedInstance = container.resolveFromType(binding)
		if resolvedInstance == nil {
			return nil
		}
	}

	container.resolved[binding.concreteType] = resolvedInstance

	return resolvedInstance
}

// resolveFromType - Create a new instance of the concreteType and return it
func (container *IocContainer) resolveFromType(binding *IocContainerBinding) any {
	instance := reflect.New(binding.concreteType)

	return instance.Interface()
}

// hasBinding - Look up a Type in the container and return whether it exists
func (container *IocContainer) hasBinding(binding reflect.Type) bool {
	_, ok := container.bindings[binding]

	return ok
}

// getBindingType - Try to get a binding type from the binding arg in a few different ways
// We'll first assume we're checking for an abstract type binding...
// If we didn't get it from the abstract, we'll then check for the concrete...
// Now as a last ditch effort, we'll look the bindingType up in container.concretes
// If we can't find anything, we return nil
func (container *IocContainer) getBindingType(binding any) reflect.Type {
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

// IsBound - Check if the provided value type exists in our container
func (container *IocContainer) IsBound(binding any) bool {
	return container.getBindingType(binding) != nil
}

// addBinding - Convenience function to add a IocContainerBinding for the type &
// create a reverse lookup for Concrete -> Abstract
func (container *IocContainer) addBinding(abstractType reflect.Type, binding *IocContainerBinding) {
	container.bindings[abstractType] = binding
	container.concretes[binding.concreteType] = abstractType
}
func (container *IocContainer) addSingletonBinding(singletonType reflect.Type, binding *IocContainerBinding) {
	binding.isSingleton = true
	container.addBinding(singletonType, binding)
}

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
func (container *IocContainer) resolve(binding *IocContainerBinding, parameters ...any) any {
	if binding.isSingleton {
		return container.resolveSingleton(binding, parameters...)
	}

	if binding.isFunctionResolver {
		return container.resolveFromFunctionResolver(binding, parameters...)
	}

	return container.resolveFromType(binding)
}

// makeFromBinding - Once we've obtained our binding type from
// Make, we'll then check the containers bindings
// If it doesn't exist, and we have a parent container we'll then call makeFromBinding on the
// parent container. Which will either recurse until a resolve is made, or return nil
func (container *IocContainer) makeFromBinding(binding reflect.Type, parameters ...any) any {
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

// Make - Try to make a new instance of the provided value and return it
// This requires a type cast to work nicely...
// For example:
//  service := IocContainer.Make((*ServiceAbstract)(nil))
func (container *IocContainer) Make(abstract any, parameters ...any) any {
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
//  IocContainer.MakeTo(&service)
func (container *IocContainer) MakeTo(makeTo any, parameters ...any) {
	makeToVal := reflect.ValueOf(makeTo)
	if makeToVal.Kind() != reflect.Pointer {
		log.Printf("Call to IocContainer.MakeTo(), the makeTo arg must be a pointer to your receiving var. Ex; var service ServiceAbstract; IocContainer.MakeTo(&service)...")
		return
	}

	makeToValIndirect := reflect.Indirect(makeToVal)
	if !makeToValIndirect.CanSet() {
		log.Printf("Call to IocContainer.MakeTo(), the makeTo arg cannot be set?")
		return
	}

	bindingType := reflect.Indirect(makeToVal).Type()

	resolved := container.Make(bindingType, parameters...)
	if resolved == nil {
		return
	}

	makeToValIndirect.Set(reflect.ValueOf(resolved))
}

// Tag - When we've bound to the container, we can then tag the abstracts with a string
// This is useful when we want to obtain a "category" of implementations
//
// For example; Imagine we have a few different "statistic gathering" services
//
//  // Bind our individual services
//  Container.Bind(new(NewUserPostViewsStatService), func () {})
//  Container.Bind(new(NewPageViewsStatService), func () {})
//
//  // Add the services to the "StatServices" "Category"
//  Container.Tag("StatServices", new(NewUserPostViewsStatService), new(NewPageViewsStatService))
//
//  // Now we can obtain them all
//  Container.Tagged("StatServices")
//
func (container *IocContainer) Tag(tag string, bindings ...any) bool {
	if len(bindings) == 0 {
		return false
	}

	taggedTypes := []reflect.Type{}

	// Get the types of the provided bindings and create a new array
	for _, b := range bindings {
		binding := container.getBindingType(b)
		if binding == nil {
			continue
		}
		taggedTypes = append(taggedTypes, binding)
	}

	// If we couldn't get binding types and our array is empty... return
	if len(taggedTypes) == 0 {
		return false
	}

	// If we don't have any tagged types already with this tag, we'll just set and return
	if _, ok := container.tagged[tag]; !ok {
		container.tagged[tag] = taggedTypes
		return true
	}

	// We have types tagged with this tag already, so we need to merge, but make sure they're unique
	for _, taggedType := range taggedTypes {
		if !slices.Contains(container.tagged[tag], taggedType) {
			container.tagged[tag] = append(container.tagged[tag], taggedType)
		}
	}

	return len(container.tagged[tag]) > 0
}

// Tagged - Resolve the instances from the container using the specified tag
// Refer to Tag to see how adding tagged bindings works
func (container *IocContainer) Tagged(tag string) []any {
	resolved := []any{}

	if _, ok := container.tagged[tag]; !ok {
		return resolved
	}

	taggedTypes := container.tagged[tag]

	for _, taggedType := range taggedTypes {
		resolvedBinding := container.makeFromBinding(taggedType)
		if resolvedBinding == nil {
			continue
		}
		resolved = append(resolved, resolvedBinding)
	}

	return resolved
}
