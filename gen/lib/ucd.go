package lib

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// The UCD in XML (Unicode Character Database in XML)
const ucdURL = "https://unicode.org/Public/13.0.0/ucdxml/ucd.nounihan.flat.zip"

// UCD is the top level element of the UCD in XML
type UCD struct {
	Chars []Char `xml:"repertoire>char"`
}

// Char is the most important element of the UCD in XML
type Char struct {
	Cp  string `xml:"cp,attr"`
	Age string `xml:"age,attr"`
	Na  string `xml:"na,attr"`
	Gc  string `xml:"gc,attr"`
	Dt  string `xml:"dt,attr"`
	Dm  string `xml:"dm,attr"`
	Blk string `xml:"blk,attr"`
	Suc string `xml:"suc,attr"`
	Slc string `xml:"slc,attr"`
}

func readUCD() (*UCD, error) {
	const filename = "ucd.nounihan.flat.xml"

	// download
	resp, err := http.Get(ucdURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// unzip
	zipReader, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		return nil, err
	}
	if len(zipReader.File) != 1 {
		return nil, fmt.Errorf("want: 1 zip member, have: %d members", len(zipReader.File))
	}
	zipMember := zipReader.File[0]
	if zipMember.Name != filename {
		return nil, fmt.Errorf("want: %s, have: %s", filename, zipMember.Name)
	}
	f, err := zipMember.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	buf, err = ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// ucd
	ucd := UCD{}
	err = xml.Unmarshal(buf, &ucd)
	if err != nil {
		return nil, err
	}
	return &ucd, nil
}

func uplus(s string) string {
	if len(s) == 0 || s == "#" {
		return ""
	}
	return fmt.Sprintf("U+%s", s)
}

func writeUCD(f io.Writer, ucd *UCD) {
	fmt.Fprint(f, "cp,age,na,gc,dt,dm,blk,suc,slc")
	fmt.Fprintln(f, "")
	for _, c := range ucd.Chars {
		if len(c.Cp) > 0 {
			fmt.Fprintf(f, "%s", uplus(c.Cp))
			fmt.Fprintf(f, ",%s", c.Age)
			fmt.Fprintf(f, ",%s", c.Na)
			fmt.Fprintf(f, ",%s", c.Gc)
			fmt.Fprintf(f, ",%s", c.Dt)
			fmt.Fprintf(f, ",%s", uplus(c.Dm))
			fmt.Fprintf(f, ",%s", c.Blk)
			fmt.Fprintf(f, ",%s", uplus(c.Suc))
			fmt.Fprintf(f, ",%s", uplus(c.Slc))
			fmt.Fprintln(f, "")
		}
	}
}

// GenUCD generates the UCD in CSV
func GenUCD(f io.Writer) error {
	ucd, err := readUCD()
	if err != nil {
		return err
	}
	writeUCD(f, ucd)
	return nil
}
