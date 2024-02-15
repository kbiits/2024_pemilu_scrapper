package kawalpemiluscrapper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kbiits/scrap_pemilu/pkg/cache"
	pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"
)

const (
	baseUrl = "https://kp24-fd486.et.r.appspot.com/h"
)

type KawalPemiluScrapper struct {
	cache *cache.Cache
}

func NewKawalPemiluScrapper(cache *cache.Cache) *KawalPemiluScrapper {
	return &KawalPemiluScrapper{
		cache: cache,
	}
}

func (s *KawalPemiluScrapper) StartScrapping(writer pemiluscrapper.TPSResultWriter, tpsChan <-chan pemiluscrapper.AreaWithUrl) {
	for area := range tpsChan {
		resultSubDistrict := s.getFromArea(area)
		tpsNumInStr := []rune(area.Code)[len(area.Code)-3:]
		tpsNumInInt, err := strconv.ParseInt(string(tpsNumInStr), 10, 64)
		if err != nil {
			panic(err)
		}

		resultTPS, ex := resultSubDistrict.Result.Aggregated[int(tpsNumInInt)]
		if !ex || len(resultTPS) == 0 {
			log.Default().Printf("not exists tps data with code %s", area.Code)
			continue
		}

		firstResult := resultTPS[0]

		compatibleKPUResult := firstResult.ConvertToCompatibleKPUType(buildUrlFromAreaCode(area.Code), area.Code)
		err = writer.Write(compatibleKPUResult)
		if err != nil {
			log.Default().Printf("failed to write tps %s from kawal pemilu", area.Code)
			writer.WriteError(compatibleKPUResult)
		}
	}
	writer.Close()
}

func (s *KawalPemiluScrapper) getFromArea(area pemiluscrapper.AreaWithUrl) KawalPemiluTPSData {
	areaCodeSubDistrict := string([]rune(area.Code)[:len(area.Code)-3])
	result := s.cache.SetOrGet(areaCodeSubDistrict, time.Second*60, func() interface{} {
		req, err := http.NewRequest(http.MethodGet, baseUrl+"?id="+areaCodeSubDistrict, nil)
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

		var result KawalPemiluTPSData
		err = json.Unmarshal(respBytes, &result)
		if err != nil {
			panic(err)
		}

		return result
	})
	if val, isOk := result.(KawalPemiluTPSData); !isOk {
		panic("invalid type from cache")
	} else {
		return val
	}
}

func buildUrlFromAreaCode(code string) string {
	tpsNumInStr := []rune(code)[len(code)-3:]
	tpsNumInInt, err := strconv.ParseInt(string(tpsNumInStr), 10, 64)
	if err != nil {
		panic(err)
	}
	areaCodeSubDistrict := string([]rune(code)[:len(code)-3])

	return fmt.Sprintf("https://kawalpemilu.org/h/%s%d", areaCodeSubDistrict, tpsNumInInt)
}
