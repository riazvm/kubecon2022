package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	pb "example.com/go-edgemgmt-grpc/edgemgmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	address = "localhost:443"
)

func loadTLSCredentials() (credentials.TransportCredentials, error) {
	fmt.Printf("Recieved in load")
	clientCert, err := tls.LoadX509KeyPair("cert/client/tls.crt", "cert/client/tls.key")
	if err != nil {
		return nil, err
	}
	pemServerCA, err := ioutil.ReadFile("cert/client/ca.crt")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Recieved in load 2")
	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, fmt.Errorf("failed to add server CA's certificate")
	}

	config := &tls.Config{
		Certificates:       []tls.Certificate{clientCert},
		RootCAs:            certPool,
		InsecureSkipVerify: true,
	}
	return credentials.NewTLS(config), nil
}

func main() {

	http.HandleFunc("/edgeActivate", activateEdgeDevice)
	http.HandleFunc("/edgePay", processPayment)
	fmt.Printf("Starting server for testing HTTP POST...\n")
	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatal(err)
	}

	/*var activate_device = make(map[string]int32)
		activate_device["12345"] = 114477
		for merchantID, deviceID := range activate_device {
			r, err := c.ActivateNewDevice(ctx, &pb.NewDevice{MerchantID: merchantID, DeviceID: deviceID})
			if err != nil {
				log.Fatalf("could not create user: %v", err)
			}
			log.Printf(`Device Details:
	MerchantID: %s
	DeviceID: %d
	DeviceStatus: %s`, r.GetMerchantID(), r.GetDeviceID(), r.GetDeviceStatus())*/

}

func activateEdgeDevice(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	/*tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}*/
	//conn, err := grpc.Dial(address, grpc.WithTransportCredentials(tlsCredentials), grpc.WithBlock())
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEdgeManagementClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var activate_device = make(map[string]int32)
	activate_device["12345"] = 114477
	for merchantID, deviceID := range activate_device {
		r, err := c.ActivateNewDevice(ctx, &pb.NewDevice{MerchantID: merchantID, DeviceID: deviceID})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf(`Device Details:
	MerchantID: %s
	DeviceID: %d
	DeviceStatus: %s`, r.GetMerchantID(), r.GetDeviceID(), r.GetDeviceStatus())

		json.NewEncoder(w).Encode(r)

	}
}

func processPayment(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	tlsCredentials, err := loadTLSCredentials()
	if err != nil {
		log.Fatal("cannot load TLS credentials: ", err)
	}

	conn, err := grpc.Dial("localhost", grpc.WithTransportCredentials(tlsCredentials), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEdgePaymentClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var payment_process = make(map[string]string)
	payment_process["12345"] = "this is a payment"
	for paymentRefNo, paymentDetails := range payment_process {
		r, err := c.ProcessPayment(ctx, &pb.Payment{PaymentRefNo: paymentRefNo, PaymentDetails: paymentDetails})
		if err != nil {
			log.Fatalf("could not create user: %v", err)
		}
		log.Printf(`Payment Details:
	PaymentRefNo: %s
	PaymentStatus: %s`, r.GetPaymentRefNo(), r.GetPaymentStatus())
		json.NewEncoder(w).Encode(r)

	}

}
