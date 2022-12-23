package handler

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"

	pb "github.com/finest08/go-connect/internal/gen/messaging/v1"
	pbcnn "github.com/finest08/go-connect/internal/gen/messaging/v1/messagingv1connect"
	"github.com/finest08/go-connect/store"
)

type MessagingServer struct {
	Store store.Storer
	pbcnn.UnimplementedMessagingServiceHandler
}

func (s MessagingServer) Create(ctx context.Context, req *connect.Request[pb.CreateRequest]) (*connect.Response[pb.CreateResponse], error) {
	reqMsg := req.Msg
	msg := reqMsg.MessageThread

	msg.Id = uuid.NewString()
	msg.Messages[0].Date = time.Now().Unix()

	// store functions
	if err := s.Store.CreateMsg(ctx, msg); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.CreateResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) Query(ctx context.Context, req *connect.Request[pb.QueryRequest]) (*connect.Response[pb.QueryResponse], error) {
	reqMsg := req.Msg

	if reqMsg.SearchText != "" {
		pattern, err := regexp.Compile(`^[a-zA-Z@. ]+$`)
		if err != nil {
			return nil, connect.NewError(connect.CodeAborted, err)
		}
		if !pattern.MatchString(reqMsg.SearchText) {
			return nil, connect.NewError(connect.CodeInvalidArgument,
				errors.New("invalid search text format"))
		}
	}

	cur, mat, err := s.Store.QueryMsg(ctx, reqMsg)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	rsp := &pb.QueryResponse{
		Cursor:  cur,
		Matches: mat,
	}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) Get(ctx context.Context, req *connect.Request[pb.GetRequest]) (*connect.Response[pb.GetResponse], error) {
	reqMsg := req.Msg.MessageId

	// store functions
	msg, err := s.Store.GetMsg(ctx, reqMsg)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.GetResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) Update(ctx context.Context, req *connect.Request[pb.UpdateRequest]) (*connect.Response[pb.UpdateResponse], error) {
	reqMsg := req.Msg
	msg := reqMsg.MessageThread

	msg.Messages[len(msg.Messages)-1].Date = time.Now().Unix()

	// store functions
	if err := s.Store.UpdateMsg(reqMsg.MessageId, ctx, msg); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.UpdateResponse{MessageThread: msg}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) CreateUser(ctx context.Context, req *connect.Request[pb.CreateUserRequest]) (*connect.Response[pb.CreateUserResponse], error) {
	reqMsg := req.Msg
	usr := reqMsg.User

	usr.Id = uuid.NewString()
	usr.Date = time.Now().Unix()

	// store functions
	if err := s.Store.CreateUser(ctx, usr); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.CreateUserResponse{User: usr}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) GetUser(ctx context.Context, req *connect.Request[pb.GetUserRequest]) (*connect.Response[pb.GetUserResponse], error) {
	id := req.Msg.UserId

	// store functions
	usr, err := s.Store.GetUser(ctx, id)
	if err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.GetUserResponse{UserId: usr}
	return connect.NewResponse(rsp), nil
}

func (s MessagingServer) DeleteUser(ctx context.Context, req *connect.Request[pb.DeleteUserRequest]) (*connect.Response[pb.DeleteUserResponse], error) {
	reqMsg := req.Msg
	id := reqMsg.UserId

	// store functions
	if err := s.Store.DeleteUser(id); err != nil {
		return nil, connect.NewError(connect.CodeAborted, err)
	}

	// response
	rsp := &pb.DeleteUserResponse{}
	return connect.NewResponse(rsp), nil
}

