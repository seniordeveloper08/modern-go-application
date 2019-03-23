package tododriver

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	todov1beta1 "github.com/sagikazarmark/modern-go-application/.gen/api/proto/todo/v1beta1"
	"github.com/sagikazarmark/modern-go-application/internal/todo"
)

type grpcServer struct {
	createTodo grpctransport.Handler
	listTodos  grpctransport.Handler
	markAsDone grpctransport.Handler
}

// MakeGRPCServer makes a set of endpoints available as a gRPC server.
func MakeGRPCServer(endpoints Endpoints, errorHandler todo.ErrorHandler) todov1beta1.TodoListServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerFinalizer(func(ctx context.Context, err error) {
			if err != nil {
				errorHandler.Handle(err)
			}
		}),
	}

	return &grpcServer{
		createTodo: grpctransport.NewServer(
			endpoints.Create,
			decodeCreateTodoGRPCRequest,
			encodeCreateTodoGRPCResponse,
			options...,
		),
		listTodos: grpctransport.NewServer(
			endpoints.List,
			decodeListTodosGRPCRequest,
			encodeListTodosGRPCResponse,
			options...,
		),
		markAsDone: grpctransport.NewServer(
			endpoints.MarkAsDone,
			decodeMarkAsDoneGRPCRequest,
			encodeMarkAsDoneGRPCResponse,
			options...,
		),
	}
}

func (s *grpcServer) CreateTodo(
	ctx context.Context,
	req *todov1beta1.CreateTodoRequest,
) (*todov1beta1.CreateTodoResponse, error) {
	_, rep, err := s.createTodo.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*todov1beta1.CreateTodoResponse), nil
}

func decodeCreateTodoGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*todov1beta1.CreateTodoRequest)

	return createTodoRequest{
		Text: req.GetText(),
	}, nil
}

func encodeCreateTodoGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(createTodoResponse)

	return &todov1beta1.CreateTodoResponse{
		Id: resp.ID,
	}, nil
}

func (s *grpcServer) ListTodos(
	ctx context.Context,
	req *todov1beta1.ListTodosRequest,
) (*todov1beta1.ListTodosResponse, error) {
	_, rep, err := s.listTodos.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*todov1beta1.ListTodosResponse), nil
}

func decodeListTodosGRPCRequest(_ context.Context, _ interface{}) (interface{}, error) {
	return nil, nil
}

func encodeListTodosGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(listTodosResponse)

	grpcResp := &todov1beta1.ListTodosResponse{
		Todos: make([]*todov1beta1.Todo, len(resp.Todos)),
	}

	for i, t := range resp.Todos {
		grpcResp.Todos[i] = &todov1beta1.Todo{
			Id:   t.ID,
			Text: t.Text,
			Done: t.Done,
		}
	}

	return grpcResp, nil
}

func (s *grpcServer) MarkAsDone(
	ctx context.Context,
	req *todov1beta1.MarkAsDoneRequest,
) (*todov1beta1.MarkAsDoneResponse, error) {
	_, rep, err := s.markAsDone.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*todov1beta1.MarkAsDoneResponse), nil
}

func decodeMarkAsDoneGRPCRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*todov1beta1.MarkAsDoneRequest)

	return markAsDoneRequest{
		ID: req.GetId(),
	}, nil
}

func encodeMarkAsDoneGRPCResponse(_ context.Context, response interface{}) (interface{}, error) {
	if f, ok := response.(endpoint.Failer); ok && f.Failed() != nil {
		err := f.Failed()
		code := codes.Internal

		if e, ok := err.(*todoError); ok && e.Code() == codeNotFound {
			code = codes.NotFound
		}

		return nil, status.Error(code, err.Error())
	}

	return &todov1beta1.MarkAsDoneResponse{}, nil
}
