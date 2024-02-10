package opentelemetry

import (
	"fmt"

	"go.opentelemetry.io/otel/attribute"
)

func parseAttributeKeyValue(key string, value any) (attribute.KeyValue, error) {
	attr := attribute.KeyValue{
		Key: attribute.Key(key),
	}

	switch x := value.(type) {
	case string:
		attr.Value = attribute.StringValue(x)
	case []string:
		attr.Value = attribute.StringSliceValue(x)
	case float64:
		attr.Value = attribute.Float64Value(x)
	case float32:
		attr.Value = attribute.Float64Value(float64(x))
	case int64:
		attr.Value = attribute.Int64Value(x)
	case int32:
		attr.Value = attribute.Int64Value(int64(x))
	case int16:
		attr.Value = attribute.Int64Value(int64(x))
	case int8:
		attr.Value = attribute.Int64Value(int64(x))
	case int:
		attr.Value = attribute.Int64Value(int64(x))
	case bool:
		attr.Value = attribute.BoolValue(x)
	default:
		return attribute.KeyValue{}, fmt.Errorf("type is not supported: %+v", value)
	}

	return attr, nil
}
