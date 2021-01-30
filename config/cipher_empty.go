package config

func NewEmptyCipher() EmptyCipher {
	return EmptyCipher{}
}

type EmptyCipher struct{}

func (c EmptyCipher) Encrypt(textToEncrypt []byte) ([]byte, error) {
	return textToEncrypt, nil
}

func (c EmptyCipher) Decrypt(textToDecrypt []byte) ([]byte, error) {
	return textToDecrypt, nil
}
