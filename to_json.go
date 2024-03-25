package main

import (
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

type Person struct {
	Name       string `json:"name"`
	Age        int    `json:"age"`
	IsMarraied bool   `json:"is_marraied"`
}

func ToJson() {
	k := map[int]Person{
		1: {Name: "uccs", Age: 18, IsMarraied: false},
		2: {Name: "uccs", Age: 18, IsMarraied: true},
		3: {Name: "uccs", Age: 18, IsMarraied: true},
	}

	fmt.Println(JsonMarshal(k))

	fmt.Println(JsonMarshal(&[...]interface{}{
		1,
		&Person{Name: "uccs", Age: 18, IsMarraied: false},
		Person{Name: "uccs", Age: 18, IsMarraied: true},
		true,
		Person{Name: "uccs", Age: 18, IsMarraied: true},
	}))
}

func JsonMarshal(v interface{}) (string, error) {
	rv := reflect.ValueOf(v)
	rt := rv.Type()

	switch rt.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%v", rv.Int()), nil
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%v", rv.Float()), nil
	case reflect.String:
		return fmt.Sprintf("%q", rv.String()), nil
	case reflect.Bool:
		return fmt.Sprintf("%v", rv.Bool()), nil
	case reflect.Slice:
		return marshalSlice(rv)
	case reflect.Struct:
		return marshalStruct(rv)
	case reflect.Array:
		return marshalSlice(rv.Slice(0, rv.Len()))
	case reflect.Map:
		return marshalMap(rv)
	case reflect.Pointer:
		if rv.Elem().Kind() == reflect.Array {
			return marshalSlice(rv.Elem().Slice(0, rv.Len()))
		}
		if rv.Elem().Kind() == reflect.Struct {
			return JsonMarshal(rv.Elem().Interface())
		}
		return JsonMarshal(rv.Elem().Interface())
	default:
		return "", fmt.Errorf("unsupported type: %s", rt)
	}

}

func marshalMap(rv reflect.Value) (string, error) {
	var items []string
	for _, mapKey := range rv.MapKeys() {
		mapValue := rv.MapIndex(mapKey)
		keyJsonString, err := JsonMarshal(mapKey.Interface())
		if err != nil {
			return "", err
		}
		valueJsonString, err := JsonMarshal(mapValue.Interface())
		if err != nil {
			return "", err
		}
		items = append(items, fmt.Sprintf("%v:%v", keyJsonString, valueJsonString))
	}
	return "{" + strings.Join(items, ",") + "}", nil
}

func marshalSlice(rv reflect.Value) (string, error) {
	var items []string
	for i := 0; i < rv.Len(); i++ {
		value, err := JsonMarshal(rv.Index(i).Interface())
		if err != nil {
			return "", err
		}
		items = append(items, value)
	}
	return "[" + strings.Join(items, ",") + "]", nil
}

func marshalStruct(rv reflect.Value) (string, error) {
	var items []string
	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		jsonTag := rv.Type().Field(i).Tag.Get("json")
		key := rv.Type().Field(i).Name
		if !isFieldExported(key) {
			continue
		}
		if jsonTag != "" {
			key = jsonTag
		}
		value, err := JsonMarshal(fieldValue.Interface())
		if err != nil {
			return "", err
		}

		items = append(items, fmt.Sprintf("%q:%v", key, value))
	}
	return "{" + strings.Join(items, ",") + "}", nil
}

func isFieldExported(name string) bool {
	r, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(r)
}
