include .env

up:
	@echo "Starting containers..."
	docker-compose up --build -d --remove-orphans

down:
	@echo "Stopping containers..."
	docker-compose down

build:
	go build -o ${BINARY} ./api/

start:
	@env PORT=${PORT} MONGO_DB_USERNAME=${MONGO_DB_USERNAME} MONGO_DB_PASSWORD=${MONGO_DB_PASSWORD} JWT_SECRET=${JWT_SECRET} ./${BINARY}

restart: build start

# connection string
# mongodb://admin:password@localhost:27017/todos_db?authSource=admin&readPreference=primary&appname=MongDB%20Compass&directConnection=true&ssl=false
# mongodb://admin:password@localhost:65510/todos_db?authSource=admin&readPreference=primary&appname=MongDB%20Compass&directConnection=true&ssl=false