package logger

func init() {
	RegisterWriter("RotateFile", NewRotateFileWriterWithOptions)
	RegisterWriter("Stdout", NewStdoutWriterWithOptions)

	RegisterFormatter("Json", NewJsonFormatter)
	RegisterFormatter("Text", NewTextFormatterWithOptions)
}
