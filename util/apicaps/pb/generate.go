package moby_devkit_v1_apicaps //nolint:revive

//go:generate protoc -I=. -I=../../../vendor/ -I=../../../../../../ --gogo_out=plugins=grpc:. caps.proto
