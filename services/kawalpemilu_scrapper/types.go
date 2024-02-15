package kawalpemiluscrapper

import pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"

type KawalPemiluTPSData struct {
	Result Result `json:"result"`
}

type Result struct {
	Aggregated map[int][]TPSAggregatedData `json:"aggregated"`
}

type TPSAggregatedData struct {
	TotalLaporTps     int    `json:"totalLaporTps"`
	IDLokasi          string `json:"idLokasi"`
	Pas2              int    `json:"pas2"`
	TotalTps          int    `json:"totalTps"`
	Pas3              int    `json:"pas3"`
	TotalCompletedTps int    `json:"totalCompletedTps"`
	Dpt               int    `json:"dpt"`
	TotalJagaTps      int    `json:"totalJagaTps"`
	TotalPendingTps   int    `json:"totalPendingTps"`
	Name              string `json:"name"`
	TotalErrorTps     int    `json:"totalErrorTps"`
	Pas1              int    `json:"pas1"`
	UpdateTs          int64  `json:"updateTs"`
}

func (data *TPSAggregatedData) ConvertToCompatibleKPUType(url string, code string) pemiluscrapper.TPSResultWithMetadata {
	var pas1Val *int
	var pas2Val *int
	var pas3Val *int

	if data.TotalCompletedTps > 0 {
		pas1Val = &data.Pas1
		pas2Val = &data.Pas2
		pas3Val = &data.Pas3
	}

	return pemiluscrapper.TPSResultWithMetadata{
		Url:  url,
		Code: code,
		TPSResult: pemiluscrapper.TPSResult{
			Chart: &pemiluscrapper.Chart{
				Num100025: pas1Val,
				Num100026: pas2Val,
				Num100027: pas3Val,
			},
		},
	}
}
