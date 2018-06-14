package json

import "encoding/json"
import "reflect"
import "testing"

type A struct {
	A string `json:"a"`
	B int64 `json:"b"`
	C bool `json:"c"`
}

func Test1(t *testing.T) {

	check(t,  A{"test", 99, true})
	check(t, map[string]interface{}{"a": []interface{}{"a", "b", "c"}, "b": float64(1), "c": true})
	check(t , 1)
	check(t , true)
	check(t , "a")

}

func check(t *testing.T, v interface{}) {
	json_ := Encode(v)
	i := Decode(json_)
	bs, _ := json.Marshal(i)
	v2 := reflect.New(reflect.TypeOf(v))
	json.Unmarshal(bs, v2.Interface())

	if !reflect.DeepEqual(v, v2.Elem().Interface()) {
		t.Fatal("test json failed",  v, v2.Elem().Interface())
	}
}

type F struct {
	A int64 `json:"a"`
	B bool `json:"b"`
}

type B struct {
	A string `json:"a"`
	B int64 `json:"b"`
	C float64 `json:"c"`
	D bool `json:"d"`
	E []int64 `json:"e"`
	F F `json:"f"`
}

func Test2(t *testing.T) {

	type_ := EncodeType(map[string]interface{}{
				"a": "STRING",
				"b": "LONG",
				"c": "DOUBLE",
				"d": "BOOL",
				"e": []interface{}{"LONG"},
				"f": map[interface{}]interface{}{"a": "LONG", "b": "BOOL"}})

	if type_.Kind() != OBJECT ||
		type_.(JSONObjectType).Get("a").Kind() != STRING ||
		type_.(JSONObjectType).Get("b").Kind() != LONG ||
		type_.(JSONObjectType).Get("c").Kind() != DOUBLE ||
		type_.(JSONObjectType).Get("d").Kind() != BOOL ||
		type_.(JSONObjectType).Get("e").Kind() != ARRAY ||
		type_.(JSONObjectType).Get("f").Kind() != OBJECT ||
		type_.(JSONObjectType).Get("e").(JSONArrayType).GetType().Kind() != LONG ||
		type_.(JSONObjectType).Get("f").(JSONObjectType).Get("a").Kind() != LONG ||
		type_.(JSONObjectType).Get("f").(JSONObjectType).Get("b").Kind() != BOOL {
		t.Fatal("encode type failed", type_)
	}

	check2(t, type_, B{"A", 15, 15.0, true, []int64{1, 2}, F{100, true}},
		B{"A", 15, 15.0, true, []int64{1, 2}, F{100, true}})

	check2(t, type_, B{"A", 15, 15.0, true, []int64{1, 2}, F{100, true}},
		B{"A", 15, 15.0, true, []int64{1, 2}, F{100, true}})



}

func check2(t *testing.T, type_ JSONType, d interface{}, b B) {
	json_ := EncodeByType(d, type_)
	i := Decode(json_)
	bs, _ := json.Marshal(i)
	v2 := reflect.New(reflect.TypeOf(b))
	json.Unmarshal(bs, v2.Interface())

	if !reflect.DeepEqual(b, v2.Elem().Interface()) {
		t.Fatal("test json failed", b, v2.Elem().Interface())
	}
}

