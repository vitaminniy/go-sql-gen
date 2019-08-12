package main

import (
	"os"
	"reflect"
	"regexp"
	"testing"
)

func TestParseFile(t *testing.T) {
	cases := []struct {
		file             string
		expected         parsedFile
		ignoreUnexported bool
		re               *regexp.Regexp
	}{
		{
			file: "testdata/parse_file_1",
			expected: parsedFile{
				Package: "main",
				Structs: []parsedStruct{{Name: "ParseFile", Fields: []string{"A", "B", "C", "D"}}},
			},
			re:               regexp.MustCompile("."),
			ignoreUnexported: false,
		},
		{
			file: "testdata/parse_file_2",
			expected: parsedFile{
				Package: "main",
				Structs: []parsedStruct{{Name: "ParseFile", Fields: []string{"A", "C", "D"}}},
			},
			re:               regexp.MustCompile("."),
			ignoreUnexported: true,
		},
		{
			file: "testdata/parse_file_3",
			expected: parsedFile{
				Package: "main",
				Structs: []parsedStruct{{Name: "ParseFile", Fields: []string{"A", "D"}}},
			},
			re:               regexp.MustCompile("."),
			ignoreUnexported: true,
		},
		{
			file: "testdata/parse_file_3",
			expected: parsedFile{
				Package: "main",
				Structs: []parsedStruct{{Name: "ParseFile", Fields: []string{"A", "D"}}},
			},
			re:               regexp.MustCompile("File"),
			ignoreUnexported: true,
		},
		{
			file: "testdata/parse_file_3",
			expected: parsedFile{
				Package: "main",
				Structs: []parsedStruct{},
			},
			re:               regexp.MustCompile("existedStructName"),
			ignoreUnexported: true,
		},
	}

	testFn := func(t *testing.T, file string, expected parsedFile, re *regexp.Regexp, ignoreUnexported bool) func(*testing.T) {
		return func(t *testing.T) {
			f, err := os.Open(file)
			if err != nil {
				t.Fatalf("couldn't open %s: %v\n", file, err)
			}
			defer f.Close()

			actual, err := parseFile(f, re, ignoreUnexported)
			if err != nil {
				t.Fatalf("parseFile failed: %v\n", err)
			}

			if expected.Package != actual.Package {
				t.Fatalf("package name mismatch: want %s; got %s\n", expected.Package, actual.Package)
			}

			if len(expected.Structs) != len(actual.Structs) {
				t.Fatalf("structs amount mismatch: want %d; got %d\n", len(expected.Structs), len(actual.Structs))
			}

			for i := 0; i < len(expected.Structs); i++ {
				if expected.Structs[i].Name != actual.Structs[i].Name {
					t.Fatalf("struct #%d name mismatch: want %s; got %s\n", i, expected.Structs[i].Name, actual.Structs[i].Name)
				}

				if !reflect.DeepEqual(expected.Structs[i].Fields, actual.Structs[i].Fields) {
					t.Fatalf("structs #%d fields mismatch: want %v; got %v\n", i, expected.Structs[i].Fields, actual.Structs[i].Fields)
				}
			}
		}
	}

	for _, c := range cases {
		t.Run(c.file, testFn(t, c.file, c.expected, c.re, c.ignoreUnexported))
	}
}
