// gendoc generates tracking DOCX or XLSX files for CyberPhish attachment campaigns.
// The generated files embed a remote image URL that triggers when opened in
// Microsoft Word/Excel, logging an "Opened Attachment" event.
//
// Usage:
//   go run ./cmd/gendoc -url https://phish.example.com/track-open?rid=XXXXX -o Invoice.docx
//   go run ./cmd/gendoc -url https://phish.example.com/track-open?rid=XXXXX -o Report.xlsx
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gophish/gophish/util"
)

func main() {
	trackURL := flag.String("url", "", "Tracking URL (e.g. https://phish.example.com/track-open?rid=XXXXX)")
	output := flag.String("o", "", "Output file path (e.g. Invoice.docx or Report.xlsx)")
	flag.Parse()

	if *trackURL == "" || *output == "" {
		fmt.Fprintln(os.Stderr, "Usage: gendoc -url <tracking-url> -o <output-file>")
		fmt.Fprintln(os.Stderr, "  Supported extensions: .docx, .xlsx")
		os.Exit(1)
	}

	ext := strings.ToLower(filepath.Ext(*output))
	var err error
	switch ext {
	case ".docx":
		err = util.GenerateTrackingDocx(*trackURL, *output)
	case ".xlsx":
		err = util.GenerateTrackingXlsx(*trackURL, *output)
	default:
		fmt.Fprintf(os.Stderr, "Unsupported extension: %s (use .docx or .xlsx)\n", ext)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("[+] Generated: %s\n", *output)
	fmt.Printf("[+] Tracking URL embedded: %s\n", *trackURL)
	fmt.Println("[+] Drop the file in ./static/attachments/ and use in email template:")
	fmt.Printf("    {{.URL}}/download?rid={{.RId}}&f=%s\n", filepath.Base(*output))
}
