#!/bin/bash
#
# Generate a new X.509 certificate and private key.
# Source:
# http://www.reddit.com/r/linux/comments/2nkyou/til_one_liner_to_generate_a_selfsigned/

################################################################################

DAYS=3650
COMMON_NAME="test.example.com"
KEYFILE="secret.key"
CERTFILE="Public.crt"

openssl req -new -newkey rsa:4096 -days "$(DAYS)" -nodes -x509 -utf8 -sha256 -subj "/CN=$(COMMON_NAME)" -keyout "$(KEYFILE)" -out "$(CERTFILE)"
