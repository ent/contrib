// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entproto

import (
	"errors"
	"fmt"
	"math"
	"path"
	"path/filepath"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/builder"
	"google.golang.org/protobuf/types/descriptorpb"
	_ "google.golang.org/protobuf/types/known/emptypb"
	_ "google.golang.org/protobuf/types/known/timestamppb"
	_ "google.golang.org/protobuf/types/known/wrapperspb" // needed to load wkt to global proto registry
)

const (
	DefaultProtoPackageName = "entpb"
	IDFieldNumber           = 1
)

var (
	ErrSchemaSkipped   = errors.New("entproto: schema not annotated with Generate=true")
	repeatedFieldLabel = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
	wktsPaths          = map[string]string{
		// TODO: handle more Well-Known proto types
		"google.protobuf.Timestamp":   "google/protobuf/timestamp.proto",
		"google.protobuf.Empty":       "google/protobuf/empty.proto",
		"google.protobuf.Int32Value":  "google/protobuf/wrappers.proto",
		"google.protobuf.Int64Value":  "google/protobuf/wrappers.proto",
		"google.protobuf.UInt32Value": "google/protobuf/wrappers.proto",
		"google.protobuf.UInt64Value": "google/protobuf/wrappers.proto",
		"google.protobuf.FloatValue":  "google/protobuf/wrappers.proto",
		"google.protobuf.DoubleValue": "google/protobuf/wrappers.proto",
		"google.protobuf.StringValue": "google/protobuf/wrappers.proto",
		"google.protobuf.BoolValue":   "google/protobuf/wrappers.proto",
		"google.protobuf.BytesValue":  "google/protobuf/wrappers.proto",
	}
)

// LoadAdapter takes a *gen.Graph and parses it into protobuf file descriptors
func LoadAdapter(graph *gen.Graph) (*Adapter, error) {
	a := &Adapter{
		graph:            graph,
		descriptors:      make(map[string]*desc.FileDescriptor),
		schemaProtoFiles: make(map[string]string),
		errors:           make(map[string]error),
	}
	if err := a.parse(); err != nil {
		return nil, err
	}
	return a, nil
}

// Adapter facilitates the transformation of ent gen.Type to desc.FileDescriptors
type Adapter struct {
	graph            *gen.Graph
	descriptors      map[string]*desc.FileDescriptor
	schemaProtoFiles map[string]string
	errors           map[string]error
}

// AllFileDescriptors returns a file descriptor per proto package for each package that contains
// a successfully parsed ent.Schema
func (a *Adapter) AllFileDescriptors() map[string]*desc.FileDescriptor {
	return a.descriptors
}

// GetMessageDescriptor retrieves the protobuf message descriptor for `schemaName`, if an error was returned
// while trying to parse that error they are returned
func (a *Adapter) GetMessageDescriptor(schemaName string) (*desc.MessageDescriptor, error) {
	fd, err := a.GetFileDescriptor(schemaName)
	if err != nil {
		return nil, err
	}
	findMessage := fd.FindMessage(fd.GetPackage() + "." + schemaName)
	if findMessage != nil {
		return findMessage, nil
	}
	return nil, errors.New("entproto: couldnt find message descriptor")
}

// parse transforms the ent gen.Type objects into file descriptors
func (a *Adapter) parse() error {
	var dpbDescriptors []*descriptorpb.FileDescriptorProto

	protoPackages := make(map[string]*descriptorpb.FileDescriptorProto)

	for _, genType := range a.graph.Nodes {
		messageDescriptor, err := a.toProtoMessageDescriptor(genType)

		// store specific message parse failures
		if err != nil {
			a.errors[genType.Name] = err
			continue
		}

		protoPkg, err := protoPackageName(genType)
		if err != nil {
			a.errors[genType.Name] = err
			continue
		}

		if _, ok := protoPackages[protoPkg]; !ok {
			goPkg := a.goPackageName(protoPkg)
			protoPackages[protoPkg] = &descriptorpb.FileDescriptorProto{
				Name:    relFileName(protoPkg),
				Package: &protoPkg,
				Syntax:  strptr("proto3"),
				Options: &descriptorpb.FileOptions{
					GoPackage: &goPkg,
				},
			}
		}
		fd := protoPackages[protoPkg]
		fd.MessageType = append(fd.MessageType, messageDescriptor)
		a.schemaProtoFiles[genType.Name] = *fd.Name

		depPaths, err := a.extractDepPaths(messageDescriptor)
		if err != nil {
			a.errors[genType.Name] = err
			continue
		}
		fd.Dependency = append(fd.Dependency, depPaths...)

		svcAnnotation, err := extractServiceAnnotation(genType)
		if errors.Is(err, errNoServiceDef) {
			continue
		}
		if err != nil {
			return err
		}
		if svcAnnotation.Generate {
			svcResources, err := a.createServiceResources(genType, svcAnnotation.Methods)
			if err != nil {
				return err
			}
			fd.Service = append(fd.Service, svcResources.svc)
			fd.MessageType = append(fd.MessageType, svcResources.svcMessages...)
			fd.Dependency = append(fd.Dependency, "google/protobuf/empty.proto")
		}
	}

	// Append the well known types to the context.
	for _, wktPath := range wktsPaths {
		typeDesc, err := desc.LoadFileDescriptor(wktPath)
		if err != nil {
			return err
		}
		dpbDescriptors = append(dpbDescriptors, typeDesc.AsFileDescriptorProto())
	}

	for _, fd := range protoPackages {
		fd.Dependency = dedupe(fd.Dependency)
		dpbDescriptors = append(dpbDescriptors, fd)
	}

	descriptors, err := desc.CreateFileDescriptors(dpbDescriptors)
	if err != nil {
		return err
	}

	// cleanup the WKT protos from the map
	for _, wp := range wktsPaths {
		delete(descriptors, wp)
	}

	for dp, fd := range descriptors {
		fbuild, err := builder.FromFile(fd)
		if err != nil {
			return err
		}
		fbuild.SetSyntaxComments(builder.Comments{
			LeadingComment: " Code generated by entproto. DO NOT EDIT.",
		})
		fd, err = fbuild.Build()
		if err != nil {
			return err
		}
		descriptors[dp] = fd
	}

	a.descriptors = descriptors

	return nil
}

func (a *Adapter) goPackageName(protoPkgName string) string {
	// TODO(rotemtam): make this configurable from an annotation
	entBase := a.graph.Config.Package
	slashed := strings.ReplaceAll(protoPkgName, ".", "/")
	return path.Join(entBase, "proto", slashed)
}

// GetFileDescriptor returns the proto file descriptor containing the transformed proto message descriptor for
// `schemaName` along with any other messages in the same protobuf package.
func (a *Adapter) GetFileDescriptor(schemaName string) (*desc.FileDescriptor, error) {
	if err, ok := a.errors[schemaName]; ok {
		return nil, err
	}
	fn, ok := a.schemaProtoFiles[schemaName]
	if !ok {
		return nil, fmt.Errorf("entproto: could not find file descriptor for schema %s", schemaName)
	}

	dsc, ok := a.descriptors[fn]
	if !ok {
		return nil, fmt.Errorf("entproto: could not find file descriptor for schema %s", schemaName)
	}

	return dsc, nil
}

func protoPackageName(genType *gen.Type) (string, error) {
	msgAnnot, err := extractMessageAnnotation(genType)
	if err != nil {
		return "", err
	}

	if msgAnnot.Package != "" {
		return msgAnnot.Package, nil
	}
	return DefaultProtoPackageName, nil
}

func relFileName(packageName string) *string {
	parts := strings.Split(packageName, ".")
	fileName := parts[len(parts)-1] + ".proto"
	parts = append(parts, fileName)
	joined := filepath.Join(parts...)
	return &joined
}

func (a *Adapter) extractDepPaths(m *descriptorpb.DescriptorProto) ([]string, error) {
	var out []string
	for _, fld := range m.Field {
		if *fld.Type == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE { //nolint
			fieldTypeName := *fld.TypeName
			if wp, ok := wktsPaths[fieldTypeName]; ok { //nolint
				out = append(out, wp)
			} else if graphContainsDependency(a.graph, fieldTypeName) {
				fieldTypeName = extractLastFqnPart(fieldTypeName)
				depType, err := extractGenTypeByName(a.graph, fieldTypeName)
				if err != nil {
					return nil, err
				}
				depPackageName, err := protoPackageName(depType)
				if err != nil {
					return nil, err
				}
				selfType, err := extractGenTypeByName(a.graph, *m.Name)
				if err != nil {
					return nil, err
				}
				selfPackageName, _ := protoPackageName(selfType)
				if depPackageName != selfPackageName {
					importPath := relFileName(depPackageName)
					out = append(out, *importPath)
				}
			} else {
				return nil, fmt.Errorf("entproto: failed extracting deps, unknown path for %s", fieldTypeName)
			}
		}
	}
	return out, nil
}

func graphContainsDependency(graph *gen.Graph, fieldTypeName string) bool {
	gt, err := extractGenTypeByName(graph, extractLastFqnPart(fieldTypeName))
	if err != nil {
		return false
	}
	return gt != nil
}

func extractLastFqnPart(fqn string) string {
	parts := strings.Split(fqn, ".")
	return parts[len(parts)-1]
}

type unsupportedTypeError struct {
	Type *field.TypeInfo
}

func (e unsupportedTypeError) Error() string {
	return fmt.Sprintf("unsupported field type %q", e.Type.ConstName())
}

func (a *Adapter) toProtoMessageDescriptor(genType *gen.Type) (*descriptorpb.DescriptorProto, error) {
	msgAnnot, err := extractMessageAnnotation(genType)
	if err != nil || !msgAnnot.Generate {
		return nil, ErrSchemaSkipped
	}
	msg := &descriptorpb.DescriptorProto{
		Name:     &genType.Name,
		EnumType: []*descriptorpb.EnumDescriptorProto(nil),
	}

	if !genType.ID.UserDefined {
		genType.ID.Annotations = map[string]interface{}{FieldAnnotation: Field(IDFieldNumber)}
	}

	all := []*gen.Field{genType.ID}
	all = append(all, genType.Fields...)

	for _, f := range all {
		if _, ok := f.Annotations[SkipAnnotation]; ok {
			continue
		}

		protoField, err := toProtoFieldDescriptor(f)
		if err != nil {
			return nil, err
		}
		// If the field is an enum type, we need to create the enum descriptor as well.
		if f.Type.Type == field.TypeEnum {
			dp, err := toProtoEnumDescriptor(f)
			if err != nil {
				return nil, err
			}
			msg.EnumType = append(msg.EnumType, dp)
		}
		msg.Field = append(msg.Field, protoField)
	}

	for _, e := range genType.Edges {
		if _, ok := e.Annotations[SkipAnnotation]; ok {
			continue
		}

		descriptor, err := a.extractEdgeFieldDescriptor(genType, e)
		if err != nil {
			return nil, err
		}
		if descriptor != nil {
			msg.Field = append(msg.Field, descriptor)
		}
	}

	if err := verifyNoDuplicateFieldNumbers(msg); err != nil {
		return nil, err
	}

	return msg, nil
}

func verifyNoDuplicateFieldNumbers(msg *descriptorpb.DescriptorProto) error {
	mem := make(map[int32]struct{})
	for _, fld := range msg.Field {
		if _, seen := mem[fld.GetNumber()]; seen {
			return fmt.Errorf("entproto: field %d already defined on message %q",
				fld.GetNumber(), msg.GetName())
		} else {
			mem[fld.GetNumber()] = struct{}{}
		}
	}
	return nil
}

func (a *Adapter) extractEdgeFieldDescriptor(source *gen.Type, e *gen.Edge) (*descriptorpb.FieldDescriptorProto, error) {
	t := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	msgTypeName := pascal(e.Type.Name)

	edgeAnnotation, err := extractEdgeAnnotation(e)
	if err != nil {
		return nil, fmt.Errorf("entproto: failed extracting proto field number annotation: %w", err)
	}

	if edgeAnnotation.Number == 1 {
		return nil, fmt.Errorf("entproto: edge %q has number 1 which is reserved for id", e.Name)
	}

	if num := int64(edgeAnnotation.Number); num > math.MaxInt32 || num < math.MinInt32 {
		return nil, fmt.Errorf("value %v overflows int32", num)
	}
	fieldNum := int32(edgeAnnotation.Number) //nolint:gosec
	fieldDesc := &descriptorpb.FieldDescriptorProto{
		Number: &fieldNum,
		Name:   &e.Name,
		Type:   &t,
	}

	if !e.Unique {
		fieldDesc.Label = &repeatedFieldLabel
	}

	relType, err := extractGenTypeByName(a.graph, msgTypeName)
	if err != nil {
		return nil, err
	}
	dstAnnotation, err := extractMessageAnnotation(relType)
	if err != nil || !dstAnnotation.Generate {
		return nil, fmt.Errorf("entproto: message %q is not generated", msgTypeName)
	}

	sourceAnnotation, err := extractMessageAnnotation(source)
	if err != nil {
		return nil, err
	}
	if sourceAnnotation.Package == dstAnnotation.Package {
		fieldDesc.TypeName = &msgTypeName
	} else {
		fqn := dstAnnotation.Package + "." + msgTypeName
		fieldDesc.TypeName = &fqn
	}

	return fieldDesc, nil
}

func toProtoEnumDescriptor(fld *gen.Field) (*descriptorpb.EnumDescriptorProto, error) {
	enumAnnotation, err := extractEnumAnnotation(fld)
	if err != nil {
		return nil, err
	}
	if err := enumAnnotation.Verify(fld); err != nil {
		return nil, err
	}
	enumName := pascal(fld.Name)
	dp := &descriptorpb.EnumDescriptorProto{
		Name:  strptr(enumName),
		Value: []*descriptorpb.EnumValueDescriptorProto{},
	}
	if !fld.Default {
		dp.Value = append(dp.Value, &descriptorpb.EnumValueDescriptorProto{
			Number: int32ptr(0),
			Name:   strptr(strings.ToUpper(snake(fld.Name)) + "_UNSPECIFIED"),
		})
	}
	for _, opt := range fld.Enums {
		n := strings.ToUpper(snake(NormalizeEnumIdentifier(opt.Value)))
		if !enumAnnotation.OmitFieldPrefix {
			n = strings.ToUpper(snake(fld.Name)) + "_" + n
		}
		dp.Value = append(dp.Value, &descriptorpb.EnumValueDescriptorProto{
			Number: int32ptr(enumAnnotation.Options[opt.Value]),
			Name:   strptr(n),
		})
	}
	return dp, nil
}

func toProtoFieldDescriptor(f *gen.Field) (*descriptorpb.FieldDescriptorProto, error) {
	fieldDesc := &descriptorpb.FieldDescriptorProto{
		Name: &f.Name,
	}
	fann, err := extractFieldAnnotation(f)
	if err != nil {
		return nil, err
	}
	if num := int64(fann.Number); num > math.MaxInt32 || num < math.MinInt32 {
		return nil, fmt.Errorf("value %v overflows int32", num)
	}
	fieldNumber := int32(fann.Number) //nolint:gosec
	if fieldNumber == 1 && strings.ToUpper(f.Name) != "ID" {
		return nil, fmt.Errorf("entproto: field %q has number 1 which is reserved for id", f.Name)
	}
	fieldDesc.Number = &fieldNumber
	if fann.Type != descriptorpb.FieldDescriptorProto_Type(0) {
		fieldDesc.Type = &fann.Type
		if len(fann.TypeName) > 0 {
			fieldDesc.TypeName = &fann.TypeName
		}
		return fieldDesc, nil
	}
	typeDetails, err := extractProtoTypeDetails(f)
	if err != nil {
		return nil, err
	}
	fieldDesc.Type = &typeDetails.protoType
	if typeDetails.messageName != "" {
		fieldDesc.TypeName = &typeDetails.messageName
	}
	if typeDetails.repeated {
		fieldDesc.Label = &repeatedFieldLabel
	}
	return fieldDesc, nil
}

func extractProtoTypeDetails(f *gen.Field) (fieldType, error) {
	if f.Type.Type == field.TypeJSON {
		return extractJSONDetails(f)
	}
	cfg, ok := typeMap[f.Type.Type]
	if !ok || cfg.unsupported {
		return fieldType{}, unsupportedTypeError{Type: f.Type}
	}
	if f.Optional {
		if cfg.optionalType == "" {
			return fieldType{}, unsupportedTypeError{Type: f.Type}
		}
		return fieldType{
			protoType:   descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
			messageName: cfg.optionalType,
		}, nil
	}
	name := cfg.msgTypeName
	if cfg.namer != nil {
		name = cfg.namer(f)
	}
	return fieldType{
		protoType:   cfg.pbType,
		messageName: name,
	}, nil
}

func extractJSONDetails(f *gen.Field) (fieldType, error) {
	switch f.Type.Ident {
	case "[]string":
		return fieldType{
			protoType: descriptorpb.FieldDescriptorProto_TYPE_STRING,
			repeated:  true,
		}, nil
	case "[]int32":
		return fieldType{
			protoType: descriptorpb.FieldDescriptorProto_TYPE_INT32,
			repeated:  true,
		}, nil
	case "[]int64":
		return fieldType{
			protoType: descriptorpb.FieldDescriptorProto_TYPE_INT64,
			repeated:  true,
		}, nil
	case "[]uint32":
		return fieldType{
			protoType: descriptorpb.FieldDescriptorProto_TYPE_UINT32,
			repeated:  true,
		}, nil
	case "[]uint64":
		return fieldType{
			protoType: descriptorpb.FieldDescriptorProto_TYPE_UINT64,
			repeated:  true,
		}, nil
	}
	return fieldType{}, unsupportedTypeError{Type: f.Type}
}

type fieldType struct {
	messageName string
	protoType   descriptorpb.FieldDescriptorProto_Type
	repeated    bool
}

func strptr(s string) *string {
	return &s
}

func int32ptr(i int32) *int32 {
	return &i
}

func extractGenTypeByName(graph *gen.Graph, name string) (*gen.Type, error) {
	for _, sch := range graph.Nodes {
		if sch.Name == name {
			return sch, nil
		}
	}
	return nil, fmt.Errorf("entproto: could not find schema %q in graph", name)
}

func dedupe(s []string) []string {
	out := make([]string, 0, len(s))
	seen := make(map[string]struct{})
	for _, item := range s {
		if _, skip := seen[item]; skip {
			continue
		}
		out = append(out, item)
		seen[item] = struct{}{}
	}
	return out
}
