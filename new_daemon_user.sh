#!/bin/sh
# Quickly add a deamon user

adduser --disabled-password --disabled-login --shell /usr/sbin/nologin $1
