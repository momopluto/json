package json_test

import (
	officialjson "encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/momopluto/json"
	"reflect"
	"strconv"
	"testing"
)

/*
// 模拟 proto 文件
// mock.proto

enum MockPbEnum
{
  MC_MNG    = 0;
  MC_MEMBER = 1;
  MC_MM     = 2;
}
message MockInnerPbStruct
{
  optional string   key    = 1;
  optional string   value  = 2;
  optional string   def    = 3 [default = "def-str"];
}
message MockPbStruct
{
  optional int32              int32_with_def      = 1 [default = 10];
  optional int32              int32_with_no_def   = 2;
  optional bool               bool_with_def       = 3 [default = true];
  optional bool               bool_with_no_def    = 4;
  optional string             string_with_def     = 5 [default = "test-string"];
  optional string             string_with_no_def  = 6;
  optional MockPbEnum         enum_with_def       = 7 [default = MC_MEMBER];
  optional MockPbEnum         enum_with_no_def    = 8;

  repeated int64              int64_slice         = 9;
  repeated MockInnerPbStruct  struct_slice        = 10;
  optional MockInnerPbStruct  struct2             = 11;
}
*/

// -----------------------------------------------------------------------------------------

// 模拟 proto 编译结果文件
// mock.pb.go

type MockPbEnum int32

const (
	MockPbEnum_MC_MNG    MockPbEnum = 0
	MockPbEnum_MC_MEMBER MockPbEnum = 1
	MockPbEnum_MC_MM     MockPbEnum = 2
)

var MockPbEnum_name = map[int32]string{
	0: "MC_MNG",
	1: "MC_MEMBER",
	2: "MC_MM",
}
var MockPbEnum_value = map[string]int32{
	"MC_MNG":    0,
	"MC_MEMBER": 1,
	"MC_MM":     2,
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

const Default_MockInnerPbStruct_Def string = "def-str"

func (m *MockInnerPbStruct) GetKey() string {
	if m != nil && m.Key != nil {
		return *m.Key
	}
	return ""
}

func (m *MockInnerPbStruct) GetValue() string {
	if m != nil && m.Value != nil {
		return *m.Value
	}
	return ""
}

func (m *MockInnerPbStruct) GetDef() string {
	if m != nil && m.Def != nil {
		return *m.Def
	}
	return Default_MockInnerPbStruct_Def
}

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

const Default_MockPbStruct_Int32WithDef int32 = 10
const Default_MockPbStruct_BoolWithDef bool = true
const Default_MockPbStruct_StringWithDef string = "test-string"
const Default_MockPbStruct_EnumWithDef MockPbEnum = MockPbEnum_MC_MNG

func (m *MockPbStruct) GetInt32WithDef() int32 {
	if m != nil && m.Int32WithDef != nil {
		return *m.Int32WithDef
	}
	return Default_MockPbStruct_Int32WithDef
}

func (m *MockPbStruct) GetInt32WithNoDef() int32 {
	if m != nil && m.Int32WithNoDef != nil {
		return *m.Int32WithNoDef
	}
	return 0
}

func (m *MockPbStruct) GetBoolWithDef() bool {
	if m != nil && m.BoolWithDef != nil {
		return *m.BoolWithDef
	}
	return Default_MockPbStruct_BoolWithDef
}

func (m *MockPbStruct) GetBoolWithNoDef() bool {
	if m != nil && m.BoolWithNoDef != nil {
		return *m.BoolWithNoDef
	}
	return false
}

func (m *MockPbStruct) GetStringWithDef() string {
	if m != nil && m.StringWithDef != nil {
		return *m.StringWithDef
	}
	return Default_MockPbStruct_StringWithDef
}

func (m *MockPbStruct) GetStringWithNoDef() string {
	if m != nil && m.StringWithNoDef != nil {
		return *m.StringWithNoDef
	}
	return ""
}

func (m *MockPbStruct) GetEnumWithDef() MockPbEnum {
	if m != nil && m.EnumWithDef != nil {
		return *m.EnumWithDef
	}
	return Default_MockPbStruct_EnumWithDef
}

func (m *MockPbStruct) GetEnumWithNoDef() MockPbEnum {
	if m != nil && m.EnumWithNoDef != nil {
		return *m.EnumWithNoDef
	}
	return MockPbEnum_MC_MEMBER
}

func (m *MockPbStruct) GetInt64Slice() []int64 {
	if m != nil {
		return m.Int64Slice
	}
	return nil
}

func (m *MockPbStruct) GetStructSlice() []*MockInnerPbStruct {
	if m != nil {
		return m.StructSlice
	}
	return nil
}

func (m *MockPbStruct) GetStruct2() *MockInnerPbStruct {
	if m != nil {
		return m.Struct2
	}
	return nil
}

// -----------------------------------------------------------------------------------------

func TestNilSliceIgnoreOmitempty(t *testing.T) {
	testCases := []string{
		`{}`,
		`{"int64_slice":[],"struct_slice":[],"struct2":null}`, // 和上述结果一致
		// 注意: "struct2":{} 转成 pb 会分配内存空间不算空, "struct2":null 才算空
		`{"struct2":{}}`,
	}

	json.Init(true)

	for _, str := range testCases {
		fmt.Printf("           input:\t %s\n", str)
		rsp := &MockPbStruct{}
		_ = officialjson.Unmarshal([]byte(str), rsp)

		rspByte1, _ := officialjson.Marshal(rsp)
		fmt.Printf("omitempty work  :\t %s\n", string(rspByte1))

		json.Init(false) // no effect because of sync.Once

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

func TestFillPbDefaultVal(t *testing.T) {
	testCases := []string{
		`{}`,

		// 测试 slice 和 struct 填充 default 值
		`{"struct_slice":[{"key":"key1"},{"value":"value2"}],"struct2":{}}`,
	}

	for _, str := range testCases {
		fmt.Printf("              input:\t %s\n", str)
		rsp := &MockPbStruct{}
		_ = officialjson.Unmarshal([]byte(str), rsp)
		var pm proto.Message
		pm = rsp

		rspByte1, _ := officialjson.Marshal(rsp)
		fmt.Printf("before fill default:\t %s\n", string(rspByte1))

		FillPbDefaultVal(reflect.ValueOf(pm))

		rspByte2, _ := officialjson.Marshal(rsp)
		fmt.Printf(" after fill default:\t %s\n", string(rspByte2))

		fmt.Println("---------------------------")
	}
}

func FillPbDefaultVal(val reflect.Value) {
	// traverse until done

	sv := reflect.Indirect(val)
	if sv.Kind() == reflect.Struct {
		st := sv.Type()
		sprops := proto.GetProperties(st)
		for i := 0; i < st.NumField(); i++ {
			sf := sv.Field(i)
			switch reflect.Indirect(sf).Kind() {
			case reflect.Struct:
				// entry next level
				FillPbDefaultVal(sf)
			case reflect.Slice:
				for j := 0; j < sf.Len(); j++ {
					sfj := sf.Index(j)
					if reflect.Indirect(sfj).Kind() == reflect.Struct {
						// entry next level
						FillPbDefaultVal(sfj)
					}
				}
			default:
				props := sprops.Prop[i]
				if props.HasDefault && props.Default != "" {
					trySetDefaultValue(sf, props.Default)
				}
			}
		}
	}
}

func trySetDefaultValue(x reflect.Value, def string) bool {
	if x.IsValid() && x.Kind() == reflect.Ptr && x.IsZero() == true && def != "" && x.CanSet() {
		x.Set(reflect.New(x.Type().Elem()))
		x = reflect.Indirect(x)

		defV := reflect.Value{}

		switch x.Kind() { // 目前 proto2 协议中只处理这3种
		case reflect.Bool:
			bv := false
			if def == "1" {
				bv = true
			}
			defV = reflect.ValueOf(bv)
		case reflect.String:
			defV = reflect.ValueOf(def)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			i64, _ := strconv.ParseInt(def, 10, 64)
			defV = reflect.ValueOf(i64)
		default:
			fmt.Printf("unhandled kind = %v\n", x.Kind())
			return false
		}

		x.Set(defV.Convert(x.Type()))
		return true
	}

	return false
}

/*
=== RUN   TestFillPbDefaultVal
              input:	 {}
before fill default:	 {}
 after fill default:	 {"int32_with_def":10,"bool_with_def":true,"string_with_def":"test-string","enum_with_def":1}
---------------------------
              input:	 {"struct_slice":[{"key":"key1"},{"value":"value2"}],"struct2":{}}
before fill default:	 {"struct_slice":[{"key":"key1"},{"value":"value2"}],"struct2":{}}
 after fill default:	 {"int32_with_def":10,"bool_with_def":true,"string_with_def":"test-string","enum_with_def":1,"struct_slice":[{"key":"key1","def":"def-str"},{"value":"value2","def":"def-str"}],"struct2":{"def":"def-str"}}
---------------------------
--- PASS: TestFillPbDefaultVal (0.00s)
PASS
*/
