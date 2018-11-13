package cmd

import (
	"bytes"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"github.com/lanfang/gcurl/config"
	"github.com/lanfang/gcurl/source"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"strings"
	"time"
)

//default command is invoke rpc metheod
func commandDefault(cmd *cobra.Command, args []string) error {
	var retErr error
	for loop := true; loop; loop = false {
		s, err := source.NewReflectSource(config.G_Conf.Addr)
		if err != nil {
			retErr = err
			break
		}
		srv := ""
		if len(config.G_Conf.SymbolList) > 0 {
			srv = strings.Trim(config.G_Conf.SymbolList[0], ".")
		}
		if srv == "" {
			retErr = fmt.Errorf("service method is empty")
			break
		}
		srvDs, err := s.FindSymbol(srv)
		if err != nil {
			retErr = err
			break
		}
		method, ok := srvDs.(*desc.MethodDescriptor)
		if !ok {
			retErr = fmt.Errorf("method %s not exist", srv)
			break
		}
		if resp, err := callRPC(config.G_Conf.Addr, method, config.G_Conf.Data); err == nil {
			jsonResp, _ := PbToJson(resp)
			fmt.Println(fmt.Sprintf("Response:\n%s", jsonResp))
		} else {
			retErr = err
		}
	}
	return retErr
}
func newPBRequest(input io.Reader, pb proto.Message) error {
	return (&jsonpb.Unmarshaler{}).Unmarshal(input, pb)

}

func callRPC(addr string, method *desc.MethodDescriptor, input string) (proto.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	msgFactory := dynamic.NewMessageFactoryWithDefaults()
	req := msgFactory.NewMessage(method.GetInputType())
	var retErr error
	var resp proto.Message
	for loop := true; loop; loop = false {
		if retErr = newPBRequest(bytes.NewReader([]byte(input)), req); retErr != nil {
			break
		}
		conn, _ := source.NewConn(ctx, addr)
		stub := grpcdynamic.NewStubWithMessageFactory(conn, msgFactory)

		resp, retErr = stub.InvokeRpc(ctx, method, req)
		if stat, ok := status.FromError(retErr); !ok {
			break
		} else {
			if stat.Code() == codes.OK {
				break
			} else {
				retErr = fmt.Errorf("call grpc %s", stat.Code())
			}
		}
	}
	return resp, retErr

}
