package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	pb "protobuf-test/proto"
	ds "protobuf-test/structures"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/protobuf/proto"
)

func getAvtechData() ds.AvtechResponseData {
	avtechUrl := os.Getenv("AVTECH_URL")

	resp, err := http.Get(avtechUrl)
	if err != nil {
		fmt.Println("Error fetching data from Avtech")
	}
	defer resp.Body.Close()

	var avtechResponse ds.AvtechResponseData
	err = json.NewDecoder(resp.Body).Decode(&avtechResponse)
	if err != nil {
		fmt.Println("Error decoding Avtech response")
	}

	// fmt.Println("Avtech response:", avtechResponse)
	now := time.Now().UTC().String()
	for i := 0; i < len(avtechResponse.Sensor); i++ {
		avtechResponse.Sensor[i].Time = now
	}
	return avtechResponse
}

func testJson(n int, avtechResponse ds.AvtechResponseData) {

	start := time.Now()
	for i := 0; i < n; i++ {
		serializeJson(avtechResponse)
	}
	elapsed := time.Since(start)

	fmt.Printf("JSON serialization took %s\n", elapsed)
}

func serializeJson(d ds.AvtechResponseData) {
	_, err := json.Marshal(d)
	if err != nil {
		fmt.Println("Error marshalling Avtech response")
	}
	// fmt.Println("Avtech response:", string(b))
}

func testProtobuf(n int, avtechResponse ds.AvtechResponseData) {
	var dto *pb.SensorUpdateDto = &pb.SensorUpdateDto{
		Label: avtechResponse.Sensor[0].Label,
		TempF: avtechResponse.Sensor[0].TempF,
		TempC: avtechResponse.Sensor[0].TempC,
		HighF: avtechResponse.Sensor[0].HighF,
		HighC: avtechResponse.Sensor[0].HighC,
		LowF:  avtechResponse.Sensor[0].LowF,
		LowC:  avtechResponse.Sensor[0].LowC,
	}
	now := time.Now().UTC().String()
	dto.Time = now
	start := time.Now()
	for i := 0; i < n; i++ {
		serializeProtobuf(dto)
	}
	elapsed := time.Since(start)

	fmt.Printf("Protobuf serialization took %s\n", elapsed)
}

func serializeProtobuf(dto *pb.SensorUpdateDto) {

	_, err := proto.Marshal(dto)
	if err != nil {
		fmt.Println("Error proto marshalling Avtech response")
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	avtechResponse := getAvtechData()

	iterations := 10000000
	testJson(iterations, avtechResponse)
	testProtobuf(iterations, avtechResponse)
}
