package Container

import "testing"

type ServiceAbstract interface {
	Message() string
}
type AnotherServiceAbstract interface {
	Message() string
}

type ServiceConcreteTwo struct {
	message string
}

func (service *ServiceConcreteTwo) Message() string {
	if service.message != "" {
		return service.message
	}
	return "Hello World!"
}

type ServiceConcrete struct {
	message        string
	anotherService AnotherServiceAbstract
}

func (service *ServiceConcrete) Message() string {
	if service.message != "" {
		return service.message
	}
	return "Hello World!"
}

func createSingletonServiceOne() *ServiceConcrete {
	service := new(ServiceConcrete)
	service.message = "TestBindingSingletonValueToContainer"
	return service
}
func createSingletonServiceTwo() *ServiceConcrete {
	service := new(ServiceConcrete)
	service.message = "TestBindingSingletonTypeAndResolverFuncToContainer"
	return service
}
func createSingletonServiceThree() *ServiceConcrete {
	service := new(ServiceConcrete)
	service.message = "TestBindingSingletonFuncToContainer"
	return service
}
func newServiceConcrete() *ServiceConcrete {
	return &ServiceConcrete{
		message: "plain service concrete",
	}
}
func newServiceConcreteTwo() *ServiceConcreteTwo {
	return &ServiceConcreteTwo{
		message: "plain service concrete#2",
	}
}
func newServiceConcreteWithMessageArg(message string) ServiceAbstract {
	return &ServiceConcrete{
		message: message,
	}
}
func newServiceConcreteWithMessageArgAndService(message string, service AnotherServiceAbstract) ServiceAbstract {
	return &ServiceConcrete{
		message:        message,
		anotherService: service,
	}
}
func newAnotherService() AnotherServiceAbstract {
	return &ServiceConcrete{message: "Another service"}
}

//
// BINDING SERVICES TO CONTAINER CONTAINER
//

func TestBindingAbstractInterfaceToConcreteImplementation(t *testing.T) {
	container := CreateContainer()
	if !container.Bind((*ServiceAbstract)(nil), ServiceConcrete{}) {
		t.Fatal("Could not bind ServiceAbstract to ServiceConcrete")
	}

	if !container.IsBound((*ServiceAbstract)(nil)) {
		t.Fatal("Could not find binding for ServiceAbstract")
	}
	if !container.IsBound(ServiceConcrete{}) {
		t.Fatal("Could not find binding for ServiceConcrete")
	}
}

func TestBindingAbstractInterfaceToConcreteImplementationUsingPtr(t *testing.T) {
	container := CreateContainer()
	if !container.Bind((*ServiceAbstract)(nil), &ServiceConcrete{}) {
		t.Fatal("Could not bind ServiceAbstract to ServiceConcrete")
	}

	if !container.IsBound((*ServiceAbstract)(nil)) {
		t.Fatal("Could not find binding for ServiceAbstract")
	}
	if !container.IsBound(ServiceConcrete{}) {
		t.Fatal("Could not find binding for ServiceConcrete")
	}
}

func TestBindingAbstractToConcreteViaFunction(t *testing.T) {
	container := CreateContainer()
	bound := container.Bind(func() ServiceAbstract {
		return &ServiceConcrete{}
	})

	if !bound {
		t.Fatal("Could not bind ServiceAbstract to ServiceConcrete via function")
	}
	if !container.IsBound((*ServiceAbstract)(nil)) {
		t.Fatal("Could not find binding for ServiceAbstract")
	}
}

func TestBindingAbstractToConcreteViaFunctionWithAbstractToFunctionArgs(t *testing.T) {
	container := CreateContainer()
	bound := container.Bind((*ServiceAbstract)(nil), func() *ServiceConcrete {
		return &ServiceConcrete{}
	})

	if !bound {
		t.Fatal("Could not bind ServiceAbstract to ServiceConcrete via function")
	}
	if !container.IsBound((*ServiceAbstract)(nil)) {
		t.Fatal("Could not find binding for ServiceAbstract")
	}
	if !container.IsBound(ServiceConcrete{}) {
		t.Fatal("Could not find binding for ServiceConcrete")
	}
}

//
// RESOLVING SERVICES FROM CONTAINER
//

func TestResolvingWithTypeCast(t *testing.T) {
	container := CreateContainer()
	if !container.Bind((*ServiceAbstract)(nil), ServiceConcrete{}) {
		t.Fatal("Could not bind ServiceAbstract to ServiceConcrete")
	}

	serviceResolved := container.Make((*ServiceAbstract)(nil))

	if serviceResolved == nil {
		t.Fatal("Could not resolve ServiceConcrete via ServiceAbstract")
	}

	service, ok := serviceResolved.(ServiceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of ServiceConcrete is not ServiceAbstract")
	}
	expected := "Hello World!"
	received := service.Message()
	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingWithVar(t *testing.T) {
	container := CreateContainer()
	if !container.Bind((*ServiceAbstract)(nil), ServiceConcrete{}) {
		t.Fatal("Could not bind ServiceAbstract to ServiceConcrete")
	}

	var service ServiceAbstract
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve ServiceConcrete via ServiceAbstract")
	}

	_, ok := service.(ServiceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of ServiceConcrete is not ServiceAbstract")
	}

	expected := "Hello World!"
	received := service.Message()
	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingWithProvidedArgs(t *testing.T) {
	container := CreateContainer()
	if !container.Bind(newServiceConcreteWithMessageArg) {
		t.Fatal("Could not bind via newServiceConcreteWithMessageArg func")
	}
	expected := "A custom message"

	var service ServiceAbstract
	container.MakeTo(&service, expected)

	if service == nil {
		t.Fatal("Could not resolve ServiceConcrete via ServiceAbstract")
	}

	_, ok := service.(ServiceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of ServiceConcrete is not ServiceAbstract")
	}

	received := service.Message()
	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingWithSingleProvidedArgAndRestFromContainer(t *testing.T) {
	container := CreateContainer()
	if !container.Bind(newServiceConcreteWithMessageArgAndService) {
		t.Fatal("Could not bind via newServiceConcreteWithMessageArg func")
	}
	if !container.Bind(newAnotherService) {
		t.Fatal("Could not bind newAnotherService")
	}
	expectedMessage := "A custom message"

	var service ServiceAbstract
	container.MakeTo(&service, expectedMessage, (*AnotherServiceAbstract)(nil))

	if service == nil {
		t.Fatal("Could not resolve ServiceConcrete via ServiceAbstract")
	}

	_, ok := service.(ServiceAbstract)
	if !ok {
		t.Fatal("Attempted resolve of ServiceConcrete is not ServiceAbstract")
	}

	received := service.Message()
	if received != expectedMessage {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expectedMessage)
	}
}

//
// BINDING SINGLETONS TO CONTAINER
//

func TestBindingSingletonInstanceToContainer(t *testing.T) {
	container := CreateContainer()

	didBind := container.Instance(createSingletonServiceOne())
	if didBind == false {
		t.Fatalf("Failed to bind singleton `new(ServiceConcrete)` to the container.")
	}

	if !container.IsBound(new(ServiceConcrete)) {
		t.Fatalf("Singleton `new(ServiceConcrete)` is not bound to the container.")
	}
}

func TestBindingSingletonTypeAndResolverFuncToContainer(t *testing.T) {
	container := CreateContainer()

	didBind := container.Singleton(new(ServiceConcrete), createSingletonServiceTwo)

	if didBind == false {
		t.Fatalf("Failed to bind singleton `new(ServiceConcrete)` with resolver func to the container.")
	}

	if !container.IsBound(new(ServiceConcrete)) {
		t.Fatalf("Singleton `new(ServiceConcrete)` with resolver func is not bound to the container.")
	}
}

func TestBindingSingletonFuncToContainer(t *testing.T) {
	container := CreateContainer()

	didBind := container.Singleton(createSingletonServiceThree)

	if didBind == false {
		t.Fatalf("Failed to bind singleton resolver func for ServiceConcrete to the container.")
	}

	if !container.IsBound(new(ServiceConcrete)) {
		t.Fatalf("Singleton resolver func ServiceConcrete is not bound to the container.")
	}
}

//
// RESOLVING SINGLETONS FROM CONTAINER
//

func TestResolvingSingletonInstance(t *testing.T) {
	container := CreateContainer()
	singleton := createSingletonServiceOne()
	container.Instance(singleton)

	var service *ServiceConcrete
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve singleton ServiceConcrete")
	}

	expected := singleton.message
	received := service.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}

	service.message = "Im now modified?"
	var newService *ServiceConcrete
	container.MakeTo(&newService)

	expected = "Im now modified?"
	received = newService.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingSingletonTypeAndResolverFunc(t *testing.T) {
	container := CreateContainer()
	singleton := createSingletonServiceTwo()
	container.Singleton(new(ServiceConcrete), createSingletonServiceTwo)

	var service *ServiceConcrete
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve singleton ServiceConcrete")
	}

	expected := singleton.message
	received := service.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

func TestResolvingSingletonFuncFunc(t *testing.T) {
	container := CreateContainer()
	singleton := createSingletonServiceThree()
	container.Singleton(createSingletonServiceThree)

	var service *ServiceConcrete
	container.MakeTo(&service)

	if service == nil {
		t.Fatal("Could not resolve singleton ServiceConcrete")
	}

	expected := singleton.message
	received := service.Message()

	if received != expected {
		t.Fatalf("Return of Message() on service is invalid.\nGot: '%s'.\nExpected: '%s'", received, expected)
	}
}

//
// CHILD CONTAINERS
//

func TestChildBindingDoesNotAffectParent(t *testing.T) {
	container := CreateContainer()

	childContainer := container.CreateChildContainer()
	childContainer.Singleton(createSingletonServiceOne)

	if !childContainer.IsBound(new(ServiceConcrete)) {
		t.Fatal("Could not find binding in child container for ServiceConcrete")
	}

	if container.IsBound(new(ServiceConcrete)) {
		t.Fatal("Service Concrete is bound to parent container and shouldn't be")
	}

}

func TestResettingContainer(t *testing.T) {
	container := CreateContainer()
	container.Singleton(createSingletonServiceOne)

	if !container.IsBound(new(ServiceConcrete)) {
		t.Fatal("Could not find binding in container")
	}

	container.Reset()

	if container.IsBound(new(ServiceConcrete)) {
		t.Fatal("Binding is still set in container after reset?")
	}
}

func TestClearingContainerInstances(t *testing.T) {
	container := CreateContainer()
	container.Singleton(createSingletonServiceOne)

	if !container.IsBound(new(ServiceConcrete)) {
		t.Fatal("Could not find binding in container")
	}

	var service *ServiceConcrete
	container.MakeTo(&service)
	service.message = "before clear"

	container.ClearInstances()

	if !container.IsBound(new(ServiceConcrete)) {
		t.Fatal("Could not find binding in container")
	}

	var nextService *ServiceConcrete
	container.MakeTo(&nextService)

	if nextService.Message() == "before clear" {
		t.Fatal("ServiceConcrete resolved after ClearInstances() call, still has value from previous singleton")
	}

}

//
// TAGGED BINDINGS
//

func TestAddingTaggedBindings(t *testing.T) {
	container := CreateContainer()
	container.Bind(newServiceConcrete)
	container.Bind(newServiceConcreteTwo)

	if !container.IsBound(new(ServiceConcrete)) || !container.IsBound(new(ServiceConcreteTwo)) {
		t.Fatal("Could not find bindings in container")
	}

	didTag := container.Tag("SingletonServices", new(ServiceConcrete), new(ServiceConcreteTwo))
	if !didTag {
		t.Fatal("Failed to tag services")
	}

	tagged := container.Tagged("SingletonServices")

	if len(tagged) == 0 {
		t.Fatal("Failed to get tagged services")
	}
	if len(tagged) != 2 {
		t.Fatalf("We should only have two tagged services, we have %d", len(tagged))
	}

	if _, ok := tagged[0].(*ServiceConcrete); !ok {
		t.Fatal("First tagged service is not ServiceConcrete")
	}

	if _, ok := tagged[1].(*ServiceConcreteTwo); !ok {
		t.Fatal("First tagged service is not ServiceConcreteTwo")
	}
}
