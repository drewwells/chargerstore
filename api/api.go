package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/drewwells/chargerstore"
	"github.com/drewwells/chargerstore/types"
	"github.com/gorilla/mux"
)

type options struct {
	*mux.Router
}

func New() (http.Handler, error) {

	r := mux.NewRouter()
	o := &options{
		Router: r,
	}

	r.HandleFunc("/", o.index)

	r.HandleFunc("/api/v1/summary", o.summaryHandler)
	r.HandleFunc("/pubsub/push", o.pushHandler)
	return o, nil

}

func (o *options) index(w http.ResponseWriter, r *http.Request) {
	const (
		//master  = `Names:{{block "list" .}}{{"\n"}}{{range .}}{{println "-" .}}{{end}}{{end}}`
		overlay = `
<html>
<body>
<ul>{{range .}}<li><a href="/api/v1/{{.}}">{{.}}</li>{{end}}</ul>
</body>
</html>
`
	)
	var (
		// funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{
			"summary",
			"car/id/laststatus",
			"car/id/chargerate",
			"car/id/battery",
		}
	)
	ctx := appengine.NewContext(r)
	// masterTmpl, err := template.New("master").Funcs(funcs).Parse(master)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	overlayTmpl, err := template.New("overlay").Parse(overlay)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}
	// if err := masterTmpl.Execute(os.Stdout, guardians); err != nil {
	// 	log.Fatal(err)
	// }
	if err := overlayTmpl.Execute(w, guardians); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func (o *options) summaryHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "This worker has processed %d events.", chargerstore.Count())
}

func (o *options) pushHandler(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	var req types.PushRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		err := fmt.Errorf("Could not decode body: %v", err)
		log.Errorf(ctx, "%s", err)
		// http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Infof(ctx, "request % #v\n", req)
	msg, err := chargerstore.Process(ctx, req.Message)
	if err != nil {
		err := fmt.Errorf("failed to process msg: %s", err)
		log.Errorf(ctx, "%s", err)
		// http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Infof(ctx, "processed message from: %s", msg.DeviceID)
}
