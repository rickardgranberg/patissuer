// ==========================================================================================================
// <copyright>COPYRIGHT © InfoVista Sweden AB</copyright>
//
// The copyright of the computer program herein is the property of InfoVista Sweden AB.
// The program may be used and/or copied only with the written permission from InfoVista Sweden AB
// or in the accordance with the terms and conditions stipulated in the agreement/contract under which
// the program has been supplied.
// ==========================================================================================================
//go:build mage
// +build mage

// The SiteVerification Microservice build targets.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Install the binaries.
func Build() error {
	mg.Deps(ToolInstall)
	version := os.Getenv("BUILDVERSION")
	if version == "" {
		version = "dev"
	}

	buildTime, err := sh.Output("date", "--iso-8601=seconds")
	if err != nil {
		return err
	}

	commit, err := sh.Output("git", "rev-parse", "--short", "HEAD")
	if err != nil {
		return err
	}

	return sh.Run("go",
		"install",
		fmt.Sprintf("-ldflags= -X main.version=%s -X main.commit=%s -X main.buildTime=%s", version, commit, buildTime),
		"./...")
}

// Performs a lint analysis on the Go code
func Lint() error {
	mg.Deps(ToolInstall)
	return sh.RunV("golangci-lint", "run")
}

// Performs a Vulnerability check on the Go code
func Vuln() error {
	mg.Deps(ToolInstall)
	return sh.RunV("govulncheck", "./...")
}

// Run the tests.
func Test() error {
	mg.Deps(ToolInstall)
	return sh.Run("ginkgo", "--timeout", "1m", "--no-color", "--race", "-v", "./...")
}

// Cleans the binaries.
func Clean() error {
	return sh.Run("go", "clean", "./...")
}

// Run the tests in watch mode.
func Watch() error {
	mg.Deps(ToolInstall)
	return sh.Run("ginkgo", "watch", "./...")
}

// Updates the module dependencies.
func Update() error {
	return sh.Run("go", "get", "-u", "./...")
}

// Creates a release using goreleaser
func Release() error {
	return sh.Run("goreleaser", "release", "--snapshot", "--skip-publish", "--skip-sign", "--rm-dist")
}

// Creates a release using goreleaser
func ReleaseCI() error {
	return sh.Run("goreleaser", "release", "--rm-dist")
}

// Install all tool dependencies
func ToolInstall() error {
	tools, err := findTools()
	if err != nil {
		return err
	}

	for _, t := range tools {
		if err := sh.Run("go", "install", t); err != nil {
			return err
		}
	}
	return nil
}

func findTools() ([]string, error) {
	f, err := os.Open("tools.go")
	if err != nil {
		return nil, err
	}

	defer f.Close()

	var tools []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.Trim(scanner.Text(), " \t")
		if strings.HasPrefix(line, "_") {
			tokens := strings.Split(line, " ")
			tools = append(tools, strings.Trim(tokens[1], "\""))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return tools, nil
}
