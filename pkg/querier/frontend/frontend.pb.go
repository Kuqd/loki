// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: frontend.proto

package frontend

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	queryrange "github.com/grafana/loki/pkg/querier/queryrange"
	httpgrpc "github.com/weaveworks/common/httpgrpc"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion2 // please upgrade the proto package

type ProcessRequest struct {
	HttpRequest       *httpgrpc.HTTPRequest `protobuf:"bytes,1,opt,name=httpRequest,proto3" json:"httpRequest,omitempty"`
	QueryRangeRequest *queryrange.Request   `protobuf:"bytes,2,opt,name=queryRangeRequest,proto3" json:"queryRangeRequest,omitempty"`
}

func (m *ProcessRequest) Reset()      { *m = ProcessRequest{} }
func (*ProcessRequest) ProtoMessage() {}
func (*ProcessRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_eca3873955a29cfe, []int{0}
}
func (m *ProcessRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProcessRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessRequest.Merge(m, src)
}
func (m *ProcessRequest) XXX_Size() int {
	return m.Size()
}
func (m *ProcessRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessRequest proto.InternalMessageInfo

func (m *ProcessRequest) GetHttpRequest() *httpgrpc.HTTPRequest {
	if m != nil {
		return m.HttpRequest
	}
	return nil
}

func (m *ProcessRequest) GetQueryRangeRequest() *queryrange.Request {
	if m != nil {
		return m.QueryRangeRequest
	}
	return nil
}

type ProcessResponse struct {
	HttpResponse *httpgrpc.HTTPResponse  `protobuf:"bytes,1,opt,name=httpResponse,proto3" json:"httpResponse,omitempty"`
	ApiResponse  *queryrange.APIResponse `protobuf:"bytes,2,opt,name=apiResponse,proto3" json:"apiResponse,omitempty"`
}

func (m *ProcessResponse) Reset()      { *m = ProcessResponse{} }
func (*ProcessResponse) ProtoMessage() {}
func (*ProcessResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_eca3873955a29cfe, []int{1}
}
func (m *ProcessResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ProcessResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ProcessResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ProcessResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ProcessResponse.Merge(m, src)
}
func (m *ProcessResponse) XXX_Size() int {
	return m.Size()
}
func (m *ProcessResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ProcessResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ProcessResponse proto.InternalMessageInfo

func (m *ProcessResponse) GetHttpResponse() *httpgrpc.HTTPResponse {
	if m != nil {
		return m.HttpResponse
	}
	return nil
}

func (m *ProcessResponse) GetApiResponse() *queryrange.APIResponse {
	if m != nil {
		return m.ApiResponse
	}
	return nil
}

func init() {
	proto.RegisterType((*ProcessRequest)(nil), "frontend.ProcessRequest")
	proto.RegisterType((*ProcessResponse)(nil), "frontend.ProcessResponse")
}

func init() { proto.RegisterFile("frontend.proto", fileDescriptor_eca3873955a29cfe) }

var fileDescriptor_eca3873955a29cfe = []byte{
	// 362 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xb1, 0x4e, 0xf3, 0x30,
	0x18, 0xb4, 0xff, 0xe1, 0xa7, 0x72, 0x51, 0x11, 0x46, 0x40, 0xe9, 0x60, 0xa1, 0x4e, 0x5d, 0x88,
	0x51, 0x41, 0x42, 0x30, 0x80, 0x8a, 0x10, 0x82, 0x2d, 0x8a, 0x3a, 0xb1, 0xa5, 0xc1, 0x4d, 0xa3,
	0xd2, 0x38, 0x75, 0x12, 0x2a, 0x36, 0x46, 0x06, 0x06, 0x1e, 0x83, 0x47, 0x61, 0xec, 0xd8, 0x91,
	0xba, 0x0b, 0x63, 0x1f, 0x01, 0xd5, 0x4e, 0x2d, 0x53, 0x98, 0xf2, 0x9d, 0xbe, 0xbb, 0xef, 0x2e,
	0x3a, 0xa3, 0x4a, 0x57, 0xf0, 0x38, 0x63, 0xf1, 0xbd, 0x93, 0x08, 0x9e, 0x71, 0x5c, 0x5a, 0xe2,
	0xda, 0x41, 0x18, 0x65, 0xbd, 0xbc, 0xe3, 0x04, 0x7c, 0x40, 0x43, 0x1e, 0x72, 0xaa, 0x08, 0x9d,
	0xbc, 0xab, 0x90, 0x02, 0x6a, 0xd2, 0xc2, 0xda, 0xb1, 0x45, 0x1f, 0x31, 0xff, 0x91, 0x8d, 0xb8,
	0xe8, 0xa7, 0x34, 0xe0, 0x83, 0x01, 0x8f, 0x69, 0x2f, 0xcb, 0x92, 0x50, 0x24, 0x81, 0x19, 0x0a,
	0xd5, 0x85, 0x6d, 0x22, 0xfc, 0xae, 0x1f, 0xfb, 0xf4, 0x81, 0xf7, 0x23, 0x9a, 0xf4, 0x43, 0x3a,
	0xcc, 0x99, 0x88, 0x98, 0x50, 0xdf, 0x27, 0xe1, 0xc7, 0x21, 0xb3, 0x46, 0x7d, 0xa0, 0xfe, 0x0a,
	0x51, 0xc5, 0x15, 0x3c, 0x60, 0x69, 0xea, 0xb1, 0x61, 0xce, 0xd2, 0x0c, 0x9f, 0xa0, 0xf2, 0xc2,
	0xa5, 0x80, 0x55, 0xb8, 0x0f, 0x1b, 0xe5, 0xe6, 0xb6, 0x63, 0x9c, 0x6f, 0xda, 0x6d, 0xb7, 0x58,
	0x7a, 0x36, 0x13, 0xb7, 0xd0, 0xa6, 0xba, 0xef, 0x2d, 0xee, 0x2f, 0xe5, 0xff, 0x94, 0x7c, 0xcb,
	0xb1, 0x9c, 0x97, 0xe2, 0xdf, 0xec, 0xfa, 0x0b, 0x44, 0x1b, 0x26, 0x4e, 0x9a, 0xf0, 0x38, 0x65,
	0xf8, 0x0c, 0xad, 0x6b, 0x17, 0x8d, 0x8b, 0x40, 0x3b, 0xab, 0x81, 0xf4, 0xd6, 0xfb, 0xc1, 0xc5,
	0xa7, 0xa8, 0xec, 0x27, 0x91, 0x91, 0xea, 0x30, 0xbb, 0x76, 0x98, 0x96, 0x7b, 0x6b, 0xb4, 0x36,
	0xb7, 0xe9, 0xa2, 0xd2, 0x75, 0xd1, 0x25, 0xbe, 0x42, 0x6b, 0x45, 0x2a, 0xbc, 0xe7, 0x98, 0xc6,
	0x57, 0x82, 0xd6, 0xaa, 0x7f, 0xac, 0xf4, 0x6f, 0x81, 0x06, 0x3c, 0x84, 0x97, 0xe7, 0xe3, 0x29,
	0x01, 0x93, 0x29, 0x01, 0xf3, 0x29, 0x81, 0xcf, 0x92, 0xc0, 0x77, 0x49, 0xe0, 0x87, 0x24, 0x70,
	0x2c, 0x09, 0xfc, 0x94, 0x04, 0x7e, 0x49, 0x02, 0xe6, 0x92, 0xc0, 0xb7, 0x19, 0x01, 0xe3, 0x19,
	0x01, 0x93, 0x19, 0x01, 0x77, 0xe6, 0x45, 0x75, 0xfe, 0xab, 0xca, 0x8e, 0xbe, 0x03, 0x00, 0x00,
	0xff, 0xff, 0x37, 0x97, 0xfe, 0x9d, 0x74, 0x02, 0x00, 0x00,
}

func (this *ProcessRequest) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProcessRequest)
	if !ok {
		that2, ok := that.(ProcessRequest)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.HttpRequest.Equal(that1.HttpRequest) {
		return false
	}
	if !this.QueryRangeRequest.Equal(that1.QueryRangeRequest) {
		return false
	}
	return true
}
func (this *ProcessResponse) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*ProcessResponse)
	if !ok {
		that2, ok := that.(ProcessResponse)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.HttpResponse.Equal(that1.HttpResponse) {
		return false
	}
	if !this.ApiResponse.Equal(that1.ApiResponse) {
		return false
	}
	return true
}
func (this *ProcessRequest) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&frontend.ProcessRequest{")
	if this.HttpRequest != nil {
		s = append(s, "HttpRequest: "+fmt.Sprintf("%#v", this.HttpRequest)+",\n")
	}
	if this.QueryRangeRequest != nil {
		s = append(s, "QueryRangeRequest: "+fmt.Sprintf("%#v", this.QueryRangeRequest)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func (this *ProcessResponse) GoString() string {
	if this == nil {
		return "nil"
	}
	s := make([]string, 0, 6)
	s = append(s, "&frontend.ProcessResponse{")
	if this.HttpResponse != nil {
		s = append(s, "HttpResponse: "+fmt.Sprintf("%#v", this.HttpResponse)+",\n")
	}
	if this.ApiResponse != nil {
		s = append(s, "ApiResponse: "+fmt.Sprintf("%#v", this.ApiResponse)+",\n")
	}
	s = append(s, "}")
	return strings.Join(s, "")
}
func valueToGoStringFrontend(v interface{}, typ string) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("func(v %v) *%v { return &v } ( %#v )", typ, typ, pv)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// FrontendClient is the client API for Frontend service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type FrontendClient interface {
	Process(ctx context.Context, opts ...grpc.CallOption) (Frontend_ProcessClient, error)
}

type frontendClient struct {
	cc *grpc.ClientConn
}

func NewFrontendClient(cc *grpc.ClientConn) FrontendClient {
	return &frontendClient{cc}
}

func (c *frontendClient) Process(ctx context.Context, opts ...grpc.CallOption) (Frontend_ProcessClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Frontend_serviceDesc.Streams[0], "/frontend.Frontend/Process", opts...)
	if err != nil {
		return nil, err
	}
	x := &frontendProcessClient{stream}
	return x, nil
}

type Frontend_ProcessClient interface {
	Send(*ProcessResponse) error
	Recv() (*ProcessRequest, error)
	grpc.ClientStream
}

type frontendProcessClient struct {
	grpc.ClientStream
}

func (x *frontendProcessClient) Send(m *ProcessResponse) error {
	return x.ClientStream.SendMsg(m)
}

func (x *frontendProcessClient) Recv() (*ProcessRequest, error) {
	m := new(ProcessRequest)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// FrontendServer is the server API for Frontend service.
type FrontendServer interface {
	Process(Frontend_ProcessServer) error
}

// UnimplementedFrontendServer can be embedded to have forward compatible implementations.
type UnimplementedFrontendServer struct {
}

func (*UnimplementedFrontendServer) Process(srv Frontend_ProcessServer) error {
	return status.Errorf(codes.Unimplemented, "method Process not implemented")
}

func RegisterFrontendServer(s *grpc.Server, srv FrontendServer) {
	s.RegisterService(&_Frontend_serviceDesc, srv)
}

func _Frontend_Process_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FrontendServer).Process(&frontendProcessServer{stream})
}

type Frontend_ProcessServer interface {
	Send(*ProcessRequest) error
	Recv() (*ProcessResponse, error)
	grpc.ServerStream
}

type frontendProcessServer struct {
	grpc.ServerStream
}

func (x *frontendProcessServer) Send(m *ProcessRequest) error {
	return x.ServerStream.SendMsg(m)
}

func (x *frontendProcessServer) Recv() (*ProcessResponse, error) {
	m := new(ProcessResponse)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Frontend_serviceDesc = grpc.ServiceDesc{
	ServiceName: "frontend.Frontend",
	HandlerType: (*FrontendServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Process",
			Handler:       _Frontend_Process_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "frontend.proto",
}

func (m *ProcessRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProcessRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.QueryRangeRequest != nil {
		{
			size, err := m.QueryRangeRequest.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintFrontend(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.HttpRequest != nil {
		{
			size, err := m.HttpRequest.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintFrontend(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ProcessResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ProcessResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ProcessResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ApiResponse != nil {
		{
			size, err := m.ApiResponse.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintFrontend(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.HttpResponse != nil {
		{
			size, err := m.HttpResponse.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintFrontend(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintFrontend(dAtA []byte, offset int, v uint64) int {
	offset -= sovFrontend(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *ProcessRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.HttpRequest != nil {
		l = m.HttpRequest.Size()
		n += 1 + l + sovFrontend(uint64(l))
	}
	if m.QueryRangeRequest != nil {
		l = m.QueryRangeRequest.Size()
		n += 1 + l + sovFrontend(uint64(l))
	}
	return n
}

func (m *ProcessResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.HttpResponse != nil {
		l = m.HttpResponse.Size()
		n += 1 + l + sovFrontend(uint64(l))
	}
	if m.ApiResponse != nil {
		l = m.ApiResponse.Size()
		n += 1 + l + sovFrontend(uint64(l))
	}
	return n
}

func sovFrontend(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozFrontend(x uint64) (n int) {
	return sovFrontend(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *ProcessRequest) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProcessRequest{`,
		`HttpRequest:` + strings.Replace(fmt.Sprintf("%v", this.HttpRequest), "HTTPRequest", "httpgrpc.HTTPRequest", 1) + `,`,
		`QueryRangeRequest:` + strings.Replace(fmt.Sprintf("%v", this.QueryRangeRequest), "Request", "queryrange.Request", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ProcessResponse) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ProcessResponse{`,
		`HttpResponse:` + strings.Replace(fmt.Sprintf("%v", this.HttpResponse), "HTTPResponse", "httpgrpc.HTTPResponse", 1) + `,`,
		`ApiResponse:` + strings.Replace(fmt.Sprintf("%v", this.ApiResponse), "APIResponse", "queryrange.APIResponse", 1) + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringFrontend(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *ProcessRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFrontend
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProcessRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HttpRequest", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFrontend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFrontend
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFrontend
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.HttpRequest == nil {
				m.HttpRequest = &httpgrpc.HTTPRequest{}
			}
			if err := m.HttpRequest.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field QueryRangeRequest", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFrontend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFrontend
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFrontend
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.QueryRangeRequest == nil {
				m.QueryRangeRequest = &queryrange.Request{}
			}
			if err := m.QueryRangeRequest.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFrontend(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFrontend
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthFrontend
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ProcessResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowFrontend
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ProcessResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ProcessResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field HttpResponse", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFrontend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFrontend
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFrontend
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.HttpResponse == nil {
				m.HttpResponse = &httpgrpc.HTTPResponse{}
			}
			if err := m.HttpResponse.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ApiResponse", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowFrontend
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthFrontend
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthFrontend
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.ApiResponse == nil {
				m.ApiResponse = &queryrange.APIResponse{}
			}
			if err := m.ApiResponse.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipFrontend(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthFrontend
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthFrontend
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipFrontend(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowFrontend
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFrontend
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
			return iNdEx, nil
		case 1:
			iNdEx += 8
			return iNdEx, nil
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowFrontend
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthFrontend
			}
			iNdEx += length
			if iNdEx < 0 {
				return 0, ErrInvalidLengthFrontend
			}
			return iNdEx, nil
		case 3:
			for {
				var innerWire uint64
				var start int = iNdEx
				for shift := uint(0); ; shift += 7 {
					if shift >= 64 {
						return 0, ErrIntOverflowFrontend
					}
					if iNdEx >= l {
						return 0, io.ErrUnexpectedEOF
					}
					b := dAtA[iNdEx]
					iNdEx++
					innerWire |= (uint64(b) & 0x7F) << shift
					if b < 0x80 {
						break
					}
				}
				innerWireType := int(innerWire & 0x7)
				if innerWireType == 4 {
					break
				}
				next, err := skipFrontend(dAtA[start:])
				if err != nil {
					return 0, err
				}
				iNdEx = start + next
				if iNdEx < 0 {
					return 0, ErrInvalidLengthFrontend
				}
			}
			return iNdEx, nil
		case 4:
			return iNdEx, nil
		case 5:
			iNdEx += 4
			return iNdEx, nil
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
	}
	panic("unreachable")
}

var (
	ErrInvalidLengthFrontend = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowFrontend   = fmt.Errorf("proto: integer overflow")
)
