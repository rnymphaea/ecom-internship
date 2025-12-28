IMAGE := server
PORT := 8080
API_URL := http://localhost:$(PORT)

.PHONY: all build run test api-test

all: build run

build:
	docker build -t $(IMAGE) .

run:
	docker run --rm -p $(PORT):$(PORT) $(IMAGE)

test:
	go test ./... -v

api-test:
	@echo "=== Testing API endpoints ==="
	@echo ""
	@echo "1. GET /todos (empty)"
	@curl -v $(API_URL)/todos
	@echo ""
	@echo "2. POST /todos (create)"
	@curl -v -X POST $(API_URL)/todos -H "Content-Type: application/json" -d '{"caption":"Test","description":"Test todo"}'
	@echo ""
	@echo "3. GET /todos (list)"
	@curl -v $(API_URL)/todos
	@echo ""
	@echo "4. GET /todos/1"
	@curl -v $(API_URL)/todos/1
	@echo ""
	@echo "5. POST /todos (empty caption)"
	@curl -v -X POST $(API_URL)/todos -H "Content-Type: application/json" -d '{"caption":"","description":"Empty"}'
	@echo ""
	@echo "6. GET /todos/999 (not found)"
	@curl -v $(API_URL)/todos/999
	@echo ""
	@echo "7. PUT /todos/1 (update existing - success)"
	@curl -v -X PUT $(API_URL)/todos/1 -H "Content-Type: application/json" -d '{"caption":"Updated caption","description":"Updated description","is_completed":true}'
	@echo ""
	@echo "8. PUT /todos/999 (update non-existent)"
	@curl -v -X PUT $(API_URL)/todos/999 -H "Content-Type: application/json" -d '{"caption":"Not exists","description":"Wont work","is_completed":false}'
	@echo ""
	@echo "9. PUT /todos/1 (update with empty caption - should fail)"
	@curl -v -X PUT $(API_URL)/todos/1 -H "Content-Type: application/json" -d '{"caption":"","description":"Empty caption","is_completed":false}'
	@echo ""
	@echo "10. GET /todos"
	@curl -v $(API_URL)/todos
	@echo ""
	@echo ""
	@echo "11. DELETE /todos/1 (delete existing)"
	@curl -v -X DELETE $(API_URL)/todos/1
	@echo ""
	@echo "12. DELETE /todos/2 (delete non-existent)"
	@curl -v -X DELETE $(API_URL)/todos/2
	@echo ""
	@echo "13. Final check: GET /todos (should be empty)"
	@curl -v $(API_URL)/todos
	@echo ""
	@echo "=== API test completed ==="

