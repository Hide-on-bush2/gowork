package MyMarshal

import "testing"

type Person struct {
	Id        int    `json:"id"` //id作为该JSON字段的key值
	FirstName string //不设置标签则默认使用FirstName为字段key值
	LastName  string `json:"-"`                //字段被本包忽略，即使有值也不输出
	Age       int    `json:"age,omitempty"`    //含omitempty选项的字段如果为空值会省略，如果存在Age作为该JSON字段的key值
	Height    int    `json:"height,omitempty"` //含omitempty选项的字段如果为空值会省略，如果存在height作为该JSON字段的key值
}

func Test_readini(t *testing.T) {
	v := Person{Id: 13, FirstName: "John", LastName: "Doe", Age: 42}

	expectedJson := "{\"id\":13,\"FirstName\":\"John\",\"age\":42}"
	json_byte, _ := JsonMarshal(v)
	realJson := string(json_byte)
	if realJson != expectedJson {
		t.Errorf("\ngot '%s' \nwant '%s'", realJson, expectedJson)
	}
}
