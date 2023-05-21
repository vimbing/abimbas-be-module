package adidas_backend

import (
	"errors"
	"strings"
	profilesreader "umbrella/internal/file_readers/profiles_reader"
	"umbrella/internal/file_readers/proxy"
	tasksreader "umbrella/internal/file_readers/tasks_reader"
	quicktaskshandler "umbrella/internal/quicktasks_handler"
	clititle "umbrella/internal/utils/cli_title"
	definederrors "umbrella/internal/utils/defined_errors"
)

func AddressFiller(proxy *proxy.Proxy, profile *profilesreader.Profile, task *tasksreader.TaskData, id int) {
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

			if errors.Is(AkamaiNeedsSolveError{}, err) {
				err = SolveAkamai(&config)

				if err != nil {
					continue
				}
			}
		}

		err = SavedSessionHandler(&config)

		if err != nil {
			continue
		}

		err = AddressFill(&config)

		if err != nil {
			continue
		}

		if strings.ToLower(task.ResetAfterSuccess) != "true" {
			return
		}
	}
}
