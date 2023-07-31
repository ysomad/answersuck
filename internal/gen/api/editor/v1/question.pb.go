// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.4
// source: editor/v1/question.proto

package editorv1

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Answer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       int32  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Text     string `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	MediaUrl string `protobuf:"bytes,3,opt,name=media_url,json=mediaUrl,proto3" json:"media_url,omitempty"`
}

func (x *Answer) Reset() {
	*x = Answer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_editor_v1_question_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Answer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Answer) ProtoMessage() {}

func (x *Answer) ProtoReflect() protoreflect.Message {
	mi := &file_editor_v1_question_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Answer.ProtoReflect.Descriptor instead.
func (*Answer) Descriptor() ([]byte, []int) {
	return file_editor_v1_question_proto_rawDescGZIP(), []int{0}
}

func (x *Answer) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Answer) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Answer) GetMediaUrl() string {
	if x != nil {
		return x.MediaUrl
	}
	return ""
}

type Question struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int32                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Text       string                 `protobuf:"bytes,2,opt,name=text,proto3" json:"text,omitempty"`
	Answer     *Answer                `protobuf:"bytes,3,opt,name=answer,proto3" json:"answer,omitempty"`
	Author     string                 `protobuf:"bytes,4,opt,name=author,proto3" json:"author,omitempty"`
	MediaUrl   string                 `protobuf:"bytes,5,opt,name=media_url,json=mediaUrl,proto3" json:"media_url,omitempty"`
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,50,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
}

func (x *Question) Reset() {
	*x = Question{}
	if protoimpl.UnsafeEnabled {
		mi := &file_editor_v1_question_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Question) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Question) ProtoMessage() {}

func (x *Question) ProtoReflect() protoreflect.Message {
	mi := &file_editor_v1_question_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Question.ProtoReflect.Descriptor instead.
func (*Question) Descriptor() ([]byte, []int) {
	return file_editor_v1_question_proto_rawDescGZIP(), []int{1}
}

func (x *Question) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Question) GetText() string {
	if x != nil {
		return x.Text
	}
	return ""
}

func (x *Question) GetAnswer() *Answer {
	if x != nil {
		return x.Answer
	}
	return nil
}

func (x *Question) GetAuthor() string {
	if x != nil {
		return x.Author
	}
	return ""
}

func (x *Question) GetMediaUrl() string {
	if x != nil {
		return x.MediaUrl
	}
	return ""
}

func (x *Question) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

type CreateQuestionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Question         string `protobuf:"bytes,1,opt,name=question,proto3" json:"question,omitempty"` // required
	QuestionMediaUrl string `protobuf:"bytes,2,opt,name=question_media_url,json=questionMediaUrl,proto3" json:"question_media_url,omitempty"`
	Answer           string `protobuf:"bytes,3,opt,name=answer,proto3" json:"answer,omitempty"` // required
	AnswerMediaUrl   string `protobuf:"bytes,4,opt,name=answer_media_url,json=answerMediaUrl,proto3" json:"answer_media_url,omitempty"`
}

func (x *CreateQuestionRequest) Reset() {
	*x = CreateQuestionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_editor_v1_question_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateQuestionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateQuestionRequest) ProtoMessage() {}

func (x *CreateQuestionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_editor_v1_question_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateQuestionRequest.ProtoReflect.Descriptor instead.
func (*CreateQuestionRequest) Descriptor() ([]byte, []int) {
	return file_editor_v1_question_proto_rawDescGZIP(), []int{2}
}

func (x *CreateQuestionRequest) GetQuestion() string {
	if x != nil {
		return x.Question
	}
	return ""
}

func (x *CreateQuestionRequest) GetQuestionMediaUrl() string {
	if x != nil {
		return x.QuestionMediaUrl
	}
	return ""
}

func (x *CreateQuestionRequest) GetAnswer() string {
	if x != nil {
		return x.Answer
	}
	return ""
}

func (x *CreateQuestionRequest) GetAnswerMediaUrl() string {
	if x != nil {
		return x.AnswerMediaUrl
	}
	return ""
}

type CreateQuestionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QuestionId int32 `protobuf:"varint,1,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"`
}

func (x *CreateQuestionResponse) Reset() {
	*x = CreateQuestionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_editor_v1_question_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateQuestionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateQuestionResponse) ProtoMessage() {}

func (x *CreateQuestionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_editor_v1_question_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateQuestionResponse.ProtoReflect.Descriptor instead.
func (*CreateQuestionResponse) Descriptor() ([]byte, []int) {
	return file_editor_v1_question_proto_rawDescGZIP(), []int{3}
}

func (x *CreateQuestionResponse) GetQuestionId() int32 {
	if x != nil {
		return x.QuestionId
	}
	return 0
}

type GetQuestionRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	QuestionId int32 `protobuf:"varint,1,opt,name=question_id,json=questionId,proto3" json:"question_id,omitempty"` // required
}

func (x *GetQuestionRequest) Reset() {
	*x = GetQuestionRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_editor_v1_question_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuestionRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuestionRequest) ProtoMessage() {}

func (x *GetQuestionRequest) ProtoReflect() protoreflect.Message {
	mi := &file_editor_v1_question_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuestionRequest.ProtoReflect.Descriptor instead.
func (*GetQuestionRequest) Descriptor() ([]byte, []int) {
	return file_editor_v1_question_proto_rawDescGZIP(), []int{4}
}

func (x *GetQuestionRequest) GetQuestionId() int32 {
	if x != nil {
		return x.QuestionId
	}
	return 0
}

type GetQuestionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Question *Question `protobuf:"bytes,1,opt,name=question,proto3" json:"question,omitempty"`
}

func (x *GetQuestionResponse) Reset() {
	*x = GetQuestionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_editor_v1_question_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetQuestionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetQuestionResponse) ProtoMessage() {}

func (x *GetQuestionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_editor_v1_question_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetQuestionResponse.ProtoReflect.Descriptor instead.
func (*GetQuestionResponse) Descriptor() ([]byte, []int) {
	return file_editor_v1_question_proto_rawDescGZIP(), []int{5}
}

func (x *GetQuestionResponse) GetQuestion() *Question {
	if x != nil {
		return x.Question
	}
	return nil
}

var File_editor_v1_question_proto protoreflect.FileDescriptor

var file_editor_v1_question_proto_rawDesc = []byte{
	0x0a, 0x18, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x64, 0x69, 0x74,
	0x6f, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x49, 0x0a, 0x06, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78,
	0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x1b, 0x0a,
	0x09, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x55, 0x72, 0x6c, 0x22, 0xcb, 0x01, 0x0a, 0x08, 0x51,
	0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x65, 0x78, 0x74, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x65, 0x78, 0x74, 0x12, 0x29, 0x0a, 0x06, 0x61,
	0x6e, 0x73, 0x77, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x65, 0x64,
	0x69, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x52, 0x06,
	0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x75, 0x74, 0x68, 0x6f, 0x72, 0x12, 0x1b,
	0x0a, 0x09, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x55, 0x72, 0x6c, 0x12, 0x3b, 0x0a, 0x0b, 0x63,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x32, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xd4, 0x01, 0x0a, 0x15, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x26, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x0a, 0xfa, 0x42, 0x07, 0x72, 0x05, 0x10, 0x03, 0x18, 0xc8, 0x01,
	0x52, 0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x12, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x08, 0x72, 0x06, 0xd0, 0x01, 0x01,
	0x88, 0x01, 0x01, 0x52, 0x10, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x65, 0x64,
	0x69, 0x61, 0x55, 0x72, 0x6c, 0x12, 0x21, 0x0a, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x09, 0xfa, 0x42, 0x06, 0x72, 0x04, 0x10, 0x03, 0x18, 0x64,
	0x52, 0x06, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x12, 0x35, 0x0a, 0x10, 0x61, 0x6e, 0x73, 0x77,
	0x65, 0x72, 0x5f, 0x6d, 0x65, 0x64, 0x69, 0x61, 0x5f, 0x75, 0x72, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x09, 0x42, 0x0b, 0xfa, 0x42, 0x08, 0x72, 0x06, 0xd0, 0x01, 0x01, 0x88, 0x01, 0x01, 0x52,
	0x0e, 0x61, 0x6e, 0x73, 0x77, 0x65, 0x72, 0x4d, 0x65, 0x64, 0x69, 0x61, 0x55, 0x72, 0x6c, 0x22,
	0x39, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x35, 0x0a, 0x12, 0x47, 0x65,
	0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x1f, 0x0a, 0x0b, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x49,
	0x64, 0x22, 0x46, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2f, 0x0a, 0x08, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x65, 0x64, 0x69,
	0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52,
	0x08, 0x71, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x32, 0xb6, 0x01, 0x0a, 0x0f, 0x51, 0x75,
	0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a,
	0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x20, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x21, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1d, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e,
	0x47, 0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x1e, 0x2e, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x47,
	0x65, 0x74, 0x51, 0x75, 0x65, 0x73, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x14, 0x5a, 0x12, 0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x3b,
	0x65, 0x64, 0x69, 0x74, 0x6f, 0x72, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_editor_v1_question_proto_rawDescOnce sync.Once
	file_editor_v1_question_proto_rawDescData = file_editor_v1_question_proto_rawDesc
)

func file_editor_v1_question_proto_rawDescGZIP() []byte {
	file_editor_v1_question_proto_rawDescOnce.Do(func() {
		file_editor_v1_question_proto_rawDescData = protoimpl.X.CompressGZIP(file_editor_v1_question_proto_rawDescData)
	})
	return file_editor_v1_question_proto_rawDescData
}

var file_editor_v1_question_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_editor_v1_question_proto_goTypes = []interface{}{
	(*Answer)(nil),                 // 0: editor.v1.Answer
	(*Question)(nil),               // 1: editor.v1.Question
	(*CreateQuestionRequest)(nil),  // 2: editor.v1.CreateQuestionRequest
	(*CreateQuestionResponse)(nil), // 3: editor.v1.CreateQuestionResponse
	(*GetQuestionRequest)(nil),     // 4: editor.v1.GetQuestionRequest
	(*GetQuestionResponse)(nil),    // 5: editor.v1.GetQuestionResponse
	(*timestamppb.Timestamp)(nil),  // 6: google.protobuf.Timestamp
}
var file_editor_v1_question_proto_depIdxs = []int32{
	0, // 0: editor.v1.Question.answer:type_name -> editor.v1.Answer
	6, // 1: editor.v1.Question.create_time:type_name -> google.protobuf.Timestamp
	1, // 2: editor.v1.GetQuestionResponse.question:type_name -> editor.v1.Question
	2, // 3: editor.v1.QuestionService.CreateQuestion:input_type -> editor.v1.CreateQuestionRequest
	4, // 4: editor.v1.QuestionService.GetQuestion:input_type -> editor.v1.GetQuestionRequest
	3, // 5: editor.v1.QuestionService.CreateQuestion:output_type -> editor.v1.CreateQuestionResponse
	5, // 6: editor.v1.QuestionService.GetQuestion:output_type -> editor.v1.GetQuestionResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_editor_v1_question_proto_init() }
func file_editor_v1_question_proto_init() {
	if File_editor_v1_question_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_editor_v1_question_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Answer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_editor_v1_question_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Question); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_editor_v1_question_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateQuestionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_editor_v1_question_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateQuestionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_editor_v1_question_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuestionRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_editor_v1_question_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetQuestionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_editor_v1_question_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_editor_v1_question_proto_goTypes,
		DependencyIndexes: file_editor_v1_question_proto_depIdxs,
		MessageInfos:      file_editor_v1_question_proto_msgTypes,
	}.Build()
	File_editor_v1_question_proto = out.File
	file_editor_v1_question_proto_rawDesc = nil
	file_editor_v1_question_proto_goTypes = nil
	file_editor_v1_question_proto_depIdxs = nil
}