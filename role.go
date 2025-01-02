package main

import (
	"bytes"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"text/template"
)

const (
	ShellRole = "Provide only {{ .Shell }} commands for {{ .OS }} without any description.\n" +
		"If there is a lack of details, provide most logical solution.\n" +
		"Ensure the output is a valid shell command.\n" +
		"If multiple steps required try to combine them together using &&.\n" +
		"Provide only plain text without Markdown formatting.\n" +
		"Do not provide markdown formatting such as ```."

	DescribeShellRole = `Provide a terse, single sentence description of the given shell command.
Describe each argument and option of the command.
Provide short responses in about 80 words.
APPLY MARKDOWN formatting when possible.`

	CodeRole = "Provide only code as output without any description.\n" +
		"Provide only code in plain text format without Markdown formatting.\n" +
		"Do not include symbols such as ``` or ```python.\n" +
		"If there is a lack of details, provide most logical solution.\n" +
		"You are not allowed to ask for more details.\n" +
		"For example if the prompt is \"Hello world Python\", you should return \"print('Hello world')\"."

	DefaultRole = `You are programming and system administration assistant.
You are managing {{ .OS }} operating system with {{ .Shell }} shell.
Provide short responses in about 100 words, unless you are specifically asked for more details.
If you need to store any data, assume it will be stored in the conversation.
APPLY MARKDOWN formatting when possible.`
)

const RoleTemplate = "You are {{ .Name }}\n{{ .Role }}"

type DefaultRoleName string

const (
	Default       DefaultRoleName = "ShellGPT"
	Shell         DefaultRoleName = "Shell Command Generator"
	DescribeShell DefaultRoleName = "Shell Command Descriptor"
	Code          DefaultRoleName = "Code Generator"
)

type SystemRole struct {
	name string
	role string
}

func NewRole(name string, role string, variables map[string]string) (*SystemRole, error) {
	roleString, err := execRole(role, variables)
	if err != nil {
		return nil, err
	}

	data := map[string]string{
		"Name": name,
		"Role": roleString,
	}

	var b bytes.Buffer
	tpl, err := template.New("tpl").Parse(RoleTemplate)
	if err != nil {
		return nil, err
	}
	if err := tpl.Execute(&b, data); err != nil {
		return nil, err
	}
	return &SystemRole{name, b.String()}, nil
}

func execRole(role string, variables map[string]string) (string, error) {
	if !checkVariables(variables) {
		return role, nil
	}

	var b bytes.Buffer
	tpl, err := template.New("role").Parse(role)
	if err != nil {
		return "", err
	}
	if err := tpl.Execute(&b, variables); err != nil {
		return "", err
	}
	return b.String(), nil
}

func osName() string {
	// todo: add distro name if needed
	slog.Debug(runtime.GOOS)
	return runtime.GOOS
}

func shellName() string {
	// todo: support windows shell
	slog.Debug(filepath.Base(os.Getenv("SHELL")))
	return filepath.Base(os.Getenv("SHELL"))
}

func checkVariables(variables map[string]string) bool {
	_, ok := variables["Shell"]
	if !ok {
		return false
	}
	_, ok = variables["OS"]
	return ok
}

func CheckGet(shell bool, describeShell bool, code bool) (*SystemRole, error) {
	if shell {
		variables := map[string]string{"OS": osName(), "Shell": shellName()}
		role, err := NewRole(string(Shell), ShellRole, variables)
		if err != nil {
			return nil, err
		}
		return role, nil
	}
	if describeShell {
		role, err := NewRole(string(DescribeShell), DescribeShellRole, map[string]string{})
		if err != nil {
			return nil, err
		}
		return role, nil
	}
	if code {
		role, err := NewRole(string(Code), CodeRole, map[string]string{})
		if err != nil {
			return nil, err
		}
		return role, nil
	}

	variables := map[string]string{"OS": osName(), "Shell": shellName()}
	role, err := NewRole(string(Default), DefaultRole, variables)
	if err != nil {
		return nil, err
	}
	return role, nil
}
