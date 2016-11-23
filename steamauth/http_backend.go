package steamauth

import (
	"log"
	"net/http"

	"github.com/yohcop/openid-go"

	"fmt"

	"github.com/gorilla/mux"
)

var nonceStore = openid.NewSimpleNonceStore()
var discoveryCache = openid.NewSimpleDiscoveryCache()

// Maps secret keys to Discord users.
var secretMap map[string]DiscordUser

type DiscordUser struct {
	ID   string
	Name string
}

func Initialise() {
	m := mux.NewRouter()

	// Add example secret.
	secretMap = make(map[string]DiscordUser)
	secretMap["example"] = DiscordUser{ID: "123", Name: "Alex"}

	m.HandleFunc("/auth/{secret}", AuthHandler)
	m.HandleFunc("/steam/{secret}", CallbackHandler)
	http.ListenAndServe(":3002", m)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret := vars["secret"]

	if _, ok := secretMap[secret]; !ok {
		// Secret does not exist.
		fmt.Fprintln(w, "Invalid URL, please try linking your account again.")
		return
	}

	if url, err := openid.RedirectURL(
		"http://steamcommunity.com/openid",
		fmt.Sprintf("http://localhost:3002/steam/%s", secret),
		"http://localhost:3002/",
	); err == nil {
		http.Redirect(w, r, url, 303)
	} else {
		log.Println(err)
	}
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret := vars["secret"]
	fullURL := "http://localhost:3002" + r.URL.String()

	if _, ok := secretMap[secret]; !ok {
		// Secret does not exist.
		fmt.Fprintln(w, "Invalid URL, please try linking your account again.")
		return
	}

	id, err := openid.Verify(
		fullURL,
		discoveryCache,
		nonceStore,
	)

	if err != nil {
		log.Println(err)
	} else {
		var steamID string
		fmt.Sscanf(id, "http://steamcommunity.com/openid/id/%s", &steamID)

		user, _ := secretMap[secret]

		fmt.Fprintf(w, "Thanks, the SteamID %s has now been linked to the Discord account %s! You can close this page now.", steamID, user.Name)
		log.Println(id)
	}
}
