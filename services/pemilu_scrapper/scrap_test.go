package pemiluscrapper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_buildUrl(t *testing.T) {
	type args struct {
		baseUrl string
		parents []Area
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "get province by code",
			args: args{
				parents: []Area{
					{Code: "10"},
				},
				baseUrl: baseUrlArea,
			},
			want: baseUrlArea + "10.json",
		},
		{
			name: "get city by code",
			args: args{
				parents: []Area{
					{Code: "10"},
					{Code: "1021"},
				},
				baseUrl: baseUrlArea,
			},
			want: baseUrlArea + "10/1021.json",
		},
		{
			name: "get district by code",
			args: args{
				parents: []Area{
					{Code: "10"},
					{Code: "1021"},
					{Code: "102130"},
				},
				baseUrl: baseUrlArea,
			},
			want: baseUrlArea + "10/1021/102130.json",
		},
		{
			name: "get tps result url",
			args: args{
				parents: []Area{
					{Code: "10"},
					{Code: "1021"},
					{Code: "102130"},
					{Code: "10213012"},
					{Code: "1021301232"},
				},
				baseUrl: baseUrlTPSResult,
			},
			want: baseUrlTPSResult + "10/1021/102130/10213012/1021301232.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildUrl(tt.args.baseUrl, tt.args.parents); got != tt.want {
				t.Errorf("buildUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestScrapperSvc_BuildStackAreaByAreaCodes(t *testing.T) {
	type args struct {
		code string
	}
	tests := []struct {
		name          string
		args          args
		wantStackArea []Area
	}{
		{
			name: "stack area province",
			args: args{
				code: "31",
			},
			wantStackArea: []Area{{Level: 1, Code: "31"}},
		},
		{
			name: "stack area city",
			args: args{
				code: "3174",
			},
			wantStackArea: []Area{{Level: 2, Code: "3174"}, {Level: 1, Code: "31"}},
		},
		{
			name: "stack area district",
			args: args{
				code: "317402",
			},
			wantStackArea: []Area{{Level: 3, Code: "317402"}, {Level: 2, Code: "3174"}, {Level: 1, Code: "31"}},
		},
		{
			name: "stack area subdistrict",
			args: args{
				code: "3174031001",
			},
			wantStackArea: []Area{{Level: 4, Code: "3174031001"}, {Level: 3, Code: "317403"}, {Level: 2, Code: "3174"}, {Level: 1, Code: "31"}},
		},
		{
			name: "stack area tps",
			args: args{
				code: "3174031001065",
			},
			wantStackArea: []Area{{Level: 5, Code: "3174031001065"}, {Level: 4, Code: "3174031001"}, {Level: 3, Code: "317403"}, {Level: 2, Code: "3174"}, {Level: 1, Code: "31"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &ScrapperSvc{}
			gotStackArea := s.BuildStackAreaByAreaCodes(tt.args.code)
			require.EqualValues(t, tt.wantStackArea, gotStackArea.ToSlice())
		})
	}
}

func Test_extractTPSAreaCodeFromUrl(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "success with tps url",
			args: args{
				url: "https://sirekap-obj-data.kpu.go.id/pemilu/hhcw/ppwp/31/3174/317403/3174031001/3174031001050.json",
			},
			want: "3174031001050",
		},
		{
			name: "success with district url",
			args: args{
				url: "https://sirekap-obj-data.kpu.go.id/wilayah/pemilu/ppwp/31/3174/317402.json",
			},
			want: "317402",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractTPSAreaCodeFromUrl(tt.args.url); got != tt.want {
				t.Errorf("extractTPSAreaCodeFromUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
