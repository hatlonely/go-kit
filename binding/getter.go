package binding

type Getter interface {
	Get(key string) (value interface{}, exists bool)
}
