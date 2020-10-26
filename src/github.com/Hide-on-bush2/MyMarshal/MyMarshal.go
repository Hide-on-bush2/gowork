package MyMarshal

import (
	"fmt"
	"reflect"
)

func JsonMarshal(v interface{}) ([]byte, error) {
	json, _ := struct2str(v)
	return []byte(json), nil
}

func struct2str(v interface{}) (string, error) {
	str := ""
	obj := reflect.ValueOf(v)
	obj_type := reflect.TypeOf(v)
	count := obj.NumField()

	switch obj_type.Kind() {
	case reflect.String:
		str = fmt.Sprintf("%v", v)
		// fmt.Printf(str)
		// fmt.Printf("%v:%s\n", obj_type.Name, obj.String())
		// str = obj_type.Name + " : " + obj.String()
	case reflect.Int:
		str = fmt.Sprintf("%v", v)
		// fmt.Printf(str)
		// fmt.Printf("%v:%s\n", obj_type.Name, strconv.FormatInt(obj.Int(), 10))
		// str = obj_type.Name + " : " + strconv.FormatInt(obj.Int(), 10)
	case reflect.Struct:
		str += "{"
		for i := 0; i < count; i++ {
			// field := obj_type.Field(i)
			val := obj.Field(i)
			str += fmt.Sprintf("%v", val)
			// fmt.Printf(t_str)
		}
		str += "}"
	}

	return str, nil
}
