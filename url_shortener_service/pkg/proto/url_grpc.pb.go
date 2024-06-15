// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: pkg/proto/url.proto

package url

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

// UrlClient is the client API for Url service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UrlClient interface {
	ShortenUrl(ctx context.Context, in *LongUrlRequest, opts ...grpc.CallOption) (*UrlDataResponse, error)
}

type urlClient struct {
	cc grpc.ClientConnInterface
}

func NewUrlClient(cc grpc.ClientConnInterface) UrlClient {
	return &urlClient{cc}
}

func (c *urlClient) ShortenUrl(ctx context.Context, in *LongUrlRequest, opts ...grpc.CallOption) (*UrlDataResponse, error) {
	out := new(UrlDataResponse)
	err := c.cc.Invoke(ctx, "/url.Url/ShortenUrl", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UrlServer is the server API for Url service.
// All implementations must embed UnimplementedUrlServer
// for forward compatibility
type UrlServer interface {
	ShortenUrl(context.Context, *LongUrlRequest) (*UrlDataResponse, error)
	mustEmbedUnimplementedUrlServer()
}

// UnimplementedUrlServer must be embedded to have forward compatible implementations.
type UnimplementedUrlServer struct {
}

func (UnimplementedUrlServer) ShortenUrl(context.Context, *LongUrlRequest) (*UrlDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShortenUrl not implemented")
}
func (UnimplementedUrlServer) mustEmbedUnimplementedUrlServer() {}

// UnsafeUrlServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UrlServer will
// result in compilation errors.
type UnsafeUrlServer interface {
	mustEmbedUnimplementedUrlServer()
}

func RegisterUrlServer(s grpc.ServiceRegistrar, srv UrlServer) {
	s.RegisterService(&Url_ServiceDesc, srv)
}

func _Url_ShortenUrl_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LongUrlRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UrlServer).ShortenUrl(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/url.Url/ShortenUrl",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UrlServer).ShortenUrl(ctx, req.(*LongUrlRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Url_ServiceDesc is the grpc.ServiceDesc for Url service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Url_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "url.Url",
	HandlerType: (*UrlServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ShortenUrl",
			Handler:    _Url_ShortenUrl_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/proto/url.proto",
}
