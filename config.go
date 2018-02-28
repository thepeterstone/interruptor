package interruptor

// Config holds app configuration
type Config struct {
	APIKey            string   `env:"ACCESS_TOKEN,required"`
	VerificationToken string   `env:"VERIFICATION_TOKEN,required"`
	MessagePrefix     string   `env:"MESSAGE_PREFIX"`
	Channels          []string `env:"CHANNELS,required"`
}
