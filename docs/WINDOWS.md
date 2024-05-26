# SmarterSmarterChild Quickstart for Windows 10/11

This guide explains how to download, configure and connect SmarterSmarterChild running on Windows 10/11 to an AIM/OSCAR
service.

1. **Register an Account**

   First create an account for your chatbot on your preferred AIM revival platform—SmarterSmarterChild connects to AIM
   just like any other AIM client. If you don't know where to start, set up [Retro AIM Server](https://github.com/mk6i/retro-aim-server).
   Keep the username and password handy for the configuration step. Note that the screen names `SmarterSmarterChild` and
   `SmarterChild` may already be taken by someone else.

2. **Download SmarterSmarterChild**

   Grab the latest macOS release from the [Releases page](https://github.com/mk6i/smarter-smarter-child/releases) for
   your platform (Intel or Apple Silicon).

   Because the SmarterSmarterChild executable has not been blessed by Microsoft, browsers such as Chrome may think it's
   a "suspicious" file and block the download, in which case you need to explicitly opt in to downloading the untrusted
   file.

    <p align="center">
      <img alt="screenshot of a chrome prompt showing a blocked download" src="https://github.com/mk6i/retro-aim-server/assets/2894330/90af40bd-262d-4e0f-a769-06943c7acd18">
    </p>

   > While the binaries are 100% safe, you can avoid the security concern by [building the application yourself](./BUILD.md).
   We do not provide signed binaries because of the undue cost and complexity.

   Once downloaded, extract the `.zip` archive, which contains the application and a configuration file `settings.bat`.

3. **Configure AIM Client Settings**

   Open `settings.bat` (right-click, `edit in notepad`).

   - Set `OSCAR_HOST` to your AIM provider's hostname. If you're running Retro AIM Server locally, keep the default
     value.
   - Set the `SCREEN_NAME` and `PASSWORD` values for the account you created in the first step.

4. **Configure ChatGPT Backend**

   Finally, configure the Chatbot backend that generates responses to user input. SmarterSmartChild is configured by
   default to use a mock ChatGPT service that serves canned responses.

   Configure the following values in `settings.bat`:

   - If you don't want to use ChatGPT, keep `OFFLINE_MODE=true`.
   - To enable ChatGPT for AI-generated responses, set `OFFLINE_MODE=false` and set `OPEN_AI_KEY` to your OpenAI
     service account key. If you don't have one already, register for an
     [OpenAPI Platform Account](https://platform.openai.com/) (you'll need to spend a few dollars on credits) and
     create a service account.
   
5. **Start the Application**

   Open `run.cmd` to launch SmarterSmarterChild.

   Because SmarterSmarterChild has not been blessed by Microsoft, Windows will flag the application as a security risk
   the first time you run it. You'll be presented with a `Microsoft Defender SmartScreen` warning prompt that gives you
   the option to run the blocked application.

   To proceed, click `More Options`, then `Run anyway`.

    <p align="center">
      <img alt="of screenshot microsoft defender smartscreen prompt" src="https://github.com/mk6i/retro-aim-server/assets/2894330/9ab0966b-d5dd-4b70-ba16-483e5c458f89">
      <img alt="of screenshot microsoft defender smartscreen prompt" src="https://github.com/mk6i/retro-aim-server/assets/2894330/5d4106c6-0ce6-4d4f-9260-e9bbb777c770">
    </p>

   > While the binaries are 100% safe, you can avoid the security concern by [building the application yourself](./BUILD.md).
   We do not provide signed binaries because of the undue cost and complexity.

   SmarterSmarterChild will run in the terminal, ready to chat.