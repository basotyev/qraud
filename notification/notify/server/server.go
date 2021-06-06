package main

import (
	notify "armanbimak/notify/protos"
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
	"net/smtp"
)


type Server struct {
	notify.UnimplementedNotificationServiceServer
}


func (s *Server) NotifyCreate(stream notify.NotificationService_NotifyCreateServer)  error {
	fmt.Printf("Greet function was invoked with %v \n", stream)
	var isDone bool = true
	for {
		req, err := stream.Recv()
		if err == io.EOF {

			return stream.SendAndClose(&notify.NotificationResponse{
				IsDone: isDone,
			})
		}
		if err != nil {
			log.Println("Error from client stream: %v", err)
		}

		from := req.GetFrom()
		to := req.GetTo()

		message := from + " posted new article!!!"

		if !send(message, to){
				isDone = false
		}
	}
	return nil
}


func send(body string, to string) bool{
	from := "armanbimak27@gmail.com"
	pass := "ARASH@rash27122001"


	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return false
	}else{
		log.Print("sent, to " + to + " from " + from )
		return true
	}
}

func (s *Server) NotifyFollow( ctx context.Context, req *notify.NotificationRequest)  (*notify.NotificationResponse, error) {
	email_to := req.GetTo()
	email_from := req.GetFrom()

	message := email_from + " just started following you!!!"


	res := notify.NotificationResponse{IsDone: send(message, email_to)}
	if res.IsDone{
		return &res, nil
	} else{
		return &res, errors.New("can not send email")
	}

}


func main() {
	l, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen:%v", err)
	}
	s := grpc.NewServer()
	notify.RegisterNotificationServiceServer(s, &Server{})
	log.Println("Server is running on port:50051")
	if err := s.Serve(l); err != nil {
		log.Fatalf("failed to serve:%v", err)
	}

}
