package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/Models"
	"main/grpc_client"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	Models.InitEnv()
}

const (
	port = ":50051"
)

// server is used to implement helloworld.proto.GreeterServer.
/*type server struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.proto.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.HelloReply{
		Message: "Hello " + in.GetName(),
	}, nil
}*/


func main() {

	mediaMetadataClient := grpc_client.InitMediaMetadataGrpcClient()
	mediaMetadataClient.CreateNewMediaMetadata()

	/*worker := Worker.InitWorker()
	defer worker.RabbitMQ.Conn.Close()
	defer worker.RabbitMQ.Ch.Close()
	worker.Work()*/


	/*fmt.Println("GRPC");
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc_client.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}*/

}