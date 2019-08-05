package main

import (
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "../protobuf"
)

const (
	port = ":50051"
)

type feedbackService struct {
	dictFeedback map[string]pb.Feedback
}

//PassengerFeedback service
func (fs *feedbackService) AddPassengerFeedback(ctx context.Context, req *pb.Feedback) (*pb.AddFeedbackResponse, error) {
	log.Println("Added feedback")
	fs.dictFeedback[req.BookingCode] = *req
	return &pb.AddFeedbackResponse{BookingCode: req.BookingCode, PassengerID: req.PassengerID}, status.Errorf(codes.OK, "Add Passenger Feedback successful")
}

func (fs *feedbackService) GetFeedbackByPassengerID(ctx context.Context, req *pb.PassengerID) (*pb.ArrayFeedback, error) {
	log.Println("Feedback from passenger id")
	var arrayFeedback []*pb.Feedback
	for _, feedback := range fs.dictFeedback {
		if feedback.PassengerID == req.PassengerID {
			arrayFeedback = append(arrayFeedback, &pb.Feedback{BookingCode: feedback.BookingCode, PassengerID: feedback.PassengerID, Feedback: feedback.Feedback})
		}
	}
	if len(arrayFeedback) == 0 {
		return nil, status.Errorf(codes.NotFound, "Feedback not found")
	}

	return &pb.ArrayFeedback{ArrayFeedback: arrayFeedback}, status.Errorf(codes.OK, "Get Feedback successful")
}

func (fs *feedbackService) GetFeedbackByBookingCode(ctx context.Context, req *pb.BookingCode) (*pb.Feedback, error) {
	log.Println("Feedback from booking code")
	if value, ok := fs.dictFeedback[req.BookingCode]; ok == true {
		return &value, status.Errorf(codes.OK, "Get successful")
	}

	return nil, status.Errorf(codes.NotFound, "Feedback not found")
}

func (fs *feedbackService) DeleteByPassengerID(ctx context.Context, req *pb.PassengerID) (*pb.DeleteByPassengerIDResponse, error) {
	log.Println("Delete feedback from passenger id")
	var count int32
	for key, fb := range fs.dictFeedback {
		if fb.PassengerID == req.PassengerID {
			delete(fs.dictFeedback, key)
			count++
		}
	}

	if count == 0 {
		return nil, status.Errorf(codes.NotFound, "Feedback not found")
	}

	return &pb.DeleteByPassengerIDResponse{DeletedCount: count}, status.Errorf(codes.OK, "Delete successful")
}

func main() {
	// Set up a connection to the server.
	listen, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	pb.RegisterPassengerFeedbackServer(server, &feedbackService{dictFeedback: make(map[string]pb.Feedback)})
	if err := server.Serve(listen); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
