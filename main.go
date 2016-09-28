package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/queue"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"github.com/streadway/amqp"
	"notification/setting"
	"notification/service"
)

func main() {
	//初始化消息队列
	queue.SetupAMQP(config.GetValue("amqp_url").ToString())

	//订阅订单事件
	queue.ConsumeOrderEvent(func(event *queue.OrderEvent,dv amqp.Delivery) {
		//订单已付款事件
		if event.EventKey==queue.ORDER_EVENT_PAID {
			tmpId :=setting.GetYunTongXunSetting()["morder_template_id"]

			service.SendSMSOfYunTongXun()
		}
	})

}

//发送短信
func SendMorderNotify()  {

}
