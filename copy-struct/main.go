package main

import (
	"errors"
	"fmt"
	"reflect"
)

type Person struct {
	Name         string
	BirthYear    int
	NumberOfCars *int
}

type Target struct {
	Name      string
	BirthYear int
}

func CopyStruct(dst interface{}, src interface{}, lenient bool) error {
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Struct {
		return errors.New("Source type is not struct")
	}

	dstType := reflect.TypeOf(dst)
	if dstType.Kind() != reflect.Ptr {
		return errors.New("Destination is not a pointer")
	}

	// Since dst is a pointer to a struct, we need to resolve the ptr
	// indirection and get the type of the pointed data.
	dstPtrValue := reflect.ValueOf(dst)
	dstValue := reflect.Indirect(dstPtrValue)
	dstType = reflect.TypeOf(dstValue.Interface())

	if dstType.Kind() != reflect.Struct {
		return errors.New("Destination is not pointing to a struct")
	}

	srcValue := reflect.ValueOf(src)

	for i := 0; i < dstType.NumField(); i++ {
		dstField := dstType.Field(i)

		_, found := srcType.FieldByName(dstField.Name)
		if !found && !lenient {
			return fmt.Errorf("Field %s not found in source struct %s", dstField.Name, srcType.Name())
		}

		srcFieldValue := srcValue.FieldByName(dstField.Name)
		dstValue.Field(i).Set(srcFieldValue)
	}

	return nil
}

func main() {
	cars := 2
	p := Person{
		Name:         "Joe King",
		BirthYear:    1965,
		NumberOfCars: &cars,
	}
	pcopy := Target{}

	CopyStruct(&pcopy, p, false)

	fmt.Printf("Target struct %v\n", pcopy)
}
