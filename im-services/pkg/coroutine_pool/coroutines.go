package coroutine_pool

import (
	"im-services/internal/config"

	"github.com/panjf2000/ants/v2"
)

var AntsPool *ants.Pool

func ConnectPool() *ants.Pool {
	//设置数量
	AntsPool, _ = ants.NewPool(config.Conf.Server.CoroutinePool)
	return AntsPool
}
