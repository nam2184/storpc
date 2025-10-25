package driver

import (
	"fmt"

	"github.com/nam2184/storpc/driver/types"
)

// ------------------------
// BTree
// ------------------------

type MemoryBTree struct {
	root  types.PageNode // root node
	t     int            // minimum degree
	size  int            // total number of keys
	pager *MemoryPager
}

func NewMemoryBTree(t int) *MemoryBTree {
	pager := NewMemoryPager()
	rootID := pager.AllocatePage()
	root := &MemoryPage{
		id:   rootID,
		keys: []int{},
		leaf: true,
	}
	pager.WritePage(root)
	return &MemoryBTree{
		root:  root,
		t:     t,
		size:  0,
		pager: pager,
	}
}

func (bt *MemoryBTree) Insert(key int) {
	fmt.Println("Insert called on MemoryBTree with key:", key)
	// TODO: implement split/insert logic

	bt.size++
}

func (bt *MemoryBTree) Search(key int) bool {
	fmt.Println("Search called on MemoryBTree for key:", key)
	return false
}

type DiskBTree struct {
	root  types.PageNode
	t     int
	size  int
	pager *DiskPager
}

func NewDiskBTree(t int) *DiskBTree {
	pager := NewDiskPager()
	rootID := pager.AllocatePage()
	root := &DiskPage{
		id:   rootID,
		keys: []int{},
		leaf: true,
	}
	pager.WritePage(root)
	return &DiskBTree{
		root:  root,
		t:     t,
		size:  0,
		pager: pager,
	}
}

func (bt *DiskBTree) Insert(key int) {
	fmt.Println("Insert called on DiskBTree with key:", key)
	// TODO: implement split/insert logic
	bt.size++
}

func (bt *DiskBTree) Search(key int) bool {
	fmt.Println("Search called on DiskBTree for key:", key)
	return false
}

// ------------------------
// Page Node
// ------------------------

type MemoryPage struct {
	id       types.PageID
	keys     []int
	children []*MemoryPage
	leaf     bool
}

func (p *MemoryPage) ID() types.PageID {
	return p.id
}

func (p *MemoryPage) Header() types.PageHeader {
	return &PageHeader{}
}

func (p *MemoryPage) Content() types.PageContent {
	return &PageContent{}
}

func (p *MemoryPage) Size() uint8 {
	return uint8(len(p.keys))
}

func (p *MemoryPage) Write() error {
	// placeholder: in-memory write is trivial
	fmt.Println("MemoryPage written:", p.id)
	return nil
}

type DiskPage struct {
	id       types.PageID
	keys     []int
	children []*DiskPage
	leaf     bool
}

func (p *DiskPage) ID() types.PageID {
	return p.id
}

func (p *DiskPage) Header() types.PageHeader {
	return &PageHeader{}
}

func (p *DiskPage) Content() types.PageContent {
	return &PageContent{}
}

func (p *DiskPage) Size() uint8 {
	return uint8(len(p.keys))
}

func (p *DiskPage) Write() error {
	// placeholder: disk write simulation
	fmt.Println("DiskPage written:", p.id)
	return nil
}

// ------------------------
// Pagers
// ------------------------

type MemoryPager struct {
	nextID types.PageID
	pages  map[types.PageID]*MemoryPage
}

func NewMemoryPager() *MemoryPager {
	return &MemoryPager{
		nextID: 1,
		pages:  make(map[types.PageID]*MemoryPage),
	}
}

func (mp *MemoryPager) ReadPage(id types.PageID) (*MemoryPage, error) {
	page, ok := mp.pages[id]
	if !ok {
		return nil, fmt.Errorf("page %d not found", id)
	}
	return page, nil
}

func (mp *MemoryPager) WritePage(node *MemoryPage) error {
	mp.pages[node.ID()] = node
	return nil
}

func (mp *MemoryPager) AllocatePage() types.PageID {
	id := mp.nextID
	mp.nextID++
	return id
}

func (mp *MemoryPager) Type() types.PageType {
	return types.InMemory
}

type DiskPager struct {
	nextID types.PageID
	pages  map[types.PageID]*DiskPage
}

func NewDiskPager() *DiskPager {
	return &DiskPager{
		nextID: 1,
		pages:  make(map[types.PageID]*DiskPage),
	}
}

func (dp *DiskPager) ReadPage(id types.PageID) (*DiskPage, error) {
	page, ok := dp.pages[id]
	if !ok {
		return nil, fmt.Errorf("page %d not found", id)
	}
	return page, nil
}

func (dp *DiskPager) WritePage(node *DiskPage) error {
	dp.pages[node.ID()] = node
	return nil
}

func (dp *DiskPager) AllocatePage() types.PageID {
	id := dp.nextID
	dp.nextID++
	return id
}

func (dp *DiskPager) Type() types.PageType {
	return types.InDisk
}
