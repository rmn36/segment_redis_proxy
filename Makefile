# HELP
# This will output the help for each task
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help test

help: ## This is helpful
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

build: ## Build the container
	docker-compose build --no-cache redis-proxy

run: ## Build and run the container
	docker-compose up -d redis-proxy 

stop: ## Stop and remove a running container
	docker-compose down -v

rm: stop ## Stop and remove running containers
	docker-compose rm redis-proxy

rmi: stop rm ## Removes image
	docker rmi segment_redis_proxy

clean: rmi ## Cleans out Docker artifacts

test: ## Runs test env	
	cd scripts/ ; ./test.sh

test_down:
	cd test/ ; docker-compose down -v