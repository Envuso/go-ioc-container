package container

import (
	"log"
	"reflect"
)

func isInvocable(typ reflect.Type) (bool, string) {
	kind := typ.Kind()

	if kind == reflect.Ptr {
		typ = indirectType(typ)
	}

	if kind == reflect.Struct {
		return true, "struct"
	}

	if kind == reflect.Func {
		return true, "func"
	}

	return false, ""
}

func indirectType(typ reflect.Type) reflect.Type {
	switch typ.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		return typ.Elem()
	}
	return typ
}

// getType - Allows us to get the type if it's not already a type. If it is a type...
// return it, just prevents us constantly doing .TypeOf() or this everywhere
func getType(t any) reflect.Type {
	if _, ok := t.(reflect.Type); ok {
		return t.(reflect.Type)
	}
	return reflect.TypeOf(t)
}

// getVal - Allows us to get the val if it's not already a value. If it is a value...
// return it, just prevents us constantly doing .ValueOf() or this everywhere
func getVal(t any) reflect.Value {
	if _, ok := t.(reflect.Value); ok {
		return t.(reflect.Value)
	}
	return reflect.ValueOf(t)
}

// getAbstractReturnType - Convenience function, if our type is a pointer, we'll get the underlying type
// If it's also not an interface, nil will be returned
func getAbstractReturnType(abstractType reflect.Type) reflect.Type {
	var abstract = abstractType

	if abstractType.Kind() == reflect.Ptr {
		abstract = abstractType.Elem()
	}

	if abstract.Kind() != reflect.Interface {
		return nil
	}

	return abstract
}

// getConcreteReturnType - Allows us to pass a function and get it's first
// return arg or pass a struct and get the type of that
func getConcreteReturnType(concrete reflect.Type) reflect.Type {
	returnType := concrete

	if concrete.Kind() == reflect.Struct {
		return returnType
	}

	if concrete.Kind() == reflect.Func {
		numOut := concrete.NumOut()
		if numOut == 0 {
			log.Printf("Trying to get function return type for binding but it doesnt have a return type...")
			return nil
		}
		if numOut > 1 {
			log.Printf("Getting a function return type, but the function has > 1 return args. Only the first arg is handled.")
		}

		returnType = concrete.Out(0)
	}

	if returnType.Kind() == reflect.Pointer {
		return indirectType(returnType)
	}

	if returnType != nil {
		return returnType
	}

	return nil
}
