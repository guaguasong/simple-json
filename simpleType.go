package json

import (
	"strconv"
)

//----------------------------------------------------------------------------//

type JSONLong struct {
	value int64
}

func NewJSONLong(value int64) *JSONLong {
	return &JSONLong{value}
}

func (r *JSONLong) Value() int64 {
	return r.value
}

func (r *JSONLong) Set(value int64) {
	r.value = value
}

func (r *JSONLong) Type() int {
	return LONG
}

func (r *JSONLong) ToString() string {
	return strconv.FormatInt(r.value, 10)
}

//----------------------------------------------------------------------------//

type JSONDouble struct {
	value float64
}

func NewJSONDouble(value float64) *JSONDouble {
	return &JSONDouble{value}
}

func (r *JSONDouble) Value() float64 {
	return r.value
}

func (r *JSONDouble) Set(value float64) {
	r.value = value
}

func (r *JSONDouble) Type() int {
	return DOUBLE
}

func (r *JSONDouble) ToString() string {
	return strconv.FormatFloat(r.value, 'f', -1, 64)
}

//----------------------------------------------------------------------------//

type JSONBool struct {
	value bool
}

func NewJSONBool(value bool) *JSONBool {
	return &JSONBool{value}
}

func (r *JSONBool) Value() bool {
	return r.value
}

func (r *JSONBool) Set(value bool) {
	r.value = value
}

func (r *JSONBool) Type() int {
	return BOOL
}

func (r *JSONBool) ToString() string {
	return strconv.FormatBool(r.value)
}

//----------------------------------------------------------------------------//

type JSONString struct {
	value string
}

func NewJSONString(value string) *JSONString {
	return &JSONString{value}
}

func (r *JSONString) Value() string {
	return r.value
}

func (r *JSONString) Set(value string) {
	r.value = value
}

func (r *JSONString) Type() int {
	return STRING
}

func (r *JSONString) ToString() string {
	return r.value
}

