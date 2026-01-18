package ratelimit

import "time"

type Config struct {
	Default DefaultConfig `yaml:"__default__"`
	Rules   []MainConfig  `yaml:"rules"`
}

type DefaultConfig struct {
	MainConfig `yaml:",inline"`
}

type MainConfig struct {
	// Clients клиенты, к которым относится данное правило
	//  если список пуст - применяется для любых запросов
	//  если имя клиента не попадает в список - имя клиента трактуется как unknown
	Clients []string `yaml:"clients"`
	// Handlers Ручки, к которым относится правило
	//  если пустой список - применяется для любых запросов
	Handlers []string `yaml:"handlers"`
	// Limit установленный rps,
	//  если не установлен или равен 0 - отклоняем ВСЕ запросы
	Limit float64 `yaml:"limit"`
	// Burst допускаемый прирост запросов в секунду в интервал времени
	// если не установлен - равен лимиту
	Burst int `yaml:"burst"`
	//Timeout время, которое будет ожидать запрос при получении токена
	Timeout time.Duration `yaml:"timeout"`
}

func withDefaults(cfg Config) Config {
	cfg.Default.MainConfig = defaultConfig(cfg.Default.MainConfig)

	return cfg
}

func enrich(cfg, enrichment MainConfig) MainConfig {
	cfg.Handlers = overrideIfNilSlice(cfg.Handlers, enrichment.Handlers)
	cfg.Clients = overrideIfNilSlice(cfg.Clients, enrichment.Clients)
	cfg.Limit = overrideIfZero(cfg.Limit, enrichment.Limit)
	cfg.Timeout = overrideIfZero(cfg.Timeout, enrichment.Timeout)
	cfg.Burst = overrideIfZero(cfg.Burst, enrichment.Burst)
	return cfg
}

func defaultConfig(cfg MainConfig) MainConfig {
	defaultCfg := MainConfig{}
	return enrich(cfg, defaultCfg)
}
