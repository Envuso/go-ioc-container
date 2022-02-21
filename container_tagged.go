package container

import (
	"reflect"

	"golang.org/x/exp/slices"
)

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
func (container *ContainerInstance) Tag(tag string, bindings ...any) bool {
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
func (container *ContainerInstance) Tagged(tag string) []any {
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
