#!/bin/sh

#Using protoc version 3.0.0
cd gauge-proto
protoc --go_out=../gauge_messages/ api.proto api_v2.proto messages.proto spec.proto
