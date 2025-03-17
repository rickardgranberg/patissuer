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
	"fmt"
	"os"

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
	mg.Deps(ToolInstall)
	return sh.Run("goreleaser", "release", "--snapshot", "--skip=publish", "--skip=sign", "--clean")
}

// Creates a release using goreleaser
func ReleaseCI() error {
	mg.Deps(ToolInstall)
	return sh.Run("goreleaser", "release", "--clean")
}

// Install all tool dependencies
func ToolInstall() error {

	return sh.Run("go", "install", "tool")
}
