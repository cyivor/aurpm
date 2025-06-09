package aur

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

type PKGBUILD struct {
	PkgName       string
	PkgVer        string
	PkgRel        string
	Depends       []string
	MakeDepends   []string
	Source        []string
	Sha256Sums    []string
	BuildScript   string // build()
	PackageScript string // package()
}

// reads the pkgbuild file and extracts the fields
func ParsePKGBUILD(filePath string) (*PKGBUILD, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	pkg := &PKGBUILD{}
	var currentFunc string
	var funcLines []string
	inFunc := false

	scanner := bufio.NewScanner(file)
	arrayVar := regexp.MustCompile(`(\w+)=(\(.*?\))`)
	singleVar := regexp.MustCompile(`(\w+)='(.*?)'|(\w+)="(.*?)"|(\w+)=(.*)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// array variables (e.g., source, depends, makedepends, sha256sums)
		if matches := arrayVar.FindStringSubmatch(line); len(matches) > 0 {
			key := matches[1]
			value := strings.Trim(matches[2], "()")
			values := strings.Fields(value)
			switch key {
			case "depends":
				pkg.Depends = values
			case "makedepends":
				pkg.MakeDepends = values
			case "source":
				pkg.Source = values
			case "sha256sums":
				pkg.Sha256Sums = values
			}
			continue
		}

		// single variables (pkgname, pkgver, pkgrel)
		if matches := singleVar.FindStringSubmatch(line); len(matches) > 0 {
			for i := 1; i < len(matches); i += 2 {
				if matches[i] != "" {
					key := matches[i]
					value := matches[i+1]
					switch key {
					case "pkgname":
						pkg.PkgName = value
					case "pkgver":
						pkg.PkgVer = value
					case "pkgrel":
						pkg.PkgRel = value
					}
				}
			}
			continue
		}

		// function definitions
		if strings.HasPrefix(line, "build() {") {
			inFunc = true
			currentFunc = "build"
			funcLines = []string{}
			continue
		} else if strings.HasPrefix(line, "package() {") {
			inFunc = true
			currentFunc = "package"
			funcLines = []string{}
			continue
		} else if inFunc && line == "}" {
			inFunc = false
			if currentFunc == "build" {
				pkg.BuildScript = strings.Join(funcLines, "\n")
			} else if currentFunc == "package" {
				pkg.PackageScript = strings.Join(funcLines, "\n")
			}
			continue
		}
		if inFunc {
			funcLines = append(funcLines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return pkg, nil
}
