package main

import (
	pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"
	tpswriter "github.com/kbiits/scrap_pemilu/services/tps_writer"
)

func main() {

	csvWriter := tpswriter.NewCSVWriter("./pasar_minggu.csv")

	scrapperService := pemiluscrapper.NewScrapper()
	// urlsChan := scrapperService.GenerateTPSUri()
	stacks := scrapperService.BuildStackAreaByAreaCodes("317404")
	urlsChan := scrapperService.GenerateTPSUriInArea(stacks)
	scrapperService.StartScrapping(csvWriter, urlsChan)
}
