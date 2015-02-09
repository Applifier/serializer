package serializer

import "testing"

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
