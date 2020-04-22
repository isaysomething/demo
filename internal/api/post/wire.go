package post

import "github.com/google/wire"

var Set = wire.NewSet(
	New,
	NewService,
)
