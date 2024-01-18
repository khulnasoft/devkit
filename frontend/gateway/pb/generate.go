package khulnasoft_devkit_v1_frontend //nolint:revive

//go:generate protoc -I=. -I=../../../vendor/ -I=../../../../../../ --gogo_out=plugins=grpc:. gateway.proto
