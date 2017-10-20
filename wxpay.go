package wxpay

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"
)

type WXPay struct {
	appId     string
	apiKey    string
	mchId     string
	Client    *http.Client
	NotifyURL string
}

func New(appId, apiKey, mchId string) (client *WXPay) {
	client = &WXPay{}
	client.appId = appId
	client.mchId = mchId
	client.apiKey = apiKey
	client.Client = http.DefaultClient
	return client
}

func (this *WXPay) doRequest(method, url string, param map[string]interface{}, results interface{}) (err error) {
	param["appid"] = this.appId
	param["mch_id"] = this.mchId
	param["nonce_str"] = getNonceStr()
	if _, ok := param["notify_url"]; ok == false {
		if len(this.NotifyURL) > 0 {
			param["notify_url"] = this.NotifyURL
		}
	}

	var sign = signMD5(param, this.apiKey)
	param["sign"] = sign

	req, err := http.NewRequest(method, url, strings.NewReader(mapToXML(param)))
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	resp, err := this.Client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, results)

	return err
}

func (this *WXPay) DoRequest(method, url string, param map[string]interface{}, results interface{}) (err error) {
	return this.doRequest(method, url, param, results)
}
