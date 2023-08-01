package mdtable

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/tabwriter"
	"text/template"
)

type Row struct {
	Cols []string
	// ...
}

type Data struct {
	Header Row
	Rows   []Row
}

type Opts struct {
	HeaderFormat string
	BodyFormat   string
}

func Generate(d Data, o Opts) (string, error) {
	// validate parameters
	o.HeaderFormat = strings.ReplaceAll(o.HeaderFormat, `\t`, "\t")
	o.BodyFormat = strings.ReplaceAll(o.BodyFormat, `\t`, "\t")
	headerTmpl, err := template.New("header").Parse(o.HeaderFormat)
	if err != nil {
		return "", fmt.Errorf("parsing header format: %w", err)
	}
	bodyTmpl, err := template.New("body").Parse(o.BodyFormat)
	if err != nil {
		return "", fmt.Errorf("parsing body format: %w", err)
	}

	// generate the table
	var sb strings.Builder
	w := tabwriter.NewWriter(&sb, 0, 8, 1, ' ', tabwriter.Debug) // Debug is a way to have tabwriter adding '|' after the columns.

	// header
	var rawHeader bytes.Buffer
	if err := headerTmpl.Execute(&rawHeader, d.Header); err != nil {
		return "", fmt.Errorf("executing header template: %w", err)
	}
	fmt.Fprintf(w, "%s\n", rawHeader.String())
	// post-header delimiter
	re := regexp.MustCompile("[^|\\s]")
	delimLine := re.ReplaceAllString(rawHeader.String(), "-")
	fmt.Fprintf(w, "%s\n", delimLine)

	for _, row := range d.Rows {
		if err := bodyTmpl.Execute(w, row); err != nil {
			return "", fmt.Errorf("executing body template: %w", err)
		}
		fmt.Fprintln(w)
	}
	w.Flush()
	tw := sb.String()
	// with tabwriter.Debug, we automatically have pipes on the right of fields, but we miss the first pipe on the left.

	lines := strings.Split(tw, "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}
		lines[i] = "|" + line
	}
	md := strings.Join(lines, "\n")
	return md, nil
}
