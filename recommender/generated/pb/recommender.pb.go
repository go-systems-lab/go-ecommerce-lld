// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: recommender.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProductInteraction struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	UserId          string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	ProductId       string                 `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	InteractionType string                 `protobuf:"bytes,3,opt,name=interaction_type,json=interactionType,proto3" json:"interaction_type,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *ProductInteraction) Reset() {
	*x = ProductInteraction{}
	mi := &file_recommender_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductInteraction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductInteraction) ProtoMessage() {}

func (x *ProductInteraction) ProtoReflect() protoreflect.Message {
	mi := &file_recommender_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductInteraction.ProtoReflect.Descriptor instead.
func (*ProductInteraction) Descriptor() ([]byte, []int) {
	return file_recommender_proto_rawDescGZIP(), []int{0}
}

func (x *ProductInteraction) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *ProductInteraction) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *ProductInteraction) GetInteractionType() string {
	if x != nil {
		return x.InteractionType
	}
	return ""
}

type RecommendationRequestForUserId struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	UserId        string                 `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Skip          uint64                 `protobuf:"varint,2,opt,name=skip,proto3" json:"skip,omitempty"`
	Take          uint64                 `protobuf:"varint,3,opt,name=take,proto3" json:"take,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecommendationRequestForUserId) Reset() {
	*x = RecommendationRequestForUserId{}
	mi := &file_recommender_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecommendationRequestForUserId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecommendationRequestForUserId) ProtoMessage() {}

func (x *RecommendationRequestForUserId) ProtoReflect() protoreflect.Message {
	mi := &file_recommender_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecommendationRequestForUserId.ProtoReflect.Descriptor instead.
func (*RecommendationRequestForUserId) Descriptor() ([]byte, []int) {
	return file_recommender_proto_rawDescGZIP(), []int{1}
}

func (x *RecommendationRequestForUserId) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *RecommendationRequestForUserId) GetSkip() uint64 {
	if x != nil {
		return x.Skip
	}
	return 0
}

func (x *RecommendationRequestForUserId) GetTake() uint64 {
	if x != nil {
		return x.Take
	}
	return 0
}

type RecommendationRequestOnViews struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Ids           []string               `protobuf:"bytes,1,rep,name=ids,proto3" json:"ids,omitempty"`
	Skip          uint64                 `protobuf:"varint,2,opt,name=skip,proto3" json:"skip,omitempty"`
	Take          uint64                 `protobuf:"varint,3,opt,name=take,proto3" json:"take,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RecommendationRequestOnViews) Reset() {
	*x = RecommendationRequestOnViews{}
	mi := &file_recommender_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecommendationRequestOnViews) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecommendationRequestOnViews) ProtoMessage() {}

func (x *RecommendationRequestOnViews) ProtoReflect() protoreflect.Message {
	mi := &file_recommender_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecommendationRequestOnViews.ProtoReflect.Descriptor instead.
func (*RecommendationRequestOnViews) Descriptor() ([]byte, []int) {
	return file_recommender_proto_rawDescGZIP(), []int{2}
}

func (x *RecommendationRequestOnViews) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *RecommendationRequestOnViews) GetSkip() uint64 {
	if x != nil {
		return x.Skip
	}
	return 0
}

func (x *RecommendationRequestOnViews) GetTake() uint64 {
	if x != nil {
		return x.Take
	}
	return 0
}

type ProductReplica struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Description   string                 `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Price         float64                `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProductReplica) Reset() {
	*x = ProductReplica{}
	mi := &file_recommender_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductReplica) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductReplica) ProtoMessage() {}

func (x *ProductReplica) ProtoReflect() protoreflect.Message {
	mi := &file_recommender_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductReplica.ProtoReflect.Descriptor instead.
func (*ProductReplica) Descriptor() ([]byte, []int) {
	return file_recommender_proto_rawDescGZIP(), []int{3}
}

func (x *ProductReplica) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ProductReplica) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProductReplica) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *ProductReplica) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

type RecommendationResponse struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	RecommendedProducts []*ProductReplica      `protobuf:"bytes,1,rep,name=recommended_products,json=recommendedProducts,proto3" json:"recommended_products,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *RecommendationResponse) Reset() {
	*x = RecommendationResponse{}
	mi := &file_recommender_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RecommendationResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecommendationResponse) ProtoMessage() {}

func (x *RecommendationResponse) ProtoReflect() protoreflect.Message {
	mi := &file_recommender_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecommendationResponse.ProtoReflect.Descriptor instead.
func (*RecommendationResponse) Descriptor() ([]byte, []int) {
	return file_recommender_proto_rawDescGZIP(), []int{4}
}

func (x *RecommendationResponse) GetRecommendedProducts() []*ProductReplica {
	if x != nil {
		return x.RecommendedProducts
	}
	return nil
}

var File_recommender_proto protoreflect.FileDescriptor

const file_recommender_proto_rawDesc = "" +
	"\n" +
	"\x11recommender.proto\x12\x02pb\x1a\x1bgoogle/protobuf/empty.proto\"w\n" +
	"\x12ProductInteraction\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x1d\n" +
	"\n" +
	"product_id\x18\x02 \x01(\tR\tproductId\x12)\n" +
	"\x10interaction_type\x18\x03 \x01(\tR\x0finteractionType\"a\n" +
	"\x1eRecommendationRequestForUserId\x12\x17\n" +
	"\auser_id\x18\x01 \x01(\tR\x06userId\x12\x12\n" +
	"\x04skip\x18\x02 \x01(\x04R\x04skip\x12\x12\n" +
	"\x04take\x18\x03 \x01(\x04R\x04take\"X\n" +
	"\x1cRecommendationRequestOnViews\x12\x10\n" +
	"\x03ids\x18\x01 \x03(\tR\x03ids\x12\x12\n" +
	"\x04skip\x18\x02 \x01(\x04R\x04skip\x12\x12\n" +
	"\x04take\x18\x03 \x01(\x04R\x04take\"l\n" +
	"\x0eProductReplica\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12 \n" +
	"\vdescription\x18\x03 \x01(\tR\vdescription\x12\x14\n" +
	"\x05price\x18\x04 \x01(\x01R\x05price\"_\n" +
	"\x16RecommendationResponse\x12E\n" +
	"\x14recommended_products\x18\x01 \x03(\v2\x12.pb.ProductReplicaR\x13recommendedProducts2\x93\x02\n" +
	"\x12RecommenderService\x12]\n" +
	"\x1bGetRecommendationsForUserId\x12\".pb.RecommendationRequestForUserId\x1a\x1a.pb.RecommendationResponse\x12Y\n" +
	"\x19GetRecommendationsOnViews\x12 .pb.RecommendationRequestOnViews\x1a\x1a.pb.RecommendationResponse\x12C\n" +
	"\x11RecordInteraction\x12\x16.pb.ProductInteraction\x1a\x16.google.protobuf.EmptyB\x06Z\x04./pbb\x06proto3"

var (
	file_recommender_proto_rawDescOnce sync.Once
	file_recommender_proto_rawDescData []byte
)

func file_recommender_proto_rawDescGZIP() []byte {
	file_recommender_proto_rawDescOnce.Do(func() {
		file_recommender_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_recommender_proto_rawDesc), len(file_recommender_proto_rawDesc)))
	})
	return file_recommender_proto_rawDescData
}

var file_recommender_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_recommender_proto_goTypes = []any{
	(*ProductInteraction)(nil),             // 0: pb.ProductInteraction
	(*RecommendationRequestForUserId)(nil), // 1: pb.RecommendationRequestForUserId
	(*RecommendationRequestOnViews)(nil),   // 2: pb.RecommendationRequestOnViews
	(*ProductReplica)(nil),                 // 3: pb.ProductReplica
	(*RecommendationResponse)(nil),         // 4: pb.RecommendationResponse
	(*emptypb.Empty)(nil),                  // 5: google.protobuf.Empty
}
var file_recommender_proto_depIdxs = []int32{
	3, // 0: pb.RecommendationResponse.recommended_products:type_name -> pb.ProductReplica
	1, // 1: pb.RecommenderService.GetRecommendationsForUserId:input_type -> pb.RecommendationRequestForUserId
	2, // 2: pb.RecommenderService.GetRecommendationsOnViews:input_type -> pb.RecommendationRequestOnViews
	0, // 3: pb.RecommenderService.RecordInteraction:input_type -> pb.ProductInteraction
	4, // 4: pb.RecommenderService.GetRecommendationsForUserId:output_type -> pb.RecommendationResponse
	4, // 5: pb.RecommenderService.GetRecommendationsOnViews:output_type -> pb.RecommendationResponse
	5, // 6: pb.RecommenderService.RecordInteraction:output_type -> google.protobuf.Empty
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_recommender_proto_init() }
func file_recommender_proto_init() {
	if File_recommender_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_recommender_proto_rawDesc), len(file_recommender_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_recommender_proto_goTypes,
		DependencyIndexes: file_recommender_proto_depIdxs,
		MessageInfos:      file_recommender_proto_msgTypes,
	}.Build()
	File_recommender_proto = out.File
	file_recommender_proto_goTypes = nil
	file_recommender_proto_depIdxs = nil
}
