package spider

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/cihub/seelog"
)

type Processer interface {
	Process(string, bool, *http.Response) ([]*http.Request, error)
	Finish()
}

const (
	SelectorDefault = "script, link, a, img, frame, iframe, area, base, blockquote, body, del, head, ins, object, q"
)

type OptionDomProcesser func(*DomProcesser)

func OptionDomProcesserSelectors(selectors []string) OptionDomProcesser {
	return func(processer *DomProcesser) {
		processer.selectors = strings.Join(selectors, ",")
	}
}

type DomProcesser struct {
	selectors string
}

func NewDomProcesser(options ...OptionDomProcesser) *DomProcesser {
	dp := &DomProcesser{}
	for _, option := range options {
		option(dp)
	}
	if dp.selectors == "" {
		dp.selectors = SelectorDefault
	}
	return dp
}

func (dp *DomProcesser) Finish() {
	return
}

func (dp *DomProcesser) Process(charSet string, certain bool, rsp *http.Response) ([]*http.Request, error) {
	defer rsp.Body.Close()

	/*
		utfReader, err := iconv.NewReader(rsp.Body, charSet, "utf-8")
		if err != nil {
			seelog.Errorf("DomProcesser::Process | ioconv new reader charset: %s, err: %s", charSet, err)
			return nil, err
		}
		dom, err := goquery.NewDocumentFromReader(utfReader)
		if err != nil {
			seelog.Errorf("DomProcesser::Process | new document from reader err: %s", err)
			return nil, err
		}
	*/
	dom, err := goquery.NewDocumentFromResponse(rsp)
	if err != nil {
		seelog.Errorf("DomProcesser::Process | new document from reader err: %s", err)
		return nil, err
	}

	base := rsp.Request.URL.String()

	links := extractLinks(base, dom)

	var reqs []*http.Request
	for _, link := range links {
		req, err := http.NewRequest(http.MethodGet, link, nil)
		if err != nil {
			continue
		}
		//防防盗链
		req.Header.Add("Refer", rsp.Request.URL.String())
		reqs = append(reqs, req)
	}
	return reqs, nil
}

func extractLinks(base string, doc *goquery.Document) []string {
	internalUrls := []string{}
	if doc != nil {
		doc.Find(SelectorDefault).Each(func(i int, s *goquery.Selection) {
			sub, exists := s.Attr("href")
			if !exists {

				sub, exists = s.Attr("src")
				if !exists {

					sub, exists = s.Attr("action")
					if !exists {

						sub, exists = s.Attr("codebase")
						if !exists {

							sub, exists = s.Attr("cite")
							if !exists {

								sub, exists = s.Attr("longdesc")
								if !exists {

									sub, exists = s.Attr("usemap")
									if !exists {

										sub, exists = s.Attr("profile")
										if !exists {
											return
										}
									}
								}
							}
						}
					}
				}
			}

			u := mergeUrl(base, sub)
			if u != "" {
				internalUrls = append(internalUrls, u)
			}
			seelog.Infof("Spider::extractLinks | base: %s, sub: %s, after merge: %s", base, sub, u)
		})
		return internalUrls
	}
	return internalUrls
}

func mergeUrl(base, sub string) string {
	subU, err := url.Parse(sub)
	if err != nil {
		return ""
	}

	if subU.IsAbs() {
		//绝对路径且同域
		if strings.HasPrefix(sub, base) {
			return sub
		}
		return ""
	}

	baseU, err := url.Parse(base)
	if err != nil {
		return ""
	}

	mergeU := baseU.ResolveReference(subU)
	return mergeU.String()
}
