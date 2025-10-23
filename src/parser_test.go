package storpc

import (
	"os"
	"path/filepath"
	"testing"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

func createTestFileDescriptorSet(t *testing.T) string {
	// define a simple message type
	msg := &descriptorpb.DescriptorProto{
		Name: proto.String("LoginRequest"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   proto.String("username"),
				Number: proto.Int32(1),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
			{
				Name:   proto.String("password"),
				Number: proto.Int32(2),
				Label:  descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL.Enum(),
				Type:   descriptorpb.FieldDescriptorProto_TYPE_STRING.Enum(),
			},
		},
	}

	svc := &descriptorpb.ServiceDescriptorProto{
		Name: proto.String("AuthService"),
		Method: []*descriptorpb.MethodDescriptorProto{
			{
				Name:       proto.String("Login"),
				InputType:  proto.String(".testpkg.LoginRequest"),
				OutputType: proto.String(".testpkg.LoginRequest"), // just reuse same type for simplicity
			},
		},
	}

	fd := &descriptorpb.FileDescriptorProto{
		Name:        proto.String("auth.proto"),
		Package:     proto.String("testpkg"),
		MessageType: []*descriptorpb.DescriptorProto{msg},
		Service:     []*descriptorpb.ServiceDescriptorProto{svc},
	}

	set := &descriptorpb.FileDescriptorSet{
		File: []*descriptorpb.FileDescriptorProto{fd},
	}

	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "fds.bin")
	data, err := proto.Marshal(set)
	if err != nil {
		t.Fatalf("failed to marshal descriptor set: %v", err)
	}

	if err := os.WriteFile(file, data, 0644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	return file
}

func TestParseWithRealMessage(t *testing.T) {
	file := createTestFileDescriptorSet(t)
	opts := &ProtoParserOptions{Filepath: file}
	parser := NewProtoParser(opts)

	if err := parser.Parse(); err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// access the message descriptor
	fd := parser.fileDesc
	msgDesc := fd.Messages().Get(0) // LoginRequest
	if string(msgDesc.Name()) != "LoginRequest" {
		t.Errorf("expected message LoginRequest, got %s", msgDesc.Name())
	}
}

func TestFilterServices_Streaming(t *testing.T) {
	fd := &descriptorpb.FileDescriptorProto{
		Name:    proto.String("stream.proto"),
		Package: proto.String("streampkg"),
		Service: []*descriptorpb.ServiceDescriptorProto{
			{
				Name: proto.String("StreamService"),
				Method: []*descriptorpb.MethodDescriptorProto{
					{
						Name:            proto.String("ServerStream"),
						InputType:       proto.String(".google.protobuf.Empty"),
						OutputType:      proto.String(".google.protobuf.Empty"),
						ServerStreaming: proto.Bool(true),
					},
				},
			},
		},
	}

	set := &descriptorpb.FileDescriptorSet{File: []*descriptorpb.FileDescriptorProto{fd}}
	data, _ := proto.Marshal(set)
	tmpFile := filepath.Join(t.TempDir(), "stream.bin")
	os.WriteFile(tmpFile, data, 0644)

	opts := &ProtoParserOptions{Filepath: tmpFile}
	parser := NewProtoParser(opts)

	err := parser.Parse()
	if err == nil {
		t.Fatalf("expected error for streaming RPC, got nil")
	}
}

func TestNewProtoParserOptions(t *testing.T) {
	args := map[ParseArgs]string{
		Input:   "/tmp/file.proto",
		Verbose: "-v",
		Group:   "users",
		Case:    "snake",
	}

	opts := NewProtoParserOptions(args)

	if opts.Filepath != "/tmp/file.proto" {
		t.Errorf("filepath mismatch")
	}
	if !opts.Verbose {
		t.Errorf("verbose expected true")
	}
	if opts.KeyGroup != "users" {
		t.Errorf("KeyGroup mismatch")
	}
	if opts.WordCase != "snake" {
		t.Errorf("WordCase mismatch")
	}
}
