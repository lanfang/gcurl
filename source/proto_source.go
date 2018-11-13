package source

import "github.com/jhump/protoreflect/desc"

type protoSource struct {
	proto   string
	imports []string
}

func NewProtoSource(proto string, depPath []string) (Source, error) {

	return nil, nil
}

func (src *protoSource) ListServices() ([]string, error) {
	return nil, nil
}
func (src *protoSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	return nil, nil
}
func (src *protoSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	return nil, nil
}
