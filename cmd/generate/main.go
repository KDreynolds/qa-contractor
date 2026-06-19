package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/a-h/templ"
	"qa-contractor/templates"
)

var basePath = flag.String("base-path", "", "Base path for GitHub Pages (e.g. /qa-contractor)")

type page struct {
	path   string
	render func() templ.Component
}

func main() {
	flag.Parse()

	outDir := "out"
	os.RemoveAll(outDir)

	pages := []page{
		{"index.html", func() templ.Component { return templates.Home() }},
		{"services/index.html", func() templ.Component { return templates.Services() }},
		{"portfolio/index.html", func() templ.Component { return templates.Portfolio() }},
		{"about/index.html", func() templ.Component { return templates.About() }},
		{"contact/index.html", func() templ.Component { return templates.Contact() }},
	}

	for _, p := range pages {
		fullPath := filepath.Join(outDir, p.path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR mkdir %s: %v\n", filepath.Dir(fullPath), err)
			continue
		}

		f, err := os.Create(fullPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR create %s: %v\n", fullPath, err)
			continue
		}

		var buf strings.Builder
		p.render().Render(context.Background(), &buf)
		html := buf.String()

		if *basePath != "" {
			html = applyBasePath(html, *basePath)
		}

		f.WriteString(html)
		f.Close()
		fmt.Printf("  %s\n", p.path)
	}

	// Copy static files
	srcStatic := "static"
	dstStatic := filepath.Join(outDir, "static")
	fmt.Println("  copying static/ ...")
	copyDir(srcStatic, dstStatic)

	fmt.Println("\nDone — output in", outDir)
}

func applyBasePath(html, base string) string {
	base = strings.TrimSuffix(base, "/")

	// Replace absolute internal paths with base-prefixed versions.
	// Include closing quotes to avoid breaking attributes like href="/" → href="/qa-contractor/ class=".
	repl := strings.NewReplacer(
		`="/static/`, fmt.Sprintf(`="%s/static/`, base),
		`="/services"`, fmt.Sprintf(`="%s/services/"`, base),
		`="/portfolio"`, fmt.Sprintf(`="%s/portfolio/"`, base),
		`="/about"`, fmt.Sprintf(`="%s/about/"`, base),
		`="/contact"`, fmt.Sprintf(`="%s/contact/"`, base),
		`="/"`, fmt.Sprintf(`="%s/"`, base),
	)
	return repl.Replace(html)
}

func copyDir(src, dst string) {
	srcInfo, err := os.Stat(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR stat %s: %v\n", src, err)
		return
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR mkdir %s: %v\n", dst, err)
		return
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR readdir %s: %v\n", src, err)
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			copyDir(srcPath, dstPath)
			continue
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR open %s: %v\n", srcPath, err)
			continue
		}
		defer srcFile.Close()

		dstFile, err := os.Create(dstPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR create %s: %v\n", dstPath, err)
			continue
		}
		defer dstFile.Close()

		if _, err := io.Copy(dstFile, srcFile); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR copy %s: %v\n", dstPath, err)
		}
	}
}
