package main

import (
	"context"
	"encoding/json"
	"log"
	api "tf-serving/protos/github.com/tensorflow/serving/tensorflow_serving/apis"
	"time"

	"google.golang.org/grpc"

	tfs_api_pb "tf-serving/protos/github.com/tensorflow/serving/tensorflow_serving/apis"
	tf_framework "tf-serving/protos/github.com/tensorflow/tensorflow/tensorflow/go/core/framework"
)

const (
	HOST                 = "0.0.0.0:8500"
	MODEL_NAME           = "half_plus_two"
	MODEL_SIGNATURE_NAME = "serving_default"
)

func main() {
	conn, err := grpc.Dial(HOST, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := api.PredictionLog{}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	predictRequest := tfs_api_pb.PredictRequest{
		ModelSpec: &tfs_api_pb.ModelSpec{
			Name:          MODEL_NAME,
			SignatureName: MODEL_SIGNATURE_NAME,
		},
		Inputs: map[string]*tf_framework.TensorProto{
			"x": &tf_framework.TensorProto{
				Dtype: tf_framework.DataType_DT_FLOAT,
				TensorShape: &tf_framework.TensorShapeProto{
					Dim: []*tf_framework.TensorShapeProto_Dim{
						{
							Size: 1,
						},
						{
							Size: 3,
						},
					},
				},
				FloatVal: []float32{
					1.0, 2.0, 5.0,
				},
			},
		},
	}

	for key, value := range predictRequest.Inputs {
		log.Printf("Input %s", key)

		for _, element := range value.FloatVal {
			log.Printf("\t%f", element)
		}
	}

	predictResponse, err := c.Predict(ctx, &predictRequest)
	if err != nil {
		log.Fatalf("could not get response: %v", err)
	}

	jsonResponse, err := json.Marshal(predictResponse)
	if err != nil {
		log.Fatalf("could not marshal: %v", err)
	}
	log.Printf("Response: %s", jsonResponse)

	for key, value := range predictResponse.Outputs {
		log.Printf("Output %s", key)

		for _, element := range value.FloatVal {
			log.Printf("\t%f", element)
		}
	}
}
