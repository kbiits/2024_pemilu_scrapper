package tpswriter

import (
	"encoding/csv"
	"fmt"
	"os"

	pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"
)

var _ pemiluscrapper.TPSResultWriter = (*CSVWriter)(nil)
var (
	cols        = []string{"Area code", "URL", "Suara Total", "Suara Sah", "Suara Tidak Sah", "Paslon 01", "Paslon 02", "Paslon 03"}
	paslonIdMap = map[int]string{
		1: "100025",
		2: "100026",
		3: "100027",
	}
)

func NewCSVWriter(pathToCsv string) pemiluscrapper.TPSResultWriter {
	file, err := os.OpenFile(pathToCsv, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		panic(err)
	}

	csvWriter := csv.NewWriter(file)

	err = csvWriter.Write(cols)
	if err != nil {
		panic(err)
	}

	return &CSVWriter{
		writer: csvWriter,
	}
}

type CSVWriter struct {
	writer *csv.Writer
}

// Write implements pemiluscrapper.TPSResultWriter.
func (w *CSVWriter) Write(result pemiluscrapper.TPSResultWithMetadata) error {
	data := []string{result.Code, result.Url}

	if result.TPSResult.Administrasi != nil {
		adm := result.TPSResult.Administrasi
		data = append(data, fmt.Sprintf("%d", adm.SuaraTotal))
		data = append(data, fmt.Sprintf("%d", adm.SuaraSah))
		data = append(data, fmt.Sprintf("%d", adm.SuaraTidakSah))
	} else {
		data = append(data, "null", "null", "null")
	}

	if result.TPSResult.Chart != nil {
		chart := result.TPSResult.Chart
		var (
			paslon1 string
			paslon2 string
			paslon3 string
		)

		if chart.Num100025 == nil {
			paslon1 = "null"
		} else {
			paslon1 = fmt.Sprintf("%d", *chart.Num100025)
		}

		if chart.Num100026 == nil {
			paslon2 = "null"
		} else {
			paslon2 = fmt.Sprintf("%d", *chart.Num100026)
		}

		if chart.Num100027 == nil {
			paslon3 = "null"
		} else {
			paslon3 = fmt.Sprintf("%d", *chart.Num100027)
		}

		data = append(data, paslon1)
		data = append(data, paslon2)
		data = append(data, paslon3)
	} else {
		data = append(data, "null", "null", "null")
	}

	return w.writer.Write(data)
}

// WriteError implements pemiluscrapper.TPSResultWriter.
func (w *CSVWriter) WriteError(result pemiluscrapper.TPSResultWithMetadata) {
	data := []string{result.Url}
	w.writer.Write(data)
}

func (w *CSVWriter) Close() {
	w.writer.Flush()
}
