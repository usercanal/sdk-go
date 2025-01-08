// convert/convert.go
package convert

import (
	"fmt"
	"unicode/utf8"

	pb "github.com/usercanal/sdk-go/proto"
	"github.com/usercanal/sdk-go/types"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// EventToProto converts a types.Event to a protobuf Event
func EventToProto(e *types.Event) (*pb.Event, error) {
	var props map[string]*pb.Value
	if e.Properties != nil {
		var err error
		props, err = ToProperties(e.Properties)
		if err != nil {
			// Initialize empty map instead of nil
			props = make(map[string]*pb.Value)
		}
	} else {
		// Initialize empty map instead of nil
		props = make(map[string]*pb.Value)
	}

	base := &pb.MessageBase{
		DistinctId: e.UserId,
		UserId:     e.UserId,
		Timestamp:  timestamppb.New(e.Timestamp),
	}

	track := &pb.TrackEvent{
		Base:       base,
		EventName:  e.Name.String(), // Convert EventName to string
		Properties: props,           // Always pass initialized map
	}

	return &pb.Event{
		Type: pb.Event_TRACK,
		Event: &pb.Event_Track{
			Track: track,
		},
	}, nil
}

// IdentityToProto converts a types.Identity to a protobuf Event
func IdentityToProto(i *types.Identity) (*pb.Event, error) {
	props, _ := ToProperties(i.Properties)
	return &pb.Event{
		Type: pb.Event_IDENTIFY,
		Event: &pb.Event_Identify{
			Identify: &pb.IdentifyEvent{
				Base: &pb.MessageBase{
					UserId: i.UserId,
				},
				Traits:  props,
				Context: nil,
			},
		},
	}, nil
}

// GroupToProto converts a types.GroupInfo to a protobuf Event
func GroupToProto(g *types.GroupInfo) (*pb.Event, error) {
	props, _ := ToProperties(g.Properties)
	return &pb.Event{
		Type: pb.Event_GROUP,
		Event: &pb.Event_Group{
			Group: &pb.GroupEvent{
				Base: &pb.MessageBase{
					UserId: g.UserId,
				},
				GroupId:   g.GroupId,
				GroupType: "organization",
				Traits:    props,
			},
		},
	}, nil
}

// ToValue converts a Go interface{} to a protobuf Value
func ToValue(v interface{}) (*pb.Value, error) {
	if v == nil {
		return nil, nil
	}

	switch val := v.(type) {
	case string:
		// Try to encode as UTF-8 string first
		if utf8.ValidString(val) {
			return &pb.Value{Value: &pb.Value_StringValue{StringValue: val}}, nil
		}
		// Fall back to bytes if not valid UTF-8
		return &pb.Value{Value: &pb.Value_BytesValue{BytesValue: []byte(val)}}, nil
	case []byte:
		// Raw bytes should be stored as bytes value
		return &pb.Value{Value: &pb.Value_BytesValue{BytesValue: val}}, nil
	case int:
		return &pb.Value{Value: &pb.Value_IntValue{IntValue: int64(val)}}, nil
	case int32:
		return &pb.Value{Value: &pb.Value_IntValue{IntValue: int64(val)}}, nil
	case int64:
		return &pb.Value{Value: &pb.Value_IntValue{IntValue: val}}, nil
	case float32:
		return &pb.Value{Value: &pb.Value_DoubleValue{DoubleValue: float64(val)}}, nil
	case float64:
		return &pb.Value{Value: &pb.Value_DoubleValue{DoubleValue: val}}, nil
	case bool:
		return &pb.Value{Value: &pb.Value_BoolValue{BoolValue: val}}, nil
	case types.EventName: // Add support for EventName type
		return &pb.Value{Value: &pb.Value_StringValue{StringValue: val.String()}}, nil
	case types.Currency: // Add support for Currency type
		return &pb.Value{Value: &pb.Value_StringValue{StringValue: string(val)}}, nil
	case types.RevenueType: // Add support for RevenueType type
		return &pb.Value{Value: &pb.Value_StringValue{StringValue: string(val)}}, nil
	case types.AuthMethod: // Add support for AuthMethod type
		return &pb.Value{Value: &pb.Value_StringValue{StringValue: string(val)}}, nil
	case types.PaymentMethod: // Add support for PaymentMethod type
		return &pb.Value{Value: &pb.Value_StringValue{StringValue: string(val)}}, nil
	case []interface{}:
		return toArrayValue(val)
	case map[string]interface{}:
		return toObjectValue(val)
	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}

// ToProperties converts a map of Go values to protobuf Values
func ToProperties(props map[string]interface{}) (map[string]*pb.Value, error) {
	// Always initialize the map, even if props is nil
	result := make(map[string]*pb.Value)

	// If props is nil, return empty initialized map
	if props == nil {
		return result, nil
	}

	for k, v := range props {
		val, err := ToValue(v)
		if err != nil {
			return nil, fmt.Errorf("property %q: %w", k, err)
		}
		if val != nil {
			result[k] = val
		}
	}
	return result, nil
}

func toArrayValue(arr []interface{}) (*pb.Value, error) {
	values := make([]*pb.Value, 0, len(arr))
	for i, v := range arr {
		val, err := ToValue(v)
		if err != nil {
			return nil, fmt.Errorf("array index %d: %w", i, err)
		}
		if val != nil {
			values = append(values, val)
		}
	}
	return &pb.Value{
		Value: &pb.Value_ArrayValue{
			ArrayValue: &pb.Array{Values: values},
		},
	}, nil
}

func toObjectValue(obj map[string]interface{}) (*pb.Value, error) {
	fields, err := ToProperties(obj)
	if err != nil {
		return nil, err
	}
	return &pb.Value{
		Value: &pb.Value_ObjectValue{
			ObjectValue: &pb.Object{Fields: fields},
		},
	}, nil
}

// FromValue converts a protobuf Value to a Go interface{}
func FromValue(v *pb.Value) interface{} {
	if v == nil {
		return nil
	}

	switch val := v.Value.(type) {
	case *pb.Value_StringValue:
		return val.StringValue
	case *pb.Value_BytesValue:
		// Try to convert bytes back to string if it's valid UTF-8
		if utf8.Valid(val.BytesValue) {
			return string(val.BytesValue)
		}
		// Return raw bytes if not valid UTF-8
		return val.BytesValue
	case *pb.Value_IntValue:
		return val.IntValue
	case *pb.Value_DoubleValue:
		return val.DoubleValue
	case *pb.Value_BoolValue:
		return val.BoolValue
	case *pb.Value_ArrayValue:
		return fromArrayValue(val.ArrayValue)
	case *pb.Value_ObjectValue:
		return fromObjectValue(val.ObjectValue)
	default:
		return nil
	}
}

func fromArrayValue(arr *pb.Array) []interface{} {
	if arr == nil {
		return nil
	}
	result := make([]interface{}, len(arr.Values))
	for i, v := range arr.Values {
		result[i] = FromValue(v)
	}
	return result
}

func fromObjectValue(obj *pb.Object) map[string]interface{} {
	if obj == nil {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range obj.Fields {
		result[k] = FromValue(v)
	}
	return result
}
