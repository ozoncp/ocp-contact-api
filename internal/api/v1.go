package api

import (
	"context"

	desc "github.com/ozoncp/ocp-contact-api/pkg/ocp-contact-api"

	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type contactApiServer struct {
	log zerolog.Logger
	desc.UnimplementedOcpContactApiServer
}

func NewOcpContactApiServer(log zerolog.Logger) desc.OcpContactApiServer {
	return &contactApiServer{log: log}
}

func (s *contactApiServer) ListContactsV1(
	context context.Context,
	req *desc.ListContactsV1Request,
) (*desc.ListContactsV1Response, error) {
	s.log.Info().Msg("list")
	return &desc.ListContactsV1Response{}, nil
}

func (s *contactApiServer) DescribeContactV1(
	context context.Context,
	req *desc.DescribeContactV1Request,
) (*desc.DescribeContactV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.log.Info().Msgf("describe contact with id: %v", req.ContactId)
	return &desc.DescribeContactV1Response{}, nil
}

func (s *contactApiServer) CreateContactV1(
	context context.Context,
	req *desc.CreateContactV1Request,
) (*desc.CreateContactV1Response, error) {
	s.log.Info().Msgf("create contact with userId: %v, type: %v, text: %v",
		req.UserId, req.Type, req.Text)
	return &desc.CreateContactV1Response{}, nil
}

func (s *contactApiServer) RemoveContactV1(
	context context.Context,
	req *desc.RemoveContactV1Request,
) (*desc.RemoveContactV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	s.log.Info().Msgf("remove contact with id: %v", req.ContactId)
	return &desc.RemoveContactV1Response{}, nil
}
