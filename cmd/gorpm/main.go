// Create binary rpm package with ease
package main

import (
	"fmt"
	"github.com/greenpau/gorpm/pkg/gorpm"
	"github.com/greenpau/versioned"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

var (
	app        *versioned.PackageManager
	appVersion string
	gitBranch  string
	gitCommit  string
	buildUser  string
	buildDate  string
)

func init() {
	app = versioned.NewPackageManager("gorpm")
	app.Description = "RPM utilities in Go."
	app.Documentation = "https://github.com/greenpau/gorpm/"
	app.SetVersion(appVersion, "")
	app.SetGitBranch(gitBranch, "")
	app.SetGitCommit(gitCommit, "")
	app.SetBuildUser(buildUser, "")
	app.SetBuildDate(buildDate, "")
}

func main() {
	cliApp := cli.NewApp()
	cliApp.Name = app.Name
	cliApp.Version = app.Version
	cliApp.Usage = "RPM utilities in Go"
	cliApp.UsageText = fmt.Sprintf("%s <cmd> <options>", app.Name)
	cliApp.Commands = []cli.Command{
		{
			Name:   "generate-spec",
			Usage:  "Generate the SPEC file",
			Action: generateSpec,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "file, f",
					Value: "config.json",
					Usage: "Path to the config.json file",
				},
				&cli.StringFlag{
					Name:  "arch, a",
					Value: "",
					Usage: "Target CPU architecture of the build, e.g. amd64",
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
				&cli.StringFlag{
					Name:  "distro",
					Value: "",
					Usage: "Target distribution of the build",
				},
				&cli.StringFlag{
					Name:  "cpu",
					Value: "",
					Usage: "Target CPU Instruction Set Architecture (ISA) of the build, e.g. x86_64",
				},
				&cli.StringFlag{
					Name:  "output, o",
					Value: "",
					Usage: "File path to the resulting RPM .spec file",
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
					Value: "config.json",
					Usage: "Path to the config.json file",
				},
				&cli.StringFlag{
					Name:  "build-area, b",
					Value: "pkg-build",
					Usage: "Path to the build area",
				},
				&cli.StringFlag{
					Name:  "arch, a",
					Value: "",
					Usage: "Target CPU architecture of the build, e.g. amd64",
				},
				&cli.StringFlag{
					Name:  "release",
					Value: "",
					Usage: "Target release of the build",
				},
				&cli.StringFlag{
					Name:  "distro",
					Value: "",
					Usage: "Target distribution of the build",
				},
				&cli.StringFlag{
					Name:  "cpu",
					Value: "",
					Usage: "Target CPU Instruction Set Architecture (ISA) of the build, e.g. x86_64",
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
					Value: "config.json",
					Usage: "Path to the config.json file",
				},
			},
		},
	}

	err := cliApp.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func generateSpec(c *cli.Context) error {
	cliInput := make(map[string]string)
	cliInput["file"] = c.String("file")
	cliInput["arch"] = c.String("arch")
	cliInput["version"] = c.String("version")
	cliInput["release"] = c.String("release")
	cliInput["distro"] = c.String("distro")
	cliInput["cpu"] = c.String("cpu")

	if cliInput["file"] == "" {
		return cli.NewExitError("--file,-f argument is required", 1)
	}

	output := c.String("output")
	rpmJSON := gorpm.Package{}

	if err := rpmJSON.Load(cliInput["file"]); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err := rpmJSON.Normalize(cliInput); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	spec, err := rpmJSON.GenerateSpecFile("")
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if output != "" {
		if err := ioutil.WriteFile(output, []byte(spec), 0644); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}
	} else {
		fmt.Printf("%s", spec)
	}

	return nil
}

func generatePkg(c *cli.Context) error {
	var err error

	cliInput := make(map[string]string)
	cliInput["file"] = c.String("file")
	cliInput["arch"] = c.String("arch")
	cliInput["version"] = c.String("version")
	cliInput["release"] = c.String("release")
	cliInput["distro"] = c.String("distro")
	cliInput["cpu"] = c.String("cpu")

	buildArea := c.String("build-area")
	output := c.String("output")
	if output == "" {
		return cli.NewExitError("--output,-o argument is required", 1)
	}

	rpmJSON := gorpm.Package{}

	if err = rpmJSON.Load(cliInput["file"]); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if buildArea, err = filepath.Abs(buildArea); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	if err = rpmJSON.Normalize(cliInput); err != nil {
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

	rpmJSON := gorpm.Package{}

	if err := rpmJSON.Load(file); err != nil {
		return cli.NewExitError(err.Error(), 1)
	}

	fmt.Println("File is correct")

	return nil
}
