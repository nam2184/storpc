package storpc

const STORPC_VERSION_MAJOR uint8 = 1
const STORPC_VERSION_MINOR uint8 = 0
const STORPC_VERSION_PATCH uint8 = 0

type GenIR struct {
	Header *GenHeader
	Body   *GenBody
}

func NewGenIR(header *GenHeader, body *GenBody) *GenIR {
	return &GenIR{
		Header: header,
		Body:   body,
	}
}

type MethodIR struct {
	Header *MethodHeader
	Body   *MethodBody
}

func NewMethodIR(header *MethodHeader, body *MethodBody) *MethodIR {
	return &MethodIR{
		Header: header,
		Body:   body,
	}
}

// 11 bytes header
type GenHeader struct {
	VersionMajor uint8  // 1 byte
	VersionMinor uint8  // 1 byte
	VersionPatch uint8  // 1 byte
	NumMessages  uint32 // 4 bytes
	NumEnums     uint32 // 4 bytes
}

func NewGenHeader(numMessages, numEnums uint32) *GenHeader {
	return &GenHeader{
		VersionMajor: STORPC_VERSION_MAJOR,
		VersionMinor: STORPC_VERSION_MINOR,
		VersionPatch: STORPC_VERSION_PATCH,
		NumMessages:  numMessages,
		NumEnums:     numEnums,
	}
}

type GenBody struct {
	Messages []Message
	Enums    []Enum
	Group    string
}

func NewGenBody() *GenBody {
	return &GenBody{}
}

// 8 bytes header
type MethodHeader struct {
	VersionMajor uint8  // 1 byte
	VersionMinor uint8  // 1 byte
	VersionPatch uint8  // 1 byte
	Group        uint32 // 4 bytes
	Operation    uint8  // 1 byte operation type
}

func NewMethodHeader(operation uint8) *MethodHeader {
	return &MethodHeader{
		VersionMajor: STORPC_VERSION_MAJOR,
		VersionMinor: STORPC_VERSION_MINOR,
		VersionPatch: STORPC_VERSION_PATCH,
		Operation:    operation,
	}
}

type MethodBody struct {
	Type    string
	Message map[string]interface{}
}

func NewMethodBody(typeof string, payload map[string]interface{}) *MethodBody {
	return &MethodBody{
		Type:    typeof,
		Message: payload,
	}
}

type Message struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name   string
	Type   string
	Number int32
	Nested *Message
}

type Enum struct {
	Name   string
	Values []EnumValue
}

type EnumValue struct {
	Name  string
	Value int32
}
