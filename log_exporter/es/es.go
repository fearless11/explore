package es

import (
	"context"
	"fmt"
	"io"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/fearless11/explore/log_exporter/prome"
	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/olivere/elastic.v5"
)

type ES struct {
	Client *elastic.Client
}

func NewESClient(url string) *ES {
	client, err := elastic.NewClient(elastic.SetSniff(false), elastic.SetHealthcheckInterval(5*time.Second), elastic.SetURL(url))
	if err != nil {
		log.Fatal("connect es fail ", err)
	}
	return &ES{Client: client}
}

// NgxLog nginx logfromat
type NgxLog struct {
	Status      string `json:"status"`
	ServerName  string `json:"server_name"`
	RequestTime string `json:"request_time"`
}

// SearchScroll specify interval consumption elastic data to prometheus
func (e *ES) SearchScroll(index string, intervalT string, lastT *int64) {
	// set the end time to the last pause time, increase by interval
	dt, _ := time.ParseDuration("-" + intervalT)
	endT := time.Unix(*lastT, 0)
	startT := endT.Add(dt)
	// move current consumption time by interval
	dt, _ = time.ParseDuration(intervalT)
	*lastT = endT.Add(dt).Unix()

	// consumption elastic data by interval
	queryL := elastic.NewRangeQuery("@timestamp").Gte(startT).Lte(endT)
	scrollC := elastic.NewScrollService(e.Client)
	searchR, err := scrollC.Index(index).Query(queryL).KeepAlive("1m").Size(500).Do(context.Background())
	if err != nil && err != io.EOF {
		log.Println("scroll search error ", err)
		return
	}

	// search total data
	fmt.Println("total", searchR.Hits.TotalHits)

	var nlog NgxLog
	for {
		for _, item := range searchR.Each(reflect.TypeOf(nlog)) {
			t := item.(NgxLog)
			httpRtt, err := strconv.ParseFloat(t.RequestTime, 64)
			if err != nil {
				log.Println("httpRtt conversion error ", err)
				continue
			}
			prome.HttpResponseStatus.With(prometheus.Labels{"url": t.ServerName, "code": t.Status}).Inc()
			prome.HttpResponseDuration.With(prometheus.Labels{"url": t.ServerName}).Observe(float64(httpRtt))
		}

		// batch search total
		// fmt.Println("scrollID: ", len(searchR.Hits.Hits))
		searchR, err = scrollC.ScrollId(searchR.ScrollId).Do(context.Background())
		if err != nil && err != io.EOF {
			log.Println("scrollID search  error", err)
			continue
		}

		if len(searchR.Hits.Hits) == 0 {
			break
		}
	}

	err = scrollC.Clear(context.Background())
	if err != nil {
		log.Println("scrollID search  error", err)
		return
	}
}

// SearchTerm specify query condition
func (e *ES) SearchTerm(index string, qCondition string, intervalT string) (total int64) {

	queryC := strings.Split(qCondition, ":")
	termQuery := elastic.NewTermQuery(queryC[0], queryC[1])
	dt, _ := time.ParseDuration("-" + intervalT)
	endT := time.Now()
	startT := endT.Add(dt)
	rangeQuery := elastic.NewRangeQuery("@timestamp").Gte(startT).Lte(endT)
	query := elastic.NewBoolQuery().Must(termQuery, rangeQuery)

	scrollC := elastic.NewScrollService(e.Client)
	searchR, err := scrollC.Index(index).Query(query).KeepAlive("1m").Size(500).Do(context.Background())
	if err != nil && err != io.EOF {
		log.Println("scroll search error ", err)
		return
	}
	return searchR.Hits.TotalHits
}

// SearchRange specify query condition range
func (e *ES) SearchRange(index string, boundary string, intervalT string) (total int64) {

	dt, _ := time.ParseDuration("-" + intervalT)
	endT := time.Now()
	startT := endT.Add(dt)
	timeRangeQuery := elastic.NewRangeQuery("@timestamp").Gte(startT).Lte(endT)
	conditionRangeQuery := elastic.NewRangeQuery("request_time").Gte(boundary)
	termsQuery := elastic.NewTermsQuery("status", "101", "499")
	query := elastic.NewBoolQuery().Must(conditionRangeQuery, timeRangeQuery).MustNot(termsQuery)

	scrollC := elastic.NewScrollService(e.Client)
	searchR, err := scrollC.Index(index).Query(query).KeepAlive("1m").Size(500).Do(context.Background())
	if err != nil && err != io.EOF {
		log.Println("scroll search error ", err)
		return
	}
	return searchR.Hits.TotalHits
}
