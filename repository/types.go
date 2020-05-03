package repository

type Sort struct {
	Name      string
	Direction string
}

type Pageable struct {
	Offset int
	Limit  int
}

// NewPageable returns a new default pageable with offset 0 and limit 20
func NewPageable() *Pageable {
	return &Pageable{
		Offset: 0,
		Limit:  20,
	}
}

type ArticleFilter struct {
	Pageable
	Tag       string
	Author    string
	Favorited string
}
