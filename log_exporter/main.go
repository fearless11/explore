package main

import (
	"flag"
	"strings"
	"time"

	"github.com/fearless11/explore/log_exporter/es"
	"github.com/fearless11/explore/log_exporter/prome"
)

var (
	u        string
	i        string
	t        string
	lastT    int64
	ESClient *es.ES
)

func init() {
	flag.StringVar(&u, "u", "http://192.168.56.101:9200", "elastic url")
	flag.StringVar(&i, "i", "filebeat-6.5.1-2019.05.07", "index prefix")
	flag.StringVar(&t, "t", "1m", "query time interval")
}

// consumption timer
func consumption(indexs []string, t string) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for {
		for _, index := range indexs {
			// lastT is a competition resources. How To solve it ?
			go ESClient.SearchScroll(index, t, &lastT)
			// specify condition
			condition := "status:101"
			go ESClient.SearchTerm(index, condition, "24h")
		}
		<-ticker.C
	}
}

func main() {
	flag.Parse()
	go prome.Start()
	// lastT = time.Now().Unix()
	lastT = 1557201600
	// index := i + time.Now().Format("2006.01.02")
	index := strings.Split(i, ",")
	ESClient = es.NewESClient(u)
	consumption(index, t)
}
