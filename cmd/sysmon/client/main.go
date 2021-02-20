package main

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/artofey/sysmon/pkg/server/pb"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewMonitorClient(conn)

	inReader := bufio.NewReader(os.Stdin)

	for {
		req, err := getRequest(inReader)
		if err != nil {
			log.Printf("request error: %v", err)
			continue
		}
		monitorClient, err := client.GetStats(context.Background(), req)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("MonRequest submitted")
		for {
			ss, err := monitorClient.Recv()
			if err != nil {
				log.Printf("response error: %v", err)
				return
			}
			fmt.Println(pbSSToString(ss))
		}
	}
}

func pbSSToString(ss *pb.StatSnapshot) string {
	lavg := fmt.Sprintf(
		"LoadAVG: %v, %v, %v",
		ss.Lavg.Load1, ss.Lavg.Load5, ss.Lavg.Load15,
	)
	lcpu := fmt.Sprintf(
		"LoadCPU: %v, %v, %v",
		ss.Lcpu.User, ss.Lcpu.System, ss.Lcpu.Idle,
	)
	return fmt.Sprintf("%v \n%v", lavg, lcpu)
}

func getRequest(reader *bufio.Reader) (*pb.MonRequest, error) {
	log.Printf("write MonRequest <Timeout> <AveragedOver>:")
	text, err := reader.ReadString('\n')
	if err != nil {
		return nil, errors.New("wrong input, try again")
	}

	parts := strings.Split(strings.TrimSpace(text), " ")
	if len(parts) < 2 {
		return nil, errors.New("wrong input, try again")
	}
	timeout, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, errors.New("wrong input, try again")
	}
	avgOver, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, errors.New("wrong input, try again")
	}

	return &pb.MonRequest{
		Timeout:      uint32(timeout),
		AveragedOver: uint32(avgOver),
	}, nil
}
