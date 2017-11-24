package raweb

type Service interface {
	Start(config *Config) error
}
