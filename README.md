# image-dataset

## face images

### collect images from Twitter

```sh
export CONSUMER_KEY='*********************'
export CONSUMER_SECRET='****************************************'
go run cmd/collect_twitter_images/main.go -screen_name <Twitter ScreenName> > python/data.tsv
```


### detect faces and save metadata

```sh
cd python
pip install -r requirements.txt
# download a trained facial shape predictor for dlib
python create_data.py data.tsv
```


### upload data

```sh
DEVELOPMENT=true go run cmd/upload_images/*.go -datadir python/data -projectID <Project ID>
```


## development

```sh
gcloud beta emulators firestore --project <Project ID> start --host-port :8081
```

```sh
export FIRESTORE_EMULATOR_HOST=localhost:8081
export GOOGLE_CLOUD_PROJECT=<Project ID>
export GOOGLE_APPLICATION_CREDENTIALS=<path/to/serviceAccountKey.json>
export SESSION_KEY=<secret key>
go run web/main.go
```

```sh
cd frontend
npm start
```


## deployment

```sh
go run cmd/generate_index/main.go > firestore.indexes.json
firebase deploy --only firestore:indexes
gcloud app deploy web
```

## dump images

```sh
go run cmd/dump_data/main.go -projectID <Project ID> -num 10000
```
