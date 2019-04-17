package config

type WebServerConfig struct {
	_prefix    string `prefix:"webServer"`
	Host 	   string `val:"host"`
	Port       int    `val:"port"`
	CertFile   string `val:"certFile"`
	KeyFile    string `val:"keyFile"`
}
