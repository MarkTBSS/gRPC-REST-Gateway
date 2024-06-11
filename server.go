package main

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/MarkTBSS/gRPC-REST-Gateway/proto_src/proto_dst/email"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

type server struct {
	email.UnimplementedEmailServiceServer
	db *pgx.Conn
}

func (s *server) SendEmail(ctx context.Context, req *email.EmailRequest) (*email.EmailResponse, error) {
	// Implement your logic to send an email or insert into PostgreSQL
	// Assuming you are inserting into the database for the example
	_, err := s.db.Exec(ctx, "INSERT INTO emails (recipient, subject, body) VALUES ($1, $2, $3)",
		req.Recipient, req.Subject, req.Body)
	if err != nil {
		log.Printf("Failed to insert: %v", err)
		return nil, err
	}

	return &email.EmailResponse{
		Recipient: req.Recipient,
		Subject:   req.Subject,
		Body:      req.Body,
	}, nil
}

func main() {
	go runGRPCServer()
	runHTTPServer()
}

func runGRPCServer() {
	// Connect to PostgreSQL
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:Pass1234@localhost:5432/postgres")
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	// Create a gRPC server object
	s := grpc.NewServer()
	email.RegisterEmailServiceServer(s, &server{db: conn})

	// Register reflection service on gRPC server.
	reflection.Register(s)

	// Create a listener on TCP port 8080
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Serve gRPC server
	log.Println("Serving gRPC on 0.0.0.0:8080")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func runHTTPServer() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:  true,
				UseEnumNumbers: false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		}),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}

	err := email.RegisterEmailServiceHandlerFromEndpoint(ctx, mux, "localhost:8080", opts)
	if err != nil {
		log.Fatalf("Failed to register gRPC-Gateway: %v", err)
	}

	log.Println("Serving gRPC-Gateway on http://0.0.0.0:8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
	}
}
