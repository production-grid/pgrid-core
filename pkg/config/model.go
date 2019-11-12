package config

// CoreConfiguration models the basic configuration of a pgrid application.
type CoreConfiguration struct {
	ApplicationName       string                `yaml:"name"`
	RootTenantHost        string                `yaml:"root_tenant_host"`
	SecureCookies         bool                  `yaml:"secure_cookies"`
	PortNumber            int                   `yaml:"port"`
	DatabaseConfiguration DatabaseConfiguration `yaml:"database"`
	Notifications         NotificationConfig    `yaml:"notifications"`
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

// NotificationConfig encapsulates transport specific notification settings.
type NotificationConfig struct {
	SendGrid SendGridConfig `yaml:"sendgrid"`
	Twilio   TwilioConfig   `yaml:"twilio"`
}

// SendGridConfig models sendgrid configuration
type SendGridConfig struct {
	APIKey      string `yaml:"api_key"`
	SenderEMail string `yaml:"sender_email"`
	SenderName  string `yaml:"sender_name"`
}

// TwilioConfig models Twilio configuration
type TwilioConfig struct {
	Number    string `yaml:"number"`
	SID       string `yaml:"sid"`
	AuthToken string `yaml:"auth_token"`
}
