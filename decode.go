package json

func Decode(json_ JSONInterface) interface{} {

	if json_ == nil {
		return nil
	}

	switch json_.Type() {
	case LONG:
		return json_.(*JSONLong).Value()
	case DOUBLE:
		return json_.(*JSONDouble).Value()
	case BOOL:
		return json_.(*JSONBool).Value()
	case STRING:
		return json_.(*JSONString).Value()
	case OBJECT:
		obj := json_.(*JSONObject)
		if obj.IsNil() { return nil }
		keys := obj.Keys()
		m := make(map[string]interface{})
		for _, key := range keys {
			m[key] = Decode(obj.Get(key))
		}
		return m
	case ARRAY:
		array := json_.(*JSONArray)
		if array.IsNil() { return nil }
		a := make([]interface{}, array.Length())
		for i:=0; i<len(a);i++ {
			a[i] = Decode(array.Get(i))
		}
		return a
	default:
		return nil
	}
	return nil
}

