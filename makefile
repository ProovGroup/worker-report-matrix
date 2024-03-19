REGION = eu-west-1
ENV ?= test
SERVICE_NAME = worker-report-matrix

lambda: clean
	@env GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -o ./bin/bootstrap ./cmd/worker/worker.go
	@mkdir ./bin/assets && cp -r ./assets/* ./bin/assets
	@cd ./bin && zip $(SERVICE_NAME).zip ./assets/* bootstrap

deploy: lambda
	$(eval VERSION = $(shell aws lambda update-function-code --function-name ${SERVICE_NAME} --region ${REGION} --zip-file fileb://bin/$(SERVICE_NAME).zip --publish |  python -c "import sys, json; print json.load(sys.stdin)['Version']"))
	@aws lambda update-alias --function-name ${SERVICE_NAME} --region ${REGION} --name ${ENV} --function-version $(VERSION)

clean:
	@echo "clean ${SERVICE_NAME}"
	@rm -rf bin
