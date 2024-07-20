.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/RegisterTrip/bootstrap ./RegisterTrip
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/ListTrips/bootstrap ./ListTrips
	env GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o ./bin/TripIA/bootstrap ./TripIA

clean:
	rm -rf ./bin ./vendor go.sum

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
