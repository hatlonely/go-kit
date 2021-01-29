package logger

func init() {
	RegisterWriter("RotateFile", NewRotateFileWriterWithOptions)
	RegisterWriter("Stdout", NewStdoutWriterWithOptions)
	RegisterWriter("ElasticSearch", NewElasticSearchWriterWithOptions)
}
