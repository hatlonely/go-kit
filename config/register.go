package config

func init() {
	RegisterCipher("", EmptyCipher{})
	RegisterCipher("Empty", EmptyCipher{})
	RegisterCipher("AES", NewAESCipherWithOptions)
	RegisterCipher("Base64", NewBase64CipherWithOptions)
	RegisterCipher("Group", NewCipherGroupWithOptions)
	RegisterCipher("KMS", NewKMSCipherWithOptions)

	RegisterDecoder("", Json5Decoder{})
	RegisterDecoder("Json", Json5Decoder{})
	RegisterDecoder("Json5", Json5Decoder{})
	RegisterDecoder("Prop", PropDecoder{})
	RegisterDecoder("Toml", TomlDecoder{})
	RegisterDecoder("Yaml", YamlDecoder{})
	RegisterDecoder("Ini", IniDecoder{})

	RegisterProvider("", NewLocalProviderWithOptions)
	RegisterProvider("Local", NewLocalProviderWithOptions)
	RegisterProvider("Memory", NewMemoryProviderWithOptions)
	RegisterProvider("OTS", NewOTSProviderWithOptions)
	RegisterProvider("OTSLegacy", NewOTSLegacyProviderWithOptions)
	RegisterProvider("ACM", NewACMProviderWithOptions)
}
