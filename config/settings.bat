rem Set logging granularity. Possible values: 'debug', 'info', 'warn', 'error'.
set LOG_LEVEL=info

rem Specifies the maximum number of messages a user can send to the bot per
rem minute before rate limiting is applied.
set MAX_MSG_PER_MIN=10

rem The OSCAR hostname to connect to.
set OSCAR_HOST=127.0.0.1

rem The OSCAR port to connect to.
set OSCAR_PORT=5190

rem Use a static chat bot that serves canned responses instead of OpenAI for
rem testing.
set OFFLINE_MODE=true

rem Key required to connect to the OpenAI API.
set OPEN_AI_KEY=

rem The bot's account password.
set PASSWORD=

rem The bot's screen name.
set SCREEN_NAME=smartersmarterchild

rem The maximum number of words sent to the bot in a single message.
set WORD_COUNT_LIMIT=25

rem The maximum length of any word sent to the bot in a single message.
set WORD_LENGTH_LIMIT=15

rem The bot's HTML profile information.
set PROFILE_HTML='<HTML><BODY BGCOLOR="#CDFFFE"><FONT FACE="Courier New" COLOR="#000080" LANG="0">Hello, %n!<BR>Send me an IM to get started!</FONT><BR><BR><HR><FONT SIZE=1>Powered by <A HREF="https://github.com/mk6i/smarter-smarter-child">SmarterSmarterChild</A>.</FONT></BODY></HTML>'

rem The bot's message response. @MsgContent@ will be replaced with the content
rem of the bot's response.
set MSG_FORMAT='<HTML><BODY BGCOLOR="#CDFFFE"><FONT FACE="Courier New" COLOR="#000080" LANG="0">@MsgContent@</FONT></BODY></HTML>'

