package Models

type MongoPagination struct {
	Limit *int64
	Skip  *int64
}

type PostgresPagination struct {
	Limit string
	Skip  string
}
