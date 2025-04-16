package grpc

import (
	"context"

	"ipcheck/internal/geoip"
	pb "ipcheck/internal/grpc/ipcheckpb"

	"google.golang.org/grpc"
)

type grpcServer struct {
	pb.UnimplementedIpCheckerServer
	geo *geoip.Server
}

func Register(s *grpc.Server, g *geoip.Server) {
	pb.RegisterIpCheckerServer(s, &grpcServer{geo: g})
}

func (s *grpcServer) CheckIP(ctx context.Context, req *pb.CheckRequest) (*pb.CheckResponse, error) {
	ok, country, err := s.geo.CheckIP(req.GetIp(), req.GetAllowedCountries())
	if err != nil {
		return nil, err
	}
	return &pb.CheckResponse{Allowed: ok, Country: country}, nil
}
