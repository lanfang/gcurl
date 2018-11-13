package source

import "github.com/jhump/protoreflect/desc"

type SourceType int

const (
	SourceProto SourceType = iota + 1
	SourceReflect
)

type Source interface {
	// ListServices asks the server for the fully-qualified names of all exposed
	// services.
	ListServices() ([]string, error)
	// FindSymbol returns symbol info
	FindSymbol(fullyQualifiedName string) (desc.Descriptor, error)
	// AllExtensionsForType returns all known extension fields that extend the given message type name.
	AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error)
}
