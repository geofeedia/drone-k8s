PROJECT=drone-k8s
VERSION=1.0.0

all: build

build:
	docker build -f Dockerfile.build -t quay.io/geofeedia/${PROJECT}-build:${VERSION} .
	docker run --rm quay.io/geofeedia/${PROJECT}-build:${VERSION} | docker build -t quay.io/geofeedia/${PROJECT}:${VERSION} -
	docker rmi -f quay.io/geofeedia/${PROJECT}-build:${VERSION}
