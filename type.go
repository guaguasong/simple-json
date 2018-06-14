package json

import (
	"reflect"
)

const (
	NULL = iota
	LONG
	DOUBLE
	BOOL
	STRING
	OBJECT
	ARRAY
)

type JSONType interface {
	Kind() int
}

type JSONLongType struct {}
func (r JSONLongType) Kind() int { return LONG }

type JSONDoubleType struct {}
func (r JSONDoubleType) Kind() int { return DOUBLE }

type JSONBoolType struct {}
func (r JSONBoolType) Kind() int { return BOOL }

type JSONStringType struct {}
func (r JSONStringType) Kind() int { return STRING }

type JSONObjectType struct {
	m map[string]JSONType
}
func (r JSONObjectType) Kind() int { return OBJECT }

func (r JSONObjectType) Get(name string) JSONType {
	return r.m[name]
}

func (r JSONObjectType) Set(name string, type_ JSONType) {
	if r.m == nil { r.m = make(map[string]JSONType) }
	r.m[name] = type_
}

func (r JSONObjectType) Keys() []string {
	keys := make([]string, len(r.m))
	i := 0
	for k, _ := range r.m {
		keys[i] = k
		i++
	}
	return keys
}

type JSONArrayType struct {
	t JSONType
}
func (r JSONArrayType) Kind() int { return ARRAY }

func (r JSONArrayType) GetType() JSONType { return r.t }

func EncodeType(d interface{}) JSONType {
	return encodeType(reflect.ValueOf(d))
}

func encodeType(v reflect.Value) JSONType {
	switch v.Kind() {
	case reflect.Interface:
		return encodeType(v.Elem())
	case reflect.String:
		switch v.String() {
		case "LONG": return JSONLongType{}
		case "DOUBLE": return JSONDoubleType{}
		case "BOOL": return JSONBoolType{}
		case "STRING": return JSONStringType{}
		}
		return nil
	case reflect.Map:
		m := make(map[string]JSONType)
		vs := v.MapKeys()
		for _, k := range vs {
			key := ""
			if k.Kind() == reflect.String {
				key = k.String()
			} else if k.Kind() == reflect.Interface && k.Elem().Kind() == reflect.String {
				key = k.Elem().String()
			} else {
				continue
			}
			m[key] = encodeType(v.MapIndex(k))
		}
		return JSONObjectType{m}
	case reflect.Slice, reflect.Array:
		n := v.Len()
		for i := 0; i < n; i++ {
			return JSONArrayType{encodeType(v.Index(i))}
		}
		return nil
	}
	return nil
}

////////////////////////////////////////////////////////////////////////////////

type JSONInterface interface {
	Type() int
}

//----------------------------------------------------------------------------//

type JSONObject struct {
	value map[string]JSONInterface
}

func NewJSONObject(value map[string]JSONInterface) *JSONObject {
	return &JSONObject{value}
}

func (r *JSONObject) Value() map[string]interface{} {
	return nil
}

func (r *JSONObject) Reset(value map[string]JSONInterface) {
	r.value = value
}

func (r *JSONObject) Get(name string) JSONInterface {
	return r.value[name]
}

func (r *JSONObject) Set(name string, value JSONInterface) {
	r.value[name] = value
}

func (r *JSONObject) IsNil() bool {
	return r.value == nil
}

func (r *JSONObject) Keys() []string {
	keys := make([]string, len(r.value))
	i := 0
	for key, _ := range r.value {
		keys[i] = key
		i++
	}
	return keys
}

func (r *JSONObject) Type() int {
	return OBJECT
}

func (r *JSONObject) ToString() string {
	return ""
}

//----------------------------------------------------------------------------//

type JSONArray struct {
	value []JSONInterface
}

func NewJSONArray(value []JSONInterface) *JSONArray {
	return &JSONArray{value}
}

func (r *JSONArray) Value() interface{} {
	return nil
}

func (r *JSONArray) Reset(value []JSONInterface) {
	r.value = value
}

func (r *JSONArray) Get(index int) JSONInterface {
	return r.value[index]
}

func (r *JSONArray) Set(index int, value JSONInterface) {
	r.value[index] = value
}

func (r *JSONArray) IsNil() bool {
	return r.value == nil
}

func (r *JSONArray) Length() int {
	return len(r.value)
}

func (r *JSONArray) Type() int {
	return ARRAY
}

func (r *JSONArray) ToString() string {
	return ""
}

