package mm7utils

import (
	"encoding/xml"
	"fmt"
)

// XMLLayout XML structure
type XMLLayout struct {
	InnerXML string `xml:",innerxml"`
}

// XMLHead XML structure
type XMLHead struct {
	Layout XMLLayout `xml:"layout"`
}

// XMLSmilBody XML structure
type XMLSmilBody struct {
	InnerXML string `xml:",innerxml"`
}

// XMLSmil XML structure
type XMLSmil struct {
	XMLName xml.Name    `xml:"smil"`
	Head    XMLHead     `xml:"head"`
	Body    XMLSmilBody `xml:"body"`
}

type SMILMedia struct {
	ContentID string
	MediaType string
}

func renderSMIL(attachments []SMILMedia) ([]byte, error) {
	var err error
	var medias string
	var regions string

	for i, attachment := range attachments {
		id := fmt.Sprintf("%s-%d", attachment.MediaType, i)
		medias += fmt.Sprintf(`<par><%s src="cid:%s" region="%s"/></par>`, attachment.MediaType, attachment.ContentID, id)
		regions += fmt.Sprintf(`<region id="%s" top="50%%" left="0" height="50%%" width="100%%" fit="hidden"/>`, id)
	}

	smil := XMLSmil{
		Head: XMLHead{
			Layout: XMLLayout{
				InnerXML: `<root-layout width="100%" height="100%"/>` + regions,
			},
		},
		Body: XMLSmilBody{
			InnerXML: medias,
		},
	}

	body, err := xml.Marshal(smil)
	if err != nil {
		return nil, err
	}

	body = append([]byte(xml.Header), body...)

	return body, err
}
