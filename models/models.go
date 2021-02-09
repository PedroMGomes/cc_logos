package models

// Currency ID is the currency symbol.
type Currency struct {
	// When calling gorm.Model inside a struct it adds the properties:
	// ID        uint `gorm:"primary_key"`
	// CreatedAt time.Time
	// UpdatedAt time.Time
	// DeletedAt *time.Time
	// gorm.Model
	ID      string `json:"ID" binding:"required" gorm:"primary_key"`
	Name    string `json:"name" binding:"required"`
	LogoURL string `json:"logo_url" binding:"required"`
}

// AppConfig - Access is the passphrase to establish connections.
type AppConfig struct {
	// mapstructure used by Viper to read from config files.
	Access      string `mapstructure:"ACCESS"`
	GinMode     string `mapstructure:"GIN_MODE"`
	Addr        string `mapstructure:"ADDR"`
	Port        string `mapstructure:"PORT"`
	DatabaseURL string `mapstructure:"DATABASE_URL"`
}
