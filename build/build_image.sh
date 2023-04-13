#!/usr/bin/env bash

go install

buildah from --name puzzlepassstrengthserver-working-container scratch
buildah copy puzzlepassstrengthserver-working-container $HOME/go/bin/puzzlepassstrengthserver /bin/puzzlepassstrengthserver
buildah config --env SERVICE_PORT=50051 puzzlepassstrengthserver-working-container
buildah config --port 50051 puzzlepassstrengthserver-working-container
buildah config --entrypoint '["/bin/puzzlepassstrengthserver"]' puzzlepassstrengthserver-working-container
buildah commit puzzlepassstrengthserver-working-container puzzlepassstrengthserver
buildah rm puzzlepassstrengthserver-working-container

buildah push puzzlepassstrengthserver docker-daemon:puzzlepassstrengthserver:latest
