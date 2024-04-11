package file

import (
	"github.com/cavaliergopher/grab/v3"
)

func Download(srcUrl, dst string) (size int64, err error) {
	rsp, err := grab.Get(dst, srcUrl)
	if err != nil {
		return 0, err
	}
	return rsp.Size(), nil
}
