package pemiluscrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/kbiits/scrap_pemilu/pkg/stack"
)

const (
	baseUrl          = "https://pemilu2024.kpu.go.id"
	baseUrlTPSResult = "https://sirekap-obj-data.kpu.go.id/pemilu/hhcw/ppwp/"
	// example url to get tps result: https://sirekap-obj-data.kpu.go.id/pemilu/hhcw/ppwp/31/3174/317402/3174021004/3174021004010.json

	baseUrlArea = "https://sirekap-obj-data.kpu.go.id/wilayah/pemilu/ppwp/"

	getAllProvincesUrl = "https://sirekap-obj-data.kpu.go.id/wilayah/pemilu/ppwp/0.json"
)

type ScrapperSvc struct {
}
type TPSResultWriter interface {
	Write(result TPSResultWithMetadata) error
	WriteError(result TPSResultWithMetadata)
	Close()
}

func NewScrapper() *ScrapperSvc {
	return &ScrapperSvc{}
}

func (s *ScrapperSvc) GenerateTPSUri() <-chan AreaWithUrl {
	chanResult := make(chan AreaWithUrl, 100)

	provinces := getAllProvinceCodes()

	go func() {
		defer close(chanResult)

		parentStack := stack.New[Area]()
		for _, province := range provinces {
			parentStack.Push(province)
			s.GenerateTPSUriRecursively(parentStack, chanResult)
			parentStack.Pop()
		}
	}()

	return chanResult

}

func (s *ScrapperSvc) GenerateTPSUriInArea(stackArea *stack.Stack[Area]) <-chan AreaWithUrl {
	chanResult := make(chan AreaWithUrl)

	go func() {
		s.GenerateTPSUriRecursively(stackArea, chanResult)
		close(chanResult)
	}()

	return chanResult
}

func (s *ScrapperSvc) GenerateTPSUriRecursively(parentStack *stack.Stack[Area], chanUri chan<- AreaWithUrl) {
	parent, err := parentStack.Peek()
	if err != nil {
		panic(err)
	}

	parentsInOrder := parentStack.ToSlice()
	slices.Reverse(parentsInOrder)

	if parent.Level == 5 {
		tpsResultUrl := buildUrl(baseUrlTPSResult, parentsInOrder)
		areaWithUrl := AreaWithUrl{
			Area: Area{
				Code:  extractTPSAreaCodeFromUrl(tpsResultUrl),
				Level: 5,
			},
			UrlJson: tpsResultUrl,
		}
		chanUri <- areaWithUrl
		return
	}

	url := buildUrl(baseUrlArea, parentsInOrder)
	areas := getAreaByCode(url)

	for _, area := range areas {
		parentStack.Push(area)
		s.GenerateTPSUriRecursively(parentStack, chanUri)
		parentStack.Pop()
	}
}

func (s *ScrapperSvc) StartScrapping(writer TPSResultWriter, chanTps <-chan AreaWithUrl) {

	i := 1
	for tps := range chanTps {
		req, err := http.NewRequest(http.MethodGet, tps.UrlJson, nil)
		if err != nil {
			panic(err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		}

		respBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		resp.Body.Close()

		var tpsResult TPSResult
		err = json.Unmarshal(respBytes, &tpsResult)
		if err != nil {
			panic(err)
		}

		var tpsResultWithMeta = TPSResultWithMetadata{
			TPSResult: tpsResult,
			Url:       tps.UrlJson,
			Code:      tps.Code,
		}

		err = writer.Write(tpsResultWithMeta)
		if err != nil {
			writer.WriteError(tpsResultWithMeta)
			log.Default().Printf("error write tps result. url %s", tpsResultWithMeta.Url)
		}

		fmt.Printf("[NUM-%d] done scrapping %s\n", i, tps.UrlJson)
		i++
	}

	writer.Close()
}

func (s *ScrapperSvc) BuildStackAreaByAreaCodes(code string) (stackArea *stack.Stack[Area]) {
	stack := stack.New[Area]()

	var maxLevel int
	switch {
	case len(code) <= 6:
		// province -> city -> district use 2 char code
		maxLevel = len(code) / 2
	case len(code) == 10:
		// subdistrict use 4 char code
		maxLevel = 4
	case len(code) == 13:
		// tps use 3 char code
		maxLevel = 5
	default:
		panic("invalid length of code")
	}

	for i := 0; i < maxLevel; i++ {
		codeRune := []rune(code)
		var area Area
		switch {
		case i < 3:
			// for province, city and district
			area = Area{
				Code:  string(codeRune[:i*2+2]),
				Level: i + 1,
			}
		case i < 4:
			area = Area{
				Code:  string(codeRune[:10]),
				Level: i + 1,
			}
		case i < 5:
			area = Area{
				Code:  string(codeRune[:]),
				Level: i + 1,
			}
		default:
			panic("invalid i value")
		}

		stack.Push(area)
	}

	return stack
}

func getAllProvinceCodes() []Area {

	req, err := http.NewRequest(http.MethodGet, getAllProvincesUrl, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var respJson []Area
	err = json.Unmarshal(respBytes, &respJson)
	if err != nil {
		panic(err)
	}

	return respJson
}

func getAreaByCode(url string) []Area {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var respJson []Area
	err = json.Unmarshal(respBytes, &respJson)
	if err != nil {
		panic(err)
	}

	return respJson
}

func buildUrl(baseUrl string, parents []Area) string {
	url := baseUrl

	if len(parents) == 0 {
		return url + ".json"
	}

	lastIdx := len(parents) - 1

	for i, v := range parents {

		if i != lastIdx {
			url += fmt.Sprintf("%s/", v.Code)
		} else {
			url += fmt.Sprintf("%s.json", v.Code)
		}
	}

	return url
}

func extractTPSAreaCodeFromUrl(url string) string {
	urlSplitted := strings.Split(url, "/")
	if len(urlSplitted) == 0 {
		panic("invalid url")
	}

	lastPart := urlSplitted[len(urlSplitted)-1]
	if !strings.Contains(lastPart, ".json") {
		panic("invalid url, not contains .json")
	}

	return strings.TrimSuffix(lastPart, ".json")
}
