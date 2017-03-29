package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"net"
	"reflect"
)

var (
	httpClient  *http.Client
	protocol    string
	duration    int
	users       int
	delay       int
	host        string
	port        int
	timeout     int
	locationLat float64
	locationLon float64
	limit       int
	cityID      int
)
var requests []Request

func randomRequest() Request {
	return requests[rand.Intn(len(requests))]
}
func randomIntElem(elems []int) int {
	return elems[rand.Intn(len(elems))]
}
func randomIntSlice(elems []int) []int {
	end := rand.Intn(len(elems)) + 1
	start := rand.Intn(end)
	return elems[start:end]
}
func intsToStrings(a []int) []string {
	var results []string
	for _, x := range a {
		results = append(results, strconv.Itoa(x))
	}
	return results
}
func randomStringElem(elems []string) string {
	return elems[rand.Intn(len(elems))]
}

type Request interface {
	URL() string
}
type MallsByShop struct {
	shopIDs []int
}

func (r *MallsByShop) URL() string {
	shop := randomIntElem(r.shopIDs)
	return fmt.Sprintf("/malls/?city=%d&sort=name&shop=%d", cityID, shop)
}

type MallsByQuery struct {
	queries []string
}

func (r *MallsByQuery) URL() string {
	query := randomStringElem(r.queries)
	params := url.Values{}
	params.Add("city", strconv.Itoa(cityID))
	params.Add("sort", "name")
	params.Add("query", query)
	return fmt.Sprintf("/malls/?%s", params.Encode())
}

type MallDetails struct {
	mallIDs []int
}

func (r *MallDetails) URL() string {
	mall := randomIntElem(r.mallIDs)
	return fmt.Sprintf("/malls/%d/", mall)
}

type ShopsByMall struct {
	mallIDs []int
}

func (r *ShopsByMall) URL() string {
	mall := randomIntElem(r.mallIDs)
	return fmt.Sprintf("/shops/?city=%d&mall=%d&sort=name", cityID, mall)
}

type ShopsByQuery struct {
	queries []string
}

func (r *ShopsByQuery) URL() string {
	query := randomStringElem(r.queries)
	params := url.Values{}
	params.Add("city", strconv.Itoa(cityID))
	params.Add("sort", "name")
	params.Add("query", query)
	return fmt.Sprintf("/shops/?%s", params.Encode())
}

type ShopsByCategory struct {
	categoryIDs []int
}

func (r *ShopsByCategory) URL() string {
	category := randomIntElem(r.categoryIDs)
	return fmt.Sprintf("/shops/?city=%d&sort=name&category=%d&limit=%d", cityID, category, limit)
}

type ShopDetails struct {
	shopIDs []int
}

func (r *ShopDetails) URL() string {
	shop := randomIntElem(r.shopIDs)
	return fmt.Sprintf("/shops/%d/?city=%d&location_lat=%f&location_lon=%f", shop, cityID, locationLat, locationLon)

}

type CurrentMall struct{}

func (r *CurrentMall) URL() string {
	return fmt.Sprintf("/current_mall/?location_lat=%f&location_lon=%f", locationLat, locationLon)
}

type cities struct{}

func (r *cities) URL() string {
	return "/cities/?sort=name"
}

type ShopsInMalls struct {
	shopIDs []int
	mallIDs []int
}

func (r *ShopsInMalls) URL() string {
	randomShops := randomIntSlice(r.shopIDs)
	randomMalls := randomIntSlice(r.mallIDs)
	params := url.Values{"shops": intsToStrings(randomShops), "malls": intsToStrings(randomMalls)}
	return fmt.Sprintf("/shops_in_malls/?%s", params.Encode())
}

type Search struct {
	shopIDs []int
}

func (r *Search) URL() string {
	randomShops := randomIntSlice(r.shopIDs)
	params := url.Values{"shops": intsToStrings(randomShops)}
	return fmt.Sprintf("/search/?city=%d&location_lat=%f&location_lon=%f&sort=distance&limit=%d&%s", cityID, locationLat, locationLon, limit, params.Encode())
}

type ListOfIDs []struct {
	ID int `json:"id"`
}
type CategoriesResponse struct {
	Data ListOfIDs `json:"data"`
}
type PaginationResponse struct {
	Data struct {
		Results ListOfIDs `json:"results"`
	} `json:"data"`
}

func getShops() []int {
	const shopsCount = 10
	requestUrl := fmt.Sprintf("%s://%s:%d/shops/?city=%d&sort=-malls_count&limit=%d", protocol, host, port, cityID, shopsCount)
	locLog := log.WithFields(log.Fields{"url": requestUrl, "location": "get shops"})
	resp, err := http.Get(requestUrl)
	if err != nil {
		locLog.Panicf("Cannot get shop ids: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		locLog.Panicf("Expect ok status, get: %d", resp.StatusCode)
	}
	d := json.NewDecoder(resp.Body)
	body := PaginationResponse{}
	err = d.Decode(&body)
	if err != nil {
		locLog.Panicf("Cannot decode response: %s", err)
	}
	var shopIDs []int
	for _, result := range body.Data.Results {
		shopIDs = append(shopIDs, result.ID)
	}
	if len(shopIDs) != shopsCount {
		locLog.WithFields(log.Fields{"expect": shopsCount, "actual": len(shopIDs)}).Panicf("number of shops does not match")
	}
	return shopIDs
}
func getMalls() []int {
	const mallsCount = 10
	requestUrl := fmt.Sprintf("%s://%s:%d/malls/?city=%d&sort=-shops_count&limit=%d", protocol, host, port, cityID, mallsCount)
	locLog := log.WithFields(log.Fields{"url": requestUrl, "location": "get malls"})
	resp, err := http.Get(requestUrl)
	if err != nil {
		locLog.Panicf("Cannot get mall ids: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		locLog.Panicf("Expect ok status, get: %d", resp.StatusCode)
	}
	d := json.NewDecoder(resp.Body)
	body := PaginationResponse{}
	err = d.Decode(&body)
	if err != nil {
		locLog.Panicf("Cannot decode response: %s", err)
	}
	var mallIDs []int
	for _, result := range body.Data.Results {
		mallIDs = append(mallIDs, result.ID)
	}
	if len(mallIDs) != mallsCount {
		locLog.WithFields(log.Fields{"expect": mallsCount, "actual": len(mallIDs)}).Panicf("number of malls does not match")
	}
	return mallIDs
}
func getCategories() []int {
	requestUrl := fmt.Sprintf("%s://%s:%d/categories/?city=%d", protocol, host, port, cityID)
	locLog := log.WithFields(log.Fields{"url": requestUrl, "location": "get categories"})
	resp, err := http.Get(requestUrl)
	if err != nil {
		locLog.Panicf("Cannot get category ids: %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		locLog.Panicf("Expect ok status, get: %d", resp.StatusCode)
	}
	d := json.NewDecoder(resp.Body)
	body := CategoriesResponse{}
	err = d.Decode(&body)
	if err != nil {
		locLog.Panicf("Cannot decode response: %s", err)
	}
	var categoryIDs []int
	for _, result := range body.Data {
		categoryIDs = append(categoryIDs, result.ID)
	}
	return categoryIDs
}

type WorkerResult struct {
	OKResponses      int
	NonOKResponses   int
	TimeoutResponses int
	ErrorResponses   int
}

func delayWorker() WorkerResult {
	end := time.After(time.Second * time.Duration(duration))
	tick := time.NewTicker(time.Millisecond * time.Duration(delay))
	result := WorkerResult{}
	defer tick.Stop()
	for {
		select {
		case <-end:
			return result
		case <-tick.C:
			doRequest(&result)
		}
	}
}
func noDelayWorker() WorkerResult {
	end := time.After(time.Second * time.Duration(duration))
	result := WorkerResult{}
	for {
		select {
		case <-end:
			return result
		default:
			doRequest(&result)
		}
	}
}
func doRequest(result *WorkerResult) {
	req := randomRequest()
	uri := fmt.Sprintf("%s://%s:%d%s", protocol, host, port, req.URL())
	locLog := log.WithField("url", uri)
	resp, err := httpClient.Get(uri)
	if err != nil {
		if netError, ok := err.(net.Error); ok && netError.Timeout() {
			result.TimeoutResponses++
		} else {
			locLog.Errorf("Cannot do request: %s", err)
			result.ErrorResponses++
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		b, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(string(b))
		result.NonOKResponses++
		return
	}
	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		if netError, ok := err.(net.Error); ok && netError.Timeout() {
			result.TimeoutResponses++
		} else {
			fmt.Println(reflect.TypeOf(err))
			locLog.Errorf("Cannot read response body: %s", err)
			result.ErrorResponses++
		}
		return
	}
	result.OKResponses++
}
func displayResults(results []WorkerResult) {
	var totalOK, totalNonOK, totalTimeout, totalError int
	for _, r := range results {
		totalOK += r.OKResponses
		totalNonOK += r.NonOKResponses
		totalTimeout += r.TimeoutResponses
		totalError += r.ErrorResponses
	}
	rps := totalOK / duration
	log.WithFields(log.Fields{"OK": totalOK, "NonOK": totalNonOK, "Timeout": totalTimeout, "Error": totalError, "RPS": rps}).Info("END")
}
func main() {
	var ssl bool
	rand.Seed(time.Now().Unix())
	flag.IntVar(&duration, "duration", 30, "duration of load test in seconds")
	flag.IntVar(&users, "users", 10, "number of simultaneous users")
	flag.IntVar(&delay, "delay", 100, "time between user request in milliseconds, 0 means no delay")
	flag.StringVar(&host, "host", "localhost", "host")
	flag.IntVar(&port, "port", 8080, "port")
	flag.IntVar(&timeout, "timeout", 1, "request timeout in seconds")
	flag.Float64Var(&locationLat, "latitude", 55.654069, "x user coordinate")
	flag.Float64Var(&locationLon, "longitude", 37.646245, "y user coordinate")
	flag.IntVar(&limit, "limit", 10, "default pagination limit")
	flag.IntVar(&cityID, "city", 1, "city id")
	flag.BoolVar(&ssl, "ssl", false, "http or https")
	flag.Parse()
	if !ssl {
		protocol = "http"
	} else {
		protocol = "https"
	}
	httpClient = &http.Client{Timeout: time.Second * time.Duration(timeout)}
	log.Info("Start building requests...")
	shopIDs := getShops()
	mallIDs := getMalls()
	categoryIDs := getCategories()
	mallsByShopReq := &MallsByShop{shopIDs: shopIDs}
	mallsByQueryReq := &MallsByQuery{queries: []string{"тц", "т", "ат", "мо", "а", "б", "в", "с", "г", "д", "о", "е", "п", "м", "л", "н"}}
	mallDetailsReq := &MallDetails{mallIDs: mallIDs}
	shopsByMallReq := &ShopsByMall{mallIDs: mallIDs}
	shopsByQueryReq := &ShopsByQuery{queries: []string{"h", "крош", "ма", "спорт", "nik", "add"}}
	shopsByCategoryReq := &ShopsByCategory{categoryIDs: categoryIDs}
	shopDetailsReq := &ShopDetails{shopIDs: shopIDs}
	currentMallReq := &CurrentMall{}
	shopsInMallsReq := &ShopsInMalls{shopIDs: shopIDs, mallIDs: mallIDs}
	searchReq := &Search{shopIDs: shopIDs}

	requests = []Request{
		mallsByShopReq,
		mallsByQueryReq,
		mallDetailsReq,
		shopsByMallReq,
		shopsByQueryReq,
		shopsByCategoryReq,
		shopDetailsReq,
		currentMallReq,
		shopsInMallsReq,
		searchReq,
		searchReq,
	}
	var worker func() WorkerResult
	if delay != 0 {
		worker = delayWorker
	} else {
		worker = noDelayWorker
	}
	wg := sync.WaitGroup{}
	m := sync.Mutex{}
	var results []WorkerResult
	log.WithFields(log.Fields{"users": users, "duration": duration, "delay": delay}).Info("Start sending")
	for i := 0; i < users; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := worker()
			m.Lock()
			results = append(results, result)
			m.Unlock()
		}()
	}
	wg.Wait()
	displayResults(results)
}
