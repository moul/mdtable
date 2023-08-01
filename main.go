package main

import (
	"context"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	ff "github.com/peterbourgon/ff/v3"
	"moul.io/banner"
	"moul.io/climan"
	"moul.io/mdtable/mdtable"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		if !errors.Is(err, flag.ErrHelp) {
			fmt.Fprintf(os.Stderr, "error: %+v\n", err)
		}
		os.Exit(1)
	}
}

type opts struct {
	mdt mdtable.Opts
	csv struct {
		noHeader bool
	}
	// json struct {}
}

var o opts

func (o *opts) commonFlagBuilder(fs *flag.FlagSet) {
	fs.StringVar(&o.mdt.BodyFormat, "md-body", o.mdt.BodyFormat, "mdtable body format")
	fs.StringVar(&o.mdt.HeaderFormat, "md-header", o.mdt.HeaderFormat, "mdtable header format")
}

func run(args []string) error {
	// default opts
	o.mdt.BodyFormat = `{{range .Cols}}{{.}}\t{{end}}`
	o.mdt.HeaderFormat = `{{range .Cols}}{{.}}\t{{end}}`

	// parse CLI
	root := &climan.Command{
		Name:           "mdtable",
		ShortUsage:     "mdtable [global flags] <subcommand> [flags] [args]",
		ShortHelp:      "More info on https://moul.io/mdtable.",
		FlagSetBuilder: o.commonFlagBuilder,
		FFOptions:      []ff.Option{ff.WithEnvVarPrefix("mdtable")},
		LongHelp:       banner.Inline("mdtable"),
		Subcommands: []*climan.Command{
			{
				Name:       "csv",
				ShortUsage: "mdtable csv [flags]",
				ShortHelp:  "CSV to markdown table",
				FlagSetBuilder: func(fs *flag.FlagSet) {
					o.commonFlagBuilder(fs)
					fs.BoolVar(&o.csv.noHeader, "csv-no-header", o.csv.noHeader, "csv: no header, just body")
				},
				Exec: doCSV,
			},
			/*{
				Name:       "json",
				ShortUsage: "mdtable json [flags]",
				ShortHelp:  "JSON to markdown table",
				FlagSetBuilder: func(fs *flag.FlagSet) {
					o.commonFlagBuilder(fs)
				},
			},*/
		},
	}
	if err := root.Parse(args); err != nil {
		return fmt.Errorf("parse error: %w", err)
	}
	if err := root.Run(context.Background()); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func doCSV(_ context.Context, args []string) error {
	if len(args) > 0 {
		return flag.ErrHelp
	}

	data := mdtable.Data{
		Header: mdtable.Row{},
		Rows:   []mdtable.Row{},
	}

	csvFile := os.Stdin
	reader := csv.NewReader(csvFile)
	// Read CSV data
	for i := 0; ; i++ {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading CSV: %w", err)
		}
		row := mdtable.Row{
			Cols: record,
		}
		if i == 0 && !o.csv.noHeader {
			data.Header = row
		} else {
			data.Rows = append(data.Rows, row)
		}
	}

	if o.csv.noHeader {
		data.Header = mdtable.Row{
			Cols: make([]string, len(data.Rows[0].Cols)),
		}
		for i := range data.Rows[0].Cols {
			data.Header.Cols[i] = "<!-- -->"
		}
	}

	output, err := mdtable.Generate(data, o.mdt)
	if err != nil {
		return fmt.Errorf("mdtable generate: %w", err)
	}
	fmt.Print(output)
	return nil
}
