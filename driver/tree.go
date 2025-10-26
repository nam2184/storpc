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
	m     int            // max keys
	size  int            // total number of keys
	pager *MemoryPager
}

func NewMemoryBTree(m int) *MemoryBTree {
	pager := NewMemoryPager()
	rootID := pager.AllocatePage()
	root := &MemoryPage{
		id:   rootID,
		keys: []types.PageContent{},
		leaf: true,
	}
	pager.WritePage(root)
	return &MemoryBTree{
		root:  root,
		m:     m,
		size:  0,
		pager: pager,
	}
}

func (bt *MemoryBTree) Root() types.PageNode {
	return bt.root
}

func (bt *MemoryBTree) Insert(key uint32) {
	fmt.Println("Insert called on MemoryBTree with key:", key)

	bt.size++
}

func (bt *MemoryBTree) Delete(key uint32) {
	fmt.Println("Delete called on MemoryBTree with key:", key)
	// TODO: implement split/insert logic

	bt.size--
}

func (bt *MemoryBTree) Search(key uint32) bool {
	fmt.Println("Search called on MemoryBTree for key:", key)
	return false
}

func (bt *MemoryBTree) Traverse(fn func(types.PageNode)) {
	fmt.Println("Traverse called")
}

func (bt *MemoryBTree) Balance() error {
	fmt.Println("Balance called")
	return nil
}

func (bt *MemoryBTree) Height() int {
	return 0
}

func (bt *MemoryBTree) Size() int {
	return 0
}

type DiskBTree struct {
	root  *DiskPage
	m     int
	size  int
	pager *DiskPager
}

func NewDiskBTree(m int) *DiskBTree {
	pager := NewDiskPager()
	rootID := pager.AllocatePage()
	root := &DiskPage{
		id:   rootID,
		keys: []types.PageContent{},
		leaf: true,
	}
	pager.WritePage(root)
	return &DiskBTree{
		root:  root,
		m:     m,
		size:  0,
		pager: pager,
	}
}

func (bt *DiskBTree) Root() types.PageNode {
	fmt.Println("Root called on DiskyBTree")
	// TODO: implement split/insert logic
	return bt.root
}

func (bt *DiskBTree) Insert(key types.PageContent) (types.PageContent, error) {
	fmt.Println("Insert called on DiskBTree with key:", key)
	if key == nil {
		return nil, fmt.Errorf("no key found")
	}
	if bt.root == nil {
		bt.root = NewDiskPageNode()
		bt.root.keys = append(bt.root.keys, key)
		bt.size++
		return nil, nil
	} else {
		if len(bt.root.keys) >= bt.maxKeys() {
			key2, second := bt.root.split(bt.maxKeys() / 2)
			oldroot := bt.root
			bt.root = NewDiskPageNode()
			bt.root.keys = append(bt.root.keys, key2)
			bt.root.children = append(bt.root.children, oldroot, second)
		}
	}
	out := bt.root.insert(key, bt.maxKeys())
	if out == nil {
		bt.size++
	}
	return out, nil
}

func (bt *DiskBTree) Delete(key types.PageContent) {
	fmt.Println("Delete called on DiskBTree with key:", key)
	// TODO: implement split/insert logic

	bt.size--
}

func (bt *DiskBTree) Search(key uint32) bool {
	fmt.Println("Search called on DiskBTree for key:", key)
	return false
}

func (bt *DiskBTree) Traverse(fn func(types.PageNode)) {
	fmt.Println("Traverse called")
}

func (bt *DiskBTree) Balance() error {
	fmt.Println("Balance called")
	return nil
}

func (bt *DiskBTree) Height() int {
	return 0
}

func (bt *DiskBTree) Size() int {
	return 0
}

func (bt *DiskBTree) maxKeys() int {
	return 0
}

// ------------------------
// Page Node
// ------------------------

type MemoryPage struct {
	id       types.PageID
	header   types.PageHeader
	keys     []types.PageContent
	children []*MemoryPage
	leaf     bool
}

func (p *MemoryPage) ID() types.PageID {
	return p.id
}

func (p *MemoryPage) Header() types.PageHeader {
	return p.header
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
	header   types.PageHeader
	keys     []types.PageContent
	children []*DiskPage
	leaf     bool
}

func NewDiskPageNode() *DiskPage {
	return &DiskPage{}
}

func (p *DiskPage) insert(key types.PageContent, maxItems int) types.PageContent {
	return nil
}

func (p *DiskPage) split(i int) (types.PageContent, *DiskPage) {
	key := p.keys[i]
	next := NewDiskPageNode()
	next.keys = append(next.keys, p.keys[i+1:]...)
	Truncate(&p.keys, i)
	if len(p.children) > 0 {
		next.children = append(next.children, p.children[i+1:]...)
		Truncate(&p.children, i+1)
	}
	return key, next
}

func (p *DiskPage) ID() types.PageID {
	return p.id
}

func (p *DiskPage) Header() types.PageHeader {
	return p.header
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
