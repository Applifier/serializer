package serializer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"strings"
)

type SecureSerializer struct {
	EncryptKey  []byte
	ValidateKey []byte
}

func md5sum(d []byte) []byte {
	h := md5.New()
	h.Write(d)
	return h.Sum(nil)
}

func evpBytesToKey(password string, keyLen int) (key []byte) {
	const md5Len = 16

	cnt := (keyLen-1)/md5Len + 1
	m := make([]byte, cnt*md5Len)
	copy(m, md5sum([]byte(password)))

	// Repeatedly call md5 until bytes generated is enough.
	// Each call to md5 uses data: prev md5 sum + password.
	d := make([]byte, md5Len+len(password))
	start := 0
	for i := 1; i < cnt; i++ {
		start += md5Len
		copy(d, m[start-md5Len:start])
		copy(d[md5Len:], password)
		copy(m[start:], md5sum(d))
	}
	return m[:keyLen]
}

func randString(n int) []byte {
	const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	symbols := big.NewInt(int64(len(alphanum)))
	states := big.NewInt(0)
	states.Exp(symbols, big.NewInt(int64(n)), nil)
	r, err := rand.Int(rand.Reader, states)
	if err != nil {
		panic(err)
	}
	var bytes = make([]byte, n)
	r2 := big.NewInt(0)
	symbol := big.NewInt(0)
	for i := range bytes {
		r2.DivMod(r, symbols, symbol)
		r, r2 = r2, r
		bytes[i] = alphanum[symbol.Int64()]
	}
	return bytes
}

func sign(data, key []byte) []byte {
	mac := hmac.New(sha1.New, key)
	mac.Write(data)
	return mac.Sum(nil)
}

func (serializer *SecureSerializer) Stringify(obj interface{}) (string, error) {
	jsonData, err := json.Marshal(obj)

	if err != nil {
		return "", err
	}

	nonceCheck := randString(8)
	if err != nil {
		return "", err
	}

	nonceCrypt := randString(8)
	if err != nil {
		return "", err
	}

	password := append(serializer.EncryptKey, nonceCrypt[:]...)
	key := evpBytesToKey(string(password), 48)

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return "", err
	}

	iv := key[32:]

	encrypter := cipher.NewCFBEncrypter(block, iv)

	encrypted := make([]byte, len(nonceCheck)+len(jsonData))
	encrypter.XORKeyStream(encrypted, append(nonceCheck, jsonData[:]...))

	digest := sign(jsonData, append(serializer.ValidateKey, nonceCheck[:]...))

	digestBase64 := base64.StdEncoding.EncodeToString(digest)
	digestBase64 = strings.Replace(digestBase64, "+", "-", -1)
	digestBase64 = strings.Replace(digestBase64, "/", "_", -1)

	return fmt.Sprint(
		digestBase64,
		string(nonceCrypt),
		hex.EncodeToString(encrypted)), nil
}

func (serializer *SecureSerializer) Parse(base64data string, obj interface{}) error {
	expectedDigest := base64data[0:28]
	nonceCrypt := base64data[28:36]
	encryptedDataHex := base64data[36:]

	password := append(serializer.EncryptKey, nonceCrypt[:]...)
	key := evpBytesToKey(string(password), 48)

	block, err := aes.NewCipher(key[:32])
	if err != nil {
		return err
	}

	iv := key[32:]

	decrypter := cipher.NewCFBDecrypter(block, iv)

	encryptedData, err := hex.DecodeString(encryptedDataHex)
	if err != nil {
		return err
	}

	decrypted := make([]byte, len(encryptedData))
	decrypter.XORKeyStream(decrypted, encryptedData)

	nonceCheck := decrypted[0:8]

	digest := sign(decrypted[8:], append(serializer.ValidateKey, nonceCheck[:]...))
	digestBase64 := base64.StdEncoding.EncodeToString(digest)
	digestBase64 = strings.Replace(digestBase64, "+", "-", -1)
	digestBase64 = strings.Replace(digestBase64, "/", "_", -1)

	if !strings.EqualFold(digestBase64, expectedDigest) {
		return errors.New("Bad digest")
	}

	return json.Unmarshal(decrypted[8:], obj)
}

func NewSecureSerializer(encryptKey []byte, validateKey []byte) *SecureSerializer {
	serializer := &SecureSerializer{encryptKey, validateKey}
	return serializer
}

func SecureStringify(obj interface{}, encryptKey []byte, validateKey []byte) (string, error) {
	serializer := NewSecureSerializer(encryptKey, validateKey)
	return serializer.Stringify(obj)
}

func SecureParse(data string, obj interface{}, encryptKey []byte, validateKey []byte) error {
	serializer := NewSecureSerializer(encryptKey, validateKey)
	return serializer.Parse(data, obj)
}
