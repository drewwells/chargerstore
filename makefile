cloudserve:
	cd app;	GCLOUD_PROJECT=particle-volt; PUBSUB_TOPIC=CAR; go run *.go

serve:
	dev_appserver.py app/app.yaml

deploy:
	gcloud app deploy app/app.yaml app/index.yaml
	#cd app; gcloud app deploy index.yaml

cloudpkgandupload:
	#cd app; GOOS=linux GOARCH=amd64 go build -i -v -o ../dist/app
	xgo --targets=linux/amd64 github.com/drewwells/chargerstore/app
	tar -c -f dist/bundle.tar app-linux-amd64
	tar -c -f dist/bundle.tar scripts/startup.sh
	gsutil cp dist/bundle.tar particle-volt

clouddeploy: #pkgandupload
	gcloud compute instances create my-app-instance \
    --image-family=debian-8 \
    --machine-type=g1-small \
    --scopes userinfo-email,cloud-platform \
    --metadata-from-file startup-script=scripts/startup.sh \
    --zone us-central1-f \
    --tags http-server
