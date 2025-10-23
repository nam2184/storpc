package storpc

const STORPC_VERSION_MAJOR uint8 = 1
const STORPC_VERSION_MINOR uint8 = 0
const STORPC_VERSION_PATCH uint8 = 0

type SerialisedGenIR struct {
	Header *SerialisedGenHeader
	Body   *SerialisedGenBody
}

func NewSerialisedGenIR(header *SerialisedGenHeader, body *SerialisedGenBody) *SerialisedGenIR {
	return &SerialisedGenIR{
		Header: header,
		Body:   body,
	}
}

type SerialisedMethodIR struct {
	Header *SerialisedMethodHeader
	Body   *SerialisedMethodBody
}

func NewSerialisedMethodIR(header *SerialisedMethodHeader, body *SerialisedMethodBody) *SerialisedMethodIR {
	return &SerialisedMethodIR{
		Header: header,
		Body:   body,
	}
}

// 15 bytes header
type SerialisedGenHeader struct {
	VersionMajor uint8  // 1 byte
	VersionMinor uint8  // 1 byte
	VersionPatch uint8  // 1 byte
	Group        uint32 // 4 bytes
	NumMessages  uint32 // 4 bytes
	NumEnums     uint32 // 4 bytes
}

func NewSerialisedGenHeader(numMessages, numEnums uint32) *SerialisedGenHeader {
	return &SerialisedGenHeader{
		VersionMajor: STORPC_VERSION_MAJOR,
		VersionMinor: STORPC_VERSION_MINOR,
		VersionPatch: STORPC_VERSION_PATCH,
		NumMessages:  numMessages,
		NumEnums:     numEnums,
	}
}

type SerialisedGenBody struct {
	Messages []SerialisedMessage
	Enums    []SerialisedEnum
}

// 8 bytes header
type SerialisedMethodHeader struct {
	VersionMajor uint8  // 1 byte
	VersionMinor uint8  // 1 byte
	VersionPatch uint8  // 1 byte
	Group        uint32 // 4 bytes
	Operation    uint8  // 1 byte operation type
}

func NewSerialisedMethodHeader(operation uint8) *SerialisedMethodHeader {
	return &SerialisedMethodHeader{
		VersionMajor: STORPC_VERSION_MAJOR,
		VersionMinor: STORPC_VERSION_MINOR,
		VersionPatch: STORPC_VERSION_PATCH,
		Operation:    operation,
	}
}

type SerialisedMethodBody struct {
	Type    string
	Message map[string]interface{}
}

func NewSerialisedMethodBody(typeof string, payload map[string]interface{}) *SerialisedMethodBody {
	return &SerialisedMethodBody{
		Type:    typeof,
		Message: payload,
	}
}

type SerialisedMessage struct {
	Name   string
	Fields []SerialisedField
}

type SerialisedField struct {
	Name string
	Type string
}

type SerialisedEnum struct {
	Name   string
	Values []SerialisedEnumValue
}

type SerialisedEnumValue struct {
	Name  string
	Value int32
}
