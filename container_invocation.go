package container

// func (container *ContainerInstance) binding(abstract any) *Binding {
// 	binding := container.getBindingType(abstract)
//
// 	if binding == nil {
// 		log.Printf("Failed to resolve binding for abstract type %s", reflect.TypeOf(abstract).String())
// 		return nil
// 	}
//
// 	containerBinding, ok := container.bindings[binding]
//
// 	if ok {
// 		return containerBinding
// 	}
//
// 	if container.parent != nil {
// 		return container.parent.Binding(abstract)
// 	}
//
// 	log.Printf("Failed to resolve container binding for abstract type %s", binding.String())
// 	return nil
// }

// Call - Call the specified function via the container, you can add parameters to your function,
// and they will be resolved from the container, if they're registered
func (container *ContainerInstance) Call(function any, parameters ...any) []any {
	invocable := CreateInvocableFunction(function)

	instanceReturnValues := invocable.CallMethodWith(container, parameters...)

	returnResult := make([]any, len(instanceReturnValues))
	for i, value := range instanceReturnValues {
		returnResult[i] = value.Interface()
	}

	return returnResult
}
