package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	pb "github.com/rodrwan/cat-grpc/categoryapi"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var (
	address     = flag.String("host", "localhost", "Server host")
	port        = flag.Int("port", 10000, "Server port")
	description = flag.String("description", "ADIDAS PARQUE ARAUCO", "Description to be categorized")
)

func main() {
	flag.Parse()
	url := fmt.Sprintf("%s:%d", *address, *port)
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	start := time.Now()
	c := pb.NewCategoryAPIClient(conn)

	r, err := c.Categorize(context.Background(), &pb.Transaction{
		Description: *description,
	})
	if err != nil {
		grpclog.Fatalf("failed to categorize: %v", err)
	}

	elapsed := time.Since(start)

	fmt.Printf("\nDescripcion: %s\n", r.Description)
	fmt.Printf("Category %d: %s\n", r.CategoryID, r.CategoryName)
	log.Printf("Prediction took %s\n\n", elapsed)
}
