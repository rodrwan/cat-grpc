package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/net/trace"

	"google.golang.org/grpc"

	_ "github.com/lib/pq"
	pb "github.com/rodrwan/cat-grpc/categoryapi"
	"github.com/rodrwan/lucky"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

func newServer(url, train, labels string) *categoryServer {
	log.Println("Running new server")
	cs := new(categoryServer)
	cs.newModel = &lucky.Lucky{
		LabelsPath:       labels,
		TrainingDataPath: train,
		URL:              url,
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
		url    = flag.String("url", "you should set database url", "Database url")
		path   = flag.String("data", "training_data.txt", "Training data path")
		labels = flag.String("cats", "labels.txt", "Labels sample path")
	)
	flag.Parse()
	log.Println("Running RPC server")

	gs := grpc.NewServer()
	cs := newServer(*url, *path, *labels)
	pb.RegisterCategoryAPIServer(gs, cs)

	hs := health.NewHealthServer()
	hs.SetServingStatus("grpc.health.v1.categoryserver", 1)
	healthpb.RegisterHealthServer(gs, hs)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	go gs.Serve(lis)

	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) { return true, true }
	log.Println("Hello service started successfully.")
	http.HandleFunc("/", handler)
	log.Fatalln(http.ListenAndServe(":8080", nil))
}
