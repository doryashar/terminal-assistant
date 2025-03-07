# Terminal assistant - written in bash
the assistant can be used to perform various tasks such as helping you with terminal commands, asking questions, writing scripts and more.
Dont use it. this is not meant for anyone. the code in this repo is not safe.

### Features:
- AI-powered assistant.
- Multiple model support.
- Configurable using easy to read json file.
- Easy to use.

### Install & Run:
```bash
curl -L https://github.com/doryashar/terminal-assistant/releases/download/v1.0.12/install | bash
# OR grep the latest: curl -L $(curl -s "https://api.github.com/repos/doryashar/terminal-assistant/releases/latest" | ./jq-linux64 -r '.assets[0].browser_download_url') | bash
```

### Usage example:
```bash
$ ai find all the files in the current directory that contain the word "hello" and are older than 10 days.
```
### Roadmap:
- [X] When installing, prompt the user for installation and config directory.
- [x] Create config file, with default values and prompt the user for values.
- [x] Add support for AUTOMATIC update.
- [ ] Add functions to the assistant. specifically, add a terminal command and put it inline for user to use (input buffer).
- [ ] Add support for multiple models.
- [x] Save conversation history.
- [ ] Use dynamic information in chat (datetime, shell type, terminal history, command history, current directory, etc).
- [ ] Suggest command into the terminal or show Text in Box.
- [ ] Allow to show execution time, token count, price and more.
- [ ] Handle errors when suggesting commands.
- [ ] Add menu for config settings.
- [ ] Add menu for selecting model.
- [x] Add support for streaming the responses.
- [ ] Add ability to log errors.
