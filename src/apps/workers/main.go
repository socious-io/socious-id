package workers

import (
	"github.com/socious-io/gomail"
	"github.com/socious-io/gomq"
)

var consumers = []gomq.AddConsumerParams{
	{
		Channel:       gomail.GetConfig().WorkerChannel,
		Consumer:      gomail.EmailWorker,
		IsCategorized: true,
	},
}

func RegisterConsumers() {
	for _, consumer := range consumers {
		gomq.AddConsumer(consumer)
	}
}
