# go CLIweather

# 1. Introduction
This is a Golang project where you can use a CLI to gather information from WeatherAPI about the forecast of your current location.

Let's start explaining the commands that the CLI have:
- **root**: By default, there are two options for users to configure various visual aspects of the forecast:
    - **no-color**: User can choose between colorful or plain information. By default, info will be with colors. In case you want switch off colors:
        ```
        cliweather --no-color forecast
        ```
    - **no-emoji**: Same case as before but this time you can decide if you want emojis or not. By default, emojis are on, but if you want to disable them:
        ```
        cliweather --no-emoji forecast
        ```
- **completion**: This command is designed to generate shell autocompletion scripts for various shells: **Bash**, **Zsh** and **Fish**. This improves the user experience by enabling tab-completion for commans, flags and args. Example of use:
    ```
    cliweather completion bash
    ```
    There si no need to add double dashes, just put directly the terminal you want.

- **version**: This command shows the user the current version of the executed binary:
    ```
    cliweather version
    ```
- **forecast**: 