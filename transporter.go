package raft

//------------------------------------------------------------------------------
//
// Typedefs
//
//------------------------------------------------------------------------------

// Transporter is the interface for allowing the host application to transport
// requests to other nodes.
type Transporter interface {
	SendVoteRequest(server *Server, peer *Peer, req *RequestVoteRequest) (*RequestVoteResponse, error)
	SendAppendEntriesRequest(server *Server, peer *Peer, req *AppendEntriesRequest) (*AppendEntriesResponse, error)
	SendSnapshotRequest(server *Server, peer *Peer, req *SnapshotRequest) (*SnapshotResponse, error)
}
