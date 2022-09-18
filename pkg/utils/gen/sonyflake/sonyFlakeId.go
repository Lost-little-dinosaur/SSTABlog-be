package sonyflake

import (
	"SSTABlog-be/internal/logger"
	"github.com/sony/sonyflake"
)

var flake *sonyflake.Sonyflake

func init() {
	flake = sonyflake.NewSonyflake(sonyflake.Settings{})
}

func GenSonyFlakeId() (int64, error) {

	id, err := flake.NextID()
	if err != nil {
		logger.Warning.Println("flake NextID failed: ", err)
		return 0, err
	}

	return int64(id), nil
}
