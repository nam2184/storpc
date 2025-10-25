package types

type PageID uint32
type PageType uint8
type PageContentType uint8

const (
	InMemory PageType = 0
	InDisk   PageType = 1
	InMap    PageType = 2
)

const (
	Index    PageContentType = 0
	Branch   PageContentType = 1
	Leaf     PageContentType = 2
	Metadata PageContentType = 3
)

type BTree interface {
	Insert(key uint32) error
	Search(id uint32) PageNode
}

type Pager interface {
	ReadPage(id PageID) (*PageNode, error)
	WritePage(node *PageNode) error
	AllocatePage() PageID
	Type() PageType
}

type PageNode interface {
	Header() PageHeader
	Content() PageContent
	ID() PageID
	Size() uint8
	Write() error
}

type PageHeader interface {
	Size() uint8
	Type() PageContentType
	Read() error
	Write() error
}

type PageContent interface {
	Next() uint8
	Read() error
	Write() error
}
