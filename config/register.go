package config

func init() {
	RegisterCipher("", NewEmptyCipher)
	RegisterCipher("Empty", NewEmptyCipher)
	RegisterCipher("AES", NewAESCipherWithOptions)
	RegisterCipher("Base64", NewBase64CipherWithOptions)
	RegisterCipher("Group", NewCipherGroupWithOptions)
	RegisterCipher("KMS", NewKMSCipherWithOptions)

	RegisterDecoder("", NewJson5Decoder)
	RegisterDecoder("Json", NewJson5Decoder)
	RegisterDecoder("Json5", NewJson5Decoder)
	RegisterDecoder("Prop", NewPropDecoder)
	RegisterDecoder("Toml", NewTomlDecoder)
	RegisterDecoder("Yaml", NewYamlDecoder)
	RegisterDecoder("Ini", NewIniDecoder)

	RegisterProvider("", NewLocalProviderWithOptions)
	RegisterProvider("Local", NewLocalProviderWithOptions)
	RegisterProvider("Memory", NewMemoryProviderWithOptions)
	RegisterProvider("OTS", NewOTSProviderWithOptions)
	RegisterProvider("OTSLegacy", NewOTSLegacyProviderWithOptions)
	RegisterProvider("ACM", NewACMProviderWithOptions)
}
