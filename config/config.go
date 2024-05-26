package config

//go:generate go run github.com/mk6i/smarter-smarter-child/cmd/config_generator windows settings.bat
//go:generate go run github.com/mk6i/smarter-smarter-child/cmd/config_generator unix settings.env
type Config struct {
	LogLevel        string `envconfig:"LOG_LEVEL" required:"true" val:"info" description:"Set logging granularity. Possible values: 'debug', 'info', 'warn', 'error'."`
	MaxMsgPerMin    int    `envconfig:"MAX_MSG_PER_MIN" required:"true" val:"10" description:"Specifies the maximum number of messages a user can send to the bot per minute before rate limiting is applied."`
	OSCARHost       string `envconfig:"OSCAR_HOST" required:"true" val:"127.0.0.1" description:"The OSCAR hostname to connect to."`
	OSCARPort       string `envconfig:"OSCAR_PORT" required:"true" val:"5190" description:"The OSCAR port to connect to."`
	OfflineMode     bool   `envconfig:"OFFLINE_MODE" required:"false" val:"true" description:"Use a static chat bot that serves canned responses instead of OpenAI for testing."`
	OpenAIKey       string `envconfig:"OPEN_AI_KEY" required:"false" val:"" description:"Key required to connect to the OpenAI API."`
	Password        string `envconfig:"PASSWORD" required:"true" val:"" description:"The bot's account password."`
	ScreenName      string `envconfig:"SCREEN_NAME" required:"true" val:"smartersmarterchild" description:"The bot's screen name."`
	WordCountLimit  int    `envconfig:"WORD_COUNT_LIMIT" required:"true" val:"25" description:"The maximum number of words sent to the bot in a single message."`
	WordLengthLimit int    `envconfig:"WORD_LENGTH_LIMIT" required:"true" val:"15" description:"The maximum length of any word sent to the bot in a single message."`
}
