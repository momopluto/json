package json_test

import (
	officialjson "encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/momopluto/json"
	"sync"
	"testing"
)

/*
// 模拟 proto 文件
// mock.proto

enum MockPbEnum
{
	PS_MNG = 0;
	PS_MEMBER = 1;
	PS_MM = 2;
}
message MockInnerPbStruct
{
	optional string 			key 	= 1;
	optional string 			value 	= 2;
	optional string 			def 	= 3 [default = "def-str"];
}
message MockPbStruct
{
	optional int32 				int32_with_def 		= 1 [default = 10];
	optional int32 				int32_with_no_def 	= 2;
	optional bool  				bool_with_def 		= 3 [default = true];
	optional bool  				bool_with_no_def 	= 4;
	optional string 			string_with_def 	= 5 [default = "test-string"];
	optional string 			string_with_no_def 	= 6;
	optional MockPbEnum 		enum_with_def 		= 7 [default = PS_MNG];
	optional MockPbEnum 		enum_with_no_def 	= 8;

	repeated int64 				int64_slice 		= 9;
	repeated MockInnerPbStruct 	struct_slice 		= 10;
	optional MockInnerPbStruct 	struct2 			= 11;
}
*/

// -----------------------------------------------------------------------------------------

// 模拟 proto 编译结果文件
// mock.pb.go

type MockPbEnum int32

const (
	MockPbEnum_PS_MNG    MockPbEnum = 0
	MockPbEnum_PS_MEMBER MockPbEnum = 1
	MockPbEnum_PS_MM     MockPbEnum = 2
)

var MockPbEnum_name = map[int32]string{
	0: "PS_MNG",
	1: "PS_MEMBER",
	2: "PS_MM",
}
var MockPbEnum_value = map[string]int32{
	"PS_MNG":    0,
	"PS_MEMBER": 1,
	"PS_MM":     2,
}

func (x MockPbEnum) Enum() *MockPbEnum {
	p := new(MockPbEnum)
	*p = x
	return p
}
func (x MockPbEnum) String() string {
	return proto.EnumName(MockPbEnum_name, int32(x))
}
func (x *MockPbEnum) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MockPbEnum_value, data, "MockPbEnum")
	if err != nil {
		return err
	}
	*x = MockPbEnum(value)
	return nil
}

type MockInnerPbStruct struct {
	Key              *string `protobuf:"bytes,1,opt,name=key" json:"key,omitempty"`
	Value            *string `protobuf:"bytes,2,opt,name=value" json:"value,omitempty"`
	Def              *string `protobuf:"bytes,3,opt,name=def,def=def-str" json:"def,omitempty"`
	XXX_unrecognized []byte  `json:"-" bson:"-"`
}

func (m *MockInnerPbStruct) Reset()         { *m = MockInnerPbStruct{} }
func (m *MockInnerPbStruct) String() string { return proto.CompactTextString(m) }
func (*MockInnerPbStruct) ProtoMessage()    {}

type MockPbStruct struct {
	// 问题1: pb default 值转 json 会丢失
	Int32WithDef    *int32      `protobuf:"varint,1,opt,name=int32_with_def,def=10" json:"int32_with_def,omitempty"`
	Int32WithNoDef  *int32      `protobuf:"varint,2,opt,name=int32_with_no_def" json:"int32_with_no_def,omitempty"`
	BoolWithDef     *bool       `protobuf:"varint,3,opt,name=bool_with_def,def=1" json:"bool_with_def,omitempty"`
	BoolWithNoDef   *bool       `protobuf:"varint,4,opt,name=bool_with_no_def" json:"bool_with_no_def,omitempty"`
	StringWithDef   *string     `protobuf:"bytes,5,opt,name=string_with_def,def=test-string" json:"string_with_def,omitempty"`
	StringWithNoDef *string     `protobuf:"bytes,6,opt,name=string_with_no_def" json:"string_with_no_def,omitempty"`
	EnumWithDef     *MockPbEnum `protobuf:"varint,7,opt,name=enum_with_def,enum=MockPbEnum,def=1" json:"enum_with_def,omitempty"`
	EnumWithNoDef   *MockPbEnum `protobuf:"varint,8,opt,name=enum_with_no_def,enum=MockPbEnum" json:"enum_with_no_def,omitempty"`

	// 问题2: 空 slice 转 json 被 omitempty 忽略, 不返回 []
	//	json.Marshal() 即使不 omitempty 也只会返回 null; 使用 json.MarshalSafeCollections() 解决返回 []
	Int64Slice       []int64              `protobuf:"varint,9,rep,name=int64_slice" json:"int64_slice,omitempty"`
	StructSlice      []*MockInnerPbStruct `protobuf:"bytes,10,rep,name=struct_slice" json:"struct_slice,omitempty"`
	Struct2          *MockInnerPbStruct   `protobuf:"bytes,11,opt,name=struct2" json:"struct2,omitempty"`
	XXX_unrecognized []byte               `json:"-" bson:"-"`
}

func (m *MockPbStruct) Reset()         { *m = MockPbStruct{} }
func (m *MockPbStruct) String() string { return proto.CompactTextString(m) }
func (*MockPbStruct) ProtoMessage()    {}

// -----------------------------------------------------------------------------------------

func TestNilSliceIgnoreOmitempty(t *testing.T) {
	testCases := []string{
		`{}`,
		`{"int64_slice":[],"struct_slice":[],"struct2":null}`, // 和上述结果一致
		// 注意: "struct2":{} 转成 pb 会分配内存空间不算空, "struct2":null 才算空
		`{"struct2":{}}`,
	}

	once := sync.Once{}
	once.Do(func() {
		json.Init(true) // 只能初始化1次. 因为 json 处理每个字段有缓存, 之后再调用此函数无法保证结果
	})

	for _, str := range testCases {
		fmt.Printf("           input:\t %s\n", str)
		rsp := &MockPbStruct{}
		_ = officialjson.Unmarshal([]byte(str), rsp)

		rspByte1, _ := officialjson.Marshal(rsp)
		fmt.Printf("omitempty work  :\t %s\n", string(rspByte1))

		rspByte3, _ := json.MarshalSafeCollections(rsp)
		fmt.Printf("ignore omitempty:\t %s\n", string(rspByte3))

		fmt.Println("---------------------------")
	}
}

/*
=== RUN   TestNilSliceIgnoreOmitempty
		  input:	 {}
omitempty work  :	 {}
ignore omitempty:	 {"int64_slice":[],"struct_slice":[]}
---------------------------
		  input:	 {"int64_slice":[],"struct_slice":[],"struct2":null}
omitempty work  :	 {}
ignore omitempty:	 {"int64_slice":[],"struct_slice":[]}
---------------------------
		  input:	 {"struct2":{}}
omitempty work  :	 {"struct2":{}}
ignore omitempty:	 {"int64_slice":[],"struct_slice":[],"struct2":{}}
---------------------------
--- PASS: TestNilSliceIgnoreOmitempty (0.00s)
PASS
*/
