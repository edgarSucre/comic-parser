test:
	export API_PORT=8000 COMIC_HOST=https://xkcd.com/ && go test ./...

build:
	docker-compose build --no-cache

run:
	docker-compose up