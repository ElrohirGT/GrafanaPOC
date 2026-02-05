#!/usr/bin/env bash

set -exu

kubectl apply -f . --server-side
