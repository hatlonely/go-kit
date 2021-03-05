package logger

func init() {
	RegisterWriter("RotateFile", NewRotateFileWriterWithOptions)
	RegisterWriter("Stdout", NewStdoutWriterWithOptions)

	RegisterFormatter("Json", JsonFormatter{})
	RegisterFormatter("Text", NewTextFormatterWithOptions)
}
