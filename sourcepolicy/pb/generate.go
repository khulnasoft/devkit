package moby_devkit_v1_sourcepolicy //nolint:revive

//go:generate protoc -I=. --gogofaster_out=plugins=grpc:. policy.proto
