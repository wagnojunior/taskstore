# CURL commands for testing

## POST request to create a new store
curl -d '{"name":"store_1"}' -H 'Content-Type: application/json' http://localhost:3030/store

## POST request to create a new task
curl -d '{"store_id":1, "text":"test_task_1", "tags":["test_tag_1"], "due":"2021-02-18T21:54:42.123Z"}' -H 'Content-Type: application/json' http://localhost:3030/task

## GET request to get a task by ID
curl -v http://localhost:3030/task/1/2
