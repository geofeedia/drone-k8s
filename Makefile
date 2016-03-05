PROJECT=drone-k8s
MAJOR_VERSION=1
VERSION=1.0.0

all: build tags

build:
	docker build --no-cache -f Dockerfile.build -t quay.io/geofeedia/${PROJECT}-build:${VERSION} .
	docker run --rm quay.io/geofeedia/${PROJECT}-build:${VERSION} | docker build -t quay.io/geofeedia/${PROJECT}:${VERSION} -
	docker rmi -f quay.io/geofeedia/${PROJECT}-build:${VERSION}

tags:
	docker tag quay.io/geofeedia/${PROJECT}:${VERSION} quay.io/geofeedia/${PROJECT}:${MAJOR_VERSION}
