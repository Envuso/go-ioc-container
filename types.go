package container

import (
	"reflect"
	"sync"
)

type Types struct {
	types sync.Map
}

var ContainerTypes = &Types{}

func nameOfType(p reflect.Type) *PkgType {
	name := p.Name()
	path := p.PkgPath()
	fullPath := ""

	if path != "" {
		fullPath = path + "/" + name
	}

	return &PkgType{
		Name:     name,
		Path:     path,
		FullName: fullPath,
		TypeStr:  p.String(),
		Kind:     p.Kind(),
		Type:     p,
	}
}

func (t *Types) Of(typ any) *PkgType {
	reflected := getType(typ)

	if reflected.Kind() == reflect.Ptr {
		return nameOfType(reflected.Elem())
	}

	return nameOfType(reflected)
}

func (t *Types) Clear() {
	t.types.Range(func(key any, value any) bool {
		t.types.Delete(key)
		return true
	})
}

func (t *Types) Has(r any) bool {
	typ := ContainerTypes.Of(r)

	_, ok := t.types.Load(typ.FullName)

	return ok
}

type PkgType struct {
	Name string
	Path string
	// FullName - The Path + Name split by a /
	FullName string

	Type    reflect.Type
	Kind    reflect.Kind
	TypeStr string
}

func (t *PkgType) Save() {
	ContainerTypes.types.LoadOrStore(t.FullName, t)
}
