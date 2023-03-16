package geecache

import "geecache/pb"

// PeerPicker is the interface that must be implemented to locate
// the peer that owns a specific key.
// PeerPicker 的 PickPeer() 方法用于根据传入的 key 选择相应节点 PeerGetter
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter is the interface that must be implemented by a peer.
// Get() 方法用于从对应 group 查找缓存值。PeerGetter 就对应于上述流程中的 HTTP 客户端。
type PeerGetter interface {
	Get(in *pb.Request, out *pb.Response) error
}
