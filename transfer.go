package json

func C(d interface{}, t JSONType) interface{} {
	return Decode(EncodeByType(d, t))
}

