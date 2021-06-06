package main

import (
	notify "armanbimak/notify/protos"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os"

)


func doLongGreet(c notify.NotificationServiceClient) {

	requests := []*notify.NotificationRequest{
		{
			From: "askar",
			To: "kanabeesss@gmail.com",
		},
		{
			From: "arman",
			To: "kanabeesss@gmail.com",
		},
		{
			From: "dauka",
			To: "kanabeesss@gmail.com",
		},
	}

	ctx := context.Background()
	stream, err := c.NotifyCreate(ctx)
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	for _, req := range requests {
		stream.Send(req)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error with getting response from server : %v", err)
	}
	fmt.Printf("Response: %v\n", res.GetIsDone())
}



func doGreet(c notify.NotificationServiceClient) {

	requests := notify.NotificationRequest{To: "kanabeesss@gmail.com", From: "Seka"}

	ctx := context.Background()
	_, err := c.NotifyFollow(ctx, &requests)
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}
}



func main() {


	fmt.Println("Hello I'm a client")

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer conn.Close()

	c := notify.NewNotificationServiceClient(conn)

	doLongGreet(c)
	doGreet(c)
	os.Exit(1)
}