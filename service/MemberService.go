package service

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
	"github.com/degary/CloudRestaurant/dao"
	"github.com/degary/CloudRestaurant/model"
	"github.com/degary/CloudRestaurant/tool"
	"math/rand"
	"time"
)

type MemberService struct {
}

func (ms *MemberService) SendCode(phone string) bool {
	//1.产生验证码
	code := fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))

	//2.调用阿里SDK,完成发送
	config := tool.GetConfig().Sms
	client, err := dysmsapi.NewClientWithAccessKey(config.RegionId, config.AppKey, config.AppSecret)
	if err != nil {
		fmt.Errorf("调用短信接口失败: %s", err.Error())
		return false
	}
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.SignName = config.SignName
	request.TemplateCode = config.TemplateCode
	request.PhoneNumbers = phone
	//把生成的验证码传到短信模板中
	par, err := json.Marshal(map[string]interface{}{
		"code": code,
	})
	if err != nil {
		fmt.Errorf("解析code失败: %s", err.Error())
		return false
	}
	request.TemplateParam = string(par)

	//3.接收返回结果,并判断状态
	response, err := client.SendSms(request)
	fmt.Println(response)
	if err != nil {
		fmt.Errorf("发送失败: %s", err.Error())
		return false
	}

	if response.Code == "OK" {
		//将验证码保存到数据库中
		smsCode := model.SmsCode{
			Phone:      phone,
			BizId:      response.BizId,
			Code:       code,
			CreateTime: time.Now().Unix(),
		}
		memberDao := dao.MemberDao{tool.DbEngine}
		result := memberDao.InsertCode(smsCode)
		return result > 0
	}
	return false
}
