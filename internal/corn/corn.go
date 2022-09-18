package corn

import (
	"SSTABlog-be/internal/logger"
	"github.com/robfig/cron"
)

func init() {
	c := cron.New()
	err := c.AddFunc("0 0/10 * * * *", func() {})
	if err != nil {
		logger.Error.Fatalln(err)
	}
	c.Start()
	logger.Info.Println("corn init SUCCESS ")
}
