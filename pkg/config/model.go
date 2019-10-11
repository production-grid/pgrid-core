package config

// CoreConfiguration models the basic configuration of a pgrid application.
type CoreConfiguration struct {
	ApplicationName       string                `yaml:"name"`
	PortNumber            int                   `yaml:"port"`
	DatabaseConfiguration DatabaseConfiguration `yaml:"database"`
}

// DatabaseConfiguration wraps database configuration settings.
type DatabaseConfiguration struct {
	Primary RelationalDatasource `yaml:"primary"`
	Replica RelationalDatasource `yaml:"replica"`
}

// RelationalDatasource describes configuration settings for a relational datasource.
type RelationalDatasource struct {
	Hostname   string `yaml:"hostname"`
	Portnumber int    `yaml:"port"`
	Schema     string `yaml:"schema"`
	Username   string `yaml:"user"`
	Password   string `yaml:"password"`
}
