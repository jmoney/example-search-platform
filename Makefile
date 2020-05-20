all: build

build: datastore

datastore:
	$(MAKE) -C datastore

run: datastore
	docker-compose up

clean:
	docker-compose down