package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

//roadNum shows the numbers of the roads that we have
const roadNum = 4

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
	rand.Seed(time.Now().Unix())

	log.SetFormatter(&log.TextFormatter{})

	inRoadMap := make(map[int]roadStruct)
	outRoadMap := make(map[int]roadStruct)
	circleIn := make(chan carStruct, 8)
	// circleOut := make(chan carStruct, 8)
	ctx, cancel := context.WithCancel(context.Background())
	// wg := &sync.WaitGroup{}

	//Describe inRoad
	inRoadMap[1] = roadStruct{time.Second * 2, rand.Intn(5) + 1, circleIn}
	inRoadMap[2] = roadStruct{time.Second * 2, 1, circleIn}
	inRoadMap[3] = roadStruct{time.Second * 3, 1, circleIn}
	inRoadMap[4] = roadStruct{time.Second * 4, 10, circleIn}

	//Describe outRoads
	outRoadMap[1] = roadStruct{time.Second * 2, rand.Intn(5) + 1, make(chan carStruct, 8)}
	outRoadMap[2] = roadStruct{time.Second * 2, 1, make(chan carStruct, 8)}
	outRoadMap[3] = roadStruct{time.Hour, 1, make(chan carStruct, 8)}
	outRoadMap[4] = roadStruct{time.Second * 4, 10, make(chan carStruct, 8)}

	for road, config := range inRoadMap {
		go inRoad(ctx, config, road)
	}
	for _, config := range outRoadMap {
		go outRoad(ctx, config)
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
			log.Info("InRoad got Context done, closing itself")
			return
		case <-ticker.C:
			destRoadNum := rand.Intn(4) + 1
			carID := fmt.Sprintf("Car#%d from road #%d to road #%d", j, roadNum, destRoadNum)
			newCar := carStruct{carID, roadNum, destRoadNum}
			rDesc.roadCh <- newCar
			log.Info(carID, "was generated and put in the circle")
		}
	}

}

func calcSleep(car carStruct) time.Duration {
	var sleep time.Duration

	if car.departure >= car.destination {
		sleep = time.Second * time.Duration(roadNum-(car.departure-car.destination))
	} else {
		sleep = time.Second * time.Duration(car.destination-car.departure)
	}

	return sleep
}

func waitingCircle(ctx context.Context, car carStruct, circleOut chan carStruct) {
	sleep := calcSleep(car)

	log.Info("car on waiting circle, sleep for ", sleep, " Car desc : ", car.id)
	ticker := time.NewTicker(sleep)
	defer ticker.Stop()

	select {
	case <-ctx.Done():
	case <-ticker.C:
		circleOut <- car
		log.Info("Car driving to the out road")
	}

	return
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
			log.Info("traffic circle took car for sleeping ", car.id)

			go waitingCircle(ctx, car, circleOut[car.destination].roadCh)
		}
	}
}

func outRoad(ctx context.Context, rDesc roadStruct) {
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
