package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"unicode"
	"unicode/utf8"
)

var (
	errDecodeError = errors.New("UTF-8 decode error")
)

type buffer struct {
	out   bytes.Buffer
	src   []byte
	index int
}

func newBuffer(src []byte) buffer {
	return buffer{
		out:   bytes.Buffer{},
		src:   src,
		index: 0,
	}
}

func (b *buffer) nextRune() (rune, error) {
	r, size := utf8.DecodeRune(b.src[b.index:])

	if r == utf8.RuneError {
		switch size {
		case 0:
			return r, io.EOF
		case 1:
			return r, errDecodeError
		}
	}
	b.index += size

	return r, nil
}

func (b *buffer) skipComment() error {
	for {
		r, err := b.nextRune()
		if err != nil {
			return err
		}
		_, err = b.out.WriteRune(r)
		if err != nil {
			return err
		}

		if r == '*' {
			r, err := b.nextRune()
			if err != nil {
				return err
			}
			_, err = b.out.WriteRune(r)
			if err != nil {
				return err
			}

			if r == '/' {
				return nil
			}
		}
	}
}

func (b *buffer) skipSection(opening, closing, escape rune,
	hasEscape, allowNesting, copyToOut bool) error {

	depth := 1
	for {
		r, err := b.nextRune()
		if err != nil {
			return err
		}
		if copyToOut {
			_, err = b.out.WriteRune(r)
			if err != nil {
				return err
			}
		}

		switch r {
		case escape:
			if hasEscape {
				r, err := b.nextRune()
				if err != nil {
					return err
				}
				if copyToOut {
					_, err = b.out.WriteRune(r)
					if err != nil {
						return err
					}
				}
			}

		case closing:
			if allowNesting {
				depth--
			}
			if !allowNesting || depth == 0 {
				return nil
			}

		case opening:
			if allowNesting {
				depth++
			}
		}
	}
}

func pegFormat(src []byte) ([]byte, error) {
	b := newBuffer(src)
	indent := 0

	for {
		r, err := b.nextRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		_, err = b.out.WriteRune(r)
		if err != nil {
			return nil, err
		}

		switch r {
		case '\n':
			indent = 0
		case '\t':
			indent++
		case '/':
			r, err = b.nextRune()
			if err != nil {
				return nil, err
			}
			_, err = b.out.WriteRune(r)
			if err != nil {
				return nil, err
			}
			switch r {
			case '/':
				err = b.skipSection('/', '\n', 0, false, false, true)
			case '*':
				err = b.skipComment()
			}
		case '\'':
			err = b.skipSection('\'', '\'', '\\', true, false, true)
		case '"':
			err = b.skipSection('"', '"', '\\', true, false, true)
		case '`':
			err = b.skipSection('`', '`', 0, false, false, true)
		case '[':
			err = b.skipSection('[', ']', '\\', true, false, true)
		case '{':
			begin := b.index
			err = b.skipSection('{', '}', 0, false, true, false)
			end := b.index

			section := b.src[begin-1 : end]
			if !bytes.ContainsRune(section, '\n') {
				content := b.src[begin : end-1]
				contentWithoutLeftSpaces :=
					bytes.TrimLeftFunc(content, unicode.IsSpace)
				contentWithoutRightSpaces :=
					bytes.TrimRightFunc(content, unicode.IsSpace)
				leftSpaceBytes :=
					len(content) - len(contentWithoutLeftSpaces)
				rightSpaceBytes :=
					len(content) - len(contentWithoutRightSpaces)

				originalContent :=
					content[leftSpaceBytes : len(content)-rightSpaceBytes]
				formattedContent, err := format.Source(originalContent)
				if err != nil {
					return nil, err
				}

				leftSpaces := content[:leftSpaceBytes]
				rightSpacesAndClosingBrace :=
					section[len(section)-rightSpaceBytes-1:]
				_, err = b.out.Write(leftSpaces)
				if err != nil {
					return nil, err
				}
				_, err = b.out.Write(formattedContent)
				if err != nil {
					return nil, err
				}
				_, err = b.out.Write(rightSpacesAndClosingBrace)
				if err != nil {
					return nil, err
				}
				break
			}

			formatted, err := format.Source(section)
			if err != nil {
				_, err = b.out.Write(b.src[begin:end])
				if err != nil {
					return nil, err
				}
			} else {
				formatted = formatted[1:]
				pattern := []byte{'\n'}
				replacement := append([]byte{'\n'},
					bytes.Repeat([]byte{'\t'}, indent)...)
				formatted = bytes.Replace(formatted, pattern, replacement, -1)
				_, err = b.out.Write(formatted)
				if err != nil {
					return nil, err
				}
			}
		}

		if err != nil {
			return nil, err
		}
	}

	return b.out.Bytes(), nil
}

func main() {
	name := filepath.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s file\n", name)
		os.Exit(0)
	}

	pegFile := os.Args[1]
	pegSource, err := ioutil.ReadFile(pegFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		os.Exit(1)
	}

	pegOutput, err := pegFormat(pegSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: format: %v\n", name, err)
		os.Exit(1)
	}

	_, err = os.Stdout.Write(pegOutput)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
	}
}
