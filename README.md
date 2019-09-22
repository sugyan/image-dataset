# image-dataset

## development

```sh
export GOOGLE_APPLICATION_CREDENTIALS=<path/to/serviceAccountKey.json>
export SESSION_KEY=<secret key>
go run web/main.go
```


## deployment

```sh
gcloud app deploy web
```
