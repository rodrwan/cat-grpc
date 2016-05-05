package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"

	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	pb "github.com/rodrwan/cat-grpc/categoryapi"
	"github.com/rodrwan/lucky"
)

type categoryServer struct {
	newModel *lucky.Lucky
}

// Categorize ...
func (c *categoryServer) Categorize(ctx context.Context, t *pb.Transaction) (*pb.Transaction, error) {
	res := c.newModel.Predict(t.Description)
	start := time.Now()
	newTrans := &pb.Transaction{
		Description:  t.Description,
		CategoryID:   uint32(res.ID),
		CategoryName: res.Name,
	}
	elapsed := time.Since(start)
	log.Printf("Prediction took %s\n\n", elapsed)
	return newTrans, nil
}

func newServer(train, labels string) *categoryServer {
	log.Println("Running new server")
	cs := new(categoryServer)
	cs.newModel = &lucky.Lucky{
		LabelsPath:       labels,
		TrainingDataPath: train,
	}

	cs.newModel.Fit()
	return cs
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hi"))
}

func main() {
	var (
		port   = flag.Int("port", 10000, "The server port")
		path   = flag.String("data", "training_data.txt", "Training data path")
		labels = flag.String("cats", "labels.txt", "Labels sample path")
	)
	flag.Parse()
	log.Println("Running RPC server")

	gs := grpc.NewServer()
	cs := newServer(*path, *labels)
	pb.RegisterCategoryAPIServer(gs, cs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs.Serve(lis)
}
