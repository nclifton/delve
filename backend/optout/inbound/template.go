package inbound

import (
	"html/template"
	"log"
	"net/http"
)

type OptOutLink struct {
	Sender string
	Link   string
}

func renderLink(w http.ResponseWriter, l OptOutLink) error {
	tmpl, err := template.New("link").Parse(`
		<html>
			<body>
				<center>
					<h3>I want to stop receiving messages from {{.Sender}}</h3>

					<form action="/{{.Link}}" method="POST">
						<input type="submit" style="font-size:16pt;" value="Opt Out"/>
					</form>
				</center>
			</body>
		</html>
	`)
	if err != nil {
		log.Println(err)
	}

	return tmpl.Execute(w, l)
}

func renderUnsubscribed(w http.ResponseWriter, l OptOutLink) error {
	tmpl, err := template.New("link").Parse(`
		<html>
			<body>
				<center>
					<h3>You have been unsubscribed from {{.Sender}}</h3>
				</center>
			</body>
		</html>
	`)
	if err != nil {
		log.Println(err)
	}

	return tmpl.Execute(w, l)
}
