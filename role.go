package main

import (
	"bytes"
	"os"
	"runtime"
	"text/template"
)

const (
	SHELL_ROLE = "Provide only {{ .Shell }} commands for {{ .OS }} without any description.\n" +
		"If there is a lack of details, provide most logical solution.\n" +
		"Ensure the output is a valid shell command.\n" +
		"If multiple steps required try to combine them together using &&.\n" +
		"Provide only plain text without Markdown formatting.\n" +
		"Do not provide markdown formatting such as ```."

	DESCRIBE_SHELL_ROLE = `Provide a terse, single sentence description of the given shell command.
Describe each argument and option of the command.
Provide short responses in about 80 words.
APPLY MARKDOWN formatting when possible.`

	CODE_ROLE = "Provide only code as output without any description.\n" +
		"Provide only code in plain text format without Markdown formatting.\n" +
		"Do not include symbols such as ``` or ```python.\n" +
		"If there is a lack of details, provide most logical solution.\n" +
		"You are not allowed to ask for more details.\n" +
		"For example if the prompt is \"Hello world Python\", you should return \"print('Hello world')\"."

	DEFAULT_ROLE = `You are programming and system administration assistant.
You are managing {{ .OS }} operating system with {{ .Shell }} shell.
Provide short responses in about 100 words, unless you are specifically asked for more details.
If you need to store any data, assume it will be stored in the conversation.
APPLY MARKDOWN formatting when possible.`
)

const ROLE_TEMPLATE = "You are {{ .Name }}\n{{ .Role }}"

type DefaultRole string

const (
	DEFAULT        DefaultRole = "ShellGPT"
	SHELL          DefaultRole = "Shell Command Generator"
	DESCRIBE_SHELL DefaultRole = "Shell Command Descriptor"
	CODE           DefaultRole = "Code Generator"
)

type SystemRole struct {
	name string
	role string
}

func NewRole(name string, role string, variables map[string]string) (*SystemRole, error) {

	var b bytes.Buffer
	if checkVariables(variables) {
		tpl := template.Must(template.New("role").Parse(role))
		if err := tpl.Execute(&b, variables); err != nil {
			return nil, err
		}
		return &SystemRole{name, b.String()}, nil
	}
	return &SystemRole{name, role}, nil
}

func osName() string {
	// todo: add distro name if needed
	return runtime.GOOS
}

func shellName() string {
	// todo: support windows shell
	return os.Getenv("SHELL")
}

func checkVariables(variables map[string]string) bool {
	_, ok := variables["shell"]
	if !ok {
		return false
	}
	_, ok = variables["os"]
	if !ok {
		return false
	}
	return true
}

func CheckGet(shell bool, describeShell bool, code bool) (*SystemRole, error) {
	if shell {
		variables := map[string]string{"os": osName(), "shell": shellName()}
		role, err := NewRole(string(SHELL), SHELL_ROLE, variables)
		if err != nil {
			return nil, err
		}
		return role, nil
	}
	if describeShell {
		role, err := NewRole(string(DESCRIBE_SHELL), DESCRIBE_SHELL_ROLE, map[string]string{})
		if err != nil {
			return nil, err
		}
		return role, nil
	}
	if code {
		role, err := NewRole(string(CODE), CODE_ROLE, map[string]string{})
		if err != nil {
			return nil, err
		}
		return role, nil
	}

	variables := map[string]string{"os": osName(), "shell": shellName()}
	role, err := NewRole(string(DEFAULT), DEFAULT_ROLE, variables)
	if err != nil {
		return nil, err
	}
	return role, nil
}