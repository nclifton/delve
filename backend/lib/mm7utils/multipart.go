package mm7utils

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"strings"
)

type MMSPart struct {
	ContentType string
	ContentID   string
	Body        []byte
}

func ProcessMultiPart(header string, body io.Reader) ([]*MMSPart, error) {
	var parts []*MMSPart
	contentType, params, err := mime.ParseMediaType(header)
	if err != nil {
		log.Printf("could not parse media type of header: %s", err)
		return parts, errors.New("Not a multipart request")
	}

	// If we are not multipart then do nothing
	if !strings.HasPrefix(contentType, "multipart/") {
		return parts, nil
	}

	mpr := multipart.NewReader(body, params["boundary"])

	for {
		part, err := mpr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("unexpected error when retrieving a part of the message: %s", err)
			break
		}
		defer func() {
			err := part.Close()
			if err != nil {
				log.Printf("unexpected error when closing a part of the message: %s", err)
			}
		}()

		if strings.HasPrefix(part.Header.Get(`Content-Type`), "multipart/") {
			subparts, err := ProcessMultiPart(part.Header.Get(`Content-Type`), part)
			if err != nil {
				log.Printf("failed to process the subpart: %s", err)
				break
			}
			parts = append(parts, subparts...)
		} else {
			partBytes, err := ioutil.ReadAll(part)
			if err != nil {
				log.Printf("failed to read content of the part: %s %s", part.Header.Get(`Content-ID`), err)
				break
			}

			parts = append(parts, &MMSPart{
				ContentType: part.Header.Get(`Content-Type`),
				ContentID:   part.Header.Get(`Content-ID`),
				Body:        partBytes,
			})

		}
	}

	return parts, nil
}
