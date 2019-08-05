package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "../protobuf"
)

const (
	adress = "localhost:50051"
)

func addFeedback(fb *pb.Feedback, client pb.PassengerFeedbackClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.AddPassengerFeedback(ctx, fb)

	if err != nil {
		return err
	}

	fmt.Println("Passenger: ", response.PassengerID, " sent feedback with booking code: ", response.BookingCode)
	return nil
}

func getFeedbackFromPassengerID(pid int32, client pb.PassengerFeedbackClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.GetFeedbackByPassengerID(ctx, &pb.PassengerID{PassengerID: pid})

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Array feedback from Passenger ID %d:\n", pid)
	for _, feedback := range response.ArrayFeedback {
		logFeedback(feedback)
	}
}

func logFeedback(fb *pb.Feedback) {
	fmt.Println("Booking code: ", fb.BookingCode, ", Feedback: ", fb.Feedback)
}

func getFeedbackFromBookingCode(bc string, client pb.PassengerFeedbackClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.GetFeedbackByBookingCode(ctx, &pb.BookingCode{BookingCode: bc})

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Array feedback from booking code %s:\n", bc)
	logFeedback(response)
}

func deleteFeedbackByPassengerID(pid int32, client pb.PassengerFeedbackClient) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	response, err := client.DeleteByPassengerID(ctx, &pb.PassengerID{PassengerID: pid})

	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(response.DeletedCount, " feedback(s) were deleted")
}

func main() {
	// Set up a connection to the server.
	connection, err := grpc.Dial(adress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer connection.Close()

	client := pb.NewPassengerFeedbackClient(connection)

	for i := 0; i < 3; i++ {
		bc := fmt.Sprintf("Free ship %d", i)
		msg := "Good"
		if i%2 != 0 {
			msg = "Bad"
		}
		fb := pb.Feedback{BookingCode: bc, PassengerID: 1, Feedback: msg}
		addFeedback(&fb, client)
	}

	getFeedbackFromPassengerID(1, client)
	getFeedbackFromBookingCode("Free ship 1", client)

	deleteFeedbackByPassengerID(1, client)

	getFeedbackFromPassengerID(1, client)
}
