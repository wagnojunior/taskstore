# CURL commands for testing

## POST request to create a new store
curl -d '{"name":"store_1"}' -H 'Content-Type: application/json' http://localhost:3030/store

## POST request to create a new task
curl -d '{"store_id":1, "text":"test_task_1", "tags":["test_tag_1"], "due":"2021-02-18"}' -H 'Content-Type: application/json' http://localhost:3030/task

## POST request to delete a task
curl -d '{"store_id":1, "id":3}' -H 'Content-Type: application/json' http://localhost:3030/delete/1/3

## GET request to get a task by ID
curl -v http://localhost:3030/task/1/1

## GET request to get all tasks from a store
curl -v http://localhost:3030/task/1


# Goose commands

## Run migrations
goose postgres "host=localhost port=5432 user=admin password=admin_taskstore dbname=taskstore sslmode=disable" up


# Database

## PSQL
docker compose exec -it db psql -U admin -d taskstore