package helper

import (
	"fmt"

	"github.com/jasonlvhit/gocron"
)

func tempTask() {
	fmt.Println("generic task")
}

// StartGenericTask ... will run once a day
func StartGenericTask() {
	s := gocron.NewScheduler()

	s.Every(1).Day().Do(tempTask)

	<-s.Start()
}
