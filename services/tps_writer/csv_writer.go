package tpswriter

import (
	"encoding/csv"
	"fmt"
	"os"

	pemiluscrapper "github.com/kbiits/scrap_pemilu/services/pemilu_scrapper"
)

var _ pemiluscrapper.TPSResultWriter = (*CSVWriter)(nil)
var (
	cols        = []string{"URL", "Suara Total", "Suara Sah", "Suara Tidak Sah", "Paslon 01", "Paslon 02", "Paslon 03"}
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
	data := []string{result.Url}

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
		data = append(data, fmt.Sprintf("%d", chart.Num100025))
		data = append(data, fmt.Sprintf("%d", chart.Num100026))
		data = append(data, fmt.Sprintf("%d", chart.Num100027))
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
