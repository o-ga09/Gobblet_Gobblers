// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: game.proto

package grpc

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// GreetingServiceClient is the client API for GreetingService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreetingServiceClient interface {
	//unary
	Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error)
	//server stream
	HelloServerStream(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (GreetingService_HelloServerStreamClient, error)
	//client stream
	HelloClientStream(ctx context.Context, opts ...grpc.CallOption) (GreetingService_HelloClientStreamClient, error)
	//Bistream
	HelloBiStream(ctx context.Context, opts ...grpc.CallOption) (GreetingService_HelloBiStreamClient, error)
	//chat
	CreateMessage(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error)
	GetMessages(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (GreetingService_GetMessagesClient, error)
	//Bidirectional Chat
	Chat(ctx context.Context, opts ...grpc.CallOption) (GreetingService_ChatClient, error)
}

type greetingServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGreetingServiceClient(cc grpc.ClientConnInterface) GreetingServiceClient {
	return &greetingServiceClient{cc}
}

func (c *greetingServiceClient) Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloResponse, error) {
	out := new(HelloResponse)
	err := c.cc.Invoke(ctx, "/game.GreetingService/Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greetingServiceClient) HelloServerStream(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (GreetingService_HelloServerStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &GreetingService_ServiceDesc.Streams[0], "/game.GreetingService/HelloServerStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetingServiceHelloServerStreamClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GreetingService_HelloServerStreamClient interface {
	Recv() (*HelloResponse, error)
	grpc.ClientStream
}

type greetingServiceHelloServerStreamClient struct {
	grpc.ClientStream
}

func (x *greetingServiceHelloServerStreamClient) Recv() (*HelloResponse, error) {
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greetingServiceClient) HelloClientStream(ctx context.Context, opts ...grpc.CallOption) (GreetingService_HelloClientStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &GreetingService_ServiceDesc.Streams[1], "/game.GreetingService/HelloClientStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetingServiceHelloClientStreamClient{stream}
	return x, nil
}

type GreetingService_HelloClientStreamClient interface {
	Send(*HelloRequest) error
	CloseAndRecv() (*HelloResponse, error)
	grpc.ClientStream
}

type greetingServiceHelloClientStreamClient struct {
	grpc.ClientStream
}

func (x *greetingServiceHelloClientStreamClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greetingServiceHelloClientStreamClient) CloseAndRecv() (*HelloResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greetingServiceClient) HelloBiStream(ctx context.Context, opts ...grpc.CallOption) (GreetingService_HelloBiStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &GreetingService_ServiceDesc.Streams[2], "/game.GreetingService/HelloBiStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetingServiceHelloBiStreamClient{stream}
	return x, nil
}

type GreetingService_HelloBiStreamClient interface {
	Send(*HelloRequest) error
	Recv() (*HelloResponse, error)
	grpc.ClientStream
}

type greetingServiceHelloBiStreamClient struct {
	grpc.ClientStream
}

func (x *greetingServiceHelloBiStreamClient) Send(m *HelloRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greetingServiceHelloBiStreamClient) Recv() (*HelloResponse, error) {
	m := new(HelloResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greetingServiceClient) CreateMessage(ctx context.Context, in *MessageRequest, opts ...grpc.CallOption) (*MessageResponse, error) {
	out := new(MessageResponse)
	err := c.cc.Invoke(ctx, "/game.GreetingService/CreateMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *greetingServiceClient) GetMessages(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (GreetingService_GetMessagesClient, error) {
	stream, err := c.cc.NewStream(ctx, &GreetingService_ServiceDesc.Streams[3], "/game.GreetingService/GetMessages", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetingServiceGetMessagesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type GreetingService_GetMessagesClient interface {
	Recv() (*MessageResponse, error)
	grpc.ClientStream
}

type greetingServiceGetMessagesClient struct {
	grpc.ClientStream
}

func (x *greetingServiceGetMessagesClient) Recv() (*MessageResponse, error) {
	m := new(MessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *greetingServiceClient) Chat(ctx context.Context, opts ...grpc.CallOption) (GreetingService_ChatClient, error) {
	stream, err := c.cc.NewStream(ctx, &GreetingService_ServiceDesc.Streams[4], "/game.GreetingService/Chat", opts...)
	if err != nil {
		return nil, err
	}
	x := &greetingServiceChatClient{stream}
	return x, nil
}

type GreetingService_ChatClient interface {
	Send(*MessageRequest) error
	Recv() (*MessageResponse, error)
	grpc.ClientStream
}

type greetingServiceChatClient struct {
	grpc.ClientStream
}

func (x *greetingServiceChatClient) Send(m *MessageRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *greetingServiceChatClient) Recv() (*MessageResponse, error) {
	m := new(MessageResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GreetingServiceServer is the server API for GreetingService service.
// All implementations must embed UnimplementedGreetingServiceServer
// for forward compatibility
type GreetingServiceServer interface {
	//unary
	Hello(context.Context, *HelloRequest) (*HelloResponse, error)
	//server stream
	HelloServerStream(*HelloRequest, GreetingService_HelloServerStreamServer) error
	//client stream
	HelloClientStream(GreetingService_HelloClientStreamServer) error
	//Bistream
	HelloBiStream(GreetingService_HelloBiStreamServer) error
	//chat
	CreateMessage(context.Context, *MessageRequest) (*MessageResponse, error)
	GetMessages(*emptypb.Empty, GreetingService_GetMessagesServer) error
	//Bidirectional Chat
	Chat(GreetingService_ChatServer) error
	mustEmbedUnimplementedGreetingServiceServer()
}

// UnimplementedGreetingServiceServer must be embedded to have forward compatible implementations.
type UnimplementedGreetingServiceServer struct {
}

func (UnimplementedGreetingServiceServer) Hello(context.Context, *HelloRequest) (*HelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Hello not implemented")
}
func (UnimplementedGreetingServiceServer) HelloServerStream(*HelloRequest, GreetingService_HelloServerStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HelloServerStream not implemented")
}
func (UnimplementedGreetingServiceServer) HelloClientStream(GreetingService_HelloClientStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HelloClientStream not implemented")
}
func (UnimplementedGreetingServiceServer) HelloBiStream(GreetingService_HelloBiStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method HelloBiStream not implemented")
}
func (UnimplementedGreetingServiceServer) CreateMessage(context.Context, *MessageRequest) (*MessageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMessage not implemented")
}
func (UnimplementedGreetingServiceServer) GetMessages(*emptypb.Empty, GreetingService_GetMessagesServer) error {
	return status.Errorf(codes.Unimplemented, "method GetMessages not implemented")
}
func (UnimplementedGreetingServiceServer) Chat(GreetingService_ChatServer) error {
	return status.Errorf(codes.Unimplemented, "method Chat not implemented")
}
func (UnimplementedGreetingServiceServer) mustEmbedUnimplementedGreetingServiceServer() {}

// UnsafeGreetingServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GreetingServiceServer will
// result in compilation errors.
type UnsafeGreetingServiceServer interface {
	mustEmbedUnimplementedGreetingServiceServer()
}

func RegisterGreetingServiceServer(s grpc.ServiceRegistrar, srv GreetingServiceServer) {
	s.RegisterService(&GreetingService_ServiceDesc, srv)
}

func _GreetingService_Hello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreetingServiceServer).Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GreetingService/Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreetingServiceServer).Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GreetingService_HelloServerStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(HelloRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreetingServiceServer).HelloServerStream(m, &greetingServiceHelloServerStreamServer{stream})
}

type GreetingService_HelloServerStreamServer interface {
	Send(*HelloResponse) error
	grpc.ServerStream
}

type greetingServiceHelloServerStreamServer struct {
	grpc.ServerStream
}

func (x *greetingServiceHelloServerStreamServer) Send(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GreetingService_HelloClientStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreetingServiceServer).HelloClientStream(&greetingServiceHelloClientStreamServer{stream})
}

type GreetingService_HelloClientStreamServer interface {
	SendAndClose(*HelloResponse) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type greetingServiceHelloClientStreamServer struct {
	grpc.ServerStream
}

func (x *greetingServiceHelloClientStreamServer) SendAndClose(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greetingServiceHelloClientStreamServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GreetingService_HelloBiStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreetingServiceServer).HelloBiStream(&greetingServiceHelloBiStreamServer{stream})
}

type GreetingService_HelloBiStreamServer interface {
	Send(*HelloResponse) error
	Recv() (*HelloRequest, error)
	grpc.ServerStream
}

type greetingServiceHelloBiStreamServer struct {
	grpc.ServerStream
}

func (x *greetingServiceHelloBiStreamServer) Send(m *HelloResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greetingServiceHelloBiStreamServer) Recv() (*HelloRequest, error) {
	m := new(HelloRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _GreetingService_CreateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GreetingServiceServer).CreateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/game.GreetingService/CreateMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GreetingServiceServer).CreateMessage(ctx, req.(*MessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GreetingService_GetMessages_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(emptypb.Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(GreetingServiceServer).GetMessages(m, &greetingServiceGetMessagesServer{stream})
}

type GreetingService_GetMessagesServer interface {
	Send(*MessageResponse) error
	grpc.ServerStream
}

type greetingServiceGetMessagesServer struct {
	grpc.ServerStream
}

func (x *greetingServiceGetMessagesServer) Send(m *MessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _GreetingService_Chat_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GreetingServiceServer).Chat(&greetingServiceChatServer{stream})
}

type GreetingService_ChatServer interface {
	Send(*MessageResponse) error
	Recv() (*MessageRequest, error)
	grpc.ServerStream
}

type greetingServiceChatServer struct {
	grpc.ServerStream
}

func (x *greetingServiceChatServer) Send(m *MessageResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *greetingServiceChatServer) Recv() (*MessageRequest, error) {
	m := new(MessageRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// GreetingService_ServiceDesc is the grpc.ServiceDesc for GreetingService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GreetingService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "game.GreetingService",
	HandlerType: (*GreetingServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Hello",
			Handler:    _GreetingService_Hello_Handler,
		},
		{
			MethodName: "CreateMessage",
			Handler:    _GreetingService_CreateMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "HelloServerStream",
			Handler:       _GreetingService_HelloServerStream_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "HelloClientStream",
			Handler:       _GreetingService_HelloClientStream_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "HelloBiStream",
			Handler:       _GreetingService_HelloBiStream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "GetMessages",
			Handler:       _GreetingService_GetMessages_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Chat",
			Handler:       _GreetingService_Chat_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "game.proto",
}
