package logger

func init() {
	RegisterWriter("RotateFile", NewRotateFileWriterWithOptions)
	RegisterWriter("Stdout", NewStdoutWriterWithOptions)

	RegisterFormatter("Json", NewJsonFormatterWithOptions)
	RegisterFormatter("Text", NewTextFormatterWithOptions)
}
