package config

// CipherGroup is also an cipher
type CipherGroup struct {
	ciphers []Cipher
}

func NewCipherGroupWithOptions(options *CipherGroupOptions) (*CipherGroup, error) {
	var ciphers []Cipher
	for _, opt := range *options {
		cipher, err := NewCipherWithOptions(&opt)
		if err != nil {
			return nil, err
		}
		ciphers = append(ciphers, cipher)
	}
	return NewCipherGroup(ciphers...), nil
}

func NewCipherGroup(ciphers ...Cipher) *CipherGroup {
	return &CipherGroup{
		ciphers: ciphers,
	}
}

func (c *CipherGroup) Encrypt(textToEncrypt []byte) ([]byte, error) {
	buf := textToEncrypt
	var err error
	for _, cipher := range c.ciphers {
		buf, err = cipher.Encrypt(buf)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

func (c *CipherGroup) Decrypt(textToDecrypt []byte) ([]byte, error) {
	buf := textToDecrypt
	var err error
	for i := len(c.ciphers) - 1; i >= 0; i-- {
		buf, err = c.ciphers[i].Decrypt(buf)
		if err != nil {
			return nil, err
		}
	}
	return buf, nil
}

type CipherGroupOptions []CipherOptions
