package ratelimit

import "time"

type Config struct {
	Rules []MainConfig `yaml:"rules"`
}

type DefaultConfig struct {
	MainConfig `yaml:",inline"`
}

type MainConfig struct {
	// Clients клиенты, к которым относится данное правило
	//  если список пуст - применяется для любых запросов
	//  если имя клиента не попадает в список - имя клиента трактуется как unknown.
	//  Применимо только для серверного rate limiter'а
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
