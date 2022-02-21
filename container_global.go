package container

// Not sure if this is the right way to do it...
// We're exposing some function which are the same as using "Container"
// but for an end user, they will have to use "container.Container"
// So we'll "proxy" those calls through to the global container?

func Bind(bindingDef ...any) bool {
	return Container.Bind(bindingDef...)
}
func Singleton(singleton any, concreteResolverFunc ...any) bool {
	return Container.Singleton(singleton, concreteResolverFunc...)
}
func Instance(instance any) bool {
	return Container.Instance(instance)
}
func IsBound(binding any) bool {
	return Container.IsBound(binding)
}
func Make(abstract any, parameters ...any) any {
	return Container.Make(abstract, parameters...)
}
func MakeTo(makeTo any, parameters ...any) {
	Container.MakeTo(makeTo, parameters...)
}
func CreateChildContainer() *ContainerInstance {
	return Container.CreateChildContainer()
}
func ClearInstances() {
	Container.ClearInstances()
}
func Reset() {
	Container.Reset()
}
func ParentContainer() *ContainerInstance {
	return Container.ParentContainer()
}
func Call(function any, parameters ...any) []any {
	return Container.Call(function, parameters...)
}
func Tag(tag string, bindings ...any) bool {
	return Container.Tag(tag, bindings...)
}
func Tagged(tag string) []any {
	return Container.Tagged(tag)
}
