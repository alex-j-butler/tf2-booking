package steamauth

import (
	"log"
	"net/http"
	"reflect"

	"fmt"

	"sync"

	"alex-j-butler.com/tf2-booking/config"
	"alex-j-butler.com/tf2-booking/database"
	"github.com/gorilla/mux"
	"github.com/yohcop/openid-go"
)

type HTTPServer struct {
	Address string
	Port    int
	RootURL string

	handlersMu     sync.RWMutex
	handlers       map[interface{}][]reflect.Value
	nonceStore     *openid.SimpleNonceStore
	discoveryCache *openid.SimpleDiscoveryCache
}

// Create a new instance of the HTTP Steam authentication server.
func New(address string, port int, rootURL string) *HTTPServer {
	httpServer := &HTTPServer{
		Address: address,
		Port:    port,
		RootURL: rootURL,

		nonceStore:     openid.NewSimpleNonceStore(),
		discoveryCache: openid.NewSimpleDiscoveryCache(),
	}

	return httpServer
}

func (s *HTTPServer) Run() {
	m := mux.NewRouter()

	m.HandleFunc("/auth/{secret}", s.AuthHandler)
	m.HandleFunc("/steam/{secret}", s.CallbackHandler)

	log.Println(fmt.Sprintf("Listening on %s:%d", s.Address, s.Port))
	http.ListenAndServe(fmt.Sprintf("%s:%d", s.Address, s.Port), m)
}

func (s *HTTPServer) AuthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret := vars["secret"]

	var authSecret database.AuthSecret
	if err := database.DB.Where("secret = ?", secret).First(&authSecret).Error; err != nil {
		// Secret does not exist.
		fmt.Fprintln(w, "Invalid URL, please try linking your account again.")
		return
	}

	if url, err := openid.RedirectURL(
		"http://steamcommunity.com/openid",
		fmt.Sprintf("%s/steam/%s", s.RootURL, secret),
		s.RootURL,
	); err == nil {
		go s.handle(&LinkAttemptEvent{
			Secret:    authSecret.Secret,
			DiscordID: authSecret.DiscordID,
		})

		http.Redirect(w, r, url, 303)
	} else {
		log.Println(err)
	}
}

func (s *HTTPServer) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret := vars["secret"]
	fullURL := config.Conf.SteamAuthServer.RootURL + r.URL.String()

	var authSecret database.AuthSecret
	if err := database.DB.Where("secret = ?", secret).First(&authSecret).Error; err != nil {
		// Secret does not exist.
		fmt.Fprintln(w, "Invalid URL, please try linking your account again.")
		return
	}

	id, err := openid.Verify(
		fullURL,
		s.discoveryCache,
		s.nonceStore,
	)

	if err != nil {
		log.Println(err)

		go s.handle(&LinkFailureEvent{
			Secret:    authSecret.Secret,
			DiscordID: authSecret.DiscordID,
		})
	} else {
		var steamID string
		fmt.Sscanf(id, "http://steamcommunity.com/openid/id/%s", &steamID)

		go s.handle(&LinkSuccessEvent{
			Secret:    authSecret.Secret,
			DiscordID: authSecret.DiscordID,
			SteamID:   steamID,
		})

		fmt.Fprintf(w, "Thanks, the SteamID %s has now been linked to the Discord account %s! You can close this page now.", steamID, authSecret.DiscordID)
		log.Println(id)
	}
}

func (s *HTTPServer) AddHandler(handler interface{}) func() {
	s.initialise()

	eventType := s.validateHandler(handler)

	s.handlersMu.Lock()
	defer s.handlersMu.Unlock()

	h := reflect.ValueOf(handler)

	s.handlers[eventType] = append(s.handlers[eventType], h)

	return func() {
		s.handlersMu.Lock()
		defer s.handlersMu.Unlock()

		handlers := s.handlers[eventType]
		for i, v := range handlers {
			if h == v {
				s.handlers[eventType] = append(handlers[:i], handlers[i+1:]...)
			}
		}
	}
}

func (s *HTTPServer) initialise() {
	s.handlersMu.Lock()
	if s.handlers != nil {
		s.handlersMu.Unlock()
		return
	}

	s.handlers = make(map[interface{}][]reflect.Value)
	s.handlersMu.Unlock()
}

func (s *HTTPServer) validateHandler(handler interface{}) reflect.Type {
	handlerType := reflect.TypeOf(handler)

	if handlerType.NumIn() != 2 {
		panic("Unable to add event handler, handler must be of type func(*steamauth.HTTPServer, *steamauth.EventType)")
	}

	if handlerType.In(0) != reflect.TypeOf(s) {
		panic("Unable to add event handler, first argument must be of type *steamauth.HTTPServer")
	}

	eventType := handlerType.In(1)

	if eventType.Kind() == reflect.Interface {
		eventType = nil
	}

	return eventType
}

func (s *HTTPServer) handle(event interface{}) {
	s.handlersMu.RLock()
	defer s.handlersMu.RUnlock()

	if s.handlers == nil {
		return
	}

	handlerParameters := []reflect.Value{reflect.ValueOf(s), reflect.ValueOf(event)}

	if handlers, ok := s.handlers[nil]; ok {
		for _, handler := range handlers {
			go handler.Call(handlerParameters)
		}
	}

	if handlers, ok := s.handlers[reflect.TypeOf(event)]; ok {
		for _, handler := range handlers {
			go handler.Call(handlerParameters)
		}
	}
}

// Old stuff

/*
var nonceStore = openid.NewSimpleNonceStore()
var discoveryCache = openid.NewSimpleDiscoveryCache()

func Initialise() {
	m := mux.NewRouter()

	m.HandleFunc("/auth/{secret}", AuthHandler)
	m.HandleFunc("/steam/{secret}", CallbackHandler)
	log.Println(fmt.Sprintf("Listening on %s:%d", config.Conf.SteamAuthServer.Address, config.Conf.SteamAuthServer.Port))
	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Conf.SteamAuthServer.Address, config.Conf.SteamAuthServer.Port), m)
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret := vars["secret"]

	var authSecret database.AuthSecret
	if err := database.DB.Where("secret = ?", secret).First(&authSecret).Error; err != nil {
		// Secret does not exist.
		fmt.Fprintln(w, "Invalid URL, please try linking your account again.")
		return
	}

	if url, err := openid.RedirectURL(
		"http://steamcommunity.com/openid",
		fmt.Sprintf("%s/steam/%s", config.Conf.SteamAuthServer.RootURL, secret),
		config.Conf.SteamAuthServer.RootURL,
	); err == nil {
		http.Redirect(w, r, url, 303)
	} else {
		log.Println(err)
	}
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	secret := vars["secret"]
	fullURL := config.Conf.SteamAuthServer.RootURL + r.URL.String()

	var authSecret database.AuthSecret
	if err := database.DB.Where("secret = ?", secret).First(&authSecret).Error; err != nil {
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

		fmt.Fprintf(w, "Thanks, the SteamID %s has now been linked to the Discord account %s! You can close this page now.", steamID, authSecret.DiscordID)
		log.Println(id)
	}
}
*/
