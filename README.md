# postal-service

Rest server for [libpostal](https://github.com/openvenues/libpostal) written in [Go](https://go.dev/).

## build

To build the docker image run:

```console
$ docker build --tag postal-service .
...
```

## run

To run the docker image run:

```console
$ docker run --rm -it  -p 9876:9876 postal-service
2024/04/30 04:27:23 INFO Starting server port=9876
```
