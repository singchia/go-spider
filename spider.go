package spider

import (
	"bytes"
	"context"
	"errors"
	"io"
	"math/rand"
	"mime"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cihub/seelog"
	"golang.org/x/net/html/charset"
)

const (
	SleepTypeNode = iota
	SleepTypeFixed
	SleepTypeRandom
)

const (
	SpiderConcuDefault = 100

	SleepTypeDefault = SleepTypeRandom
	SleepMinDefault  = 0
	SleepMaxDefault  = 1000

	DownloadPathDefault = "/tmp/tamper"
)

type OptionSpider func(*Spider)

func OptionSpiderFilter(filter Filter) OptionSpider {
	return func(spider *Spider) {
		spider.filter = filter
	}
}

func OptionSpiderProcesser(processer Processer) OptionSpider {
	return func(spider *Spider) {
		spider.processer = processer
	}
}

func OptionSpiderDownloader(downloader Downloader) OptionSpider {
	return func(spider *Spider) {
		spider.downloader = downloader
	}
}

func OptionSpiderScheduler(scheduler Scheduler) OptionSpider {
	return func(spider *Spider) {
		spider.scheduler = scheduler
	}
}

func OptionSpiderSchduler(resourceMgr ResourceManager) OptionSpider {
	return func(spider *Spider) {
		spider.resourceMgr = resourceMgr
	}
}

func OptionSpiderConcu(concu uint32) OptionSpider {
	return func(spider *Spider) {
		if concu == 0 {
			concu = 1
		}
		spider.concu = concu
	}
}

func OptionSpiderSleep(tp, min, max uint) OptionSpider {
	return func(spider *Spider) {
		if tp > SleepTypeRandom {
			return
		}
		if min >= max {
			return
		}
		spider.sleepType = tp
		spider.sleepMin = min
		spider.sleepMax = max
	}
}

func OptionSpiderRequestCheckRedirect(checkRedirect func(req *http.Request, via []*http.Request) error) OptionSpider {
	return func(spider *Spider) {
		spider.checkRedirect = checkRedirect
	}
}

func OptionSpiderRequestHeader(key, value string) OptionSpider {
	return func(spider *Spider) {
		spider.defaultHeader.Add(key, value)
	}
}

func OptionSpiderRequestTimeout(timeout time.Duration) OptionSpider {
	return func(spider *Spider) {
		spider.timeout = timeout
	}
}

func OptionSpiderResponseChunkedAllowed(allowed bool) OptionSpider {
	return func(spider *Spider) {
		spider.rspChunkedAllowed = allowed
	}
}

//不包含content-type嗅探
type Spider struct {
	results map[string]*Result
	mutex   sync.RWMutex

	defaultHeader http.Header
	checkRedirect func(req *http.Request, via []*http.Request) error
	timeout       time.Duration

	rspChunkedAllowed bool

	//对外模块
	filter     Filter
	processer  Processer
	downloader Downloader

	//下个版本可以废除
	scheduler   Scheduler
	resourceMgr ResourceManager

	concu uint32 //并发

	//sleep duration in millisecond
	sleepMin  uint
	sleepMax  uint
	sleepType uint
}

type Result struct {
	Error string `json:"error,omitempty"`

	//request result
	Req *http.Request  `json:"-"`
	Rsp *http.Response `json:"-"`

	//inner parser result
	Size    int64  `json:"size,omitempty"`
	Suffix  string `json:"suffix,omitempty"`
	CharSet string `json:"charset,omitempty"`

	//download result
	UrlPath  *string `json:"url_path,omitempty"`
	HdrPath  *string `json:"hdr_path,omitempty"`
	BodyPath *string `json:"body_path,omitempty"`

	//processer result
	Depth uint     `json:"depth"`
	Subs  []string `json:"subs,omitempty"`
}

func NewSpider(options ...OptionSpider) *Spider {
	spider := &Spider{
		defaultHeader:     make(http.Header),
		results:           make(map[string]*Result),
		rspChunkedAllowed: true,
		concu:             SpiderConcuDefault,
		sleepMin:          SleepMinDefault,
		sleepMax:          SleepMaxDefault,
		sleepType:         SleepTypeDefault,
	}

	for _, option := range options {
		option(spider)
	}
	if spider.scheduler == nil {
		spider.scheduler = NewSchedulerChan()
	}
	if spider.resourceMgr == nil {
		spider.resourceMgr = NewResourceChan(spider.concu)
	}
	if spider.processer == nil {
		spider.processer = NewDomProcesser()
	}
	if spider.downloader == nil {
		spider.downloader = NewFileDownloader(DownloadPathDefault)
	}
	return spider
}

func (spider *Spider) Run() *Spider {
	for {
		req := spider.scheduler.Poll()
		if req == nil {
			if spider.resourceMgr.Used() == uint32(0) {
				spider.processer.Finish()
				break
			}
			time.Sleep(500 * time.Millisecond)
			continue
		}

		for k, vs := range spider.defaultHeader {
			if req.Header.Get(k) == "" {
				for _, v := range vs {
					req.Header.Add(k, v)
				}
			}
		}

		spider.resourceMgr.Acquire()

		go func(req *http.Request) {
			defer spider.resourceMgr.Release()

			url := req.URL.String()
			if spider.exists(url) {
				return
			}
			result := &Result{Req: req}
			spider.record(url, result)

			defer func() {
				spider.sleep()
			}()

			depth, ok := req.Context().Value("depth").(uint)
			if !ok {
				result.Error = "no depth in http request context"
				return
			}
			result.Depth = depth

			client := &http.Client{
				CheckRedirect: spider.checkRedirect,
				Timeout:       spider.timeout,
			}
			rsp, err := client.Do(req)
			if err != nil {
				seelog.Errorf("Spider::Run | client do err: %s", err)
				result.Error = err.Error()
				return
			}
			closer := rsp.Body
			defer closer.Close()

			if httpResponseChunked(rsp.TransferEncoding) && !spider.rspChunkedAllowed {
				result.Error = "unsupported chunked transfer encoding"
				return
			}

			var suffix string
			preview := make([]byte, 1024)
			n, _ := rsp.Body.Read(preview)
			preview = preview[:n]

			if n > 0 {
				suffix, err = httpResponseContentType(
					preview,
					rsp.Header.Get("Content-Type"))
				if err != nil {
					seelog.Errorf("Spider::Run | http response content type err: %s", err)
					result.Error = err.Error()
					return
				}
				charSet, certain := httpResponseCharset(
					preview,
					rsp.Header.Get("Content-Type"))

				result.Size = rsp.ContentLength
				result.Suffix = suffix
				result.CharSet = charSet

				if !spider.filterCheck(req.Method, result.Size, result.Suffix) {
					result.Error = "filter rejected request"
					return
				}

				mergeReader := io.MultiReader(bytes.NewBuffer(preview), rsp.Body)

				//downloader and processer
				wg := sync.WaitGroup{}
				var urlPath, hdrPath, bodyPath *string
				var reqs []*http.Request
				var errD, errP error

				switch result.Suffix {
				case ContentTypeHTML, ContentTypeHTM, ContentTypeXHTML, ContentTypeXML:
					//既解析又下载
					pipeReader, pipeWriter := io.Pipe()
					teeReader := io.TeeReader(mergeReader, pipeWriter)
					rsp.Body = pipeReader
					wg.Add(2)
					go func() {
						defer wg.Done()
						reqs, errP = spider.processer.Process(
							charSet,
							certain,
							rsp)
						pipeReader.Close()
					}()

					go func() {
						defer wg.Done()
						urlPath, hdrPath, bodyPath, errD =
							spider.downloader.Download(
								req.URL,
								rsp.Header,
								teeReader,
								suffix)
						pipeWriter.Close()
					}()

				default: //仅下载
					wg.Add(1)
					go func() {
						defer wg.Done()
						urlPath, hdrPath, bodyPath, errD =
							spider.downloader.Download(
								req.URL,
								rsp.Header,
								mergeReader,
								suffix)
					}()
				}
				wg.Wait()

				if errP != nil {
					seelog.Errorf("Spider::Run | processer err: %s", errP)
					result.Error = errP.Error()
					return
				}

				if errD != nil {
					seelog.Errorf("Spider::Run | downloader err: %s", errD)
					result.Error = errD.Error()
					return
				}
				result.UrlPath = urlPath
				result.HdrPath = hdrPath
				result.BodyPath = bodyPath

				for _, subReq := range reqs {
					reqWithDepth := subReq.WithContext(context.WithValue(subReq.Context(), "depth", uint(depth+1)))
					result.Subs = append(result.Subs, subReq.URL.String())
					spider.scheduler.Push(reqWithDepth)
				}

			} else {
				result.Error = "without body assigned to the url"
				return
			}
		}(req)
	}
	return spider
}

func (spider *Spider) AddRequest(req *http.Request) *Spider {
	if req == nil {
		return spider
	}
	reqWithDepth := req.WithContext(context.WithValue(req.Context(), "depth", uint(0)))
	spider.scheduler.Push(reqWithDepth)
	return spider
}

func (spider *Spider) Result() map[string]*Result {
	spider.mutex.RLock()
	defer spider.mutex.RUnlock()

	return spider.results
}

func (spider *Spider) filterCheck(method string, size int64, suffix string) bool {
	if spider.filter != nil {
		if method == "https" && !spider.filter.HttpsAllow() {
			return false
		}
		if !spider.filter.SizeAllow(size) {
			return false
		}
		if !spider.filter.SuffixAllow(suffix) {
			return false
		}
	}
	return true
}

func (spider *Spider) exists(url string) bool {
	spider.mutex.RLock()
	defer spider.mutex.RUnlock()

	_, ok := spider.results[url]
	return ok
}

func (spider *Spider) record(url string, result *Result) {
	spider.mutex.Lock()
	defer spider.mutex.Unlock()

	spider.results[url] = result
}

func (spider *Spider) sleep() {
	switch spider.sleepType {
	case SleepTypeNode:
		return

	case SleepTypeFixed:
		time.Sleep(time.Duration(spider.sleepMin) * time.Millisecond)
		return

	case SleepTypeRandom:
		random := rand.Intn(int(spider.sleepMax-spider.sleepMin)) +
			int(spider.sleepMin)
		time.Sleep(time.Duration(random) * time.Millisecond)
		return
	}
	return
}

func httpResponseChunked(transferEncoding []string) bool {
	for _, encoding := range transferEncoding {
		if encoding == "chunked" {
			return true
		}
	}
	return false
}

//https://tools.ietf.org/html/rfc2045 #5.1
func httpResponseContentType(data []byte, contentType string) (string, error) {
	arbitrate := func(data []byte) string {
		return http.DetectContentType(data)
	}

	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil || !strings.Contains(mediaType, "/") {
		return "", errors.New("unsupported content-type")
	}
	v, ok := ContentTypes[mediaType]
	if !ok {
		mediaType = arbitrate(data)
		v, ok = ContentTypes[mediaType]
		if !ok {
			return "", errors.New("unsupported content-type")
		}
	}
	return v, nil
}

func httpResponseCharset(data []byte, contentType string) (string, bool) {
	_, name, certain := charset.DetermineEncoding(data, contentType)
	return name, certain
}
