package kafak

import (
	"fmt"
	"im-services/internal/config"
	"im-services/internal/helpers"
	dao2 "im-services/internal/service/dao"
	"sync"

	"github.com/IBM/sarama"
)

var (
	offlineMessageDao *dao2.OfflineMessageDao
)

func ConsumerInit() {

	var wg sync.WaitGroup
	consumer, err := sarama.NewConsumer([]string{fmt.Sprintf("%s:%s", config.Conf.Kafka.Host, config.Conf.Kafka.Port)}, nil)
	if err != nil {
		fmt.Printf("Failed to start consumer: %s", err)
		return
	}
	partitionList, err := consumer.Partitions("web_log") //获得该topic所有的分区
	if err != nil {
		fmt.Println("Failed to get the list of partition:, ", err)
		return
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition("web_log", int32(partition), sarama.OffsetNewest)
		if err != nil {
			fmt.Printf("Failed to start consumer for partition %d: %s\n", partition, err)
			return
		}
		wg.Add(1)
		go func(sarama.PartitionConsumer) { //为每个分区开一个go协程去取值
			for msg := range pc.Messages() { //阻塞直到有值发送过来，然后再继续等待
				offlineMessageDao.PrivateOfflineMessageSave(string(msg.Value))
			}
			defer pc.AsyncClose()
			wg.Done()
		}(pc)
	}
	wg.Wait()
	err = consumer.Close()
	helpers.ErrorHandler(err)

}
