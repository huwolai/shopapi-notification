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
			dv.Ack(false)
			return
		}

	})
	log.Info("Awaiting Message...");

	<-forever

	log.Info("程序退出...");

}

//发送商户收到订单短信
func SendMOrderNotify(event  *queue.OrderEvent)  {
	tmpId :=setting.GetYunTongXunSetting()["morder_template_id"]
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

	//用户手机号
	var mobile string
	if extData["mobile"]!=nil {
		mobile = extData["mobile"].(string)
	}
	//是否是私人订制
	var isTailor bool
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
			if resultMap!=nil{
				if resultMap["dinner_time"]!=nil{
					dinnerTime = resultMap["dinner_time"].(string)
				}
				//私人订制
				if resultMap["goods_type"]!=nil&&resultMap["goods_type"].(string)=="tailor" {
					isTailor = true
				}
			}
		}
	}

	if isTailor {
		log.Info("私人订制订单,不发送厨师短信！")
		return
	}
	if  dinnerTime=="" {
		log.Info("不是厨师订单,不发送短信")
		return
	}
	if mmobile!="" {
		err :=service.SendSMSOfYunTongXun(mmobile,tmpId,[]string{merchantName,name,address,dinnerTime,event.Content.Title,mobile})
		if err!=nil{
			log.Error("商户订单短信发送失败",err)
			return
		}
		log.Info("商户订单短信发送成功！")
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
	dinnerTime :=""
	items :=event.Content.Items
	//是否是私人订制
	var isTailor bool
	var title string
	if items!=nil&&len(items)>0{
		item :=items[0]
		if item.Json!="" {
			var resultMap map[string]interface{}
			err :=util.ReadJsonByByte([]byte(item.Json),&resultMap)
			if err!=nil{
				log.Error(err)
			}
			if resultMap!=nil{
				if resultMap["dinner_time"]!=nil{
					dinnerTime = resultMap["dinner_time"].(string)
				}

				//私人订制
				if resultMap["goods_type"]!=nil&&resultMap["goods_type"].(string)=="tailor" {
					isTailor = true
				}
			}

		}
		title = item.Title
	}

	if  dinnerTime==""&&!isTailor {
		log.Info("不是厨师订单,不发送短信")
		return
	}
	var smsData []string
	if isTailor {
		tmpId = setting.GetYunTongXunSetting()["tailor_template_id"]
		smsData = []string{title,dinnerTime,event.Content.OrderNo}
	}else{
		smsData = []string{merchantName,dinnerTime,event.Content.OrderNo}
	}
	if mobile!="" {

		err :=service.SendSMSOfYunTongXun(mobile,tmpId,smsData)
		if err!=nil{
			log.Error("用户订单短信发送失败",err)
			return
		}
		log.Info("用户订单短信发送成功！")

	}
}
