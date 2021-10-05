// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: provenance/msgfees/v1/msgfees.proto

package types

import (
	fmt "fmt"
	types "github.com/cosmos/cosmos-sdk/codec/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types1 "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/x/auth/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	_ "github.com/regen-network/cosmos-proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the set of params for the msgfees module.
type Params struct {
	// indicates if governance based controls of msgFees is allowed.
	EnableGovernance bool `protobuf:"varint,1,opt,name=enable_governance,json=enableGovernance,proto3" json:"enable_governance,omitempty"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6265859d114362, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetEnableGovernance() bool {
	if m != nil {
		return m.EnableGovernance
	}
	return false
}

// MsgFees is the core of what gets stored on the blockchain
// it consists of two parts
// 1. minimum additional fees(can be of any denom)
// 2. Fee rate which is proportional to the gas charged for processing that message.
type MsgFees struct {
	Msg *types.Any `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	// can pay in any Coin( basically a Denom and Amount, Amount can be zero)
	MinAdditionalFee github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,2,rep,name=min_additional_fee,json=minAdditionalFee,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"min_additional_fee" yaml:"min_additional_fee"`
	//  Fee rate, based on Gas used.
	FeeRate github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=fee_rate,json=feeRate,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"fee_rate,omitempty"`
}

func (m *MsgFees) Reset()         { *m = MsgFees{} }
func (m *MsgFees) String() string { return proto.CompactTextString(m) }
func (*MsgFees) ProtoMessage()    {}
func (*MsgFees) Descriptor() ([]byte, []int) {
	return fileDescriptor_0c6265859d114362, []int{1}
}
func (m *MsgFees) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgFees) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgFees.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgFees) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgFees.Merge(m, src)
}
func (m *MsgFees) XXX_Size() int {
	return m.Size()
}
func (m *MsgFees) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgFees.DiscardUnknown(m)
}

var xxx_messageInfo_MsgFees proto.InternalMessageInfo

func (m *MsgFees) GetMsg() *types.Any {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *MsgFees) GetMinAdditionalFee() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.MinAdditionalFee
	}
	return nil
}

func init() {
	proto.RegisterType((*Params)(nil), "provenance.msgfees.v1.Params")
	proto.RegisterType((*MsgFees)(nil), "provenance.msgfees.v1.MsgFees")
}

func init() {
	proto.RegisterFile("provenance/msgfees/v1/msgfees.proto", fileDescriptor_0c6265859d114362)
}

var fileDescriptor_0c6265859d114362 = []byte{
	// 466 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0x31, 0x6b, 0xdb, 0x40,
	0x14, 0xc7, 0xa5, 0x18, 0x62, 0x23, 0x77, 0x48, 0x45, 0x0a, 0x76, 0x28, 0x92, 0x71, 0xa1, 0x18,
	0x5a, 0xdf, 0xe1, 0x78, 0xcb, 0x52, 0xe2, 0x86, 0x74, 0x32, 0x04, 0x8d, 0x5d, 0xc4, 0x49, 0x7e,
	0xba, 0x1c, 0xf1, 0xdd, 0x09, 0xdd, 0x59, 0x54, 0xdf, 0xa2, 0x53, 0xe9, 0x98, 0xb9, 0x73, 0x87,
	0x7e, 0x84, 0xd0, 0x29, 0x63, 0xe9, 0xe0, 0x16, 0x7b, 0x29, 0x1d, 0xfb, 0x09, 0x8a, 0x4e, 0x92,
	0x13, 0x9a, 0xa5, 0x93, 0xee, 0xe9, 0xf7, 0x7f, 0xef, 0xff, 0xd7, 0xe9, 0x39, 0xcf, 0xd2, 0x4c,
	0xe6, 0x20, 0x88, 0x88, 0x01, 0x73, 0x45, 0x13, 0x00, 0x85, 0xf3, 0x49, 0x73, 0x44, 0x69, 0x26,
	0xb5, 0x74, 0x9f, 0xdc, 0x89, 0x50, 0x43, 0xf2, 0xc9, 0xd1, 0x21, 0x95, 0x54, 0x1a, 0x05, 0x2e,
	0x4f, 0x95, 0xf8, 0xc8, 0x8b, 0xa5, 0xe2, 0x52, 0x61, 0xb2, 0xd2, 0x97, 0x38, 0x9f, 0x44, 0xa0,
	0xc9, 0xc4, 0x14, 0x35, 0xef, 0x57, 0x3c, 0xac, 0x1a, 0xab, 0xa2, 0x41, 0x54, 0x4a, 0xba, 0x04,
	0x6c, 0xaa, 0x68, 0x95, 0x60, 0x22, 0x8a, 0x1a, 0x3d, 0xad, 0x11, 0x49, 0x19, 0x26, 0x42, 0x48,
	0x4d, 0x34, 0x93, 0x42, 0xfd, 0xe3, 0x19, 0x11, 0x05, 0x3b, 0xcf, 0x58, 0x32, 0x51, 0xf1, 0xe1,
	0x2b, 0x67, 0xff, 0x82, 0x64, 0x84, 0x2b, 0xf7, 0x85, 0xf3, 0x18, 0x04, 0x89, 0x96, 0x10, 0x52,
	0x99, 0x43, 0x66, 0xbe, 0xa9, 0x67, 0x0f, 0xec, 0x51, 0x27, 0x38, 0xa8, 0xc0, 0x9b, 0xdd, 0xfb,
	0x93, 0xce, 0xc7, 0x6b, 0xdf, 0xfa, 0x75, 0xed, 0x5b, 0xc3, 0x2f, 0x7b, 0x4e, 0x7b, 0xae, 0xe8,
	0x39, 0x80, 0x72, 0xa7, 0x4e, 0x8b, 0x2b, 0x6a, 0x9a, 0xba, 0xc7, 0x87, 0xa8, 0x0a, 0x86, 0x9a,
	0xcc, 0xe8, 0x54, 0x14, 0xb3, 0xee, 0xd7, 0xcf, 0xe3, 0xb6, 0x5a, 0x5c, 0xa1, 0xb9, 0xa2, 0x41,
	0xa9, 0x76, 0x3f, 0xd8, 0x8e, 0xcb, 0x99, 0x08, 0xc9, 0x62, 0xc1, 0xca, 0xe4, 0x64, 0x19, 0x26,
	0x00, 0xbd, 0xbd, 0x41, 0x6b, 0xd4, 0x3d, 0xee, 0xa3, 0xfa, 0x1a, 0xca, 0xfc, 0xa8, 0xce, 0x8f,
	0x5e, 0x4b, 0x26, 0x66, 0xf3, 0x9b, 0xb5, 0x6f, 0xfd, 0x59, 0xfb, 0xfd, 0x82, 0xf0, 0xe5, 0xc9,
	0xf0, 0xe1, 0x88, 0xe1, 0xa7, 0x1f, 0xfe, 0x88, 0x32, 0x7d, 0xb9, 0x8a, 0x50, 0x2c, 0x79, 0x7d,
	0xa1, 0xf5, 0x63, 0xac, 0x16, 0x57, 0x58, 0x17, 0x29, 0x28, 0x33, 0x4d, 0x05, 0x07, 0x9c, 0x89,
	0xd3, 0x5d, 0xff, 0x39, 0x80, 0x1b, 0x3a, 0x9d, 0x04, 0x20, 0xcc, 0x88, 0x86, 0x5e, 0x6b, 0x60,
	0x8f, 0x1e, 0xcd, 0xce, 0x4a, 0xcb, 0xef, 0x6b, 0xff, 0xf9, 0x7f, 0x4c, 0x3d, 0x83, 0xf8, 0xf7,
	0xda, 0x77, 0x9b, 0x09, 0x2f, 0x25, 0x67, 0x1a, 0x78, 0xaa, 0x8b, 0xa0, 0x9d, 0x00, 0x04, 0x44,
	0xc3, 0x8c, 0xdd, 0x6c, 0x3c, 0xfb, 0x76, 0xe3, 0xd9, 0x3f, 0x37, 0x9e, 0xfd, 0x7e, 0xeb, 0x59,
	0xb7, 0x5b, 0xcf, 0xfa, 0xb6, 0xf5, 0x2c, 0xa7, 0xc7, 0xcc, 0xce, 0x3c, 0xdc, 0xac, 0x0b, 0xfb,
	0xed, 0xf4, 0x9e, 0xf9, 0x9d, 0x66, 0xcc, 0xe4, 0xbd, 0x0a, 0xbf, 0xdb, 0xad, 0xac, 0x49, 0x13,
	0xed, 0x9b, 0x9f, 0x30, 0xfd, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x92, 0xf0, 0x79, 0x67, 0xd5, 0x02,
	0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.EnableGovernance {
		i--
		if m.EnableGovernance {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *MsgFees) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgFees) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgFees) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.FeeRate.Size()
		i -= size
		if _, err := m.FeeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMsgfees(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.MinAdditionalFee) > 0 {
		for iNdEx := len(m.MinAdditionalFee) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.MinAdditionalFee[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintMsgfees(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if m.Msg != nil {
		{
			size, err := m.Msg.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintMsgfees(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMsgfees(dAtA []byte, offset int, v uint64) int {
	offset -= sovMsgfees(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.EnableGovernance {
		n += 2
	}
	return n
}

func (m *MsgFees) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Msg != nil {
		l = m.Msg.Size()
		n += 1 + l + sovMsgfees(uint64(l))
	}
	if len(m.MinAdditionalFee) > 0 {
		for _, e := range m.MinAdditionalFee {
			l = e.Size()
			n += 1 + l + sovMsgfees(uint64(l))
		}
	}
	l = m.FeeRate.Size()
	n += 1 + l + sovMsgfees(uint64(l))
	return n
}

func sovMsgfees(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMsgfees(x uint64) (n int) {
	return sovMsgfees(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgfees
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EnableGovernance", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgfees
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.EnableGovernance = bool(v != 0)
		default:
			iNdEx = preIndex
			skippy, err := skipMsgfees(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsgfees
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
func (m *MsgFees) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMsgfees
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
			return fmt.Errorf("proto: MsgFees: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgFees: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgfees
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
				return ErrInvalidLengthMsgfees
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgfees
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Msg == nil {
				m.Msg = &types.Any{}
			}
			if err := m.Msg.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinAdditionalFee", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgfees
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
				return ErrInvalidLengthMsgfees
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMsgfees
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MinAdditionalFee = append(m.MinAdditionalFee, types1.Coin{})
			if err := m.MinAdditionalFee[len(m.MinAdditionalFee)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field FeeRate", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMsgfees
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMsgfees
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMsgfees
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.FeeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMsgfees(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMsgfees
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

func skipMsgfees(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMsgfees
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
					return 0, ErrIntOverflowMsgfees
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowMsgfees
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
				return 0, ErrInvalidLengthMsgfees
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMsgfees
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMsgfees
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMsgfees        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMsgfees          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMsgfees = fmt.Errorf("proto: unexpected end of group")
)