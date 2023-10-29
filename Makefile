ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

deploy:
	docker run -e GIT_SHA=$(GITHUB_SHA) -e REGISTRY -e USER -e PASS -e IMAGE -e PUBLISH -v $(ROOT_DIR)kubernetes/manifests/apps/:/manifests/ -v $(ROOT_DIR)kubernetes/templates/:/templates/ $(IMAGE)/tools

deploy-dry-run:
	docker run -e GIT_SHA=$(GITHUB_SHA) -v $(ROOT_DIR)kubernetes/manifests/apps/:/manifests/ -v $(ROOT_DIR)kubernetes/templates/:/templates/ $(IMAGE)/tools

build-image:
	docker build --tag $(IMAGE)/tools . -f "$(ROOT_DIR)hack/Dockerfile.tools"

deploy-image: build-image
	docker push $(IMAGE)/tools
