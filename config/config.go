package config

//go:generate go run github.com/mk6i/smarter-smarter-child/cmd/config_generator windows settings.bat
//go:generate go run github.com/mk6i/smarter-smarter-child/cmd/config_generator unix settings.env
type Config struct {
	LogLevel        string  `envconfig:"LOG_LEVEL" required:"true" val:"info" description:"Set logging granularity. Possible values: 'debug', 'info', 'warn', 'error'."`
	MaxMsgPerMin    int     `envconfig:"MAX_MSG_PER_MIN" required:"true" val:"10" description:"Specifies the maximum number of messages a user can send to the bot per minute before rate limiting is applied."`
	OSCARHost       string  `envconfig:"OSCAR_HOST" required:"true" val:"127.0.0.1" description:"The OSCAR hostname to connect to."`
	OSCARPort       string  `envconfig:"OSCAR_PORT" required:"true" val:"5190" description:"The OSCAR port to connect to."`
	OfflineMode     bool    `envconfig:"OFFLINE_MODE" required:"false" val:"true" description:"Use a static chat bot that serves canned responses instead of OpenAI for testing."`
	OpenAIKey       string  `envconfig:"OPEN_AI_KEY" required:"false" val:"" description:"Key required to connect to the OpenAI API."`
	Password        string  `envconfig:"PASSWORD" required:"true" val:"" description:"The bot's account password."`
	ScreenName      string  `envconfig:"SCREEN_NAME" required:"true" val:"smartersmarterchild" description:"The bot's screen name."`
	WordCountLimit  int     `envconfig:"WORD_COUNT_LIMIT" required:"true" val:"25" description:"The maximum number of words sent to the bot in a single message."`
	WordLengthLimit int     `envconfig:"WORD_LENGTH_LIMIT" required:"true" val:"15" description:"The maximum length of any word sent to the bot in a single message."`
	ProfileHTML     string  `envconfig:"PROFILE_HTML" required:"true" val:"'<HTML><BODY BGCOLOR=\"#CDFFFE\"><FONT FACE=\"Courier New\" COLOR=\"#000080\" LANG=\"0\">Hello, %n!<BR>Send me an IM to get started!</FONT><BR><BR><HR><FONT SIZE=1>Powered by <A HREF=\"https://github.com/mk6i/smarter-smarter-child\">SmarterSmarterChild</A>.</FONT></BODY></HTML>'" description:"The bot's HTML profile information."`
	MsgFormat       string  `envconfig:"MSG_FORMAT" required:"true" val:"'<HTML><BODY BGCOLOR=\"#CDFFFE\"><FONT FACE=\"Courier New\" COLOR=\"#000080\" LANG=\"0\">@MsgContent@</FONT></BODY></HTML>'" description:"The bot's message response. @MsgContent@ will be replaced with the content of the bot's response."`
	TopP            float64 `envconfig:"TOP_P" required:"true" val:"0.5" description:"The top-p value to use when querying the OpenAI API."`
	Temperature     float64 `envconfig:"TEMPERATURE" required:"true" val:"0.7" description:"The temperature value to use when querying the OpenAI API."`
	Model           string  `envconfig:"MODEL" required:"true" val:"'gpt-4o-mini'" description:"The AI model to use."`
	BotPrompt       string  `envconfig:"BOT_PROMPT" required:"true" val:"'You are SmarterChild, a dumb AIM chatbot.'" description:"The initial prompt to the OpenAI API when creating a new conversation."`
	APIUrl          string  `envconfig:"API_URL" required:"true" val:"'https://api.openai.com/v1/chat/completions'" description:"OpenAI API URL."`
}
