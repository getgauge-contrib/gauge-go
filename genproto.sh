#!/bin/sh

#Using protoc version 3.0.0
cd gauge-proto
protoc --go_out=../gauge_messages/ api.proto messages.proto spec.proto
