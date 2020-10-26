package MyMarshal

import (
	"fmt"
	"reflect"
	"strconv"
)

func JsonMarshal(v interface{}) ([]byte, error) {
	m, _ := struct2str(v)
	fmt.Println(m)
	return nil, nil
}

func struct2str(v interface{}) (string, error) {
	str := ""
	obj := reflect.ValueOf(v)
	obj_type := reflect.TypeOf(v)
	count := obj.NumField()

	switch obj_type.Kind() {
	case reflect.String:
		fmt.Printf("%v:%s\n", obj_type.Name, obj.String())
		str = obj_type.Name + " : " + obj.String()
	case reflect.Int:
		fmt.Printf("%v:%s\n", obj_type.Name, strconv.FormatInt(obj.Int(), 10))
		str = obj_type.Name + " : " + strconv.FormatInt(obj.Int(), 10)
	case reflect.Struct:
		str += "{"

	}

	return str, nil
}
