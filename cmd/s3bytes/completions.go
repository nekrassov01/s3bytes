package main

import (
	_ "embed"
	"fmt"
)

//go:embed completions/s3bytes.bash
var completionBash string

//go:embed completions/s3bytes.zsh
var completionZsh string

//go:embed completions/s3bytes.ps1
var completionPwsh string

type shell int

const (
	bash shell = iota
	zsh
	pwsh
)

func (t shell) String() string {
	switch t {
	case bash:
		return "bash"
	case zsh:
		return "zsh"
	case pwsh:
		return "pwsh"
	default:
		return ""
	}
}

func parseShell(s string) (shell, error) {
	switch s {
	case "bash":
		return bash, nil
	case "zsh":
		return zsh, nil
	case "pwsh":
		return pwsh, nil
	default:
		return 0, fmt.Errorf("unsupported shell: %q", s)
	}
}
