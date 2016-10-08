package main

import (
	"gitlab.qiyunxin.com/tangtao/utils/queue"
	"gitlab.qiyunxin.com/tangtao/utils/config"
	"github.com/streadway/amqp"
	"notification/setting"
	"notification/service"
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"gitlab.qiyunxin.com/tangtao/utils/util"
)

func main() {
	forever := make(chan bool)
	err :=config.Init(false)
	if err!=nil{
		log.Error(err)
		return
	}
	//初始化消息队列
	queue.SetupAMQP(config.GetValue("amqp_url").ToString())

	//订阅订单事件
	queue.ConsumeOrderEvent(func(event *queue.OrderEvent,dv amqp.Delivery) {
		//订单已付款事件
		if event.EventKey==queue.ORDER_EVENT_PAID {
			//发送商户订单短信
			go SendMOrderNotify(event)
			//发送用户订单短信
			go SendUOrderNotify(event)
			return
		}

		dv.Ack(false)
	})
	log.Info("Awaiting Message...");

	<-forever

	log.Info("程序退出...");

}

//发送商户收到订单短信
func SendMOrderNotify(event  *queue.OrderEvent)  {
	tmpId :=setting.GetYunTongXunSetting()["morder_template_id"]
	//服务电话
	serviceMobile := setting.GetYunTongXunSetting()["service_mobile"]
	//商户手机号
	extData :=event.Content.ExtData
	var mmobile string
	if extData["m_mobile"]!=nil{
		mmobile = extData["m_mobile"].(string)

	}
	var merchantName string
	if extData["m_name"]!=nil{
		merchantName = extData["m_name"].(string)

	}
	var name string
	if extData["name"]!=nil{
		name =extData["name"].(string)
	}
	var address string
	if extData["address"]!=nil{
		address =extData["address"].(string)
	}
	dinnerTime :=""
	items :=event.Content.Items
	if items!=nil&&len(items)>0{
		item :=items[0]
		if item.Json!="" {
			var resultMap map[string]interface{}
			err :=util.ReadJsonByByte([]byte(item.Json),&resultMap)
			if err!=nil{
				log.Error(err)
			}
			if resultMap["dinner_time"]!=nil{
				dinnerTime = resultMap["dinner_time"].(string)
			}
		}
	}
	if mmobile!="" {
		err :=service.SendSMSOfYunTongXun(mmobile,tmpId,[]string{merchantName,name,address,dinnerTime,event.Content.Title,serviceMobile})
		if err!=nil{
			log.Error("商户订单短信发送失败",err)
		}
	}

}

//发送用户订单短信
func SendUOrderNotify(event *queue.OrderEvent)  {
	tmpId :=setting.GetYunTongXunSetting()["uorder_template_id"]
	extData :=event.Content.ExtData
	//用户手机号
	var mobile string
	if extData["mobile"]!=nil {
		mobile = extData["mobile"].(string)
	}

	//厨师手机号
	//mmobile := extData["m_mobile"]
	//厨师名称
	var merchantName string
	if extData["m_name"]!=nil{
		merchantName = extData["m_name"].(string)

	}
	if mobile!="" {
		err :=service.SendSMSOfYunTongXun(mobile,tmpId,[]string{merchantName,event.Content.CreateTime,event.Content.OrderNo})
		if err!=nil{
			log.Error("用户订单短信发送失败",err)
		}

	}
}
