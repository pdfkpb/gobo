package admin

type Config struct {
	Admins  map[string]string
	Members []string
}

func (cfg *Config) IsAdmin(id string) bool {
	return cfg.Admins[id] != ""
}

var ChatPox = Config{
	Admins: map[string]string{
		"384902507383619594": "Kevin",
		"303750733700923392": "Dylan",
	},
	Members: []string{
		"303750733700923392", // Dylan
		"390693560141217814", // Seth
		"384902507383619594", // Kevin
		"482844996349984769", // Jacob
		"297889272793399296", // Doomrider
		"566875887119892510", // Faustoh
		"387121478740738050", // Kyory
		"359963174373425152", // Bahb
		"414202347137269773", // Caleb
	},
}
