package pakay

type Source struct {
	Type   string
	Labels []string
}

// FilterIn sources that should be considered in the secret evaluation
type FilterIn[T any] func(T) bool
