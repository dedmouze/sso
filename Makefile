run local:
	go run .\cmd\sso\main.go --env=env/local.env
proxy:
	go run .\cmd\proxy\main.go --config=config/local.yaml
migrate:
	go run .\cmd\migrator --storage-path=./storage/sso.db --migrations-path=./migrations