package driver

import (
	"fmt"

	"github.com/nam2184/storpc/driver/types"
)

type PageHeader struct {
	size uint8
	typ  types.PageContentType
}

func (h *PageHeader) Size() uint8 {
	return h.size
}

func (h *PageHeader) Type() types.PageContentType {
	return h.typ
}

func (h *PageHeader) Read() error {
	// placeholder read
	fmt.Println("Reading PageHeader")
	return nil
}

func (h *PageHeader) Write() error {
	// placeholder write
	fmt.Println("Writing PageHeader")
	return nil
}

// ------------------------
// Concrete PageContent
// ------------------------

type PageContent struct {
	next uint8
}

func (p *PageContent) Next() uint8 {
	return p.next
}

func (p *PageContent) Read() error {
	// placeholder read
	fmt.Println("Reading PageContent")
	return nil
}

func (p *PageContent) Write() error {
	// placeholder write
	fmt.Println("Writing PageContent")
	return nil
}
