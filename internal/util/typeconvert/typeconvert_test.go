package typeconvert

import "testing"

func TestSliceStringToInt64(t *testing.T) {
	ss := []string{"a1234", "3434"}
	vv := SliceStringToInt64(ss)
	t.Log(vv)
}
