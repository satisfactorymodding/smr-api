package validation

import (
	"context"
	"strings"

	"github.com/Vilsol/ue4pak/parser"
	"github.com/rs/zerolog/log"
)

func DecodeProperty(ctx context.Context, cleanName string, property *parser.FPropertyTag) interface{} {
	// Ignore on purpose
	if cleanName == "VertexData" {
		return nil
	}

	switch strings.Trim(property.PropertyType, "\x00") {
	case "ArrayProperty":
		return DecodeArrayProperty(ctx, property)
	case "IntProperty":
		return property.Tag
	case "Int8Property":
		return property.Tag
	case "UInt64Property":
		return property.Tag
	case "FloatProperty":
		return property.Tag
	case "BoolProperty":
		return property.TagData
	case "TextProperty":
		return trim(property.Tag.(*parser.FText).SourceString)
	case "ObjectProperty":
		return FPackageIndexToString(property.Tag)
	case "EnumProperty":
		return trim(property.Tag.(string))
	case "StrProperty":
		return trim(property.Tag.(string))
	case "NameProperty":
		return trim(property.Tag.(string))
	case "StructProperty":
		return DecodeStructProperty(ctx, property)
	case "SoftObjectProperty":
		// TODO Might need second
		return property.Tag.(*parser.FSoftObjectPath).AssetPathName
	case "ByteProperty":
		if str, ok := property.Tag.(string); ok {
			return trim(str)
		} else if b, ok := property.Tag.(byte); ok {
			return b
		} else {
			log.Error().Msgf("Unknown ByteProperty type: %#v", property)
		}
	default:
		log.Error().Msgf("Unknown property type [%s]: %s", cleanName, property.PropertyType)
	}

	return nil
}

func DecodePropertyFields(ctx context.Context, properties []*parser.FPropertyTag) map[string]interface{} {
	result := make(map[string]interface{})

	for _, property := range properties {
		cleanName := strings.Trim(property.Name, "\x00")
		result[cleanName] = DecodeProperty(ctx, cleanName, property)
	}

	return result
}

func DecodeArrayProperty(ctx context.Context, property *parser.FPropertyTag) []interface{} {
	arrayData := property.Tag.([]interface{})

	if len(arrayData) == 0 {
		return make([]interface{}, 0)
	}

	results := make([]interface{}, len(arrayData))

	switch strings.Trim(property.TagData.(string), "\x00") {
	case "StructProperty":
		for i, data := range arrayData {
			properties := data.(*parser.ArrayStructProperty).Properties
			switch properties.(type) {
			case *parser.StructType:
				break
			default:
				results[i] = DecodePropertyFields(ctx, properties.([]*parser.FPropertyTag))
			}
		}
	case "SoftObjectProperty":
		for i, data := range arrayData {
			// TODO Might need second
			results[i] = trim(data.(*parser.FSoftObjectPath).AssetPathName)
		}
	case "ObjectProperty":
		for i, data := range arrayData {
			results[i] = FPackageIndexToStringSpecial(data)
		}
	case "StrProperty":
		fallthrough
	case "EnumProperty":
		fallthrough
	case "NameProperty":
		for i, data := range arrayData {
			results[i] = trim(data.(string))
		}
	case "IntProperty":
		fallthrough
	case "FloatProperty":
		copy(results, arrayData)
	default:
		log.Error().Msgf("Unknown array property data type [%s]: %s", property.Name, property.TagData.(string))
	}

	return results
}

func DecodeStructProperty(ctx context.Context, property *parser.FPropertyTag) interface{} {
	if arr, ok := property.Tag.([]*parser.FPropertyTag); ok {
		return DecodePropertyFields(ctx, arr)
	}

	return property.Tag.(*parser.StructType).Value
}

func FPackageIndexToStringSpecial(index interface{}) string {
	fPackage := index.(*parser.FPackageIndex)
	result := ""

	if fImport, ok := fPackage.Reference.(*parser.FObjectImport); ok {
		result = trim(fImport.ObjectName)
	} else if fExport, ok := fPackage.Reference.(*parser.FObjectExport); ok {
		result = trim(fExport.ObjectName)
	}

	return result
}

func FPackageIndexToString(index interface{}) string {
	fPackage := index.(*parser.FPackageIndex)

	if fImport, ok := fPackage.Reference.(*parser.FObjectImport); ok {
		return trim(fImport.ObjectName)
	} else if fExport, ok := fPackage.Reference.(*parser.FObjectExport); ok {
		return trim(fExport.ObjectName)
	}

	return ""
}
