package wrap

type WrapperOptions struct {
	EnableTrace  bool
	EnableMetric bool
	Metric       struct {
		Buckets     []float64
		ConstLabels map[string]string
	}
}
