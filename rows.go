package memless

// Rows is a query result: the column headers and the rows of native values.
// A cell is int64, float64, string, bool, or nil for an absent value.
type Rows struct {
	Columns []string
	Rows    [][]any
}
