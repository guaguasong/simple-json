package json

import (
	"math"
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func Encode(data interface{}) JSONInterface {
	return encode(reflect.ValueOf(data))
}

type encodeField struct {
	i         int // field index in struct
	tag       string
	quoted    bool
	omitEmpty bool
}

func isValidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch c {
		case '$', '-', '_', '/', '%':
			// Acceptable
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
		}
	}
	return true
}

type tagOptions string

func parseTag(tag string) (string, tagOptions) {
	if idx := strings.Index(tag, ","); idx != -1 {
		return tag[:idx], tagOptions(tag[idx+1:])
	}
	return tag, tagOptions("")
}

func (o tagOptions) Contains(optionName string) bool {
	if len(o) == 0 {
		return false
	}
	s := string(o)
	for s != "" {
		var next string
		i := strings.Index(s, ",")
		if i >= 0 {
			s, next = s[:i], s[i+1:]
		}
		if s == optionName {
			return true
		}
		s = next
	}
	return false
}

// encodeFields returns a slice of encodeField for a given
// struct type.
func encodeFields(t reflect.Type) (fs []encodeField) {
	n := t.NumField()
	for i := 0; i < n; i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			continue
		}
		if f.Anonymous {
			// We want to do a better job with these later,
			// so for now pretend they don't exist.
			continue
		}
		var ef encodeField
		ef.i = i
		ef.tag = f.Name

		tv := f.Tag.Get("json")
		if tv != "" {
			if tv == "-" {
				continue
			}
			name, opts := parseTag(tv)
			if isValidTag(name) {
				ef.tag = name
			}
			ef.omitEmpty = opts.Contains("omitempty")
			ef.quoted = opts.Contains("string")
		}
		fs = append(fs, ef)
	}
	return fs
}

func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return v.Len() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return false
}

func encode(v reflect.Value) JSONInterface {
	if !v.IsValid() {
		return nil
	}

	switch v.Kind() {
	case reflect.Bool:
		return NewJSONBool(v.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return NewJSONLong(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return NewJSONLong(int64(v.Uint()))
	case reflect.Float32, reflect.Float64:
		return NewJSONDouble(v.Float())
	case reflect.String:
		return NewJSONString(v.String())
	case reflect.Struct:
		m := make(map[string]JSONInterface)
		for _, ef := range encodeFields(v.Type()) {
			fieldValue := v.Field(ef.i)
			if ef.omitEmpty && isEmptyValue(fieldValue) {
				continue
			}
			m[ef.tag] = encode(fieldValue)
		}
		return NewJSONObject(m)
	case reflect.Map:
		if v.IsNil() {
			return nil
		}
		m := make(map[string]JSONInterface)
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
			m[key] = encode(v.MapIndex(k))
		}
		return NewJSONObject(m)
	case reflect.Slice:
		if v.IsNil() {
			return nil
		}
		fallthrough
	case reflect.Array:
		n := v.Len()
		a := make([]JSONInterface, n)
		for i := 0; i < n; i++ {
			a[i] = encode(v.Index(i))
		}
		return NewJSONArray(a)
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return encode(v.Elem())
	default:
		return nil
	}

	return nil
}

func EncodeByType(d interface{}, t JSONType) JSONInterface {
	return encodeByType(reflect.ValueOf(d), t)
}

func encodeByType(v reflect.Value, t JSONType) JSONInterface {
	if !v.IsValid() || t == nil {
		return nil
	}

	switch v.Kind() {
	case reflect.Interface, reflect.Ptr:
		if v.IsNil() {
			return nil
		}
		return encodeByType(v.Elem(), t)
	}

	switch t.Kind() {
	case BOOL:
		switch v.Kind() {
		case reflect.Bool:
			return NewJSONBool(v.Bool())
		case reflect.String:
			b, err := strconv.ParseBool(v.String())
			if err != nil { return nil }
			return NewJSONBool(b)
		}
		return nil
	case LONG:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewJSONLong(v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return NewJSONLong(int64(v.Uint()))
		case reflect.Float32, reflect.Float64:
			return NewJSONLong(int64(math.Floor(v.Float())))
		case reflect.String:
			l, err := strconv.ParseInt(v.String(), 10, 64)
			if err != nil { return nil }
			return NewJSONLong(l)
		}
		return nil
	case DOUBLE:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewJSONDouble(float64(v.Int()))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return NewJSONDouble(float64(v.Uint()))
		case reflect.Float32, reflect.Float64:
			return NewJSONDouble(v.Float())
		case reflect.String:
			f, err := strconv.ParseFloat(v.String(), 10)
			if err != nil { return nil }
			return NewJSONDouble(f)
		}
		return nil
	case STRING:
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return NewJSONString(strconv.FormatInt(v.Int(), 10))
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return NewJSONString(strconv.FormatUint(v.Uint(), 10))
		case reflect.Float32, reflect.Float64:
			return NewJSONString(strconv.FormatFloat(v.Float(), 'f', -1, 10))
		case reflect.String:
			return NewJSONString(v.String())
		}
		return nil

	case OBJECT:
		switch v.Kind() {

		case reflect.Struct:
			m := make(map[string]JSONInterface)
			fields := encodeFields(v.Type())
			newFields := make([]encodeField, 0)
			for _, ef := range fields {
				fieldValue := v.Field(ef.i)
				if ef.omitEmpty && isEmptyValue(fieldValue) {
					continue
				}
				newFields = append(newFields, ef)
			}
			for _, k := range t.(JSONObjectType).Keys() {
				var i int = -1
				for _, ef := range newFields {
					if k == ef.tag {
						i = ef.i
						break
					}
				}
				if i == -1 {
					m[k] = nil
				} else {
					m[k] = encodeByType(v.Field(i), t.(JSONObjectType).Get(k))
				}
			}
			return NewJSONObject(m)

		case reflect.Map:
			if v.IsNil() {
				return nil
			}
			m := make(map[string]JSONInterface)
			vs := v.MapKeys()
			for _, k := range t.(JSONObjectType).Keys() {
				var ok bool = false
				for _, k_ := range vs {
					if k_.Kind() == reflect.String {
						if k == k_.String() {
							m[k] = encodeByType(v.MapIndex(k_), t.(JSONObjectType).Get(k))
							ok = true
							break
						}
					} else if k_.Kind() == reflect.Interface && k_.Elem().Kind() == reflect.String {
						if k == k_.Elem().String() {
							m[k] = encodeByType(v.MapIndex(k_), t.(JSONObjectType).Get(k))
							ok = true
							break
						}
					}
				}
				if !ok {
					m[k] = nil
				}
			}
			return NewJSONObject(m)
		}
		return nil

	case ARRAY:
		switch v.Kind() {
		case reflect.Slice:
			if v.IsNil() {
				return nil
			}
			fallthrough
		case reflect.Array:
			n := v.Len()
			a := make([]JSONInterface, n)
			for i := 0; i < n; i++ {
				a[i] = encodeByType(v.Index(i), t.(JSONArrayType).GetType())
			}
			return NewJSONArray(a)
		}
		return nil
	default:
		return nil
	}

	return nil
}

