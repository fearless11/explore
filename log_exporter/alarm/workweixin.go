package alarm

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/axgle/mahonia"
)

// workWeiXin
// golang url encode https://www.urlencoder.io/golang/
func workWeiXin(content string, phoneID string) {

	// utf-8 to gbk
	enc := mahonia.NewEncoder("gbk")
	contentGBk := enc.ConvertString(content)
	// url encode
	params := url.Values{}
	params.Add("userName", "xxx")
	params.Add("memberid", "x")
	params.Add("phone", phoneID)
	params.Add("content", contentGBk)
	params.Add("send", "x")
	sendURL := "http://xxxx/api/v1" + "?" + params.Encode()

	c := http.Client{Timeout: 10 * time.Second}
	resp, err := c.Get(sendURL)
	if err != nil {
		log.Println("send weixin error ", err)
		return
	}

	byteBody, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Println("read body response error ", err)
		return
	}
	log.Println("send success response ", string(byteBody))
}

func Event(context string, phoneID []string) {
	for _, phone := range phoneID {
		workWeiXin(context, phone)
	}
}
