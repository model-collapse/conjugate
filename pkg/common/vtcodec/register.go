// Package vtcodec registers the vtprotobuf gRPC codec.
//
// Import this package with a blank identifier to replace the default
// proto codec with vtprotobuf's faster MarshalVT/UnmarshalVT path:
//
//	import _ "github.com/conjugate/conjugate/pkg/common/vtcodec"
package vtcodec

import (
	"google.golang.org/grpc/encoding"

	vtgrpc "github.com/planetscale/vtprotobuf/codec/grpc"
)

func init() {
	encoding.RegisterCodec(vtgrpc.Codec{})
}
