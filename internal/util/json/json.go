package cjson

import (
	"github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"github.com/kissjingalex/virtpay/internal/util/typeconvert"
	"strings"
	"unsafe"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func RegisterFuzzyDecoders() {
	extra.RegisterFuzzyDecoders()
}

func RegisterInt64ArrayEncoderAsStringArray() {
	jsoniter.RegisterTypeEncoderFunc("[]int64", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
		array := *((*[]int64)(ptr))
		strArr := typeconvert.SliceInt64ToString(array)

		stream.WriteArrayStart()
		size := len(strArr)
		var sb strings.Builder
		for i := 0; i < size-1; i++ {
			sb.WriteString("\"" + strArr[i] + "\",")
		}
		sb.WriteString("\"" + strArr[size-1] + "\"")
		stream.WriteRaw(sb.String())
		stream.WriteArrayEnd()
	}, func(pointer unsafe.Pointer) bool {
		return len(*((*[]int64)(pointer))) == 0
	})
}

//func RegisterInt64EncoderAsString() {
//	jsoniter.RegisterTypeEncoderFunc("ctype.Int64s", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
//		number := *((*ctype.Int64s)(ptr))
//		stream.WriteString(strconv.FormatInt(int64(number), 10))
//	}, func(pointer unsafe.Pointer) bool {
//		return *((*int64)(pointer)) == 0
//	})
//}

//
//func RegisterPbTimeEncoderAsTimestamp() {
//	jsoniter.RegisterTypeEncoderFunc("timestamppb.Timestamp", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
//		ts := (*timestamppb.Timestamp)(ptr)
//		if ts.AsTime().UnixNano() < 0 {
//			stream.WriteInt64(0)
//			return
//		}
//		stream.WriteInt64(ts.AsTime().UnixNano() / time.Millisecond.Nanoseconds())
//	}, func(pointer unsafe.Pointer) bool {
//		ts := (*timestamppb.Timestamp)(pointer)
//		return ts.IsValid()
//	})
//}
//
//func RegisterTimeEncoderAsTimestamp() {
//	jsoniter.RegisterTypeEncoderFunc("time.Time", func(ptr unsafe.Pointer, stream *jsoniter.Stream) {
//		ts := (*time.Time)(ptr)
//		if ts.UnixNano() < 0 {
//			stream.WriteInt64(0)
//			return
//		}
//		stream.WriteInt64(ts.UnixNano() / time.Millisecond.Nanoseconds())
//
//	}, func(pointer unsafe.Pointer) bool {
//		ts := (*time.Time)(pointer)
//		return ts.UnixNano() == 0
//	})
//}
