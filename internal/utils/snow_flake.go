package utils

import (
	"fmt"
	"github.com/holdno/snowFlakeByGo"
	"strconv"
)

var worker *snowFlakeByGo.Worker

func init() {
	// 生成一个节点实例
	var (
		id  int64
		err error
	)

	// todo
	// if id, err = strconv.ParseInt(config.Common.MachineId, 0, 10); err != nil {
	if id, err = strconv.ParseInt("10", 0, 10); err != nil {
		panic(err)
	}

	worker, _ = snowFlakeByGo.NewWorker(id) // 传入当前节点id 此id在机器集群中一定要唯一 且从0开始排最多1024个节点，可以根据节点的不同动态调整该算法每毫秒生成的id上限(如何调整会在后面讲到)
}

func GenerateId() string {
	return fmt.Sprintf("%v", worker.GetId())
}
