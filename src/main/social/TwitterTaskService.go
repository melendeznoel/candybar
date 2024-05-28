package social

import (
	"github.com/jasonlvhit/gocron"
)

// StartTwitterTask will run once a day
func StartTwitterTask() {
	s := gocron.NewScheduler()

	// TODO: uncomment when tasks are ready.
	//s.Every(1).Day().Do(twitterService.RunTasks)

	s.Every(5).Seconds().Do(RunTasks)

	<-s.Start()
}
