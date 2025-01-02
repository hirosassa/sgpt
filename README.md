# sgpt

Go clone of [shell_gpt](https://github.com/TheR1D/shell_gpt) command.

A command-line productivity tool powered by AI large language models (LLM). This command-line tool offers streamlined generation of shell commands, code snippets, documentation, eliminating the need for external resources (like Google search).

> [!TIP]
> This project is still in development. The current version is a prototype.
> Many features implemented in the original shell_gpt are not yet available in this project.

## Usage (borrowed from [shell_gpt](https://github.com/TheR1D/shell_gpt) )

This tool is designed to quickly analyse and retrieve information. It's useful for straightforward requests ranging from technical configurations to general knowledge.
```shell
sgpt "What is the fibonacci sequence"
# -> The Fibonacci sequence is a series of numbers where each number ...
```

ShellGPT accepts prompt from both stdin and command line argument. Whether you prefer piping input through the terminal or specifying it directly as arguments, `sgpt` got you covered. For example, you can easily generate a git commit message based on a diff:
```shell
git diff | sgpt "Generate git commit message, for my changes"
# -> Added main feature details into README.md
```

You can analyze logs from various sources by passing them using stdin, along with a prompt. For instance, we can use it to quickly analyze logs, identify errors and get suggestions for possible solutions:
```shell
docker logs -n 20 my_app | sgpt "check logs, find errors, provide possible solutions"
```
```text
Error Detected: Connection timeout at line 7.
Possible Solution: Check network connectivity and firewall settings.
Error Detected: Memory allocation failed at line 12.
Possible Solution: Consider increasing memory allocation or optimizing application memory usage.
```

You can also use all kind of redirection operators to pass input:
```shell
sgpt "summarise" < document.txt
# -> The document discusses the impact...
sgpt << EOF
What is the best way to lear Golang?
Provide simple hello world example.
EOF
# -> The best way to learn Golang...
sgpt <<< "What is the best way to learn shell redirects?"
# -> The best way to learn shell redirects is through...
```

for more details, see [shell_gpt](https://github.com/TheR1D/shell_gpt).

## Installation

### build from source

```shell
go install github.com/hirosassa/sgpt@latest
```

### MacOS

```shell
brew install hirosassa/tap/sgpt
```
