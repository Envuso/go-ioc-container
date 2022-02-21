# IoC Container

Elegant Go IoC based on Laravels container

[![Go Reference](https://pkg.go.dev/badge/github.com/Envuso/go-ioc-container.svg)](https://pkg.go.dev/github.com/Envuso/go-ioc-container)
[![Go Report Card](https://goreportcard.com/badge/github.com/Envuso/go-ioc-container)](https://goreportcard.com/report/github.com/Envuso/go-ioc-container)

### Installation/Setup

```shell
go get github.com/envuso/go-ioc-container
```

You now have a global container at your disposal

Bind a service to an abstraction

```go
package main

type SayHelloService interface {
	SayHello() string
}

func NewHelloWorldService() SayHelloService {
	return &HelloWorldService{}
}

type HelloWorldService struct{}

func (s *HelloWorldService) SayHello() string {
	return "Hello World"
}

func main() {
	// We can do it this way
	Container.Bind(new(SayHelloService), new(HelloWorldService))
	// Binding abstract -> concrete resolver function
	Container.Bind(new(SayHelloService), NewHelloWorldService)
	// Binding abstract -> concrete via single function
	Container.Bind(NewHelloWorldService)
	
	// Now we resolve...
	
	// Option one, type casts
	service := Container.Make(new(SayHelloService)).(SayHelloService)
	
	// Option two, bind to var
	var service SayHelloService
	Container.MakeTo(&service)
}
```

#### Injection

When Make/MakeTo is called, any dependencies your service requires Will be resolved from the container... So for
example, you could do this

```go
package main

func main() {
	
	// Create and bind first thing to container
	type ServiceOne struct {
		msg string
	}
	Container.Bind(func() *ServiceOne {
		return &ServiceOne{msg: "Hello there!"}
	})
	
	// Now lets set up our service which uses this
	type ServiceTwo struct {
		ServiceOne *ServiceOne
	}
	
	// Our resolver function, ServiceOne is a dependency	
	Container.Bind(new(ServiceTwo), func(serviceOne *ServiceOne) *ServiceTwo {
		return &ServiceTwo{ServiceOne: serviceOne}
	})
	
	// Now when we resolve ServiceTwo it will have an instance of ServiceOne attached
	var serviceTwo *ServiceTwo
	Container.MakeTo(&serviceTwo)
	
	print(serviceTwo.ServiceOne.msg) // outputs "Hello There!"
}
```

#### Calling methods/Creating structs

If we wish to instantiate a struct/call a method using dependency injection, we can!

##### Calling methods with DI

```go
package main

// This is a nice simple way I like to boot my apps up
func main() {
	// Bind some service to the container
	Container.Bind(NewDatabaseService)
	Container.Call(bootApp)
}

func bootApp(database *DatabaseService) {
	// database will be automatically injected into the method & called for us.
}

```

##### Instantiating structs with DI

```go
package main

type SomeOtherServiceAbstract interface {
	DoAThing()
}
type SomeOtherService struct{}

func (s *SomeOtherService) DoAThing() {
	print("Hi :)")
}

type SomeService struct {
	someOtherService *SomeOtherService
}

func main() {
	// Bind our services to the container
	Container.Bind(func() *SomeOtherService {
		return &SomeOtherService{}
	})
	Container.Bind(func() *SomeService {
		return &SomeService{}
	})
	Container.Call(bootApp)
}

// Resolve our service from the container
func bootApp(service *SomeService) {
	// service.someOtherService will now be a fresh instance of SomeOtherService
	service.someOtherService.DoAThing()
}

```

#### Features

- Binding:
    - Abstract -> Concrete
    - Concrete
    - Abstract -> Concrete via function
    - Singletons (`` Container.Singleton(new(SingletonService)) ``)
    - Singleton Instances(pre created) (`` Container.Instance(someVarWithInstance) ``)
    - Tagging categories of bindings with a
      string (`` Container.Tag("SomeCategory", new(ServiceOne), new(ServiceTwo)) ``
      - `` Container.Tagged("SomeCategory")``)
- Resolution:
    - Finding required args to instantiate via a function and injecting them
    - Instantiating a struct and filling its fields
- Dependency Injection:
    - Ability to call a method via the container (`` Container.Call(methodReference) ``) - Type hinted parameters are resolved from the container(if bound)
    - Ability to instantiate a struct & fill the fields (atm, only for structs bound to the container)
      - This allows us to bind to the container, and have additional field level injection, rather than just the function we bind with
      - Struct tag & Config option to only inject to fields with the specified tag(basically complete, need to test & check some things)
- Child Containers - (`` Container.CreateChildContainer() ``)
    - If the binding isn't found in the child, it will be resolved from parents
    - Allowing for request based Containers, that then fall back to the main container
- "Invocation" helper:
    - This is a helper I created to make calling a method/instantiating & filling struct fields a bit cleaner
      - `` CreateInvocable(reflect.TypeOf(method or struct) `` - This will give us an instance of "Invocable" back
    - More docs to come in the future... refer to Container_test.go ^^


There's a lot more to document, but it's still a WIP, just wanted to get this out today

### TODO:

- [x] Struct field injection (if your binding has services which exist in the container, they'll be resolved and set on
  the struct during resolve) - I have the code for this, just need to add it
- [x] Ability to call a method via the container
    - [ ] Ability to call a method on a binding via the container
    - Note: This is already possible, im just unsure on the syntax I want to go with
- [ ] Container resolution events (hook into bindings being resolved)
- Probably lots more :D

### One thing to clear up

I get it that a lot of Go programmers don't feel we need these things, or think that we shouldn't try to make things "
elegant"/"framework-like". But I disagree and that's my personal opinion. We all have different ones :)