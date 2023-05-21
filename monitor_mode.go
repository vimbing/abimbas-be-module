package adidas_backend

import (
	profilesreader "umbrella/internal/file_readers/profiles_reader"
	"umbrella/internal/file_readers/proxy"
	tasksreader "umbrella/internal/file_readers/tasks_reader"
	quicktaskshandler "umbrella/internal/quicktasks_handler"
	clititle "umbrella/internal/utils/cli_title"
	definederrors "umbrella/internal/utils/defined_errors"
	waithandler "umbrella/internal/utils/wait_handler"
)

func MonitorMode(proxy *proxy.Proxy, profile *profilesreader.Profile, task *tasksreader.TaskData, id int) {
	defer clititle.DecreaseRunning()

	config, err := Init(proxy, profile, task, id)

	if err != nil {
		return
	}

	go quicktaskshandler.RegisterTaskToQuicktaskHandler(config.DefaultConfig.TaskData)

	for {
		if err != nil {
			if err.Error() == definederrors.IDENTIFIER_STOP_TASK {
				break
			}

		}

		err = waithandler.HandleUserWait(&config.DefaultConfig)

		if err != nil {
			continue
		}

		var restockedVariants []Variant
		_, restockedVariants, err = Monitor(&config)

		if err != nil {
			continue
		}

		monitorVariant := tasksreader.MonitorVariant{}
		err = monitorVariant.SetData(restockedVariants)

		if err != nil {
			config.DefaultConfig.Log.Red("Error while sending info from monitor to tasks!")
			continue
		}

		config.DefaultConfig.MonitorModeSendData(&monitorVariant)
	}
}
