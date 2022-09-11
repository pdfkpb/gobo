package admin

type Config struct {
	Admins map[string]string
}

func (cfg *Config) IsAdmin(id string) bool {
	return cfg.Admins[id] == ""
}

var ChatPox = Config{
	Admins: map[string]string{
		"<@384902507383619594>": "Kevin",
		"<@303750733700923392>": "Dylan",
	},
}
