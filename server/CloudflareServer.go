package main

import (
	pb "cfTest/cloudflareApi/cache/purge"
	"context"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"net"

	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

type server struct {
	pb.UnimplementedPurgeCloudflareServer
}

func (s *server) PurgeCloudflare(ctx context.Context, in *pb.PurgeRequestCloudflare) (*pb.PurgeReplyCloudflare, error) {
	log.Printf("purgeTask")
	log.Printf("apiKey Received: %v", in.GetApiKey())
	log.Printf("apiEmail Received: %v", in.GetApiEmail())
	log.Printf("zoneId Received: %v", in.GetZoneId())
	log.Printf("purgeList Received: %v", in.GetPurgeList())

	// Construct a new API object using a global API key
	api, err := cloudflare.New(in.GetApiKey(), in.GetApiEmail())
	// alternatively, you can use a scoped API token
	// api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
		return &pb.PurgeReplyCloudflare{Result: false}, err
	}

	// Most API calls require a Context
	ctxAPI := context.Background()

	id := in.GetZoneId()
	// set purgeRequest
	fileList := in.GetPurgeList()

	pcr := cloudflare.PurgeCacheRequest{
		Files: fileList,
	}

	response, err := api.PurgeCache(ctxAPI, id, pcr)
	if response.Response.Success == true {
		return &pb.PurgeReplyCloudflare{Result: true}, nil
	}
	return &pb.PurgeReplyCloudflare{Result: false}, nil
}

func (s *server) PurgeCloudflareEverything(ctx context.Context, in *pb.PurgeRequestCloudflareEverything) (*pb.PurgeReplyCloudflare, error) {
	log.Printf("purgeEverythingTask")
	log.Printf("apiKey Received: %v", in.GetApiKey())
	log.Printf("apiEmail Received: %v", in.GetApiEmail())
	log.Printf("zoneId Received: %v", in.GetZoneId())

	// Construct a new API object using a global API key
	api, err := cloudflare.New(in.GetApiKey(), in.GetApiEmail())
	// alternatively, you can use a scoped API token
	// api, err := cloudflare.NewWithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN"))
	if err != nil {
		log.Fatal(err)
		return &pb.PurgeReplyCloudflare{Result: false}, err
	}

	// Most API calls require a Context
	ctxAPI := context.Background()

	id := in.GetZoneId()
	// set purgeRequest

	response, err := api.PurgeEverything(ctxAPI, id)
	if response.Response.Success == true {
		return &pb.PurgeReplyCloudflare{Result: true}, nil
	}
	return &pb.PurgeReplyCloudflare{Result: false}, nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterPurgeCloudflareServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
