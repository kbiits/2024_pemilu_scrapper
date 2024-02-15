package main

import (
	"time"

	"github.com/kbiits/scrap_pemilu/pkg/cache"
	kawalpemiluscrapper "github.com/kbiits/scrap_pemilu/services/kawalpemilu_scrapper"
	pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"
	tpswriter "github.com/kbiits/scrap_pemilu/services/tps_writer"
)

func main() {

	csvWriter := tpswriter.NewCSVWriter("./pasar_minggu.csv")

	scrapperService := pemiluscrapper.NewScrapper()
	// urlsChan := scrapperService.GenerateTPSUri()
	stacks := scrapperService.BuildStackAreaByAreaCodes("317404")
	urlsChan := scrapperService.GenerateTPSUriInArea(stacks)
	// scrapperService.StartScrapping(csvWriter, urlsChan)

	kawalpemiluScrapper := kawalpemiluscrapper.NewKawalPemiluScrapper(cache.New(time.Second*30, time.Second*30))
	kawalpemiluScrapper.StartScrapping(csvWriter, urlsChan)
}
