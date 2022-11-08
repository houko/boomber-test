package main

import (
	pb "boomer-test/proto"
	"context"
	"fmt"
	"github.com/myzhan/boomer"
	"google.golang.org/grpc"
	"log"
	"time"
)

const (
	target = "localhost:50051"
)

func hello() {
	ctx := context.Background()
	conn, err := grpc.DialContext(ctx, target, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("%+v", err)
	}
	defer conn.Close()

	client := pb.NewGreeterClient(conn)

	// レスポンスタイムを計測するために、リクエストを投げる前の時刻を取得する
	start := time.Now()
	response, err := client.SayHello(ctx, &pb.HelloRequest{Name: "ミクシィ太郎"})

	if err != nil {
		// リクエストに失敗した場合はRecordFailureを呼びます
		boomer.RecordFailure(
			"grpc",     // リクエストの種別
			"SayHello", // rpc名
			time.Since(start).Nanoseconds()/int64(time.Millisecond), // レスポンスタイム
			fmt.Sprintf("failed: %+v", err),                         // エラー理由
		)
	} else {
		// リクエストに成功した場合はRecordSuccessを呼びます
		boomer.RecordSuccess(
			"grpc",     // リクエストの種別
			"SayHello", // rpc名
			time.Since(start).Nanoseconds()/int64(time.Millisecond), // レスポンスタイム
			int64(len(response.String())),                           // レスポンスサイズ
		)
	}
}

func main() {
	task := &boomer.Task{
		Name:   "sample",
		Weight: 10,
		Fn:     hello,
	}

	boomer.Run(task)
}
