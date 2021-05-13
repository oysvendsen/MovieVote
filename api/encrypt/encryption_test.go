package encrypt

import "testing"

func TestEncryption(t *testing.T) {
	secret := "This is my secret: <3"
	key1 := createKey(t)
	key2 := createKey(t)


	encryptedSecret, err := key1.Encrypt(secret, key2.encodePub())
	if err != nil {
		t.Logf("Encryption failed with error; %s", err.Error())
		t.FailNow()
	}
	assertNotEqual(t, secret, encryptedSecret)

	decryptedSecret, err2 := key2.Decrypt(encryptedSecret)
	if err2 != nil {
		t.Logf("Decryption failed with error; %s", err2.Error())
		t.FailNow()
	}
	assertEqual(t, secret, decryptedSecret)
}

func assertNotEqual(t *testing.T, s1 string, s2 string) {
	if s1 == s2 {
		t.Logf("Strings should not be equal: %s and %s", s1, s2)
		t.FailNow()
	}
}
func assertEqual(t *testing.T, s1 string, s2 string) {
	if s1 != s2 {
		t.Logf("Strings should be equal: %s and %s", s1, s2)
		t.FailNow()
	}
}

func createKey(t *testing.T) Encryption {
	key, err := New()
	if err != nil {
		t.Log("Error when creating key;" + err.Error())
		t.FailNow()
	}
	return key
}
