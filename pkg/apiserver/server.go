package apiserver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type ApiServer struct {
	bind        string
	gameService *GameService
	listService *ListService
	router      *mux.Router
}

func NewApiServer(bind string, ls *ListService, gs *GameService) (*ApiServer, error) {
	return &ApiServer{
		bind:        bind,
		gameService: gs,
		listService: ls,
	}, nil
}

func (s *ApiServer) Listen() error {
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		return err
	}

	router := mux.NewRouter().StrictSlash(true)
	router.PathPrefix("/ui").Handler(http.StripPrefix("/ui", http.FileServer(http.Dir("C:/Users/Can/Documents/Projects/q3master/ui/dist/q3party/"))))
	router.HandleFunc("/games", s.listGames).Methods("GET")
	router.HandleFunc("/games/refresh", s.refreshGames).Methods("POST")
	router.HandleFunc("/lists/{id}", s.getList).Methods("GET")
	router.HandleFunc("/lists", s.newList).Methods("POST")
	router.HandleFunc("/lists/{id}/add", s.addToList).Methods("POST")
	router.HandleFunc("/lists/{id}/remove", s.removeFromList).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})

	return http.ListenAndServe(s.bind, handlers.CORS(headers, methods, origins)(router))
}
func (s *ApiServer) listGames(w http.ResponseWriter, r *http.Request) {
	log.Trace("Started listGames")
	defer log.Trace("Exited listGames")

	res, err := s.gameService.List()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, res)
}

func (s *ApiServer) refreshGames(w http.ResponseWriter, r *http.Request) {
	log.Trace("Started refreshGames")
	defer log.Trace("Exited refreshGames")

	err := s.gameService.Refresh()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, map[string]string{"result": "ok"})
}

func (s *ApiServer) getList(w http.ResponseWriter, r *http.Request) {
	log.Trace("Started getList")
	defer log.Trace("Exited getList")

	id := mux.Vars(r)["id"]
	res, err := s.listService.Get(id)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, res)
}

func (s *ApiServer) newList(w http.ResponseWriter, r *http.Request) {
	log.Trace("Started newList")
	defer log.Trace("Exited newList")

	res, err := s.listService.Create()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	writeResponse(w, res)
}

func (s *ApiServer) addToList(w http.ResponseWriter, r *http.Request) {
	log.Trace("Started addToList")
	defer log.Trace("Exited addToList")

	id := mux.Vars(r)["id"]
	body, err := getBody(r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if v, ok := body["server"]; ok {
		if err = s.listService.AddToList(id, v.(string)); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err = errors.New("Missing server field")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *ApiServer) removeFromList(w http.ResponseWriter, r *http.Request) {
	log.Trace("Started removeFromList")
	defer log.Trace("Exited removeFromList")

	id := mux.Vars(r)["id"]
	body, err := getBody(r)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if v, ok := body["server"]; ok {
		if err = s.listService.RemoveFromList(id, v.(string)); err != nil {
			log.Error(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		err = errors.New("Missing server field")
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getBody(r *http.Request) (map[string]interface{}, error) {
	defer r.Body.Close()

	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	var body map[string]interface{}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func writeResponse(w http.ResponseWriter, res interface{}) {
	d, err := json.Marshal(res)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(d)
	if err != nil {
		log.Error(err)
	}
}
