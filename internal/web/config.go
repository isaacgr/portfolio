package web

type Configuration struct {
	Host string
	Port int
}

func NewConfiguration() (*Configuration, error) {
	cfg := Configuration{
		Host: "127.0.0.1",
		Port: 2475,
	}
	return &cfg, nil
}
