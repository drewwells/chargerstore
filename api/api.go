package api

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"google.golang.org/appengine"
	"google.golang.org/appengine/log"

	"github.com/drewwells/chargerstore"
	"github.com/drewwells/chargerstore/math"
	"github.com/drewwells/chargerstore/store"
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
<style>
p {
  padding: 5px;
  border: solid 1px green;
}
</style>
<script>
function writeDate(d) {
  var str = d.getHours() + ':' + d.getMinutes() + ':' + d.getSeconds();
  document.write(str);
}
function round(float) {
  var rounded = Math.round(float*100)/100;
  document.write(rounded);
}
</script>
<body>
<p>
Battery %: {{.battery.State.Percent}}<br/>
Last Updated: <script>writeDate(new Date('{{.battery.State.LastSOCTime}}'));</script>
</p>
<p>
Current Charging done in: <script>round({{.battery.Current.Minutes}});</script> mins
</p>
<p>
Battery %: {{.battery.State.Percent}}
</p>

  <p>
    Device ID: {{.status.DeviceID}}
  </p>
  <p>
    SOC: {{.status.LastSOC.Data}}<br/>
    Last Updated: <script>writeDate(new Date('{{.status.LastSOC.PublishTime}}'));</script>
  </p>
  <p>
    Volts: {{.status.LastVolts.Data}}<br/>
    Last Updated: <script>writeDate(new Date('{{.status.LastVolts.PublishTime}}'));</script>
  </p>
  <p>
    Amps: {{.status.LastAmps.Data}}<br/>
    Last Updated: <script>writeDate(new Date('{{.status.LastAmps.PublishTime}}'));</script>
  </p>
</body>
</html>
`
	)

	ctx := appengine.NewContext(r)
	// TODO: read deviceid from url or account
	stat, err := store.GetCarStatus(ctx, devID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	m := map[string]interface{}{
		"status": stat,
		"battery": math.BatteryCharging(
			stat.LastSOC,
			stat.LastAmps,
			stat.LastVolts,
		)}

	overlayTmpl, err := template.New("overlay").Parse(overlay)
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
