package minecraft

import (
	"io"

	"github.com/hpcloud/tail"
)

func tailLocalFile(path string) (*tail.Tail, error) {
	return tail.TailFile(path, tail.Config{
		Follow: true,
		Location: &tail.SeekInfo{
			Whence: io.SeekEnd,
			Offset: 0,
		},
		MustExist: true,
		ReOpen:    true,
	})

}
