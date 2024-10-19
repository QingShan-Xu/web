package ds

import (
	"fmt"
	"reflect"
	"strings"
)

type (
	StructReader struct {
		Name   string
		Type   reflect.Type
		Kind   reflect.Kind
		Value  reflect.Value
		Fields []*StructReader

		IsEmbedded bool
		IsList     bool

		Field reflect.StructField
		Tags  map[string]FieldTag
	}

	FieldTag struct {
		Key     string
		Value   string
		Options []string
	}
)

func NewStructReader(data interface{}) (*StructReader, error) {
	value := reflect.Indirect(reflect.ValueOf(data))
	kind := value.Kind()
	if kind != reflect.Struct {
		return nil, fmt.Errorf("NewStructReader: parameter must be a struct type")
	}

	reader, err := generateStructReader(value)
	if err != nil {
		return nil, err
	}

	return reader, nil
}

func generateStructReader(value reflect.Value) (*StructReader, error) {
	value = reflect.Indirect(value)
	structType := value.Type()

	reader := &StructReader{
		Name:   structType.Name(),
		Type:   structType,
		Kind:   value.Kind(),
		Value:  value,
		Fields: []*StructReader{},
	}

	if reader.Kind == reflect.Struct {
		for i := 0; i < value.NumField(); i++ {
			fieldValue := value.Field(i)
			field := structType.Field(i)

			fieldReader := &StructReader{
				Name:   field.Name,
				Type:   fieldValue.Type(),
				Kind:   fieldValue.Kind(),
				Value:  fieldValue,
				Fields: []*StructReader{},
				Field:  field,
				Tags:   map[string]FieldTag{},
			}
			fieldReader.parseFieldTags()

			if fieldReader.Kind == reflect.Struct {
				fieldReader.IsEmbedded = field.Anonymous
				currentReader, err := generateStructReader(fieldValue)
				if err != nil {
					return nil, err
				}
				fieldReader.Fields = currentReader.Fields
			}

			if fieldReader.Kind == reflect.Slice || fieldReader.Kind == reflect.Array {
				fieldReader.IsList = true
			}

			reader.Fields = append(reader.Fields, fieldReader)
		}
	}

	return reader, nil
}

func (r *StructReader) parseFieldName(name string) []string {
	return strings.Split(name, ".")
}

func (r *StructReader) GetValue() interface{} {
	return r.Value.Interface()
}

func (r *StructReader) GetFieldByName(name string) (*StructReader, error) {
	if name == "" {
		return nil, fmt.Errorf("GetFieldByName: name cannot be empty")
	}

	nameComponents := r.parseFieldName(name)
	currentReader := r

	for _, component := range nameComponents {
		found := false
		currentReader, found = r.searchField(currentReader, component)
		if !found {
			return nil, fmt.Errorf("field %s not found", component)
		}
	}

	return currentReader, nil
}

func (r *StructReader) searchField(reader *StructReader, fieldName string) (*StructReader, bool) {
	for _, field := range reader.Fields {
		if field.Name == fieldName {
			return field, true
		}

		if field.IsEmbedded {
			if embeddedField, found := r.searchField(field, fieldName); found {
				return embeddedField, true
			}
		}
	}

	return nil, false
}

func (r *StructReader) parseFieldTags() {
	if r.Field.Name == "" {
		return
	}

	tagStrings := strings.Fields((string(r.Field.Tag)))
	for _, tagString := range tagStrings {
		tagParts := strings.Split(tagString, ":")
		currentTag := FieldTag{
			Key: tagParts[0],
		}
		tagContent := r.Field.Tag.Get(currentTag.Key)
		tagComponents := strings.Split(tagContent, ",")
		if len(tagComponents) > 0 {
			currentTag.Value = tagComponents[0]
		}
		if len(tagComponents) > 1 {
			currentTag.Options = tagComponents[1:]
		}
		r.Tags[currentTag.Key] = currentTag
	}
}
