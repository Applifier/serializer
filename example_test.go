package serializer

import "fmt"

func ExampleNewSecureSerializer() {
	serializer := NewSecureSerializer([]byte("somesecretkey"), []byte("anothersecretstring"))

	data := map[string]string{
		"foo": "bar",
	}

	encryptedData, _ := serializer.Stringify(data)

	var returnedData map[string]string

	serializer.Parse(encryptedData, &returnedData)

	fmt.Printf("%v\n", returnedData)
	// Output: map[foo:bar]
}
