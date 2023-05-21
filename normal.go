package adidas_backend

import (
	"strings"
	profilesreader "umbrella/internal/file_readers/profiles_reader"
	"umbrella/internal/file_readers/proxy"
	tasksreader "umbrella/internal/file_readers/tasks_reader"
	quicktaskshandler "umbrella/internal/quicktasks_handler"
	clititle "umbrella/internal/utils/cli_title"
	definederrors "umbrella/internal/utils/defined_errors"
	waithandler "umbrella/internal/utils/wait_handler"
)

func Normal(proxy *proxy.Proxy, profile *profilesreader.Profile, task *tasksreader.TaskData, id int) {
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

			err = HandleTaskRestart(&config, err)

			if err != nil {
				continue
			}
		}

		err = SavedSessionHandler(&config)

		if err != nil {
			continue
		}

		err = GetAddresses(&config)

		if err != nil {
			continue
		}

		err = waithandler.HandleUserWait(&config.DefaultConfig)

		if err != nil {
			continue
		}

		var checkoutVariant Variant

		checkoutVariant, _, err = Monitor(&config)

		if err != nil {
			continue
		}

		config.DefaultConfig.CreateStartTimestamp()

		err = CheckoutId(&config, checkoutVariant)

		if err != nil {
			continue
		}

		err = Payment(&config)

		if err != nil {
			continue
		}

		err = Address(&config)

		if err != nil {
			continue
		}

		err = Discount(&config)

		if err != nil {
			continue
		}

		err = Order(&config)

		if err != nil {
			continue
		}

		if strings.ToLower(task.ResetAfterSuccess) != "true" {
			return
		}
	}
}
