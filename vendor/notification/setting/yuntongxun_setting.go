package setting

import (
	"gitlab.qiyunxin.com/tangtao/utils/log"
	"fmt"
	"io/ioutil"
	"gitlab.qiyunxin.com/tangtao/utils/util"
	"encoding/json"
	"bytes"
)


var configMap map[string]string


// 获取容联短信配置
func GetYunTongXunSetting() map[string]string {

	if configMap == nil {
		err :=LoadSettingsByFile("config/sms_yuntongxun.json",&configMap)
		log.Error(err)
	}
	return configMap
}

func LoadSettingsByFile(file string,resultMap *map[string]string)  (error) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println("Error while reading config file", err)

		util.CheckErr(err)
	}

	mdz:=json.NewDecoder(bytes.NewBuffer([]byte(content)))

	mdz.UseNumber()
	jsonErr := mdz.Decode(resultMap)

	return jsonErr
}