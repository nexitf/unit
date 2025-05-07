package call

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// StaticMethod specifies that a call is being made
// to a method that is static, which means the method is known at compile time
// and doesn't change at runtime. This can be used as a signal to stats plugins
// that this method is safe to include as a key to a measurement.
func StaticMethod(fn func(err error)) grpc.CallOption {
	return grpc.StaticMethod()
}

// Header retrieves the header metadata for a unary RPC.
func Header(md *metadata.MD) grpc.CallOption {
	return grpc.Header(md)
}

// Trailer retrieves the trailer metadata for a unary RPC.
func Trailer(md *metadata.MD) grpc.CallOption {
	return grpc.Trailer(md)
}

// Peer retrieves peer information for a unary RPC.
// The peer field will be populated *after* the RPC completes.
func Peer(peer *peer.Peer) grpc.CallOption {
	return grpc.Peer(peer)
}

// WaitForReady configures the action to take when an RPC is attempted on broken
// connections or unreachable servers. If waitForReady is false and the
// connection is in the TRANSIENT_FAILURE state, the RPC will fail
// immediately. Otherwise, the RPC client will block the call until a
// connection is available (or the call is canceled or times out) and will
// retry the call if it fails due to a transient error.  gRPC will not retry if
// data was written to the wire unless the server indicates it did not process
// the data.
func WaitForReady(waitForReady bool) grpc.CallOption {
	return grpc.WaitForReady(waitForReady)
}

// MaxCallRecvMsgSize sets the maximum message size
// in bytes the client can receive. If this is not set, gRPC uses the default
// 4MB.
func MaxCallRecvMsgSize(bytes int) grpc.CallOption {
	return grpc.MaxCallRecvMsgSize(bytes)
}

// MaxCallSendMsgSize sets the maximum message size
// in bytes the client can send. If this is not set, gRPC uses the default
// `math.MaxInt32`.
func MaxCallSendMsgSize(bytes int) grpc.CallOption {
	return grpc.MaxCallSendMsgSize(bytes)
}

// PerRPCCredentials sets credentials.PerRPCCredentials for a call.
func PerRPCCredentials(creds credentials.PerRPCCredentials) grpc.CallOption {
	return grpc.PerRPCCredentials(creds)
}

// CallContentSubtype sets the content-subtype
// for a call. For example, if content-subtype is "json", the Content-Type over
// the wire will be "application/grpc+json". The content-subtype is converted
// to lowercase before being included in Content-Type.
func CallContentSubtype(contentSubtype string) grpc.CallOption {
	return grpc.CallContentSubtype(contentSubtype)
}
