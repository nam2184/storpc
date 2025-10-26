package driver

type TableRowEntity struct {
	id      uint32
	columns []any //index is .proto defined column number
}
