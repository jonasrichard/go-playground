package main_test

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Company struct {
	Name    string
	Manager Person
	Workers []Person
}

type Person struct {
	Name              string
	Salary            int32
	BirthDate         time.Time
	Active            bool
	LastWorkTimestamp int64
	StockOptionValue  int
}

type TypeCache struct {
	// Types is a cache by reflect.TypeOf(i).String() as a key
	Types map[string]reflect.Type
}

func (t TypeCache) Unmarshal(dst interface{}, valueMap map[string]interface{}) error {
	dstPtrTypeName := reflect.TypeOf(dst).String()
	dstPtrType, ok := t.Types[dstPtrTypeName]
	if !ok {
		return fmt.Errorf("type '%s' is not in type cache", dstPtrTypeName)
	}

	if dstPtrType.Kind() != reflect.Ptr {
		return fmt.Errorf("destination '%s' is not a pointer", dstPtrTypeName)
	}

	dstPtrValue := reflect.ValueOf(dst)
	dstValue := reflect.Indirect(dstPtrValue)

	// Cut the first '*' from the name of the pointer type
	dstType := t.Types[dstPtrTypeName[1:]]

	for i := 0; i < dstType.NumField(); i++ {
		typeField := dstType.Field(i)
		valueField := dstValue.Field(i)

		v := valueMap[typeField.Name]

		log.Printf("Field %v --- Value %v\n", typeField, v)

		switch dstValue.Field(i).Kind() {
		case reflect.Int32:
			log.Printf("Struct field is int32: %s\n", typeField.Name)
			log.Printf("Value type is %T\n", v)

			switch vv := v.(type) {
			case int:
				valueField.Set(reflect.ValueOf(int32(vv)))
			case int32:
				valueField.Set(reflect.ValueOf(vv))
			default:
				return fmt.Errorf("field '%s' is int32 and cannot allow type '%T'", typeField.Name, v)
			}
		case reflect.Struct:
			// If the type is in the map, it is a custom type what we can handle with
			// Unmarshal func.
			// If it is no, the value in the valueMap will be a corresponding instance
			// of that type. Like in a time.Time type field the value in the map should
			// be a time.Time instance.
			log.Printf("Value field type: %s\n", valueField.Type().String())

			if _, ok := t.Types[valueField.Type().String()]; ok {
				if vv, ok := v.(map[string]interface{}); ok {
					subDestPtr := valueField.Addr().Interface()

					err := t.Unmarshal(subDestPtr, vv)
					if err != nil {
						return err
					}
				} else {
					return fmt.Errorf("value in the map is of type %T and not map[string]interface{}", v)
				}
			} else {
				valueField.Set(reflect.ValueOf(v))
			}
		case reflect.Slice:
			// If slice elem type is in the type cache, Unmarshal them with recursion
			elemType := typeField.Type.Elem()

			if _, ok := t.Types[elemType.String()]; ok {
				if vv, ok := v.([]map[string]interface{}); ok {
					elemSlice := reflect.MakeSlice(typeField.Type, len(vv), len(vv))

					for i := range vv {
						elemPtr := elemSlice.Index(i).Addr().Interface()

						err := t.Unmarshal(elemPtr, vv[i])
						if err != nil {
							return err
						}
					}

					valueField.Set(elemSlice)
				} else {
					return fmt.Errorf("value in the map is of type %T and not []map[string]interface{}", v)
				}
			} else {
				valueField.Set(reflect.ValueOf(v))
			}
		default:
			valueField.Set(reflect.ValueOf(v))
		}
	}

	return nil
}

func TestParsePerson(t *testing.T) {
	var person Person

	birthDate := ParsedTime("1985-12-06T12:00:00Z")
	now := time.Now().UnixMicro()

	val := PersonMap("John Smithe", 54000, birthDate, true, now, 21000)

	cache := TypeCache{
		Types: map[string]reflect.Type{
			reflect.TypeOf(person).String():  reflect.TypeOf(person),
			reflect.TypeOf(&person).String(): reflect.TypeOf(&person),
		},
	}

	err := cache.Unmarshal(&person, val)

	assert.Nil(t, err)
	assert.Equal(t, "John Smithe", person.Name)
	assert.Equal(t, int32(54000), person.Salary)
	assert.Equal(t, birthDate, person.BirthDate)
	assert.Equal(t, true, person.Active)
	assert.Equal(t, now, person.LastWorkTimestamp)
	assert.Equal(t, 21000, person.StockOptionValue)
}

func TestParseCompany(t *testing.T) {
	var (
		person  Person
		company Company
	)

	birthDate := ParsedTime("1985-12-06T12:00:00Z")
	now := time.Now().UnixMicro()

	manager := PersonMap("John Smithe", 54000, birthDate, true, now, 21000)

	val := map[string]interface{}{
		"Name":    "Apple Inc",
		"Manager": manager,
		"Workers": []map[string]interface{}{
			PersonMap("Dick Holloway", 37500, birthDate, true, now, 0),
		},
	}

	personType := reflect.TypeOf(person)
	personRefType := reflect.TypeOf(&person)
	companyType := reflect.TypeOf(company)
	companyRefType := reflect.TypeOf(&company)

	cache := TypeCache{
		Types: map[string]reflect.Type{
			personType.String():     personType,
			personRefType.String():  personRefType,
			companyType.String():    companyType,
			companyRefType.String(): companyRefType,
		},
	}

	err := cache.Unmarshal(&company, val)

	assert.Nil(t, err)
	assert.Equal(t, "Apple Inc", company.Name)
	assert.Len(t, company.Workers, 1)
	assert.Equal(t, "Dick Holloway", company.Workers[0].Name)
}

func PersonMap(name string, salary int32, birthDate time.Time, active bool, lastTS int64, stockValue int) map[string]interface{} {
	return map[string]interface{}{
		"Name":              name,
		"Salary":            salary,
		"BirthDate":         birthDate,
		"Active":            active,
		"LastWorkTimestamp": lastTS,
		"StockOptionValue":  stockValue,
	}
}

func ParsedTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}

	return t
}
