package tests

type serviceAbstract interface {
	Message() string
}
type anotherServiceAbstract interface {
	Message() string
}

type serviceConcreteTwo struct {
	message string
}

func (service *serviceConcreteTwo) Message() string {
	if service.message != "" {
		return service.message
	}
	return "Hello World!"
}

type serviceConcrete struct {
	message        string
	anotherService anotherServiceAbstract
}

func (service *serviceConcrete) Message() string {
	if service.message != "" {
		return service.message
	}
	return "Hello World!"
}

func (service *serviceConcrete) InterceptFunc(anotherService anotherServiceAbstract) string {
	return anotherService.Message()
}

func createSingletonServiceOne() *serviceConcrete {
	service := new(serviceConcrete)
	service.message = "TestBindingSingletonValueToContainer"
	return service
}
func createSingletonServiceTwo() *serviceConcrete {
	service := new(serviceConcrete)
	service.message = "TestBindingSingletonTypeAndResolverFuncToContainer"
	return service
}
func createSingletonServiceThree() *serviceConcrete {
	service := new(serviceConcrete)
	service.message = "TestBindingSingletonFuncToContainer"
	return service
}
func newServiceConcrete() *serviceConcrete {
	return &serviceConcrete{
		message: "plain service concrete",
	}
}
func newServiceConcreteTwo() *serviceConcreteTwo {
	return &serviceConcreteTwo{
		message: "plain service concrete#2",
	}
}
func newServiceConcreteWithMessageArg(message string) serviceAbstract {
	return &serviceConcrete{
		message: message,
	}
}
func newServiceConcreteWithMessageArgAndService(message string, service anotherServiceAbstract) serviceAbstract {
	return &serviceConcrete{
		message:        message,
		anotherService: service,
	}
}
func newAnotherService() anotherServiceAbstract {
	return &serviceConcrete{message: "Another service"}
}
