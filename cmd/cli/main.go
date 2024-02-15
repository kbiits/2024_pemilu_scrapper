package main

import (
	"time"

	"github.com/kbiits/scrap_pemilu/pkg/cache"
	pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"
	tpswriter "github.com/kbiits/scrap_pemilu/services/tps_writer"
)

func main() {

	csvWriter := tpswriter.NewCSVWriter("./data/dki_jakarta_kpu.csv")

	kpuCache := cache.New(time.Second*60, time.Second*60)
	scrapperService := pemiluscrapper.NewScrapper(kpuCache)
	// urlsChan := scrapperService.GenerateTPSUri()

	stacks := scrapperService.BuildStackAreaByAreaCodes("31")
	urlsChan := scrapperService.GenerateTPSUriInArea(stacks)

	scrapperService.StartScrapping(csvWriter, urlsChan)
}
