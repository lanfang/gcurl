package config

var (
	G_Conf *Configure = &Configure{}
)

type Configure struct {
	Addr       string
	ProtoFile  string
	Data       string
	SymbolList []string
	Mehtod     string
}
