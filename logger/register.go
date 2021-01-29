package logger

func init() {
	RegisterWriter("RotateFile", NewRotateFileWriterWithOptions)
	RegisterWriter("Stdout", NewStdoutWriterWithOptions)
	RegisterWriter("ElasticSearch", NewElasticSearchWriterWithOptions)

	RegisterFormatter("Json", NewJsonFormatter)
	RegisterFormatter("Text", NewTextFormatterWithOptions)
}
