# SmarterSmarterChild Quickstart for Linux (x86_64)

This guide explains how to download, configure and connect SmarterSmarterChild running on Linux (x86_64) to an AIM/OSCAR
service.

1. **Register an Account**

   First create an account for your chatbot on your preferred AIM revival platform—SmarterSmarterChild connects to AIM
   just like any other AIM client. If you don't know where to start, set up [Retro AIM Server](https://github.com/mk6i/retro-aim-server).
   Keep the username and password handy for the configuration step. Note that the screen names `SmarterSmarterChild` and
   `SmarterChild` may already be taken by someone else.

2. **Download SmarterSmarterChild**

   Grab the latest Linux release from the [Releases page](https://github.com/mk6i/smarter-smarter-child/releases) and 
   extract the archive. The extracted folder contains the application and a configuration file `settings.env`.

3. **Configure AIM Client Settings**

   Open `settings.env` in your favorite text editor.

   - Set `OSCAR_HOST` to your AIM provider's hostname. If you're running Retro AIM Server locally, keep the default
     value.
   - Set the `SCREEN_NAME` and `PASSWORD` values for the account you created in the first step.

4. **Configure ChatGPT Backend**

   Finally, configure the Chatbot backend that generates responses to user input. SmarterSmartChild is configured by
   default to use a mock ChatGPT service that serves canned responses.

   Configure the following values in `settings.env`:

   - If you don't want to use ChatGPT, keep `OFFLINE_MODE=true`.
   - To enable ChatGPT for AI-generated responses, set `OFFLINE_MODE=false` and set `OPEN_AI_KEY` to your OpenAI
     service account key. If you don't have one already, register for an
     [OpenAPI Platform Account](https://platform.openai.com/) (you'll need to spend a few dollars on credits) and
     create a service account.

5. **Start the Application**

   Run the following command to launch SmarterSmarterChild:

   ```shell
   ./run.sh
   ```

   SmarterSmarterChild will run in the terminal, ready to chat.