# IoC Container
Elegant Go IoC based on Laravels container

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

type HelloWorldService struct {}
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
When Make/MakeTo is called, any dependencies your service requires
Will be resolved from the container... So for example, you could do this

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
		return &ServiceTwo{ServiceOne:serviceOne}
	})
	
	// Now when we resolve ServiceTwo it will have an instance of ServiceOne attached
	var serviceTwo *ServiceTwo
	Container.MakeTo(&serviceTwo)
	
	print(serviceTwo.ServiceOne.msg) // outputs "Hello There!"
}
```


#### Features

- Binding:
  - Abstract -> Concrete
  - Concrete
  - Abstract -> Concrete via function
  - Singletons (`` Container.Singleton(new(SingletonService)) ``)
  - Singleton Instances(pre created) (`` Container.Instance(someVarWithInstance) ``)
  - Tagging categories of bindings with a string (`` Container.Tag("SomeCategory", new(ServiceOne), new(ServiceTwo)) `` - `` Container.Tagged("SomeCategory")``)
- Resolution:
  - Finding required args to instantiate via a function and injecting them
  - // More in todo below
- Child Containers - (`` Container.CreateChildContainer() ``)
  - If the binding isn't found in the child, it will be resolved from parents
  - Allowing for request based Containers, that then fall back to the main container



There's a lot more to document, but it's still a WIP, just wanted to get this out today

### TODO:
- [ ] Struct field injection (if your binding has services which exist in the container, they'll be resolved and set on the struct during resolve) - I have the code for this, just need to add it
- [ ] Ability to call a method via the container
  - [ ] Ability to call a method on a binding via the container
- [ ] Container resolution events (hook into bindings being resolved)
- Probably lots more :D


### One thing to clear up
I get it that a lot of Go programmers don't feel we need these things, or think that we shouldn't try to make things "elegant"/"framework-like". But I disagree and that's my personal opinion. We all have different ones :)