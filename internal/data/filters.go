package data

type Filters struct {
	Page     int    `validate:"gt=0,max=10_000_000"`
	PageSize int    `validate:"gt=0,max=100"`
	Sort     string `validate:"oneof=id title year runtime -id -title -year -runtime"`
}
