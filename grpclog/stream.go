package grpclog

import (
	"go.uber.org/zap"

	"google.golang.org/grpc"

	"github.com/golang/protobuf/proto"
)

type serverStream struct {
	grpc.ServerStream
	logger     *serverLogger
	fullMethod string
}

func newServerStream(fullMethod string, stream grpc.ServerStream, logger *serverLogger) *serverStream {
	return &serverStream{
		ServerStream: stream,
		logger:       logger,
		fullMethod:   fullMethod,
	}
}

func (s *serverStream) SendMsg(m interface{}) error {
	err := s.ServerStream.SendMsg(m)
	if err == nil && s.logger.logResponses {
		fields := []zap.Field{}
		if p, ok := m.(proto.Message); ok {
			fields = append(fields, zap.Object("grpc.response.content", &jsonpbObjectMarshaler{pb: p}))
		}
		s.logger.LogStream(s.Context(), s.fullMethod, "server response payload logged as grpc.response.content field", fields)
	}
	return err
}

func (s *serverStream) RecvMsg(m interface{}) error {
	err := s.ServerStream.RecvMsg(m)
	if err == nil && s.logger.logRequests {
		fields := []zap.Field{}
		if p, ok := m.(proto.Message); ok {
			fields = append(fields, zap.Object("grpc.request.content", &jsonpbObjectMarshaler{pb: p}))
		}
		s.logger.LogStream(s.Context(), s.fullMethod, "server request payload logged as grpc.request.content field", fields)
	}
	return err
}