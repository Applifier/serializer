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

func TestParseShouldReturnErrorOnEmptyString(t *testing.T) {
	err := serializer.Parse("", nil)

	if err == nil {
		t.Error("No error returned")
		return
	}
}

func TestParseCyptoBlocksRecover(t *testing.T) {
	badToken := `LwTUfzIlfef1395e03230974824b6643ed3ace515d5090a4498b8a9cdeefb6157a92b5e3b0bd8a4f6711e41a736c0589d42d84a230b139e3aef755da35315d06ee20d03a89aecec03f0461a68c3794bcffba89150a3920af5c6901fcbd72d0fd360abb500fbc102db9f1bfef939c5e4df49fab9f6ae765fcefe8d5c9dc0616e131ee253fa210595c05146a9cc66bc3f359f2f1fa513723015234f632f30f3f7a9a5df1e97c0fa3c202d0bc38b5ccad02219fc0d71104a2824825506ffe21bd8c380f87a74baca2e47f35ff58dd35ba22b3bab627661d53a38b6a321c421a0208d47a0b6b3532a2bee671e480f6f9acd8fe014ad563653cde79c6b53af4ac738ce90ff44a752c92b0b5144eecc7464031972a1c0d1d36289a3b3da902d17eaff6fdabf611aeead3485ba83f1f79884a91f2a911cd`
	err := serializer.Parse(badToken, nil)

	if err == nil {
		t.Error("No error recovered & returned")
	}
}

var data = map[string]interface{}{
	"foo": "bar",
}

var encryptedData, _ = serializer.Stringify(data)

func BenchmarkStringify(b *testing.B) {

	for n := 0; n < b.N; n++ {
		serializer.Stringify(data)
	}
}

func BenchmarkParse(b *testing.B) {
	for n := 0; n < b.N; n++ {
		var data map[string]interface{}
		err := serializer.Parse(encryptedData, &data)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkParseParallel(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			var data map[string]interface{}
			err := serializer.Parse(encryptedData, &data)
			if err != nil {
				b.Error(err)
			}
		}
	})
}

func BenchmarkStringifyParse(b *testing.B) {
	val := make(map[string]interface{})
	val["foo"] = "bar"
	var data map[string]interface{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		encrypted, _ := serializer.Stringify(val)
		err := serializer.Parse(encrypted, &data)
		if err != nil {
			b.Error(err)
		}
	}
}
