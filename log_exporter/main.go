package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/fearless11/explore/log_exporter/alarm"
	"github.com/fearless11/explore/log_exporter/es"
)

var (
	u               string
	i               string
	t               string
	c               string
	n               int64
	b               string
	p               string
	lastT           int64
	thresholdNumber int64
	phoneID         []string
	ESClient        *es.ES
)

func init() {
	flag.StringVar(&u, "u", "http://127.0.0.1:9200", "elastic url")
	flag.StringVar(&i, "i", "filebeat-6.5.1,filebeat-nginx", "index prefix")
	flag.StringVar(&t, "t", "5m", "query time interval")
	flag.StringVar(&c, "c", "500,502,504", "abnormal response status code")
	flag.Int64Var(&n, "n", 20, "threshold number of abnormal events")
	flag.StringVar(&b, "b", "10s", "set abnormal response time, default 10s")
	flag.StringVar(&p, "p", "112321321,1231231", "set notify to phone")
}

// consumption to the prometheus
func consumption(indexs []string, t string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		for _, index := range indexs {
			// lastT is a competition resources. How To solve it ?
			go ESClient.SearchScroll(index, t, &lastT)
		}
		<-ticker.C
	}
}

// abnormalJudgment check if the request is normal
func abnormalJudgment(index string, codes []string, boundary string, intervalT string) {

	for _, code := range codes {
		queryCondition := "status:" + code
		codeNumOfException := ESClient.SearchTerm(index, queryCondition, intervalT)
		timeNumOfException := ESClient.SearchRange(index, boundary, intervalT)

		if codeNumOfException > thresholdNumber {
			content := fmt.Sprintf("[日志告警]\n索引:%v\n时间:%v\n内容:%v分钟内%v数量为%v,超过阈值%v", index, time.Now().Format("2006-01-02 15:04:05"), intervalT, code, codeNumOfException, thresholdNumber)
			alarm.Event(content, phoneID)
		}

		if timeNumOfException > thresholdNumber {
			content := fmt.Sprintf("[日志告警]\n索引:%v\n时间:%v\n内容:%v分钟内耗时超过%vs数量为%v,超过阈值%v", index, time.Now().Format("2006-01-02 15:04:05"), intervalT, boundary, timeNumOfException, thresholdNumber)
			alarm.Event(content, phoneID)
		}
	}
}

func main() {
	flag.Parse()
	indexPrefix := strings.Split(i, ",")
	codes := strings.Split(c, ",")
	phoneID = strings.Split(p, ",")
	boundary := b
	thresholdNumber = n
	ESClient = es.NewESClient(u)

	tickerT := time.NewTicker(5 * time.Minute)
	defer tickerT.Stop()

	for {
		log.Println("check nginx log", time.Now().Format("2006-01-02 15:04:05"))
		for _, index := range indexPrefix {
			index = index + "-" + time.Now().Format("2006.01.02")
			abnormalJudgment(index, codes, boundary, t)
		}
		<-tickerT.C
	}
}
