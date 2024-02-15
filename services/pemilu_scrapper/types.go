package pemiluscrapper

type Area struct {
	Name  string `json:"nama"`
	Id    int    `json:"id"`
	Code  string `json:"kode"`
	Level int    `json:"tingkat"`
}

type AreaWithUrl struct {
	Area
	UrlJson string
}

type TPSResultWithMetadata struct {
	TPSResult
	Url  string
	Code string
}

type TPSResult struct {
	Chart        *Chart        `json:"chart"`
	Images       []string      `json:"images"`
	Administrasi *Administrasi `json:"administrasi"`
	Psu          any           `json:"psu"`
	Ts           string        `json:"ts"`
	StatusSuara  bool          `json:"status_suara"`
	StatusAdm    bool          `json:"status_adm"`
}
type Chart struct {
	Num100025 int `json:"100025"`
	Num100026 int `json:"100026"`
	Num100027 int `json:"100027"`
}
type Administrasi struct {
	SuaraSah        int `json:"suara_sah"`
	SuaraTotal      int `json:"suara_total"`
	PemilihDptJ     int `json:"pemilih_dpt_j"`
	PemilihDptL     int `json:"pemilih_dpt_l"`
	PemilihDptP     int `json:"pemilih_dpt_p"`
	PenggunaDptJ    int `json:"pengguna_dpt_j"`
	PenggunaDptL    int `json:"pengguna_dpt_l"`
	PenggunaDptP    int `json:"pengguna_dpt_p"`
	PenggunaDptbJ   int `json:"pengguna_dptb_j"`
	PenggunaDptbL   int `json:"pengguna_dptb_l"`
	PenggunaDptbP   int `json:"pengguna_dptb_p"`
	SuaraTidakSah   int `json:"suara_tidak_sah"`
	PenggunaTotalJ  int `json:"pengguna_total_j"`
	PenggunaTotalL  int `json:"pengguna_total_l"`
	PenggunaTotalP  int `json:"pengguna_total_p"`
	PenggunaNonDptJ int `json:"pengguna_non_dpt_j"`
	PenggunaNonDptL int `json:"pengguna_non_dpt_l"`
	PenggunaNonDptP int `json:"pengguna_non_dpt_p"`
}
