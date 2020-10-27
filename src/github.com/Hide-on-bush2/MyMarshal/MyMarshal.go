package MyMarshal

import (
	"fmt"
	"reflect"
	"strings"
)

func JsonMarshal(v interface{}) ([]byte, error) {
	json, _ := struct2str(v)
	return []byte(json), nil
}

func processTag(tag reflect.StructTag) int {
	field, ok := tag.Lookup("json")
	if !ok {
		//跳过，不用处理该字段
		return 0
	}
	if field == "-" {
		//忽略该字段
		return 1
	}
	if field == "-," {
		//该字段的键为"-"
		return 2
	}
	sep := ","
	arr := strings.Split(field, sep)

	if len(arr) != 2 {
		//该字段的键为field
		return 5
	}

	if arr[1] == "omitempty" {
		if arr[0] == "" {
			//该字段的键为field，如果值为空，省略该字段
			return 3
		} else {
			//该字段的键为arr[0]，如果值为空，省略该字段
			return 4
		}
	} else {
		//该标签无效，跳过
		return 0
	}
}

func isEmptyValue(v interface{}) bool {
	return v == 0 || v == "" || v == 0.0
}

func struct2str(v interface{}) (string, error) {
	str := ""
	obj := reflect.ValueOf(v)
	obj_type := reflect.TypeOf(v)

	switch obj_type.Kind() {
	case reflect.String:
		str = "\"" + fmt.Sprintf("%v", v) + "\""
		// return str, nil
		// fmt.Printf(str)
		// fmt.Printf("%v:%s\n", obj_type.Name, obj.String())
		// str = obj_type.Name + " : " + obj.String()
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Float32, reflect.Float64:
		str = fmt.Sprintf("%v", v)
		// return str, nil
		// fmt.Printf(str)
		// fmt.Printf("%v:%s\n", obj_type.Name, strconv.FormatInt(obj.Int(), 10))
		// str = obj_type.Name + " : " + strconv.FormatInt(obj.Int(), 10)
	case reflect.Struct:
		count := obj.NumField()
		str += "{"
		for i := 0; i < count; i++ {
			field_name := obj_type.Field(i)
			field_val := obj.Field(i)
			if field_name.PkgPath != "" {
				continue
			}
			tag := field_name.Tag
			json := tag.Get("json")
			sep := ","
			str_arr := strings.Split(json, sep)
			key := field_name.Name
			num := processTag(tag)

			switch num {
			case 1:
				continue
			case 2:
				key = "-"
			case 5:
				key = json
			case 3:
				if isEmptyValue(field_val.Interface()) {
					continue
				}
			case 4:
				if isEmptyValue(field_val.Interface()) {
					continue
				}
				key = str_arr[0]
			}

			str += key + ":"
			t_str, _ := struct2str(field_val.Interface())
			str += t_str
			// str += fmt.Sprintf("%v", field_val)
			if i != count-1 {
				str += ","
			}
			// fmt.Printf(t_str)
		}
		str = strings.TrimRight(str, ",")
		str += "}"
		// return str, nil
	case reflect.Slice:
		str += "["
		for i := 0; i < obj.Len(); i++ {
			t_str, _ := struct2str(obj.Index(i).Interface())
			str += t_str
			if i != obj.Len()-1 {
				str += ","
			}
		}
		str += "]"
	case reflect.Map:
		keys := obj.MapKeys()
		str += "{"
		for i := 0; i < len(keys); i++ {
			k := keys[i]
			v := obj.MapIndex(k)
			t_str1, _ := struct2str(k.Interface())
			t_str2, _ := struct2str(v.Interface())
			str += t_str1 + ":" + t_str2
			if i != len(keys)-1 {
				str += ","
			}
		}
		str += "}"
	}

	return str, nil
}
