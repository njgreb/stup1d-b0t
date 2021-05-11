#!/bin/bash

rm stup1d-b0t

go build .

docker build . -t njgreb/stup1d-b0t:$1

docker push njgreb/stup1d-b0t:$1