package grpc

import (
	"context"

	pb "github.com/xvnvdu/config-analyzer/api"
	"github.com/xvnvdu/config-analyzer/internal/checker"
	"github.com/xvnvdu/config-analyzer/internal/parser"
)

type Server struct {
	pb.UnimplementedAnalyzerServer
	checker *checker.Checker
}

func New(c *checker.Checker) *Server {
	return &Server{checker: c}
}

func (s *Server) Analyze(_ context.Context, req *pb.AnalyzeRequest) (*pb.AnalyzeResponse, error) {
	cfg, err := parser.Parse(req.Config)
	if err != nil {
		return nil, err
	}
	results := s.checker.CheckConfig(cfg)

	var issues []*pb.Issue
	for _, result := range results[0].Issues {
		issues = append(issues, &pb.Issue{
			Level:          string(result.Level),
			Message:        result.Message,
			Recommendation: result.Recommendation,
			RuleName:       result.RuleName,
		})
	}
	return &pb.AnalyzeResponse{Issues: issues}, nil
}
