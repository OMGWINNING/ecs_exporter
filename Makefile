# Ensure that 'all' is the default target otherwise it will be the first target from Makefile.common.
all::

# Needs to be defined before including Makefile.common to auto-generate targets
DOCKER_ARCHS ?= amd64 arm64
DOCKER_REPO  ?= 209202477790.dkr.ecr.us-east-1.amazonaws.com/ecs_exporter

include Makefile.common

DOCKER_IMAGE_NAME       ?= ecs-exporter
