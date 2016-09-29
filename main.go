package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/queue"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"github.com/streadway/amqp"
	"notification/setting"
	"notification/service"
	"gitlab.qiyunxin.com/tangtao/utils/log"
)

func main() {
	//初始化消息队列
	log.Info(config.GetValue("amqp_url").ToString())
	queue.SetupAMQP(config.GetValue("amqp_url").ToString())

	//订阅订单事件
	queue.ConsumeOrderEvent(func(event *queue.OrderEvent,dv amqp.Delivery) {
		//订单已付款事件
		if event.EventKey==queue.ORDER_EVENT_PAID {
			SendMOrderNotify(event)
			return
		}

		dv.Ack(false)
	})

}

//发送商户收到订单短信
func SendMOrderNotify(event  *queue.OrderEvent)  {
	tmpId :=setting.GetYunTongXunSetting()["morder_template_id"]
	//商户手机号
	mmobile := event.Content.ExtData["m_mobile"]
	if mmobile!=nil {
		service.SendSMSOfYunTongXun(mmobile.(string),tmpId,[]string{event.Content.CreateTime})
	}

}
