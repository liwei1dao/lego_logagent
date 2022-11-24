package kafka

import (
	"lego_logagent/modules/logagent/core"
	"lego_logagent/modules/logagent/sender"

	"github.com/Shopify/sarama"
	"github.com/liwei1dao/lego/sys/kafka"
)

type Sender struct {
	sender.Sender
	options IOptions
}

func (this *Sender) Init(rer core.IRunner, ser core.ISender, options core.ISenderOptions) (err error) {
	this.Sender.Init(rer, ser, options)
	this.options = options.(IOptions)
	return
}

func (this *Sender) Run(pipeId int) {
	var (
		err     error
		kafka   kafka.ISys
		data    core.ICollData
		message string
	)
	defer func() {
		kafka.Close()
		this.Wg.Done()
	}()
	if kafka, err = this.createkafkaclient(); err != nil {
		this.Runner.Log().Errorf("Run kafka pipeId:%d err:%v", pipeId, err)
		return
	} else {
		for v := range this.Cache.Out() {
			data = v.(core.ICollData)
			if message, _ = data.ToString(); err == nil {
				msg := &sarama.ProducerMessage{
					Topic: this.options.GetKafka_topic(),
					Value: sarama.StringEncoder(message),
				}
				kafka.Asyncproducer_Input() <- msg
			} else {
				this.Runner.Log().Errorf("kafka sender err:%v", err)
			}
		}
	}
}

func (this *Sender) createkafkaclient() (sys kafka.ISys, err error) {
	sys, err = kafka.NewSys(
		kafka.SetStartType(kafka.Asyncproducer),
		kafka.SetHosts(this.options.GetKafka_host()),
		kafka.SetClientID(this.options.GetKafka_client_id()),
		kafka.SetNet_DialTimeout(this.options.GetNet_DialTimeout()),
		kafka.SetNet_KeepAlive(this.options.GetNet_KeepAlive()),
		kafka.SetProducer_MaxMessageBytes(this.options.GetMax_message_bytes()),
		kafka.SetProducer_Compression(this.options.GetProducer_Compression()),
		kafka.SetProducer_Return_Errors(true),
		kafka.SetProducer_CompressionLevel(this.options.GetProducer_CompressionLevel()),
		kafka.SetProducer_Retry_Max(this.options.GetKafka_retry_max()),
	)
	return
}

func (this *Sender) kafka_msgerrhandle(kafka kafka.ISys) {
	go func() {
		for v := range kafka.Asyncproducer_Errors() {
			kafka.Asyncproducer_Input() <- &sarama.ProducerMessage{
				Topic: this.options.GetKafka_topic(),
				Value: v.Msg.Value,
			}
		}
	}()
}
