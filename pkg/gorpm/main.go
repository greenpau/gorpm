package gorpm

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver"
	//"github.com/davecgh/go-spew/spew"
	"github.com/mattn/go-zglob"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	//"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

func init() {
	log.Debugf("initialized package gorpm")
}

// Package contains the build information
type Package struct {
	Name              string            `json:"name"`
	Version           string            `json:"version,omitempty"`
	Arch              string            `json:"arch,omitempty"`
	Release           string            `json:"release,omitempty"`
	Distro            string            `json:"distro,omitempty"`
	CPU               string            `json:"cpu,omitempty"`
	Group             string            `json:"group,omitempty"`
	License           string            `json:"license,omitempty"`
	URL               string            `json:"url,omitempty"`
	Summary           string            `json:"summary,omitempty"`
	Description       string            `json:"description,omitempty"`
	ChangelogFile     string            `json:"changelog-file,omitempty"`
	ChangelogCmd      string            `json:"changelog-cmd,omitempty"`
	Files             []fileInstruction `json:"files,omitempty"`
	Sources           []string          `json:"sources,omitempty"`
	PreInstallScript  string            `json:"pre_install_script,omitempty"`
	PostInstallScript string            `json:"post_install_script,omitempty"`
	PreRemoveScript   string            `json:"pre_remove_script,omitempty"`
	PostRemoveScript  string            `json:"post_remove_script,omitempty"`
	VerifyScript      string            `json:"verify_script,omitempty"`
	CleanupScript     string            `json:"cleanup_script,omitempty"`
	BuildRequires     []string          `json:"build-requires,omitempty"`
	Requires          []string          `json:"requires,omitempty"`
	Provides          []string          `json:"provides,omitempty"`
	Conflicts         []string          `json:"conflicts,omitempty"`
	Envs              []*EnvVar         `json:"envs,omitempty"`
	Menus             []menu            `json:"menus"`
	AutoReqProv       string            `json:"auto-req-prov,omitempty"`
}

// EnvVar represents an environment variable.
type EnvVar struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

type fileInstruction struct {
	From        string `json:"from,omitempty"`
	To          string `json:"to,omitempty"`
	Base        string `json:"base,omitempty"`
	Permissions string `json:"perms,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Group       string `json:"group,omitempty"`
}

type menu struct {
	Name            string `json:"name"`           // Name of the shortcut
	GenericName     string `json:"generic-name"`   //
	Exec            string `json:"exec"`           // Exec command
	Icon            string `json:"icon"`           // Path to the installed icon
	Type            string `json:"type"`           // Type of shortcut
	StartupNotify   bool   `json:"startup-notify"` // yes/no
	Terminal        bool   `json:"terminal"`       // yes/no
	DBusActivatable bool   `json:"dbus-activable"` // yes/no
	NoDisplay       bool   `json:"no-display"`     // yes/no
	Keywords        string `json:"keywords"`       // ; separated list
	OnlyShowIn      string `json:"only-show-in"`   // ; separated list
	Categories      string `json:"categories"`     // ; separated list
	MimeType        string `json:"mime-type"`      // ; separated list
}

// Load package build information
func (p *Package) Load(file string) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return errors.Errorf("json file '%s' does not exist: %s", file, err.Error())
	}
	byt, err := ioutil.ReadFile(file)
	if err != nil {
		return errors.Errorf("error occured while reading file '%s': %s", file, err.Error())
	}
	if err := json.Unmarshal(byt, p); err != nil {
		return errors.Errorf("Invalid json file '%s': %s", file, err.Error())
	}
	return nil
}

// Normalize build information
//func (p *Package) Normalize(arch string, version string, release string) error {
func (p *Package) Normalize(params map[string]string) error {
	tokens := make(map[string]string)
	for k, v := range params {
		tokens["!"+k+"!"] = v
	}
	tokens["!name!"] = p.Name

	p.Version = replaceTokens(p.Version, tokens)
	p.Release = replaceTokens(p.Release, tokens)
	p.Arch = replaceTokens(p.Arch, tokens)
	p.Distro = replaceTokens(p.Distro, tokens)
	p.CPU = replaceTokens(p.CPU, tokens)
	p.URL = replaceTokens(p.URL, tokens)
	p.Summary = replaceTokens(p.Summary, tokens)
	p.Description = replaceTokens(p.Description, tokens)
	p.ChangelogFile = replaceTokens(p.ChangelogFile, tokens)
	p.ChangelogCmd = replaceTokens(p.ChangelogCmd, tokens)

	if p.Release == "" {
		return errors.WithStack(fmt.Errorf("release not found"))
	}
	if p.Version == "" {
		return errors.WithStack(fmt.Errorf("version not found"))
	}
	if p.Arch == "" {
		return errors.WithStack(fmt.Errorf("arch not found"))
	}
	if p.Distro == "" {
		return errors.WithStack(fmt.Errorf("distro not found"))
	}
	p.Release += "." + p.Distro
	if p.CPU == "" {
		return errors.WithStack(fmt.Errorf("cpu family not found"))
	}

	for i, v := range p.Files {
		p.Files[i].From = replaceTokens(v.From, tokens)
		p.Files[i].Base = replaceTokens(v.Base, tokens)
		p.Files[i].To = replaceTokens(v.To, tokens)
	}
	log.Infof("Arch=%s\n", p.Arch)
	log.Infof("Version=%s\n", p.Version)
	log.Infof("Release=%s\n", p.Release)
	log.Infof("Distribution=%s\n", p.Distro)
	log.Infof("CPU Family=%s\n", p.CPU)
	log.Infof("URL=%s\n", p.URL)
	log.Infof("Summary=%s\n", p.Summary)
	log.Infof("Description=%s\n", p.Description)
	log.Infof("ChangelogFile=%s\n", p.ChangelogFile)
	log.Infof("ChangelogCmd=%s\n", p.ChangelogCmd)

	shortcuts, err := p.WriteShortcutFiles()
	if err != nil {
		return errors.WithStack(err)
	}
	log.Infof("shortcuts=%s\n", shortcuts)
	for _, shortcut := range shortcuts {
		sc := fileInstruction{}
		sc.From = shortcut
		sc.To = fmt.Sprintf("%%{_datadir}/applications/")
		sc.Base = filepath.Dir(shortcut)
		p.Files = append(p.Files, sc)
		log.Infof("Added menu shortcut File=%q\n", sc)
	}
	for _, menu := range p.Menus {
		sc := fileInstruction{}
		f, err := filepath.Abs(menu.Icon)
		if err != nil {
			return errors.WithStack(err)
		}
		sc.From = f
		sc.To = fmt.Sprintf("%%{_datadir}/pixmaps/")
		sc.Base = filepath.Dir(f)
		p.Files = append(p.Files, sc)
		log.Infof("Added menu icon File=%q\n", sc)

		// desktop-file-utils is super picky.
		menu.Categories = strings.TrimSuffix(menu.Categories, ";")
		menu.Keywords = strings.TrimSuffix(menu.Keywords, ";")
		if menu.Categories != "" {
			menu.Categories += ";"
		}
		if menu.Keywords != "" {
			menu.Keywords += ";"
		}
	}

	if len(p.Menus) > 0 {
		if contains(p.BuildRequires, "desktop-file-utils") == false {
			p.BuildRequires = append(p.BuildRequires, "desktop-file-utils")
		}
	}

	if len(p.Sources) > 0 {
		for i, v := range p.Sources {
			p.Sources[i] = replaceTokens(v, tokens)
		}
	}

	if len(p.Envs) > 0 {
		for _, v := range p.Envs {
			v.Value = replaceTokens(v.Value, tokens)
		}
		envFile, err := p.WriteEnvFile()
		if err != nil {
			return errors.WithStack(err)
		}
		sc := fileInstruction{}
		sc.From = envFile
		sc.To = fmt.Sprintf("%%{_sysconfdir}/profile.d/")
		sc.Base = filepath.Dir(envFile)
		sc.Owner = "root"
		sc.Group = "root"
		sc.Permissions = "644"
		log.Infof("Added env File=%q\n", sc)
		p.Files = append(p.Files, sc)
	}
	log.Infof("p.Envs=%v\n", p.Envs)
	log.Infof("p.Requires=%s\n", p.Requires)
	log.Infof("p.BuildRequires=%s\n", p.BuildRequires)
	log.Infof("p.AutoReqProv=%s\n", p.AutoReqProv)
	return nil
}

func replaceTokens(in string, tokens map[string]string) string {
	for token, v := range tokens {
		in = strings.Replace(in, token, v, -1)
	}
	return in
}

// InitializeBuildArea intializes the build area
func (p *Package) InitializeBuildArea(buildAreaPath string) error {
	paths := make([]string, 0)
	paths = append(paths, filepath.Join(buildAreaPath, "BUILD"))
	paths = append(paths, filepath.Join(buildAreaPath, "RPMS"))
	paths = append(paths, filepath.Join(buildAreaPath, "SOURCES"))
	paths = append(paths, filepath.Join(buildAreaPath, "SPECS"))
	paths = append(paths, filepath.Join(buildAreaPath, "SRPMS"))
	paths = append(paths, filepath.Join(buildAreaPath, "RPMS", "i386"))
	paths = append(paths, filepath.Join(buildAreaPath, "RPMS", "amd64"))

	for _, p := range paths {
		if err := os.MkdirAll(p, 0755); err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}

// WriteSpecFile writes the spec file.
func (p *Package) WriteSpecFile(sourceDir string, buildAreaPath string) error {
	spec, err := p.GenerateSpecFile(sourceDir)
	if err != nil {
		return errors.WithStack(err)
	}
	path := filepath.Join(buildAreaPath, "SPECS", p.Name+".spec")
	return ioutil.WriteFile(path, []byte(spec), 0644)
}

// RunBuild executes the build of buildAreaPath.
func (p *Package) RunBuild(buildAreaPath string, output string) error {
	path := filepath.Join(buildAreaPath, "SPECS", p.Name+".spec")
	def := "_topdir " + buildAreaPath
	arch := p.Arch
	if arch == "386" {
		arch = "i386"
	}
	if arch == "amd64" {
		arch = "x86_64"
	}
	args := []string{"--target", arch, "-bb", path, "--define", def}
	log.Infof("%s %s\n", "rpmbuild", args)
	oCmd := exec.Command("rpmbuild", args...)
	oCmd.Stdout = os.Stdout
	oCmd.Stderr = os.Stderr
	if err := oCmd.Run(); err != nil {
		return errors.WithStack(err)
	}
	// if version contains a prerelease,
	// destination file generated by rpmbuild is like
	// [name]-[version]-[prerelease].[release].[arch].rpm
	// otherwise
	// [name]-[version]-[release].[arch].rpm
	pkg := fmt.Sprintf("%s/RPMS/%s/%s-%s-%s.%s.rpm", buildAreaPath, arch, p.Name, p.Version, p.Release, arch)
	v, err := semver.NewVersion(p.Version)
	if err != nil {
		return errors.WithStack(err)
	}
	if v.Prerelease() != "" {
		pkg = fmt.Sprintf("%s/RPMS/%s/%s-%s.%s.%s.rpm", buildAreaPath, arch, p.Name, p.Version, p.Release, arch)
	}
	return cp(output, pkg)
}

// GenerateSpecFile generates the spec file.
func (p *Package) GenerateSpecFile(sourceDir string) (string, error) {
	spec := ""

	// Version field of the spec file must not
	// contain non numeric characters,
	// see https://fedoraproject.org/wiki/Packaging:Naming?rd=Packaging:NamingGuidelines#Version_Tag
	// the prerelease stuff is moved into Release field
	v, err := semver.NewVersion(p.Version)
	if err != nil {
		return "", errors.WithStack(err)
	}
	okVersion := ""
	okVersion += strconv.FormatInt(int64(v.Major()), 10)
	okVersion += "." + strconv.FormatInt(int64(v.Minor()), 10)
	okVersion += "." + strconv.FormatInt(int64(v.Patch()), 10)
	preRelease := p.Release
	if v.Prerelease() != "" {
		preRelease = v.Prerelease() + "." + preRelease
	}

	if p.Name != "" {
		spec += fmt.Sprintf("Name: %s\n", p.Name)
	}
	if p.Version != "" {
		spec += fmt.Sprintf("Version: %s\n", okVersion)
	}
	if p.Release != "" {
		spec += fmt.Sprintf("Release: %s\n", preRelease)
	}
	if p.Group != "" {
		spec += fmt.Sprintf("Group: %s\n", p.Group)
	}
	if p.License != "" {
		spec += fmt.Sprintf("License: %s\n", p.License)
	}
	if p.URL != "" {
		spec += fmt.Sprintf("Url: %s\n", p.URL)
	}
	if p.Summary != "" {
		spec += fmt.Sprintf("Summary: %s\n", p.Summary)
	}
	if len(p.Sources) > 0 {
		spec += "\n"
		for _, v := range p.Sources {
			spec += fmt.Sprintf("Source0: %s\n", v)
		}
		spec += "\n"
	}
	if len(p.BuildRequires) > 0 {
		spec += fmt.Sprintf("\nBuildRequires: %s\n", strings.Join(p.BuildRequires, ", "))
	}
	if len(p.Requires) > 0 {
		spec += fmt.Sprintf("\nRequires: %s\n", strings.Join(p.Requires, ", "))
	}
	if len(p.Provides) > 0 {
		spec += fmt.Sprintf("\nProvides: %s\n", strings.Join(p.Provides, ", "))
	}
	if len(p.Conflicts) > 0 {
		spec += fmt.Sprintf("\nConflicts: %s\n", strings.Join(p.Conflicts, ", "))
	}
	if p.Description != "" {
		spec += fmt.Sprintf("\n%%description\n%s\n", p.Description)
	}
	spec += fmt.Sprintf("\n%%prep\n")
	spec += fmt.Sprintf("\n%%build\n")

	log.Warnf("Reached install section")
	spec += fmt.Sprintf("\n%%install\n")

	install, err := p.GenerateInstallSection(sourceDir)
	if err != nil {
		return "", errors.WithStack(err)
	}
	spec += fmt.Sprintf("%s\n", install)
	spec += fmt.Sprintf("\n%%files\n")
	files, err := p.GenerateFilesSection(sourceDir)
	if err != nil {
		return "", errors.WithStack(err)
	}
	spec += fmt.Sprintf("%s\n", files)
	spec += fmt.Sprintf("\n%%clean\n")

	if content := readFile(p.CleanupScript); content != "" {
		spec += fmt.Sprintf("%s\n", content)
	}

	shortcutInstall := "\n"
	for _, menu := range p.Menus {
		shortcutInstall += fmt.Sprintf("desktop-file-install --vendor='' ")
		shortcutInstall += fmt.Sprintf("--dir=%%{buildroot}%%{_datadir}/applications/%s ", p.Name)
		shortcutInstall += fmt.Sprintf("%%{buildroot}/%%{_datadir}/applications/%s.desktop", menu.Name)
		shortcutInstall += "\n"
	}
	shortcutInstall = strings.TrimSpace(shortcutInstall)
	if content := readFile(p.PreInstallScript); content != "" {
		spec += fmt.Sprintf("\n%%pre\n%s\n", content)
	}
	if content := readFile(p.PostInstallScript); content != "" {
		spec += fmt.Sprintf("\n%%post\n%s\n", content+shortcutInstall)
	} else if shortcutInstall != "" {
		spec += fmt.Sprintf("\n%%post\n%s\n", shortcutInstall)
	}
	if content := readFile(p.PreRemoveScript); content != "" {
		spec += fmt.Sprintf("\n%%preun\n%s\n", content)
	}
	if content := readFile(p.PostRemoveScript); content != "" {
		spec += fmt.Sprintf("\n%%postun\n%s\n", content)
	}
	if content := readFile(p.VerifyScript); content != "" {
		spec += fmt.Sprintf("\n%%verifyscript\n%s\n", content)
	}
	spec += fmt.Sprintf("\n%%changelog\n")
	content, err := p.GetChangelogContent()
	if err != nil {
		return "", errors.WithStack(err)
	}
	spec += fmt.Sprintf("%s\n", content)

	return spec, nil
}

// GenerateInstallSection generates the install section.
func (p *Package) GenerateInstallSection(sourceDir string) (string, error) {
	var err error
	content := ""

	log.Warnf("Generating install section")
	log.Warnf("Source directory: %s", sourceDir)

	allDirs := make([]string, 0)
	allFiles := make([]string, 0)

	if sourceDir, err = filepath.Abs(sourceDir); err != nil {
		return "", errors.WithStack(err)
	}

	for i, fileInst := range p.Files {

		log.Warnf("Processing %v", fileInst)

		if fileInst.From == "" {
			log.Infof("Skipped p.Files[%d] %q", i, fileInst)
			continue
		}

		from := fileInst.From
		to := fileInst.To
		base := fileInst.Base

		if filepath.IsAbs(from) == false {
			from = filepath.Join(sourceDir, from)
		}
		if filepath.IsAbs(base) == false {
			base = filepath.Join(sourceDir, base)
		}

		log.Infof("fileInst.From=%q\n", from)
		log.Infof("fileInst.To=%q\n", to)
		log.Infof("fileInst.Base=%q\n", base)

		items, err := zglob.Glob(from)
		if err != nil {
			log.Printf("Files not found in '%s'\n", from)
			continue
		}

		log.Infof("items=%q\n", items)

		for _, item := range items {
			n := item
			if len(item) >= len(base) && item[0:len(base)] == base {
				n = item[len(base):]
			}
			n = filepath.Join("%{buildroot}", to, n)
			dir := fmt.Sprintf("mkdir -p %s\n", filepath.Dir(n))
			if contains(allDirs, dir) == false {
				allDirs = append(allDirs, dir)
			}
			if s, err := os.Stat(item); err != nil {
				return "", err
			} else if s.IsDir() == false {
				file := fmt.Sprintf("cp %s %s\n", item, filepath.Dir(n))
				if contains(allFiles, file) == false {
					allFiles = append(allFiles, file)
				}
			}
		}
	}

	for _, d := range allDirs {
		content += d
	}
	for _, d := range allFiles {
		content += d
	}

	log.Infof("content=\n%s\n", content)

	return content, nil
}

// GenerateFilesSection generates the files section.
func (p *Package) GenerateFilesSection(sourceDir string) (string, error) {
	var err error
	content := ""
	allItems := make([]fileItem, 0)

	if sourceDir, err = filepath.Abs(sourceDir); err != nil {
		return "", errors.WithStack(err)
	}

	for _, fileInst := range p.Files {
		from := fileInst.From
		to := fileInst.To
		base := fileInst.Base
		filePerms := fileInst.Permissions

		if from == "" {
			content += fmt.Sprintf(" %s\n", filePerms)
			continue
		}

		if filepath.IsAbs(from) == false {
			from = filepath.Join(sourceDir, from)
		}
		if filepath.IsAbs(base) == false {
			base = filepath.Join(sourceDir, base)
		}

		log.Infof("fileInst.From=%q\n", from)
		log.Infof("fileInst.To=%q\n", to)
		log.Infof("fileInst.Base=%q\n", base)
		log.Infof("fileInst.Permissions=%q\n", filePerms)
		log.Infof("fileInst.Owner=%q\n", fileInst.Owner)
		log.Infof("fileInst.Group=%q\n", fileInst.Group)

		items, err := zglob.Glob(from)
		if err != nil {
			log.Printf("Files not found in '%s'\n", from)
			continue
		}

		log.Infof("items=%q\n", items)

		for _, item := range items {
			n := item
			if len(item) >= len(base) && item[0:len(base)] == base {
				n = item[len(base):]
			}
			n = filepath.Join(to, n)
			if fileItems(allItems).contains(n) == false {
				allItems = append(allItems, fileItem{
					n, filePerms,
					fileInst.Owner,
					fileInst.Group,
				})
			}
		}
	}

	for _, item := range allItems {
		content += "%"
		content += fmt.Sprintf("attr(%s, %s, %s) %s\n",
			item.Permissions, item.Owner, item.Group, item.Path,
		)
	}

	log.Infof("content=\n%s\n", content)

	return content, nil
}

// GetChangelogContent generates the changelog content.
func (p *Package) GetChangelogContent() (string, error) {
	var err error
	var c []byte
	var wd string
	var cmd *exec.Cmd
	if p.ChangelogFile != "" {
		if c, err = ioutil.ReadFile(p.ChangelogFile); err == nil {
			return string(c), nil
		}
	} else if p.ChangelogCmd != "" {
		wd, err = os.Getwd()
		if err == nil {
			cmd, err = ExecCommand(wd, p.ChangelogCmd)
			if err == nil {
				cmd.Stdout = nil
				c, err = cmd.Output()
				if err == nil {
					return string(c), nil
				}
			}
		}
	}
	return "", errors.WithStack(err)
}

// WriteShortcutFiles writes the shortcuts in the build area.
func (p *Package) WriteShortcutFiles() ([]string, error) {

	files := make([]string, 0)

	tpmDir, err := ioutil.TempDir("", "rpm-desktops")
	if err != nil {
		return files, errors.WithStack(err)
	}

	for _, m := range p.Menus {
		s := ""

		if m.Name != "" {
			s += fmt.Sprintf("%s=%s\n", "Name", m.Name)
		}

		if m.GenericName != "" {
			s += fmt.Sprintf("%s=%s\n", "GenericName", m.GenericName)
		}

		if m.Exec != "" {
			s += fmt.Sprintf("%s=%s\n", "Exec", m.Exec)
		}

		if m.Icon != "" {
			s += fmt.Sprintf("%s=%s\n", "Icon", "/usr/share/pixmaps/"+filepath.Base(m.Icon))
		}

		if m.Type != "" {
			s += fmt.Sprintf("%s=%s\n", "Type", m.Type)
		}

		if m.Categories != "" {
			s += fmt.Sprintf("%s=%s\n", "Categories", m.Categories+";")
		}

		if m.MimeType != "" {
			s += fmt.Sprintf("%s=%s\n", "MimeType", m.MimeType)
		}

		if m.OnlyShowIn != "" {
			s += fmt.Sprintf("%s=%s\n", "OnlyShowIn", m.OnlyShowIn)
		}

		if m.Keywords != "" {
			s += fmt.Sprintf("%s=%s\n", "Keywords", m.Keywords+";")
		}

		if s != "" {

			if m.StartupNotify {
				s += fmt.Sprintf("%s=%s\n", "StartupNotify", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "StartupNotify", "false")
			}

			if m.DBusActivatable {
				s += fmt.Sprintf("%s=%s\n", "DBusActivatable", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "DBusActivatable", "false")
			}

			if m.NoDisplay {
				s += fmt.Sprintf("%s=%s\n", "NoDisplay", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "NoDisplay", "false")
			}

			if m.Terminal {
				s += fmt.Sprintf("%s=%s\n", "Terminal", "true")
			} else {
				s += fmt.Sprintf("%s=%s\n", "Terminal", "false")
			}

			s = "[Desktop Entry]\n" + s

			file := filepath.Join(tpmDir, m.Name+".desktop")

			files = append(files, file)

			if err := ioutil.WriteFile(file, []byte(s), 0644); err != nil {
				return files, errors.WithStack(err)
			}
		}
	}

	return files, nil
}

// WriteEnvFile writes the env file in the build area.
func (p *Package) WriteEnvFile() (string, error) {

	file := ""

	//user, err := user.Current()
	//if err != nil {
	//	return file, errors.WithStack(err)
	//}

	//tmpDirPath := filepath.Join(os.TempDir(), user.Username, "rpm-build-env")
	tmpDirPath := filepath.Join("./etc/profile.d/")
	if _, err := os.Stat(tmpDirPath); err != nil {
		if mkdirErr := os.MkdirAll(tmpDirPath, os.ModePerm); mkdirErr != nil {
			return file, errors.WithStack(mkdirErr)
		}
	}

	//tmpDir, err := ioutil.TempDir(tmpDirPath, "")
	//if err != nil {
	//	return file, errors.WithStack(err)
	//}

	//tmpFileName := filepath.Join(tmpDir, p.Name+".sh")
	tmpFileName := filepath.Join(tmpDirPath, p.Name+".sh")

	content := fmt.Sprintf("# Global environment variables for %s\n\n", p.Name)

	for _, v := range p.Envs {
		content += fmt.Sprintf("export %s=%s\n", v.Name, v.Value)
	}

	return tmpFileName, errors.WithStack(ioutil.WriteFile(tmpFileName, []byte(content), 0644))
}

type fileItem struct {
	Path        string
	Permissions string
	Owner       string
	Group       string
}

type fileItems []fileItem

func (f fileItems) contains(path string) bool {
	for _, item := range f {
		if item.Path == path {
			return true
		}
	}
	return false
}

func contains(l []string, v string) bool {
	for _, vv := range l {
		if vv == v {
			return true
		}
	}
	return false
}

func cp(dst, src string) error {
	s, err := os.Open(src)
	if err != nil {
		return errors.WithStack(err)
	}
	defer s.Close()
	d, err := os.Create(dst)
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := io.Copy(d, s); err != nil {
		d.Close()
		return errors.WithStack(err)
	}
	return d.Close()
}

func readFile(src string) string {
	c, err := ioutil.ReadFile(src)
	if err != nil {
		return ""
	}
	return string(c)
}
