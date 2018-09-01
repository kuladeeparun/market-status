package controllers

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type context struct {
	Data []data
	Time string
}

type data struct {
	Symbol                   string
	Series                   string
	OpenPrice                string
	HighPrice                string
	LowPrice                 string
	Ltp                      string
	PreviousPrice            string
	NetPrice                 string
	TradedQuantity           string
	TurnoverInLakhs          string
	LastCorpAnnouncementDate string
	LastCorpAnnouncement     string
}

var cache = context{dat, ""}
var dat = []data{}

func gainers(templates *template.Template) {

	//Mock endpoint to test cache refresh because the market was closed over the weekend
	http.HandleFunc("/mock", func(w http.ResponseWriter, req *http.Request) {
		str := `{'data':[{'symbol':'DRREDDY','series':'EQ','openPrice':'2,397.00','highPrice':'2,499.00','lowPrice':'2,395.00','ltp':'2,492.00','previousPrice':'2,380.15','netPrice':'4.70','tradedQuantity':'19,81,320','turnoverInLakhs':'48,618.22','lastCorpAnnouncementDate':'16-Jul-2018','lastCorpAnnouncement':'Annual General Meeting-Dividend- Rs 20 Per Share'},
		{'symbol':'TECHM','series':'EQ','openPrice':'734.00','highPrice':'768.75','lowPrice':'731.25','ltp':'766.00','previousPrice':'731.90','netPrice':'4.66','tradedQuantity':'72,76,056','turnoverInLakhs':'54,806.89','lastCorpAnnouncementDate':'26-Jul-2018','lastCorpAnnouncement':'Annual General Meeting-Dividend- Rs 14 Per Share'},
		{'symbol':'LUPIN','series':'EQ','openPrice':'896.80','highPrice':'939.80','lowPrice':'895.75','ltp':'931.50','previousPrice':'893.65','netPrice':'4.24','tradedQuantity':'68,88,252','turnoverInLakhs':'63,724.60','lastCorpAnnouncementDate':'30-Jul-2018','lastCorpAnnouncement':'Annual General Meeting-Dividend Rs 5-- Per Share'},
		{'symbol':'TATAMOTORS','series':'EQ','openPrice':'259.10','highPrice':'268.25','lowPrice':'258.10','ltp':'266.95','previousPrice':'259.35','netPrice':'2.93','tradedQuantity':'1,23,60,078','turnoverInLakhs':'32,680.05','lastCorpAnnouncementDate':'18-Jul-2016','lastCorpAnnouncement':'Dividend - Re 0.20-- Per Share'},
		{'symbol':'HCLTECH','series':'EQ','openPrice':'1,021.00','highPrice':'1,055.00','lowPrice':'1,019.50','ltp':'1,044.35','previousPrice':'1,016.55','netPrice':'2.73','tradedQuantity':'32,79,336','turnoverInLakhs':'34,278.57','lastCorpAnnouncementDate':'10-Sep-2018','lastCorpAnnouncement':'Annual General Meeting'},
		{'symbol':'UPL','series':'EQ','openPrice':'695.00','highPrice':'721.90','lowPrice':'689.45','ltp':'713.35','previousPrice':'696.15','netPrice':'2.47','tradedQuantity':'58,40,817','turnoverInLakhs':'41,148.56','lastCorpAnnouncementDate':'09-Aug-2018','lastCorpAnnouncement':'Annual General Meeting-Dividend- Rs 8 Per Share'},
		{'symbol':'POWERGRID','series':'EQ','openPrice':'194.85','highPrice':'202.25','lowPrice':'194.85','ltp':'201.15','previousPrice':'196.95','netPrice':'2.13','tradedQuantity':'66,68,243','turnoverInLakhs':'13,349.16','lastCorpAnnouncementDate':'10-Sep-2018','lastCorpAnnouncement':'Annual General Meeting-Dividend Rs 2.80 Per Share'},
		{'symbol':'BAJAJ-AUTO','series':'EQ','openPrice':'2,720.00','highPrice':'2,755.00','lowPrice':'2,705.80','ltp':'2,753.95','previousPrice':'2,701.65','netPrice':'1.94','tradedQuantity':'5,12,215','turnoverInLakhs':'13,988.75','lastCorpAnnouncementDate':'05-Jul-2018','lastCorpAnnouncement':'Annual General Meeting - Dividend- Rs 60 Per Share'},
		{'symbol':'SUNPHARMA','series':'EQ','openPrice':'643.45','highPrice':'659.40','lowPrice':'643.30','ltp':'651.90','previousPrice':'639.95','netPrice':'1.87','tradedQuantity':'1,29,72,927','turnoverInLakhs':'84,779.38','lastCorpAnnouncementDate':'17-Sep-2018','lastCorpAnnouncement':'Annual General Meeting-Dividend Rs 2 Per Share'},
		{'symbol':'CIPLA','series':'EQ','openPrice':'649.55','highPrice':'664.00','lowPrice':'649.55','ltp':'661.15','previousPrice':'649.45','netPrice':'1.80','tradedQuantity':'30,56,580','turnoverInLakhs':'20,186.57','lastCorpAnnouncementDate':'13-Aug-2018','lastCorpAnnouncement':'Annual General Meeting-Dividend- Rs 3 Per Share'}],'time':'{time}'}`
		str = strings.Replace(str, "{time}", time.Now().String(), 1)
		str = strings.Replace(str, "'", "\"", -1)
		w.Header().Add("Content Type", "text/json")
		w.Write([]byte(str))
	})

	go fetch()

	//Refresh the cache every 5 minutes
	ticker := time.NewTicker(5 * time.Minute)
	go func() {
		for range ticker.C {
			fetch()
		}
	}()

	http.HandleFunc("/gainers", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Content Type", "text/html")
		templates.Lookup("gainers.html").Execute(w, cache)

		//fmt.Printf("%+v \n", cache)
		//fmt.Println(reflect.TypeOf(t))
	})

}

func fetch() {
	log.Println("Fetching fresh data")

	resp, err := http.Get("https://www.nseindia.com/live_market/dynaContent/live_analysis/gainers/niftyGainers1.json")
	//resp, err := http.Get("http://localhost:5000/mock")
	check(err, "Error while fetching JSON")

	bytes, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(bytes, &cache)
	check(err, "Error while unmarshaling")

}

func check(err error, msg string) {
	if err != nil {
		log.Println(msg)
		log.Fatal(err)
	}
}
