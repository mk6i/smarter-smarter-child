# SmarterSmarterChild Quickstart for macOS (Intel and Apple Silicon)

This guide explains how to download, configure and connect SmarterSmarterChild running on macOS (Intel and Apple
Silicon) to an AIM/OSCAR service.

1. **Register an Account**

   First create an account for your chatbot on your preferred AIM revival platform—SmarterSmarterChild connects to AIM
   just like any other AIM client. If you don't know where to start, set up [Retro AIM Server](https://github.com/mk6i/retro-aim-server).
   Keep the username and password handy for the configuration step. Note that the screen names `SmarterSmarterChild` and
   `SmarterChild` may already be taken by someone else.

2. **Download SmarterSmarterChild**

   Grab the latest macOS release from the [Releases page](https://github.com/mk6i/smarter-smarter-child/releases) for
   your platform (Intel or Apple Silicon).

   Because the SmarterSmarterChild executable has not been blessed by Apple, browsers such as Chrome may think it's a
   "suspicious" file and block the download, in which case you need to explicitly opt in to downloading the untrusted
   file.

    <p align="center">
      <img alt="screenshot of a chrome prompt showing a blocked download" src="https://github.com/mk6i/retro-aim-server/assets/2894330/90af40bd-262d-4e0f-a769-06943c7acd18">
    </p>

   > While the binaries are 100% safe, you can avoid the security concern by [building the application yourself](./BUILD.md).
   We do not provide signed binaries because of the undue cost and complexity.

   Once downloaded, extract the `.zip` archive, which contains the application and a configuration file `settings.env`.

3. **Open Terminal**

   Open a terminal and navigate to the extracted directory. This terminal will be used for the remaining steps.

   ```shell
   cd ~/Downloads/smarter_smarter_child.0.1.0.macos.intel_x86_64/
   ```

4. **Remove Quarantine**

   Because the SmarterSmarterChild `.app` has not been blessed by Apple, macOS will quarantine the application. To
   proceed, remove the quarantine flag from the `.app`. In the same terminal, run following command:

   ```shell
   sudo xattr -d com.apple.quarantine ./bin/smarter_smarter_child
   ```

   > While the binaries are 100% safe, you can avoid the security concern
   by [building the application yourself](./BUILD.md). We do not provide signed binaries because of the undue cost and
   complexity.

5. **Configure AIM Client Settings**

   Open `settings.env` in your favorite text editor.

    - Set `OSCAR_HOST` to your AIM provider's hostname. If you're running Retro AIM Server locally, keep the default
      value.
    - Set the `SCREEN_NAME` and `PASSWORD` values for the account you created in the first step.

6. **Configure ChatGPT Backend**

   Finally, configure the Chatbot backend that generates responses to user input. SmarterSmartChild is configured by
   default to use a mock ChatGPT service that serves canned responses.

   Configure the following values in `settings.env`:

    - If you don't want to use ChatGPT, keep `OFFLINE_MODE=true`.
    - To enable ChatGPT for AI-generated responses, set `OFFLINE_MODE=false` and set `OPEN_AI_KEY` to your OpenAI
      service account key. If you don't have one already, register for an
      [OpenAPI Platform Account](https://platform.openai.com/) (you'll need to spend a few dollars on credits) and
      create a service account.

7. **Start the Application**

   Run the following command to launch SmarterSmarterChild:

   ```shell
   ./run.sh
   ```

   SmarterSmarterChild will run in the terminal, ready to chat.