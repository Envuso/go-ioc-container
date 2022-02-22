package tests

import (
	"errors"
	"testing"

	container "github.com/Envuso/go-ioc-container"
	"github.com/modern-go/reflect2"
	"github.com/stretchr/testify/assert"
)

var pkgPath = "github.com/Envuso/go-ioc-container"

type TestingStr struct {
	Name string
}
type TestingInter interface {
	SomeFunc()
}

func TestTypes_Of(t *testing.T) {

	testingStr := container.ContainerTypes.Of(new(TestingStr))

	assert.Equal(t, pkgPath+"/TestingStr", testingStr.FullName)
	assert.Equal(t, pkgPath, testingStr.Path)
	assert.Equal(t, "TestingStr", testingStr.Name)

	testingInter := container.ContainerTypes.Of(new(TestingInter))
	assert.Equal(t, pkgPath+"/TestingInter", testingInter.FullName)
	assert.Equal(t, pkgPath, testingInter.Path)
	assert.Equal(t, "TestingInter", testingInter.Name)

	testingErr := container.ContainerTypes.Of(errors.New("yeet"))
	assert.Equal(t, pkgPath+"/TestingInter", testingErr.FullName)
	assert.Equal(t, pkgPath, testingErr.Path)
	assert.Equal(t, "TestingInter", testingErr.Name)

	test := reflect2.TypeByName("container.TestingStr")
	test.String()

	test2 := reflect2.TypeByPackageName("github.com/Envuso/go-ioc-container", "TestingStr")
	test2.String()
	print("")
}

func TestPkgType_Save(t *testing.T) {

	assert.False(t, container.ContainerTypes.Has(new(TestingStr)))
	assert.False(t, container.ContainerTypes.Has(TestingStr{}))

	testingStr := container.ContainerTypes.Of(new(TestingStr))
	testingStr.Save()

	assert.True(t, container.ContainerTypes.Has(new(TestingStr)))
	assert.True(t, container.ContainerTypes.Has(TestingStr{}))

}
