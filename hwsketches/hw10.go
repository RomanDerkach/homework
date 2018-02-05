package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

type carStruct struct {
	id          string
	departure   int
	destination int
}

type roadStruct struct {
	sleepDuration time.Duration
	carNum        int
	roadCh        chan carStruct
}

func main() {
	log.SetFormatter(&log.TextFormatter{})

	inRoadMap := make(map[int]roadStruct)
	outRoadMap := make(map[int]roadStruct)
	circleIn := make(chan carStruct, 8)
	// circleOut := make(chan carStruct, 8)
	ctx, cancel := context.WithCancel(context.Background())
	// wg := &sync.WaitGroup{}

	//Describe inRoad
	inRoadMap[1] = roadStruct{time.Second * 2, rand.Intn(6), circleIn}
	inRoadMap[2] = roadStruct{time.Second * 2, 1, circleIn}
	inRoadMap[3] = roadStruct{time.Second * 3, 1, circleIn}
	inRoadMap[4] = roadStruct{time.Second * 4, 10, circleIn}

	//Describe outRoads
	outRoadMap[1] = roadStruct{time.Second * 2, rand.Intn(6), make(chan carStruct, 8)}
	outRoadMap[2] = roadStruct{time.Second * 2, 1, make(chan carStruct, 8)}
	outRoadMap[3] = roadStruct{time.Hour, 1, make(chan carStruct, 8)}
	outRoadMap[4] = roadStruct{time.Second * 4, 10, make(chan carStruct, 8)}

	for road, config := range inRoadMap {
		go inRoad(ctx, config, road)
	}
	for road, config := range outRoadMap {
		go outRoad(ctx, config, road)
	}
	go trafficCircle(ctx, circleIn, outRoadMap)
	time.Sleep(time.Second * 30)
	cancel()
}

func inRoad(ctx context.Context, rDesc roadStruct, roadNum int) {
	//send cars to trafficCircle
	ticker := time.NewTicker(rDesc.sleepDuration / time.Duration(rDesc.carNum))
	defer ticker.Stop()

	for j := 0; ; j++ {
		select {
		case <-ctx.Done():
			log.Infoln("InRoad got Context done, closing itself")
			return
		case <-ticker.C:
			carID := fmt.Sprintf("Car#%d from road #%d\n", j, roadNum)
			newCar := carStruct{carID, roadNum, (rand.Intn(4) + 1)}
			rDesc.roadCh <- newCar
			log.Infoln(carID, "was generated and put in the circle")
		}
	}

}

func waitingCircle(ctx context.Context, car carStruct, sleep time.Duration, circleOut chan carStruct) {
	log.Infoln("car on waiting circle, sleep for ", sleep, "Car desc : ", car.id)
	time.Sleep(sleep)
	select {
	case <-ctx.Done():
		return
	case circleOut <- car:
		log.Infoln("Car driving to the out road")
		return
	}
}

// TrafficCircle describe a circle
func trafficCircle(ctx context.Context, circleIn chan carStruct, circleOut map[int]roadStruct) {
LoopLabel:
	for {
		select {
		case <-ctx.Done():
			return
		case car, ok := <-circleIn:
			if !ok {
				break LoopLabel
			}
			log.Infoln(len(circleIn))
			log.Info("traffic circle took car for sleeping", car.id)
			//!!!!!!need to be improved
			sleep := time.Second * time.Duration(math.Abs(float64(car.departure-car.destination)))

			go waitingCircle(ctx, car, sleep, circleOut[car.destination].roadCh)

		}
	}
}

func outRoad(ctx context.Context, rDesc roadStruct, roadNum int) {
	//get data from trafficCircle
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for i := 0; i < rDesc.carNum; i++ {
				car := <-rDesc.roadCh
				log.Info("We took something out off road ", car.id)
			}
			time.Sleep(rDesc.sleepDuration)
		}
	}
}
