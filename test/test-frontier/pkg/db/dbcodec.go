package db

import (
	"fmt"
	"github.com/golang/protobuf/ptypes"
	configV1 "github.com/nlnwa/veidemann-api-go/config/v1"
	frontierV1 "github.com/nlnwa/veidemann-api-go/frontier/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gopkg.in/rethinkdb/rethinkdb-go.v6/encoding"
	"reflect"
	"time"
)

func init() {
	encoding.SetTypeEncoding(
		reflect.TypeOf(&configV1.ConfigObject{}),
		func(value interface{}) (interface{}, error) {
			v := make(map[string]interface{})
			ProtoToRethink(value.(proto.Message), v)
			return v, nil
		},
		func(encoded interface{}, value reflect.Value) error {
			var o configV1.ConfigObject
			RethinkToProto(encoded, &o)
			value.Set(reflect.ValueOf(o))
			//fmt.Printf("CONV: %v\n%v\n", encoded, protojson.Format(&o))
			return nil
		},
	)
	encoding.SetTypeEncoding(
		reflect.TypeOf(&frontierV1.CrawlLog{}),
		func(value interface{}) (interface{}, error) {
			v := make(map[string]interface{})
			ProtoToRethink(value.(proto.Message), v)
			return v, nil
		},
		func(encoded interface{}, value reflect.Value) error {
			var o frontierV1.CrawlLog
			RethinkToProto(encoded, &o)
			value.Set(reflect.ValueOf(o))
			return nil
		},
	)
	encoding.SetTypeEncoding(
		reflect.TypeOf(&frontierV1.PageLog{}),
		func(value interface{}) (interface{}, error) {
			v := make(map[string]interface{})
			ProtoToRethink(value.(proto.Message), v)
			return v, nil
		},
		func(encoded interface{}, value reflect.Value) error {
			var o frontierV1.PageLog
			RethinkToProto(encoded, &o)
			value.Set(reflect.ValueOf(o))
			return nil
		},
	)
	encoding.SetTypeEncoding(
		reflect.TypeOf(&frontierV1.JobExecutionStatus{}),
		func(value interface{}) (interface{}, error) {
			v := make(map[string]interface{})
			ProtoToRethink(value.(proto.Message), v)
			return v, nil
		},
		func(encoded interface{}, value reflect.Value) error {
			var o frontierV1.JobExecutionStatus
			RethinkToProto(encoded, &o)
			value.Set(reflect.ValueOf(o))
			fmt.Printf("CONV: %v\n", protojson.Format(&o))
			return nil
		},
	)
}

func ProtoToRethink(src proto.Message, dest map[string]interface{}) {
	srcMessage := src.ProtoReflect()
	srcMessage.Range(func(descriptor protoreflect.FieldDescriptor, value protoreflect.Value) bool {
		key := descriptor.JSONName()
		var val interface{}
		if descriptor.IsList() {
			var l []interface{}
			list := value.List()
			for i := 0; i < list.Len(); i++ {
				l = append(l, convertValueToRethink(descriptor, list.Get(i)))
			}
			dest[key] = l
		} else if descriptor.IsMap() {
			fmt.Printf("CONVERT %v %v\n", descriptor.Name(), value)
			fmt.Printf("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!\n")
		} else {
			val = convertValueToRethink(descriptor, value)
			dest[key] = val
		}
		return true
	})
}

func RethinkToProto(src interface{}, dest proto.Message) {
	destMessage := dest.ProtoReflect()

	switch v := src.(type) {
	case map[string]interface{}:
		for f, fv := range v {
			d := destMessage.Descriptor().Fields().ByJSONName(f)
			if d.IsList() {
				m := destMessage.NewField(d)
				l := m.List()
				for _, lv := range fv.([]interface{}) {
					value := convertValueToProto(m, d, lv)
					l.Append(value)
				}
				destMessage.Set(d, m)
				continue
			}
			if d.IsMap() {
				m := destMessage.NewField(d)
				for _, lv := range fv.([]interface{}) {
					kv := lv.(map[string]interface{})
					for k, v := range kv {
						key := protoreflect.ValueOf(k).MapKey()
						value := convertValueToProto(m, d.MapValue(), v)
						m.Map().Set(key, value)
					}
				}
				destMessage.Set(d, m)
				continue
			}
			destMessage.Set(d, convertValueToProto(destMessage.NewField(d), d, fv))
		}
	case time.Time:
		if t, err := ptypes.TimestampProto(v); err == nil {
			ts := dest.(*timestamppb.Timestamp)
			*ts = *t
		}
	}
}

func convertValueToRethink(field protoreflect.FieldDescriptor, value protoreflect.Value) interface{} {
	switch field.Kind() {
	case protoreflect.StringKind:
		return value.String()
	case protoreflect.Int32Kind:
		return value.Int()
	case protoreflect.Int64Kind:
		return value.Int()
	case protoreflect.EnumKind:
		return field.Enum().Values().ByNumber(value.Enum()).Name()
	case protoreflect.MessageKind:
		v := make(map[string]interface{})
		ProtoToRethink(value.Message().Interface(), v)
		return v
	default:
		panic("Unknown Kind: " + field.Kind().String())
	}
}

func convertValueToProto(owner protoreflect.Value, field protoreflect.FieldDescriptor, value interface{}) protoreflect.Value {
	switch field.Kind() {
	case protoreflect.StringKind:
		return protoreflect.ValueOfString(value.(string))
	case protoreflect.Int32Kind:
		return protoreflect.ValueOfInt32(int32(value.(float64)))
	case protoreflect.Int64Kind:
		return protoreflect.ValueOfInt64(int64(value.(float64)))
	case protoreflect.EnumKind:
		return protoreflect.ValueOfEnum(field.Enum().Values().ByName(protoreflect.Name(value.(string))).Number())
	case protoreflect.MessageKind:
		if field.IsList() {
			v := owner.List().NewElement().Message().Interface()
			RethinkToProto(value, v)
			return protoreflect.ValueOfMessage(v.ProtoReflect())
		} else {
			v := owner.Message().Interface()
			RethinkToProto(value, v)
			return protoreflect.ValueOfMessage(v.ProtoReflect())
		}
	default:
		panic("Unknown Kind: " + field.Kind().String())
	}
}
