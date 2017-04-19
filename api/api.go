package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/drewwells/chargerstore"
	"github.com/drewwells/chargerstore/math"
	"github.com/drewwells/chargerstore/store"
	"github.com/drewwells/chargerstore/types"
	"github.com/gorilla/mux"
)

var newContext = func(r *http.Request) context.Context {
	return appengine.NewContext(r)
}

var configureContext = func(ctx context.Context) context.Context {
	return ctx
}

type options struct {
	*mux.Router
}

func New() (http.Handler, error) {

	r := mux.NewRouter()
	o := &options{
		Router: r,
	}

	r.HandleFunc("/api", o.index)
	r.HandleFunc("/{id}/status", o.tmplStatus)

	r.HandleFunc("/api/v1/summary", o.summaryHandler)
	r.HandleFunc("/api/v1/status", o.statusHandler)
	r.HandleFunc("/pubsub/push", o.pushHandler)

	// Static files
	r.PathPrefix("/public").Handler(http.StripPrefix("/public", http.FileServer(http.Dir("../public/"))))

	return o, nil

}

const devID = "520041000351353337353037"

func loadTmpl(ctx context.Context, path string) string {
	absPath, err := filepath.Abs(filepath.Join("..", "public", path))
	if err != nil {
		log.Errorf(ctx, "can not find file %s: %s", absPath, err)
		return ""
	}
	bs, err := ioutil.ReadFile(absPath)
	if err != nil {
		log.Errorf(ctx, "failed to read file %s: %s", absPath, err)
	}
	return string(bs)
}

func (o *options) statusHandler(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	// TODO: read deviceid from url or account

	stat, err := store.GetCarStatus(ctx, devID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp := make(map[string]interface{})
	resp["amps"] = stat.LastAmps
	resp["volts"] = stat.LastVolts
	resp["soc"] = stat.LastSOC
	resp["power"] = stat.LastPower
	resp["charge"] = math.BatteryCharging(
		stat.LastSOC,
		stat.LastPower,
	)
	if bc.Current.Duration > 0 {
		resp["done"] = cd
	}

	marshal(w, resp)
}

func (o *options) tmplStatus(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)

	overlay := loadTmpl(ctx, "index.html")
	// TODO: read deviceid from url or account
	stat, err := store.GetCarStatus(ctx, devID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Infof(ctx, "%#v\n", stat)

	bc := math.BatteryCharging(
		stat.LastSOC,
		stat.LastPower,
	)

	var cd time.Time
	if bc.Current.Duration > 0 {
		cd = time.Now().Add(bc.Current.Duration)
	}

	m := map[string]interface{}{
		"isCharging": stat.LastPower.Data > 0,
		"status":     stat,
		"battery":    bc,
		"done":       cd,
	}

	funcMap := template.FuncMap{
		// The name "title" is what the function will be called in the template text.
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(template.JS(a))
		},
	}
	overlayTmpl, err := template.New("overlay").
		Funcs(funcMap).Parse(overlay)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}
	if err := overlayTmpl.Execute(w, m); err != nil {
		log.Errorf(ctx, err.Error())
	}
}

func (o *options) index(w http.ResponseWriter, r *http.Request) {
	const (
		overlay = `
<html>
<meta name="viewport" content="width=device-width, initial-scale=1">
<body>
<ul>{{range .}}<li><a href="{{.}}">{{.}}</li>{{end}}</ul>
</body>
</html>
`
	)
	var (
		// funcs     = template.FuncMap{"join": strings.Join}
		guardians = []string{
			"/id/status",
			"/api/v1/status",
		}
	)
	ctx := newContext(r)
	overlayTmpl, err := template.New("overlay").Parse(overlay)
	if err != nil {
		log.Errorf(ctx, err.Error())
	}
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
