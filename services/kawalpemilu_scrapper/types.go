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
	return pemiluscrapper.TPSResultWithMetadata{
		Url:  url,
		Code: code,
		TPSResult: pemiluscrapper.TPSResult{
			Chart: &pemiluscrapper.Chart{
				Num100025: data.Pas1,
				Num100026: data.Pas2,
				Num100027: data.Pas3,
			},
		},
	}
}
