package typeconvert

import (
	"github.com/shopspring/decimal"
	"github.com/spf13/cast"
	"strconv"
)

func ConvToFloat32(d *decimal.Decimal) float32 {
	if d == nil {
		return 0
	}
	f, _ := d.Float64()
	return float32(f)
}

func StrToInt32(str *string) int32 {
	i, _ := strconv.Atoi(*str)
	return int32(i)
}

func StrToDecimal(str *string) decimal.Decimal {
	d, _ := decimal.NewFromString(*str)
	return d
}

func Int32ToStr(i int32) string {
	return strconv.Itoa(int(i))
}

func ConvToStringFixed(d *decimal.Decimal, places int32) string {
	if d == nil {
		return ""
	}
	return d.StringFixed(places)
}

func ConvToStringFixedAndTrimTrailingZeros(d *decimal.Decimal, places int32) string {
	if d == nil {
		return ""
	}
	return d.Round(places).String()
}

// StringToUint64 : string -> uint64
func StringToUint64(s string) uint64 {
	if s == "" {
		return 0
	}

	TmpInt, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return uint64(TmpInt)
}

func SliceInt64ToString(arr []int64) []string {
	strs := make([]string, 0)
	for _, v := range arr {
		strs = append(strs, cast.ToString(v))
	}
	return strs
}

func SliceStringToInt64(arr []string) []int64 {
	values := make([]int64, 0)
	for _, s := range arr {
		v, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			continue
		}
		values = append(values, v)
	}
	return values
}
