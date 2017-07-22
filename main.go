// bamfilter filters alignment records by edit distance.
package main

import (
	"os"

	"io"

	"github.com/biogo/hts/bam"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("bamfilter", "filter alignment records by edit distance.")
	app.Version("v20170721")
	maxEditDistance := app.Flag("max-edit-distance", "max edit distance (-1 for disable)").Default("-1").Int()
	bamFile := app.Arg("bamfile", "input bam file").Required().String()
	outFile := app.Arg("outfile", "outfile").Required().String()
	kingpin.MustParse(app.Parse(os.Args[1:]))

	f, err := os.Open(*bamFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bamReader, err := bam.NewReader(f, 0)
	if err != nil {
		panic(err)
	}
	h := bamReader.Header()
	bamWriter, err := bam.NewWriter(w, h, 0)
	if err != nil {
		panic(err)
	}

	for {
		rec, err := bamReader.Read()
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}
		tag, ok := rec.Tag([]byte{'N', 'M'})
		if ok {
			v := tag.Value()
			editDistance := v.(uint8)
			if int(editDistance) <= *maxEditDistance {
				bamWriter.Write(rec)
			}
		}
	}
}
