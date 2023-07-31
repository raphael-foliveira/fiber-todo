#!/bin/bash

go test -coverprofile=c.out ./pkg/todo ./pkg/common &&
go tool cover -html=c.out;
