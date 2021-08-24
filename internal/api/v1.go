package api

import (
	"context"
	"github.com/ozoncp/ocp-contact-api/internal/models"
	"github.com/ozoncp/ocp-contact-api/internal/producer"
	"github.com/ozoncp/ocp-contact-api/internal/repo"
	"github.com/ozoncp/ocp-contact-api/internal/utils"
	desc "github.com/ozoncp/ocp-contact-api/pkg/ocp-contact-api"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type contactApiServer struct {
	repo repo.Repo
	prod producer.Producer
	batchSize int
	log zerolog.Logger
	desc.UnimplementedOcpContactApiServer
}

func NewOcpContactApiServer(
	repo repo.Repo,
	prod producer.Producer,
	batchSize int,
	log zerolog.Logger,
) desc.OcpContactApiServer {
	return &contactApiServer{repo: repo, prod: prod, batchSize: batchSize, log: log}
}

func (s *contactApiServer) ListContactsV1(
	context context.Context,
	req *desc.ListContactsV1Request,
) (*desc.ListContactsV1Response, error) {
	repoContacts, err := s.repo.ListContacts(context, req.Limit, req.Offset)
	if err != nil {
		log.Error().Err(err).Msg("getting list contacts from the repo failed")
		return nil, status.Error(codes.Internal, "getting list contacts from the repo failed")
	}

	s.log.Info().Msgf("request list contacts: %v", len(repoContacts))

	contacts := make([]*desc.Contact, 0, len(repoContacts))
	for _, repoContact := range repoContacts {
		contact := &desc.Contact{
			Id:     repoContact.Id,
			UserId: repoContact.UserId,
			Type:   repoContact.Type,
			Text:   repoContact.Text,
		}

		contacts = append(contacts, contact)
	}

	response := &desc.ListContactsV1Response{
		Contacts: contacts,
	}

	return response, nil
}

func (s *contactApiServer) DescribeContactV1(
	context context.Context,
	req *desc.DescribeContactV1Request,
) (*desc.DescribeContactV1Response, error) {
	s.log.Info().Msgf("describe contact with id: %v", req.ContactId)
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	contact, err := s.repo.DescribeContact(context, req.ContactId)
	if err != nil {
		log.Error().Err(err).Msg("describe contact failed")
		return nil, status.Error(codes.NotFound, err.Error())
	}
	s.log.Info().Msgf("describe contact with id %v was successfully done", req.ContactId)
	response := &desc.DescribeContactV1Response{
		Contact: &desc.Contact{
			Id:     contact.Id,
			UserId: contact.UserId,
			Type:   contact.Type,
			Text:   contact.Text,
		},
	}

	return response, nil
}

func (s *contactApiServer) CreateContactV1(
	context context.Context,
	req *desc.CreateContactV1Request,
) (*desc.CreateContactV1Response, error) {

	s.log.Info().Msgf("creating contact with userId: %v, type: %v, text: %v",
		req.UserId, req.Type, req.Text)

	contact := models.Contact{
		UserId: req.UserId,
		Type:   req.Type,
		Text:   req.Text,
	}

	contactId, err := s.repo.CreateContact(context, contact)

	if err != nil {
		s.log.Error().Err(err).Msg("create contact failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	response := &desc.CreateContactV1Response{
		ContactId: contactId,
	}

	s.log.Info().Msgf("contact was created with id: %v", contactId)

	event := producer.EventMessage{
		Id:        contact.Id,
		Action:    producer.Create.String(),
		Timestamp: time.Now().Unix(),
	}
	msg := producer.CreateMessage(producer.Create, event)
	if err = s.prod.Send(msg); err != nil {
		s.log.Error().Err(err).Msgf("failed send message to kafka")
	}

	return response, nil
}

func (s *contactApiServer) RemoveContactV1(
	context context.Context,
	req *desc.RemoveContactV1Request,
) (*desc.RemoveContactV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.repo.RemoveContact(context, req.ContactId)

	if err != nil {
		log.Error().Err(err).Msgf("remove contact with id %v failed", req.ContactId)
		return &desc.RemoveContactV1Response{
			Result: false,
		}, nil
	}

	s.log.Info().Msgf("remove contact with id %v was removed", req.ContactId)

	event := producer.EventMessage{
		Id:        req.ContactId,
		Action:    producer.Remove.String(),
		Timestamp: time.Now().Unix(),
	}
	msg := producer.CreateMessage(producer.Remove, event)
	if err := s.prod.Send(msg); err != nil {
		s.log.Error().Err(err).Msgf("failed send message to kafka")
	}

	return &desc.RemoveContactV1Response{Result: true}, nil
}

func (s *contactApiServer) UpdateContactV1(
	ctx context.Context,
	req *desc.UpdateContactV1Request,
) (*desc.UpdateContactV1Response, error) {
	contact := models.Contact{Id: req.Contact.Id, UserId: req.Contact.UserId,
		Type: req.Contact.Type, Text: req.Contact.Text}

	if err := s.repo.UpdateContact(ctx, contact); err != nil {
		log.Error().Err(err).Msgf("update contact with id %v failed", req.Contact.Id)
		return &desc.UpdateContactV1Response{Updated: false}, err
	}

	event := producer.EventMessage{
		Id:        contact.Id,
		Action:    producer.Update.String(),
		Timestamp: time.Now().Unix(),
	}
	msg := producer.CreateMessage(producer.Update, event)
	if err := s.prod.Send(msg); err != nil {
		s.log.Error().Err(err).Msgf("failed send message to kafka")
	}

	return &desc.UpdateContactV1Response{Updated: true}, nil
}

func (s *contactApiServer) MultiCreateContactsV1(
	ctx context.Context,
	req *desc.MultiCreateContactsV1Request,
) (*desc.MultiCreateContactsV1Response, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	contacts := make([]models.Contact, 0, len(req.Contacts))

	for _, req := range req.Contacts {
		contact := models.Contact{UserId: req.UserId, Type: req.Type, Text: req.Text}
		contacts = append(contacts, contact)
	}

	batches, err := utils.Split(contacts, s.batchSize)
	if err != nil {
		log.Error().Err(err).Msgf("multiple contacts creation failed")
		return nil, status.Error(codes.Internal, err.Error())
	}

	var count uint64
	for _, batch := range batches {
		if err := s.repo.AddContacts(ctx, batch); err != nil {
			log.Error().Err(err).Msgf("multiple contacts creation failed while adding in repo")
			return nil, status.Error(codes.Internal, err.Error())
		}
		count += uint64(len(batch))
	}
	return &desc.MultiCreateContactsV1Response{
		Count: count,
	}, nil
}