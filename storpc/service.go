package storpc

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/dynamicpb"
)

type RpcMethod struct {
	md     protoreflect.MethodDescriptor
	Output *dynamicpb.Message
}

func (m RpcMethod) Operate(input *dynamicpb.Message) MethodIR {
	var operation uint8
	if m.md.Input().FullName() == "google.protobuf.Empty" {
		operation = OpGet
	} else {
		operation = OpInsert
	}
	fields := input.Descriptor().Fields()
	payload := make(map[string]interface{})

	for i := 0; i < fields.Len(); i++ {
		f := fields.Get(i)
		val := input.Get(f)                         // protoreflect.Value
		payload[string(f.Name())] = val.Interface() // convert to Go type
	}

	header := NewMethodHeader(operation)
	body := NewMethodBody(string(input.Descriptor().FullName()), payload)

	serialiasedMethod := MethodIR{
		Header: header,
		Body:   body,
	}

	return serialiasedMethod
}

func RunDynamicServer(fd protoreflect.FileDescriptor) error {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		return err
	}

	server := grpc.NewServer()

	for s := 0; s < fd.Services().Len(); s++ {
		svc := fd.Services().Get(s)

		methods := make([]grpc.MethodDesc, 0, svc.Methods().Len())

		for m := 0; m < svc.Methods().Len(); m++ {
			md := svc.Methods().Get(m)

			handler := func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
				req := dynamicpb.NewMessage(md.Input())
				if err := dec(req); err != nil {
					return nil, err
				}

				reply := dynamicpb.NewMessage(md.Output())
				return reply, nil
			}

			methods = append(methods, grpc.MethodDesc{
				MethodName: string(md.Name()),
				Handler:    handler,
			})
		}

		server.RegisterService(&grpc.ServiceDesc{
			ServiceName: string(svc.FullName()),
			HandlerType: (*interface{})(nil),
			Methods:     methods,
			Streams:     []grpc.StreamDesc{},
		}, nil)
	}

	return server.Serve(lis)
}
