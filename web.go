package main

import (
	"github.com/golang/oauth2"
	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"html/template"
	"log"
	"net/http"
	"github.com/abourget/go.rice"
	"github.com/codegangsta/negroni"
)

type Webapp struct {
	config *WebappConfig
	store  *sessions.CookieStore
}

type WebappConfigSection struct {
	Webapp WebappConfig
}

type WebappConfig struct {
	ClientID          string `json:"client_id"`
	ClientSecret      string `json:"client_secret"`
	RestrictDomain    string `json:"restrict_domain"`
	SessionAuthKey    string `json:"session_auth_key"`
	SessionEncryptKey string `json:"session_encrypt_key"`
}

var notAuthenticatedTemplate = template.Must(template.New("").Parse(`
<html><body>
You have currently not given permissions to access your data. Please authenticate this app with the Google OAuth provider.
<form action="/authorize" method="POST"><input type="submit" value="Ok, authorize this app with my id"/></form>
</body></html>
`))

var userInfoTemplate = template.Must(template.New("").Parse(`
<html><body>
This app is now authenticated to access your Google user info.  Your details are:<br />
{{.}}<br />
Name: {{.Name}}<br />
Email: {{.Email}} Verified: {{.VerifiedEmail}}<br />
Hd: {{.Hd}}<br />
</body></html>
`))

func launchWebapp() {
	var conf WebappConfigSection
	bot.LoadConfig(&conf)

	web = &Webapp{
		config: &conf.Webapp,
		store:  sessions.NewCookieStore([]byte(conf.Webapp.SessionAuthKey), []byte(conf.Webapp.SessionEncryptKey)),
	}

	var err error
	oauthCfg, err = oauth2.NewConfig(
		&oauth2.Options{
			ClientID:     conf.Webapp.ClientID,
			ClientSecret: conf.Webapp.ClientSecret,
			RedirectURL:  "http://localhost:8080/oauth2callback",
			Scopes:       []string{"openid", "profile", "email", "https://www.googleapis.com/auth/userinfo.profile"},
		},
		"https://accounts.google.com/o/oauth2/auth",
		"https://accounts.google.com/o/oauth2/token",
	)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.Handle("/static/", http.FileServer(rice.MustFindBox("static").HTTPBox()))
	// Should we put the /authorize link back ??
	mux.HandleFunc("/", handleRoot)

	n := negroni.Classic()
	n.UseHandler(NewOAuthMiddleware(context.ClearHandler(mux)))

	n.Run("localhost:8080")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	profile, _ := checkAuth(r)
	webappCode.Execute(w, profile)
	//notAuthenticatedTemplate.Execute(w, nil)
}