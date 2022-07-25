# Comic Parser
React and Go app to retrieve and parse comic meta-data.

## Running the application
This application can be run using the go tool or by running it on a container.

### Run it in a container
- Create a ```.env``` file using the **example.env** as a template.
- Install [docker](https://docs.docker.com/engine/install/)
- Compile the the **Dockerfile** ```docker build -t comic-parser:latest .```
- Run the container ``` docker run --name comic-parser -e API_PORT=8000 -e COMIC_HOST=https://xkcd.com -p 8000:8000 comic-parser:latest ```
- Alternatively install [docker-compose](https://docs.docker.com/compose/install/) and run it ``` docker-compose up --build ```

### Run it with GO
- Create the **API_PORT** and **COMIC_HOST** environment variables
- Install [Go](https://go.dev/doc/install)
- Run it directly ```go run main.go```