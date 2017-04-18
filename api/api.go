package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
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

	r.HandleFunc("/", o.index)
	r.HandleFunc("/{id}/status", o.status)

	r.HandleFunc("/api/v1/summary", o.summaryHandler)
	r.HandleFunc("/pubsub/push", o.pushHandler)
	return o, nil

}

const devID = "520041000351353337353037"

func (o *options) status(w http.ResponseWriter, r *http.Request) {
	const (
		overlay = `
<html>
<meta name="viewport" content="width=device-width, initial-scale=1">
<style>
* {
font-size: 24px;
}
p {
  padding: 5px;
  border: solid 1px green;
}
.error {
  background-color: red;
  color: #fff;
}
</style>
<script>
function writeDate(str) {
  var d = new Date(JSON.parse(str));
  document.write(d.toLocaleString());
}
function round(float) {
  var rounded = Math.round(float*100)/100;
  document.write(rounded);
}
</script>
<body>
{{if .IsCharging}}
  <p>
    Charging done: <script>writeDate({{marshal .done}});</script> (<script>round({{marshal .battery.Current.Duration.Minutes}});</script>mins)
  </p>
{{else}}
  <p class="error">Not Charging</p>
{{end}}
  <p>
    Battery %: {{.battery.State.Percent}}<br/>
    Last Updated: <script>writeDate({{marshal .battery.State.LastSOCTime}});</script>
  </p>
<div>
  <h4>Detailed Stats</h4>
  <p>
    SOC: {{.status.LastSOC.Data}}<br/>
    Last Updated: <script>writeDate({{marshal .status.LastSOC.PublishTime}});</script>
  </p>
  <p>
    Power: {{.status.LastPower.Data}}<br/>
    Last Updated: <script>writeDate({{marshal .status.LastPower.PublishTime}});</script>
  </p>
  <p>
    Volts: {{.status.LastVolts.Data}}<br/>
    Last Updated: <script>writeDate({{marshal .status.LastVolts.PublishTime}});</script>
  </p>
  <p>
    Amps: {{.status.LastAmps.Data}}<br/>
    Last Updated: <script>writeDate({{marshal .status.LastAmps.PublishTime}});</script>
  </p>
  <p>
    Device ID: {{.status.DeviceID}}
  </p>
</div>
</body>
</html>
`
	)

	ctx := appengine.NewContext(r)
	ctx = configureContext(ctx)
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
			"/api/v1/summary",
			"/api/v1/car/id/laststatus",
			"/api/v1/car/id/chargerate",
			"/api/v1/car/id/battery",
		}
	)
	ctx := appengine.NewContext(r)
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
