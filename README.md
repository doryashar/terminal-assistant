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
curl -s https://github.com/doryashar/terminal-assistant/releases/download/v1.0.1/ai | bash
```

### Usage example:
```bash
$ ai find all the files in the current directory that contain the word "hello" and are older than 10 days.
```
### Roadmap:
- [X] When installing, prompt the user for installation and config directory.
- [ ] Create config file, with default values and prompt the user for values.
- [x] Add support for AUTOMATIC update.
- [ ] Add support for multiple models.
- [ ] Save conversation history.
- [ ] Use dynamic information (datetime, shell type, terminal history, command history, current directory, etc).
- [ ] Suggest command into the terminal or show Text in Box.
- [ ] Allow to show execution time, token count, price and more.
- [ ] Handle errors when suggesting commands.
- [ ] Add menu for config settings.
- [ ] Add menu for selecting model.
- [ ] Add configuration validator.
- [ ] Add support for streaming the responses.