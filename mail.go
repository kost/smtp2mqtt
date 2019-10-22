package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	"github.com/paulrosania/go-charset/charset"
	qprintable "github.com/sloonz/go-qprintable"
)

func hasEncoding(word string) bool {
	return strings.Contains(word, "=?") && strings.Contains(word, "?=")
}

func parseSubject(subject string) string {
	if !hasEncoding(subject) {
		return subject
	}

	dec := mime.WordDecoder{}
	sub, _ := dec.DecodeHeader(subject)
	return sub
}

var headerSplitter = []byte("\r\n\r\n")

// parseBody will accept a a raw body, break it into all its parts and then convert the
// message to UTF-8 from whatever charset it may have.
func parseBody(header mail.Header, body []byte) (html []byte, text []byte, isMultipart bool, err error) {
	var mediaType string
	var params map[string]string
	mediaType, params, err = mime.ParseMediaType(header.Get("Content-Type"))
	if err != nil {
		return
	}
	partDisp := header.Get("Content-Disposition")

	if strings.HasPrefix(mediaType, "multipart/") {
		isMultipart = true
		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				break
			}

			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				// error and no results to use
				if len(slurp) == 0 {
					break
				}
			}

			partMediaType, partParams, err := mime.ParseMediaType(p.Header.Get("Content-Type"))
			if err != nil {
				break
			}

			partDisposition := p.Header.Get("Content-Disposition")
			dispo := strings.ToLower(partDisposition)

			if !strings.Contains(dispo, "attachment") {
				var htmlT, textT []byte
				htmlT, textT, err = parsePart(partMediaType, dispo, partParams["charset"], p.Header.Get("Content-Transfer-Encoding"), slurp)
				if len(htmlT) > 0 {
					html = htmlT
				} else {
					text = textT
				}
			}
		}
	} else {

		splitBody := bytes.SplitN(body, headerSplitter, 2)
		if len(splitBody) < 2 {
			isMultipart = false
			text = body
			return
		}

		body = splitBody[1]
		html, text, err = parsePart(mediaType, partDisp, params["charset"], header.Get("Content-Transfer-Encoding"), body)
	}
	return
}

func parsePart(mediaType, partDisposition, charsetStr, encoding string, part []byte) (html, text []byte, err error) {
	// deal with charset
	if strings.ToLower(charsetStr) == "iso-8859-1" {
		var cr io.Reader
		cr, err = charset.NewReader("latin1", bytes.NewReader(part))
		if err != nil {
			return
		}

		part, err = ioutil.ReadAll(cr)
		if err != nil {
			return
		}
	}

	// deal with encoding
	var body []byte
	switch strings.ToLower(encoding) {
	case "quoted-printable":
		dec := qprintable.NewDecoder(qprintable.WindowsTextEncoding, bytes.NewReader(part))
		body, err = ioutil.ReadAll(dec)
		if err != nil {
			return
		}
	case "base64":
		decoder := base64.NewDecoder(base64.StdEncoding, bytes.NewReader(part))
		body, err = ioutil.ReadAll(decoder)
		if err != nil {
			return
		}
	default:
		body = part
	}

	// deal with media type
	mediaType = strings.ToLower(mediaType)
	switch {
	case strings.Contains(mediaType, "text/html"):
		html = body
	case strings.Contains(mediaType, "text/plain"):
		text = body
	}
	return
}
