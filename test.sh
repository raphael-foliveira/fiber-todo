#!/bin/bash

go test -coverprofile=c.out ./pkg/todo ./pkg/common &&
go tool cover -func=c.out;
