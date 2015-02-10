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

func TestSecureSerializerWithEncryptedDataFromNodeSerializer(t *testing.T) {
	encrypted := "Ap8Pcq6w4OoG0A-zqFDI0vFFNkk=KdFhBDlX811f406043a74fef0ed672671c0fa6f13e71f09cba264d6662a44b18368ba8d0b40fc7cc96d3df4fdf439d7f18ceac99e8c03aa6a2b0535d1f31a31df4c7951cb12ba84fbee0cfc15b0091b86d4c02eb"

	var data [6]interface{}

	err := serializer.Parse(encrypted, &data)

	if err != nil {
		t.Error(err)
		return
	}

	if value, ok := data[0].(float64); ok {
		if value != 123 {
			t.Error("Values did not match after Stringify and Parse")
			return
		}
	} else {
		t.Error("Returned value has wrong type")
	}

	if value, ok := data[1].(string); ok {
		if value != "barfoo" {
			t.Error("Values did not match after Stringify and Parse")
			return
		}
	} else {
		t.Error("Returned value has wrong type")
	}
}
