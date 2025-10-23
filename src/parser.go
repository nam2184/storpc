package storpc

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
)

type ParseArgs string

const (
	Verbose ParseArgs = "-v"
	Quiet   ParseArgs = "-q"
	Input   ParseArgs = "--input"
	Case    ParseArgs = "--case"
	Group   ParseArgs = "--group"
)

type ProtoParser struct {
	fileDesc protoreflect.FileDescriptor
	options  *ProtoParserOptions
	logger   *slog.Logger
}

func NewProtoParser(options *ProtoParserOptions) *ProtoParser {
	var level slog.Level

	switch {
	case options.Verbose:
		level = slog.LevelDebug
	case options.Quiet:
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	logger := slog.New(handler)

	parser := &ProtoParser{
		options: options,
		logger:  logger,
	}

	return parser
}

func (p *ProtoParser) Parse() error {
	p.logger.Debug(fmt.Sprintf("parsing filepath : %v", p.options.Filepath))
	data, err := os.ReadFile(p.options.Filepath)
	if err != nil {
		return err
	}

	set := &descriptorpb.FileDescriptorSet{}
	if err := proto.Unmarshal(data, set); err != nil {
		return err
	}

	fd, err := protodesc.NewFile(set.File[0], nil)
	if err != nil {
		return err
	}

	p.fileDesc = fd

	if string(p.fileDesc.Package()) == "" {
		p.options.KeyGroup = string(p.fileDesc.Package())
	}

	err = p.filterServices()
	if err != nil {
		return err
	}

	return nil
}

func (p *ProtoParser) filterServices() error {
	fd := p.fileDesc

	for i := 0; i < fd.Services().Len(); i++ {
		svc := fd.Services().Get(i)
		for j := 0; j < svc.Methods().Len(); j++ {
			m := svc.Methods().Get(j)
			if m.IsStreamingClient() || m.IsStreamingServer() {
				return errors.New("streaming RPC detected: service " +
					string(svc.FullName()) + ", method " + string(m.Name()))
			}
		}
	}

	return nil
}

type ProtoParserOptions struct {
	Filepath string
	WordCase string //camel, snake, etc.
	Verbose  bool
	Quiet    bool
	KeyGroup string //grouping for structs
}

func NewProtoParserOptions(args map[ParseArgs]string) *ProtoParserOptions {
	return &ProtoParserOptions{
		Filepath: args[Input],
		Verbose:  args[Verbose] != "",
		Quiet:    args[Quiet] != "",
		KeyGroup: args[Group],
		WordCase: args[Case],
	}
}

/*
type FileDescriptorProto struct {
	state   protoimpl.MessageState `protogen:"open.v1"`
	Name    *string                `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`       // file name, relative to root of source tree
	Package *string                `protobuf:"bytes,2,opt,name=package" json:"package,omitempty"` // e.g. "foo", "foo.bar", etc.
	// Names of files imported by this file.
	Dependency []string `protobuf:"bytes,3,rep,name=dependency" json:"dependency,omitempty"`
	// Indexes of the public imported files in the dependency list above.
	PublicDependency []int32 `protobuf:"varint,10,rep,name=public_dependency,json=publicDependency" json:"public_dependency,omitempty"`
	// Indexes of the weak imported files in the dependency list.
	// For Google-internal migration only. Do not use.
	WeakDependency []int32 `protobuf:"varint,11,rep,name=weak_dependency,json=weakDependency" json:"weak_dependency,omitempty"`
	// Names of files imported by this file purely for the purpose of providing
	// option extensions. These are excluded from the dependency list above.
	OptionDependency []string `protobuf:"bytes,15,rep,name=option_dependency,json=optionDependency" json:"option_dependency,omitempty"`
	// All top-level definitions in this file.
	MessageType []*DescriptorProto        `protobuf:"bytes,4,rep,name=message_type,json=messageType" json:"message_type,omitempty"`
	EnumType    []*EnumDescriptorProto    `protobuf:"bytes,5,rep,name=enum_type,json=enumType" json:"enum_type,omitempty"`
	Service     []*ServiceDescriptorProto `protobuf:"bytes,6,rep,name=service" json:"service,omitempty"`
	Extension   []*FieldDescriptorProto   `protobuf:"bytes,7,rep,name=extension" json:"extension,omitempty"`
	Options     *FileOptions              `protobuf:"bytes,8,opt,name=options" json:"options,omitempty"`
	// This field contains optional information about the original source code.
	// You may safely remove this entire field without harming runtime
	// functionality of the descriptors -- the information is needed only by
	// development tools.
	SourceCodeInfo *SourceCodeInfo `protobuf:"bytes,9,opt,name=source_code_info,json=sourceCodeInfo" json:"source_code_info,omitempty"`
	// The syntax of the proto file.
	// The supported values are "proto2", "proto3", and "editions".
	//
	// If `edition` is present, this value must be "editions".
	// WARNING: This field should only be used by protobuf plugins or special
	// cases like the proto compiler. Other uses are discouraged and
	// developers should rely on the protoreflect APIs for their client language.
	Syntax *string `protobuf:"bytes,12,opt,name=syntax" json:"syntax,omitempty"`
	// The edition of the proto file.
	// WARNING: This field should only be used by protobuf plugins or special
	// cases like the proto compiler. Other uses are discouraged and
	// developers should rely on the protoreflect APIs for their client language.
	Edition       *Edition `protobuf:"varint,14,opt,name=edition,enum=google.protobuf.Edition" json:"edition,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}
*/
