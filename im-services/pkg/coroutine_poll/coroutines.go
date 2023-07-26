package coroutine_poll

import (
	"github.com/panjf2000/ants/v2"
	"im-services/internal/config"
)

var AntsPool *ants.Pool

func ConnectPool() *ants.Pool {
	//设置数量
	AntsPool, _ = ants.NewPool(config.Conf.Server.CoroutinePoll)
	return AntsPool
}
