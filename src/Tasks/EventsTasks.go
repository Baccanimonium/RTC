package Tasks

func (tm *TaskManager) RunEventTask() error {
	_, err := tm.scheduler.Every(1).Minute().Tag("EventTasks").Do(tm.eventTask)
	return err
}

func (tm *TaskManager) eventTask() {

}
