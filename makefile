cloudserve:
	cd app;	GCLOUD_PROJECT=particle-volt; PUBSUB_TOPIC=CAR; go run *.go

reflex:
	go get github.com/cespare/reflex

watch: reflex
	reflex -s -t 50ms -r 'go$$' make local

local:
	cd app; go build -i -v && ./app

start:
	yarn start

serve:
	dev_appserver.py app/app.yaml

deploy:
	yarn build
	rm -rf app/build
	mv build app/build
	gcloud app deploy app/app.yaml app/index.yaml
	#cd app; gcloud app deploy index.yaml

export:
	/opt/google-cloud-sdk/platform/google_appengine/appcfg.py download_data -A s~particle-volt --url=http://particle-volt.appspot.com/_ah/remote_api/ --filename=data.csv

import:
