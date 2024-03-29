package tb

import (
	"context"

	proto "github.com/dovydas1928/tb/pkg/proto"
)

// NewCommonServer creates a new common server using the common interface provided
func NewCommonServer(common Common) proto.CommonServer {
	return &commonServer{common: common}
}

type commonServer struct {
	common Common
}

func (s *commonServer) GetVersion(ctx context.Context, _ *proto.Void) (*proto.ResponseVersion, error) {
	version, revision, err := s.common.GetVersion(ctx)
	if err != nil {
		return nil, err
	}

	return &proto.ResponseVersion{Version: version, Revision: revision}, nil
}

func (s *commonServer) Modprobe(ctx context.Context, req *proto.RequestModprobe) (*proto.Void, error) {
	return &proto.Void{}, s.common.Modprobe(ctx, req.Module)
}
