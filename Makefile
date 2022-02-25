# Go parameters
GOCMD:=go
GORUN=$(GOCMD) run
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
MIGRATE_CONFIG:="./migration/dbconfig.yml"

run:
	$(GORUN) main.go

doc:
	swag init

test:
	$(GOCMD) clean -testcache
	$(GOTEST) ./...

testc:
	$(GOCMD) clean -testcache
	$(GOTEST) -cover ./...

migrate-up:
	sql-migrate up -config=$(MIGRATE_CONFIG) -env="localhost"

migrate-down:
	sql-migrate down -config=$(MIGRATE_CONFIG) -env="localhost"

seed-flush:
	# flush mysql
	docker exec mysql mysql -uroot -psecret -e \
	"SELECT CONCAT('TRUNCATE TABLE ', table_schema, '.', TABLE_NAME, ';') FROM INFORMATION_SCHEMA.TABLES \
	WHERE table_schema IN ('casbin_auth') AND TABLE_NAME != 'migrations'" | grep "casbin_auth*" | xargs -I {} docker exec mysql mysql -uroot -psecret -e {}
	# flush redis
	docker exec redis redis-cli flushall
	docker exec redis-cluster redis-cli -p 7000 flushall
	docker exec redis-cluster redis-cli -p 7001 flushall
	docker exec redis-cluster redis-cli -p 7002 flushall
	# exec seeder
	$(GORUN) ./pkg/seeder/main.go