all: build

build: datastore indexer

datastore:
	$(MAKE) -C datastore

indexer:
	$(MAKE) -C indexer

run: datastore indexer
	docker-compose up

clean:
	docker-compose down