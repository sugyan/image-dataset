# image-dataset

## development

```sh
export GOOGLE_APPLICATION_CREDENTIALS=<path/to/serviceAccountKey.json>
export SESSION_KEY=<secret key>
go run web/main.go
```


## deployment

```sh
go run cmd/generate_index/main.go > firestore.indexes.json
firebase deploy --only firestore:indexes
gcloud app deploy web
```
