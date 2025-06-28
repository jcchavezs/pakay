package pakay

type Source struct {
	Type   string
	Labels []string
}

// FilterIn sources that should be considered in the secret evaluation
type FilterIn func(Source) bool
