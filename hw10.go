package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

type carStruct struct {
	id           string
	onCircleTime time.Time
}

type roadStruct struct {
	sleepDuration time.Duration
	carNum        int
}

func main() {
	inRoadMap := make(map[int]roadStruct)
	outRoadMap := make(map[int]roadStruct)
	circleIn := make(chan carStruct, 8)
	circleOut := make(chan carStruct, 8)
	ctx, cancel := context.WithCancel(context.Background())
	// wg := &sync.WaitGroup{}

	//Describe inRoad
	inRoadMap[1] = roadStruct{time.Second * 2, rand.Intn(6)}
	inRoadMap[2] = roadStruct{time.Second * 2, 1}
	inRoadMap[3] = roadStruct{time.Second * 3, 1}
	inRoadMap[4] = roadStruct{time.Second * 4, 10}

	//Describe outRoads
	outRoadMap[1] = roadStruct{time.Second * 2, rand.Intn(6)}
	outRoadMap[2] = roadStruct{time.Second * 2, 1}
	outRoadMap[3] = roadStruct{time.Hour, 1}
	outRoadMap[4] = roadStruct{time.Second * 4, 10}

	for road, config := range inRoadMap {
		go inRoad(ctx, circleIn, config, road)
	}
	for road, config := range outRoadMap {
		go outRoad(ctx, circleOut, config, road)
	}
	go trafficCircle(ctx, circleIn, circleOut)
	time.Sleep(time.Second * 30)
	cancel()
	// TODO: Add wg
}

func inRoad(ctx context.Context, circleIn chan<- carStruct, rDesc roadStruct, roadNum int) {
	ticker := time.NewTicker(rDesc.sleepDuration / time.Duration(rDesc.carNum))
	defer ticker.Stop()

	//send cars to trafficCircle
	for j := 0; ; j++ {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			carID := fmt.Sprintf("Car #%d from road #%d\n", j, roadNum)
			fmt.Println(carID)
			newCar := carStruct{carID, time.Now()}
			circleIn <- newCar
		}
	}
}

//func inRoad(ctx context.Context, circleIn chan<- carStruct, rDesc roadStruct, roadNum int) {
//	//send cars to trafficCircle
//	for {
//		select {
//		case <-ctx.Done():
//			return
//		default:
//			for j := 0; j < rDesc.carNum; j++ {
//				carID := fmt.Sprintf("Car#%d from road #%d\n", j, roadNum)
//				fmt.Println(carID)
//				newCar := carStruct{carID, time.Now()}
//				circleIn <- newCar
//			}
//			fmt.Printf("Road with number %d will be sleeping for %d \n", roadNum, rDesc.sleepDuration)
//			time.Sleep(rDesc.sleepDuration)
//		}
//	}
//
//}

// TrafficCircle describe a circle
func trafficCircle(ctx context.Context, circleIn chan carStruct, circleOut chan carStruct) {
LoopLabel:
	for {
		select {
		case <-ctx.Done():
			return
		case car, ok := <-circleIn:
			if !ok {
				break LoopLabel
			}
			//fmt.Println(time.Since(car.onCircleTime))
			//fmt.Println(time.Second)

			// circleOut <- car
			//status := time.Since(car.onCircleTime) > time.Second
			//fmt.Println(status)
			if time.Since(car.onCircleTime) > time.Second {
				fmt.Println("took car from")
				circleOut <- car
			} else {
				//fmt.Println("Car is back in channel")
				circleIn <- car
			}
		}
	}
}

func outRoad(ctx context.Context, circleOut <-chan carStruct, rDesc roadStruct, roadNum int) {
	//get data from trafficCircle
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for i := 0; i < rDesc.carNum; i++ {
				car := <-circleOut
				fmt.Println("We took out off road ", car.id)
			}
			time.Sleep(rDesc.sleepDuration)
		}
	}
}
