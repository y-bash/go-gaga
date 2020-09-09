/*
gen implements gaga's auto-generator.

Usage:
	go run generate.go -output path
	go run genucdcsv.go -output path
	go run genucdexcsv.go -output path

The coommands are:
	generate.go
		print the unichar_tables (go source code).
	genucdcsv.go
		print the UCD in csv.
	genucdexcsv.go
		print the UCDEX in csv.

Examples:
	go run generate.go -output unichar_tables.go
	go run genucdcsv.go -output ucd.csv
	go run genucdexcsv.go -output ucdex.csv

*/
package main
