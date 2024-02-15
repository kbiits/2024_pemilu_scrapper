package main

import (
	pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"
	tpswriter "github.com/kbiits/scrap_pemilu/services/tps_writer"
)

func main() {

	scrapperService := pemiluscrapper.NewScrapper()
	csvWriter := tpswriter.NewCSVWriter("./mampang.csv")

	// urlsChan := scrapperService.GenerateTPSUri()
	stacks := scrapperService.BuildStackAreaByAreaCodes("3174031001")
	urlsChan := scrapperService.GenerateTPSUriInArea(stacks)

	scrapperService.StartScrapping(csvWriter, urlsChan)
}
