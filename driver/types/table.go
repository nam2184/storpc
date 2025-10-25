package types

type TableEntity interface {
	//Name of message
	Name() string

	//Name of package of .proto
	Package() string
}
