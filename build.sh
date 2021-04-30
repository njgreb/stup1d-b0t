#!/bin/bash

rm stup1d-b0t

go build .

docker build . -t njgreb/stup1d-b0t:testone

docker push njgreb/stup1d-b0t:testone