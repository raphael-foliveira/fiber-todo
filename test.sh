#!/bin/bash

go test -coverprofile=c.out ./... &&
go tool cover -func=c.out;
