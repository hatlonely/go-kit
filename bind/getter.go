package bind

type Getter interface {
	Get(key string) (value interface{}, exists bool)
}
