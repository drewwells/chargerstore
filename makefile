serve:
	cd app;	GCLOUD_PROJECT=particle-volt; PUBSUB_TOPIC=CAR; go run *.go
	

deploy:
	cd app; gcloud app deploy
