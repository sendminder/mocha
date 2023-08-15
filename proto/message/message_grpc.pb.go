// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.4
// source: proto/message.proto

package message

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MessageService_CreateMessage_FullMethodName       = "/message.MessageService/CreateMessage"
	MessageService_ReadMessage_FullMethodName         = "/message.MessageService/ReadMessage"
	MessageService_DecryptConversation_FullMethodName = "/message.MessageService/DecryptConversation"
	MessageService_PushMessage_FullMethodName         = "/message.MessageService/PushMessage"
)

// MessageServiceClient is the client API for MessageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageServiceClient interface {
	CreateMessage(ctx context.Context, in *RequestCreateMessage, opts ...grpc.CallOption) (*ResponseCreateMessage, error)
	ReadMessage(ctx context.Context, in *RequestReadMessage, opts ...grpc.CallOption) (*ResponseReadMessage, error)
	DecryptConversation(ctx context.Context, in *RequestDecryptConversation, opts ...grpc.CallOption) (*ResponseDecryptConversation, error)
	PushMessage(ctx context.Context, in *RequestPushMessage, opts ...grpc.CallOption) (*ResponsePushMessage, error)
}

type messageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageServiceClient(cc grpc.ClientConnInterface) MessageServiceClient {
	return &messageServiceClient{cc}
}

func (c *messageServiceClient) CreateMessage(ctx context.Context, in *RequestCreateMessage, opts ...grpc.CallOption) (*ResponseCreateMessage, error) {
	out := new(ResponseCreateMessage)
	err := c.cc.Invoke(ctx, MessageService_CreateMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) ReadMessage(ctx context.Context, in *RequestReadMessage, opts ...grpc.CallOption) (*ResponseReadMessage, error) {
	out := new(ResponseReadMessage)
	err := c.cc.Invoke(ctx, MessageService_ReadMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) DecryptConversation(ctx context.Context, in *RequestDecryptConversation, opts ...grpc.CallOption) (*ResponseDecryptConversation, error) {
	out := new(ResponseDecryptConversation)
	err := c.cc.Invoke(ctx, MessageService_DecryptConversation_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *messageServiceClient) PushMessage(ctx context.Context, in *RequestPushMessage, opts ...grpc.CallOption) (*ResponsePushMessage, error) {
	out := new(ResponsePushMessage)
	err := c.cc.Invoke(ctx, MessageService_PushMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageServiceServer is the server API for MessageService service.
// All implementations must embed UnimplementedMessageServiceServer
// for forward compatibility
type MessageServiceServer interface {
	CreateMessage(context.Context, *RequestCreateMessage) (*ResponseCreateMessage, error)
	ReadMessage(context.Context, *RequestReadMessage) (*ResponseReadMessage, error)
	DecryptConversation(context.Context, *RequestDecryptConversation) (*ResponseDecryptConversation, error)
	PushMessage(context.Context, *RequestPushMessage) (*ResponsePushMessage, error)
	mustEmbedUnimplementedMessageServiceServer()
}

// UnimplementedMessageServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMessageServiceServer struct {
}

func (UnimplementedMessageServiceServer) CreateMessage(context.Context, *RequestCreateMessage) (*ResponseCreateMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMessage not implemented")
}
func (UnimplementedMessageServiceServer) ReadMessage(context.Context, *RequestReadMessage) (*ResponseReadMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadMessage not implemented")
}
func (UnimplementedMessageServiceServer) DecryptConversation(context.Context, *RequestDecryptConversation) (*ResponseDecryptConversation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DecryptConversation not implemented")
}
func (UnimplementedMessageServiceServer) PushMessage(context.Context, *RequestPushMessage) (*ResponsePushMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMessage not implemented")
}
func (UnimplementedMessageServiceServer) mustEmbedUnimplementedMessageServiceServer() {}

// UnsafeMessageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageServiceServer will
// result in compilation errors.
type UnsafeMessageServiceServer interface {
	mustEmbedUnimplementedMessageServiceServer()
}

func RegisterMessageServiceServer(s grpc.ServiceRegistrar, srv MessageServiceServer) {
	s.RegisterService(&MessageService_ServiceDesc, srv)
}

func _MessageService_CreateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestCreateMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).CreateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessageService_CreateMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).CreateMessage(ctx, req.(*RequestCreateMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_ReadMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestReadMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).ReadMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessageService_ReadMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).ReadMessage(ctx, req.(*RequestReadMessage))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_DecryptConversation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestDecryptConversation)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).DecryptConversation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessageService_DecryptConversation_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).DecryptConversation(ctx, req.(*RequestDecryptConversation))
	}
	return interceptor(ctx, in, info, handler)
}

func _MessageService_PushMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RequestPushMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageServiceServer).PushMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessageService_PushMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageServiceServer).PushMessage(ctx, req.(*RequestPushMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageService_ServiceDesc is the grpc.ServiceDesc for MessageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "message.MessageService",
	HandlerType: (*MessageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMessage",
			Handler:    _MessageService_CreateMessage_Handler,
		},
		{
			MethodName: "ReadMessage",
			Handler:    _MessageService_ReadMessage_Handler,
		},
		{
			MethodName: "DecryptConversation",
			Handler:    _MessageService_DecryptConversation_Handler,
		},
		{
			MethodName: "PushMessage",
			Handler:    _MessageService_PushMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/message.proto",
}
