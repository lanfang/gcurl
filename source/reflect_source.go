package source

import (
	"fmt"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/grpcreflect"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pbclient "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"net"
	"time"
)

type reflectSource struct {
	addr   string
	client *grpcreflect.Client
}

func NewConn(ctx context.Context, addr string) (*grpc.ClientConn, error) {

	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		conn, err := (&net.Dialer{}).DialContext(ctx, "tcp", address)
		if err != nil {
			return nil, err
		}
		return conn, nil
	}

	opts := make([]grpc.DialOption, 0)
	opts = append(opts,
		grpc.WithBlock(),
		grpc.WithDialer(dialer),
		grpc.WithInsecure(),
	)
	ctxTimeout, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctxTimeout, addr, opts...)
	return conn, err
}

func NewReflectSource(addr string) (Source, error) {
	s := &reflectSource{
		addr: addr,
	}
	ctx := context.Background()
	var retErr error
	for loop := true; loop; loop = false {
		conn, err := NewConn(ctx, addr)
		if err != nil {
			retErr = err
			break
		}
		s.client = grpcreflect.NewClient(ctx, pbclient.NewServerReflectionClient(conn))
	}
	return s, retErr
}

func (s *reflectSource) ListServices() ([]string, error) {
	return s.client.ListServices()
}
func (s *reflectSource) FindSymbol(fqn string) (desc.Descriptor, error) {
	file, err := s.client.FileContainingSymbol(fqn)
	if err != nil {
		return nil, err
	}
	d := file.FindSymbol(fqn)
	if d == nil {
		return nil, fmt.Errorf("not found %s", fqn)
	}
	return d, nil
}
func (s *reflectSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	return nil, nil
}
