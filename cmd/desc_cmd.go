package cmd

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	dpb "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoprint"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/lanfang/gcurl/config"
	"github.com/lanfang/gcurl/source"
	"github.com/spf13/cobra"
	"strings"
)

var printer = &protoprint.Printer{
	Compact:                  true,
	ForceFullyQualifiedNames: true,
}

func commandDesc(cmd *cobra.Command, args []string) error {
	s, err := source.NewReflectSource(config.G_Conf.Addr)
	if err != nil {
		return err
	}
	if len(config.G_Conf.SymbolList) == 0 {
		if ss, err := s.ListServices(); err != nil {
			return err
		} else {
			fmt.Println(strings.Join(ss, "\n"))
		}
	} else {
		for _, symbol := range config.G_Conf.SymbolList {
			symbol = strings.Trim(symbol, ".")
			descirber, err := s.FindSymbol(symbol)
			if err != nil {
				fmt.Printf("get symbol for %s err:%v\n", symbol, err)
				continue
			}
			str, _ := printer.PrintProtoToString(descirber)
			if str[len(str)-1] == '\n' {
				str = str[:len(str)-1]
			}
			fmt.Println(str)
			switch d := descirber.(type) {
			case *desc.MessageDescriptor:
				msg := genDynamicMsg(d)
				str, _ := PbToJson(msg)
				fmt.Println(fmt.Sprintf("%s(json):\n%s", symbol, str))
			}
		}
	}
	return nil
}

func PbToJson(pb proto.Message) (string, error) {
	marshaler := jsonpb.Marshaler{
		EmitDefaults: true,
		Indent:       "  ",
	}
	return marshaler.MarshalToString(pb)
}

func genDynamicMsg(md *desc.MessageDescriptor) proto.Message {
	dm := dynamic.NewMessage(md)
	for _, fd := range dm.GetMessageDescriptor().GetFields() {
		if fd.IsRepeated() {
			switch fd.GetType() {
			case dpb.FieldDescriptorProto_TYPE_FIXED32,
				dpb.FieldDescriptorProto_TYPE_UINT32:
				dm.AddRepeatedField(fd, uint32(0))

			case dpb.FieldDescriptorProto_TYPE_SFIXED32,
				dpb.FieldDescriptorProto_TYPE_SINT32,
				dpb.FieldDescriptorProto_TYPE_INT32,
				dpb.FieldDescriptorProto_TYPE_ENUM:
				dm.AddRepeatedField(fd, int32(0))

			case dpb.FieldDescriptorProto_TYPE_FIXED64,
				dpb.FieldDescriptorProto_TYPE_UINT64:
				dm.AddRepeatedField(fd, uint64(0))

			case dpb.FieldDescriptorProto_TYPE_SFIXED64,
				dpb.FieldDescriptorProto_TYPE_SINT64,
				dpb.FieldDescriptorProto_TYPE_INT64:
				dm.AddRepeatedField(fd, int64(0))

			case dpb.FieldDescriptorProto_TYPE_STRING:
				dm.AddRepeatedField(fd, "")

			case dpb.FieldDescriptorProto_TYPE_BYTES:
				dm.AddRepeatedField(fd, []byte{})

			case dpb.FieldDescriptorProto_TYPE_BOOL:
				dm.AddRepeatedField(fd, false)

			case dpb.FieldDescriptorProto_TYPE_FLOAT:
				dm.AddRepeatedField(fd, float32(0))

			case dpb.FieldDescriptorProto_TYPE_DOUBLE:
				dm.AddRepeatedField(fd, float64(0))

			case dpb.FieldDescriptorProto_TYPE_MESSAGE,
				dpb.FieldDescriptorProto_TYPE_GROUP:
				dm.AddRepeatedField(fd, genDynamicMsg(fd.GetMessageType()))
			}
		} else if fd.GetMessageType() != nil {
			dm.SetField(fd, genDynamicMsg(fd.GetMessageType()))
		}
	}
	return dm
}
func toJson(fields []*desc.FieldDescriptor, buf *jsonWriter) {
	buf.Write("", "{")
	for _, fd := range fields {
		t := fd.GetType()
		switch t {
		case dpb.FieldDescriptorProto_TYPE_MESSAGE:
			toJson(fd.GetMessageType().GetFields(), buf)
		default:
			buf.Write(fd.GetJSONName(), t.String())

		}
	}
	buf.Write("", "}")
}

type jsonWriter struct {
	buf bytes.Buffer
}

func (w *jsonWriter) Write(key, val string) {
	w.buf.Write([]byte(key))
	w.buf.Write([]byte(val))
	if key != "" {
		w.buf.Write([]byte(","))
	}
}

func (w *jsonWriter) Val() string {
	return string(w.buf.Bytes())
}
