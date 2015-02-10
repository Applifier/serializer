package serializer

import (
	"reflect"
	"testing"
)

var serializer = NewSecureSerializer([]byte("somesecretkey"), []byte("anothersecretstring"))

func TestSecureSerializerWithInt(t *testing.T) {
	val := float64(135)

	result, err := serializer.Stringify(val)

	if err != nil {
		t.Error(err)
		return
	}

	var back interface{}

	err = serializer.Parse(result, &back)

	if err != nil {
		t.Error(err)
		return
	}

	if _, ok := back.(float64); ok {
		if val != back {
			t.Error("Values did not match after Stringify and Parse")
			return
		}
	} else {
		t.Error("Returned value has wrong type")
	}
}

func TestSecureSerializerWithMap(t *testing.T) {
	val := make(map[string]interface{})
	val["foo"] = "bar"

	result, err := serializer.Stringify(val)

	if err != nil {
		t.Error(err)
		return
	}

	var back interface{}

	err = serializer.Parse(result, &back)

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(val, back) {
		t.Error("Values did not match after Stringify and Parse")
	}

}

func TestSecureSerializerWithArray(t *testing.T) {
	var val [2]interface{}
	val[0] = "bar"
	val[1] = "foo"

	result, err := serializer.Stringify(val)

	if err != nil {
		t.Error(err)
		return
	}

	var back [2]interface{}

	err = serializer.Parse(result, &back)

	if err != nil {
		t.Error(err)
		return
	}

	if !reflect.DeepEqual(val, back) {
		t.Error("Values did not match after Stringify and Parse")
	}
}
