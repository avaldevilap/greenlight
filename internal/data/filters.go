package data

import "strings"

type Filters struct {
	Page     int    `validate:"gt=0,max=10_000_000"`
	PageSize int    `validate:"gt=0,max=100"`
	Sort     string `validate:"oneof=id title year runtime -id -title -year -runtime"`
}

func (f Filters) sortColumn() string {
	return strings.TrimPrefix(f.Sort, "-")
}

// Return the sort direction ("ASC" or "DESC") depending on the prefix character of the
// Sort field.
func (f Filters) sortDirection() string {
	if strings.HasPrefix(f.Sort, "-") {
		return "DESC"
	}
	return "ASC"
}
