package speedtest

type Host struct {
	IP    string
	Token string
	Port  int
}

type Version struct {
	Main   string `json:"main"`
	WebAPI string `json:"webapi"`
}
