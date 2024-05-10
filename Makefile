include services/user/Makefile
include services/auth/Makefile
include services/article/Makefile
include services/bookmark/Makefile
include services/comment/Makefile
include services/gateway/Makefile

.PHONY: migrateup
migrateup:
	cd "services/user" && make user/migrateup
	cd "services/article" && make article/migrateup
	cd "services/bookmark" && make bookmark/migrateup
	cd "services/comment" && make comment/migrateup

.PHONY: migratedown
migratedown:
	cd "services/user" && make user/migratedown
	cd "services/article" && make article/migratedown
	cd "services/bookmark" && make bookmark/migratedown
	cd "services/comment" && make comment/migratedown

.PHONY: build
build:
	cd "services/user" && make user/build
	cd "services/auth" && make auth/build
	cd "services/article" && make article/build
	cd "services/bookmark" && make bookmark/build
	cd "services/comment" && make comment/build
	cd "services/gateway" && make gateway/build

.PHONY: up
up:
	docker-compose up