// Create binary rpm package with ease
package main

import (
	"fmt"
	rpmbuilder "github.com/greenpau/go-rpm-build-lib/pkg/rpmbuilder"
	"github.com/urfave/cli"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	appName        = "go-rpm-builder"
	appVersion     = "[untracked]"
	appDocs        = "https://github.com/greenpau/go-rpm-build-lib/"
	appDescription = "RPM utilities in Go"
	gitBranch      string
	gitCommit      string
	buildUser      string // whoami
	buildDate      string // date -u
)

func main() {
	app := cli.NewApp()
	app.Name = GetAppName()
	app.Version = GetVersion()
	app.Usage = "RPM utilities in Go"
	app.UsageText = fmt.Sprintf("%s <cmd> <options>", GetAppName())
	app.Commands = []*cli.Command{
		{
			Name:   "generate-spec",
			Usage:  "Generate the SPEC file",
			Action: generateSpec,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file, f",
					Value: "rpm_config.json",
					Usage: "Path to the rpm_config.json file",
				},
				&cli.StringFlag{
					Name:  "arch, a",
					Value: "",
					Usage: "Target architecture of the build, e.g. amd64",
				},
				&cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Target version of the build",
				},
				&cli.StringFlag{
					Name:  "release",
					Value: "",
					Usage: "Target release of the build",
				},
			},
		},
		{
			Name:   "generate",
			Usage:  "Generate the package",
			Action: generatePkg,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file, f",
					Value: "rpm_config.json",
					Usage: "Path to the rpm_config.json file",
				},
				&cli.StringFlag{
					Name:  "build-area, b",
					Value: "pkg-build",
					Usage: "Path to the build area",
				},
				&cli.StringFlag{
					Name:  "arch, a",
					Value: "",
					Usage: "Target architecture of the build",
				},
				&cli.StringFlag{
					Name:  "output, o",
					Value: "",
					Usage: "File path to the resulting rpm file",
				},
				&cli.StringFlag{
					Name:  "version",
					Value: "",
					Usage: "Target version of the build",
				},
				&cli.StringFlag{
					Name:  "release",
					Value: "",
					Usage: "Target release of the build",
				},
			},
		},
		{
			Name:   "test",
			Usage:  "Test the package json file",
			Action: testPkg,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file, f",
					Value: "rpm_config.json",
					Usage: "Path to the rpm_config.json file",
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

// GetAppName returns application name
func GetAppName() string {
	return appName
}

// ShortVersion returns short version information
func GetShortVersion() string {
	var sb strings.Builder
	sb.WriteString(appName)
	if appVersion != "" {
		sb.WriteString("-" + appVersion)
	}
	sb.WriteString(fmt.Sprintf(", %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	sb.WriteString("\n")
	return sb.String()
}

// GetVersion returns version information
func GetVersion() string {
	var sb strings.Builder
	sb.WriteString(appName)
	if appVersion != "" {
		sb.WriteString("-" + appVersion)
	}
	sb.WriteString(fmt.Sprintf(", %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH))
	if gitCommit != "" {
		sb.WriteString(", commit: " + gitCommit)
	}
	if gitBranch != "" {
		sb.WriteString(", branch: " + gitBranch)
	}
	if buildDate != "" {
		sb.WriteString(", build on " + buildDate)
		if buildUser != "" {
			sb.WriteString(" by " + buildUser)
		}
	}
	sb.WriteString("\n")
	return sb.String()
}

func generateSpec(c *cli.Context) error {
	file := c.String("file")
	arch := c.String("arch")
	version := c.String("version")
	release := c.String("release")

	rpmJSON := rpmbuilder.Package{}

	if err := rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err := rpmJSON.Normalize(arch, version, release); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	spec, err := rpmJSON.GenerateSpecFile("")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	fmt.Printf("%s", spec)

	return nil
}

func generatePkg(c *cli.Context) error {
	var err error

	file := c.String("file")
	arch := c.String("arch")
	release := c.String("release")
	version := c.String("version")
	buildArea := c.String("build-area")
	output := c.String("output")

	if output == "" {
		return cli.NewExitError("--output,-o argument is required", 1)
	}

	rpmJSON := rpmbuilder.Package{}

	if err = rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if buildArea, err = filepath.Abs(buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = rpmJSON.Normalize(arch, version, release); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	rpmJSON.InitializeBuildArea(buildArea)

	if err = rpmJSON.WriteSpecFile("", buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = rpmJSON.RunBuild(buildArea, output); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("\n\nAll done!")

	return nil
}

func testPkg(c *cli.Context) error {
	file := c.String("file")

	rpmJSON := rpmbuilder.Package{}

	if err := rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("File is correct")

	return nil
}
