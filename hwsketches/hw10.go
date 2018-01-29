package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
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
	log.SetFormatter(&log.TextFormatter{})

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
}

func inRoad(ctx context.Context, circleIn chan<- carStruct, rDesc roadStruct, roadNum int) {
	//send cars to trafficCircle
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for j := 0; j < rDesc.carNum; j++ {
				carID := fmt.Sprintf("Car#%d from road #%d\n", j, roadNum)
				newCar := carStruct{carID, time.Now()}
				circleIn <- newCar
				log.Infoln(carID)
			}
			log.Infof("Road with number %d will be sleeping for %d \n", roadNum, rDesc.sleepDuration)
			time.Sleep(rDesc.sleepDuration)
		}
	}

}

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
			// circleOut <- car
			log.Info(len(circleIn))
			if time.Since(car.onCircleTime) > time.Second {
				log.Info("took car from")
				circleOut <- car
			} else {
				log.Info("Car is back in channel")
				// we will be blocked over here
				// the problem is
				// we take from channel and while we are checking a time for the machine
				// another goroutine will put some data into our channel and we will be blocked
				circleIn <- car
				log.Info("we wont get here")
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
				log.Info("trying to take a car")
				car := <-circleOut
				log.Info("We took out off road ", car.id)
			}
			time.Sleep(rDesc.sleepDuration)
		}
	}
}
