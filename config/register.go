package config

func init() {
	RegisterCipher("AES", NewAESCipherWithOptions)
	RegisterCipher("Base64", NewBase64CipherWithOptions)
	RegisterCipher("Group", NewCipherGroupWithOptions)
	RegisterCipher("KMS", NewKMSCipherWithOptions)

	RegisterDecoder("Ini", NewIniDecoder)
	RegisterDecoder("", NewJson5Decoder)
	RegisterDecoder("Json", NewJson5Decoder)
	RegisterDecoder("Json5", NewJson5Decoder)
	RegisterDecoder("Prop", NewPropDecoder)
	RegisterDecoder("Toml", NewTomlDecoder)
	RegisterDecoder("Yaml", NewYamlDecoder)

	RegisterProvider("ACM", NewACMProviderWithOptions)
	RegisterProvider("", NewLocalProviderWithOptions)
	RegisterProvider("Local", NewLocalProviderWithOptions)
	RegisterProvider("Memory", NewMemoryProviderWithOptions)
	RegisterProvider("OTS", NewOTSProviderWithOptions)
	RegisterProvider("OTSLegacy", NewOTSLegacyProviderWithOptions)
}
