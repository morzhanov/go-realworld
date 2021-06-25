// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package analytics

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

// PicturesClient is the client API for Pictures service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PicturesClient interface {
	// Get user pictures
	GetUserPictures(ctx context.Context, in *GetUserPicturesRequest, opts ...grpc.CallOption) (*PicturesMessage, error)
	// Get user picture
	GetUserPicture(ctx context.Context, in *GetUserPictureRequest, opts ...grpc.CallOption) (*PictureMessage, error)
	// Create user picture
	CreateUserPicture(ctx context.Context, in *CreateUserPictureRequest, opts ...grpc.CallOption) (*PictureMessage, error)
	// Delete user picture
	DeleteUserPicture(ctx context.Context, in *DeleteUserPictureRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type picturesClient struct {
	cc grpc.ClientConnInterface
}

func NewPicturesClient(cc grpc.ClientConnInterface) PicturesClient {
	return &picturesClient{cc}
}

func (c *picturesClient) GetUserPictures(ctx context.Context, in *GetUserPicturesRequest, opts ...grpc.CallOption) (*PicturesMessage, error) {
	out := new(PicturesMessage)
	err := c.cc.Invoke(ctx, "/main.Pictures/GetUserPictures", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *picturesClient) GetUserPicture(ctx context.Context, in *GetUserPictureRequest, opts ...grpc.CallOption) (*PictureMessage, error) {
	out := new(PictureMessage)
	err := c.cc.Invoke(ctx, "/main.Pictures/GetUserPicture", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *picturesClient) CreateUserPicture(ctx context.Context, in *CreateUserPictureRequest, opts ...grpc.CallOption) (*PictureMessage, error) {
	out := new(PictureMessage)
	err := c.cc.Invoke(ctx, "/main.Pictures/CreateUserPicture", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *picturesClient) DeleteUserPicture(ctx context.Context, in *DeleteUserPictureRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/main.Pictures/DeleteUserPicture", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PicturesServer is the server API for Pictures service.
// All implementations must embed UnimplementedPicturesServer
// for forward compatibility
type PicturesServer interface {
	// Get user pictures
	GetUserPictures(context.Context, *GetUserPicturesRequest) (*PicturesMessage, error)
	// Get user picture
	GetUserPicture(context.Context, *GetUserPictureRequest) (*PictureMessage, error)
	// Create user picture
	CreateUserPicture(context.Context, *CreateUserPictureRequest) (*PictureMessage, error)
	// Delete user picture
	DeleteUserPicture(context.Context, *DeleteUserPictureRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedPicturesServer()
}

// UnimplementedPicturesServer must be embedded to have forward compatible implementations.
type UnimplementedPicturesServer struct {
}

func (UnimplementedPicturesServer) GetUserPictures(context.Context, *GetUserPicturesRequest) (*PicturesMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPictures not implemented")
}
func (UnimplementedPicturesServer) GetUserPicture(context.Context, *GetUserPictureRequest) (*PictureMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPicture not implemented")
}
func (UnimplementedPicturesServer) CreateUserPicture(context.Context, *CreateUserPictureRequest) (*PictureMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateUserPicture not implemented")
}
func (UnimplementedPicturesServer) DeleteUserPicture(context.Context, *DeleteUserPictureRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteUserPicture not implemented")
}
func (UnimplementedPicturesServer) mustEmbedUnimplementedPicturesServer() {}

// UnsafePicturesServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PicturesServer will
// result in compilation errors.
type UnsafePicturesServer interface {
	mustEmbedUnimplementedPicturesServer()
}

func RegisterPicturesServer(s grpc.ServiceRegistrar, srv PicturesServer) {
	s.RegisterService(&Pictures_ServiceDesc, srv)
}

func _Pictures_GetUserPictures_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserPicturesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PicturesServer).GetUserPictures(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Pictures/GetUserPictures",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PicturesServer).GetUserPictures(ctx, req.(*GetUserPicturesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pictures_GetUserPicture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserPictureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PicturesServer).GetUserPicture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Pictures/GetUserPicture",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PicturesServer).GetUserPicture(ctx, req.(*GetUserPictureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pictures_CreateUserPicture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserPictureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PicturesServer).CreateUserPicture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Pictures/CreateUserPicture",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PicturesServer).CreateUserPicture(ctx, req.(*CreateUserPictureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Pictures_DeleteUserPicture_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserPictureRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PicturesServer).DeleteUserPicture(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/main.Pictures/DeleteUserPicture",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PicturesServer).DeleteUserPicture(ctx, req.(*DeleteUserPictureRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Pictures_ServiceDesc is the grpc.ServiceDesc for Pictures service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Pictures_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "main.Pictures",
	HandlerType: (*PicturesServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetUserPictures",
			Handler:    _Pictures_GetUserPictures_Handler,
		},
		{
			MethodName: "GetUserPicture",
			Handler:    _Pictures_GetUserPicture_Handler,
		},
		{
			MethodName: "CreateUserPicture",
			Handler:    _Pictures_CreateUserPicture_Handler,
		},
		{
			MethodName: "DeleteUserPicture",
			Handler:    _Pictures_DeleteUserPicture_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pictures/pictures.proto",
}