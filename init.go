package adidas_backend

import (
	"strings"
	profilesreader "umbrella/internal/file_readers/profiles_reader"
	"umbrella/internal/file_readers/proxy"
	tasksreader "umbrella/internal/file_readers/tasks_reader"
	"umbrella/internal/modules"
	globaltypes "umbrella/internal/utils/global_types"
)

func Init(proxy *proxy.Proxy, profile *profilesreader.Profile, taskData *tasksreader.TaskData, id int) (Config, error) {
	defaultConfig := modules.Init()

	err := defaultConfig.SetDefaultConfig(proxy, profile, taskData, id)

	if err != nil {
		return Config{}, err
	}

	taskStates := TaskStates{
		Login: globaltypes.TaskState{
			Name:  "LOGIN",
			Retry: 3,
		},
		SessionCheck: globaltypes.TaskState{
			Name:  "SESSION CHECK",
			Retry: 15,
		},
		Payment: globaltypes.TaskState{
			Name:  "PAYMENT",
			Retry: 15,
		},
		Delivery: globaltypes.TaskState{
			Name:  "DELIVERY",
			Retry: 15,
		},
		AddressGet: globaltypes.TaskState{
			Name:  "ADDRESS GET",
			Retry: 15,
		},
		Address: globaltypes.TaskState{
			Name:  "ADDRESS",
			Retry: 15,
		},
		CheckoutID: globaltypes.TaskState{
			Name:  "CHECKOUT ID",
			Retry: 15,
		},
		Discount: globaltypes.TaskState{
			Name:  "DISCOUNT",
			Retry: 15,
		},
		Monitor: globaltypes.TaskState{
			Name:  "MONITOR",
			Retry: 15,
		},
		Order: globaltypes.TaskState{
			Name:  "ORDER",
			Retry: 15,
		},
		Akamai: globaltypes.TaskState{
			Name:  "AKAMAI",
			Retry: 15,
		},
		Cancel: globaltypes.TaskState{
			Name:  "CANCEL",
			Retry: 15,
		},
		Register: globaltypes.TaskState{
			Name:  "REGISTER",
			Retry: 15,
		},
		SessionRefresh: globaltypes.TaskState{
			Name:  "SESSION REFRESH",
			Retry: 15,
		},
	}

	config := Config{
		DefaultConfig: defaultConfig,
		TaskStates:    taskStates,
		Resources:     Resources{},
	}

	if strings.EqualFold(config.DefaultConfig.Profile.Payment, "cod") {
		config.DefaultConfig.Profile.Payment = PAYMENT_COD
	} else if strings.EqualFold(config.DefaultConfig.Profile.Payment, "paypal") || strings.EqualFold(config.DefaultConfig.Profile.Payment, "pp") {
		config.DefaultConfig.Profile.Payment = PAYMENT_PAYPAL
	} else if strings.EqualFold(config.DefaultConfig.Profile.Payment, "cc") {
		config.DefaultConfig.Profile.Payment = PAYMENT_CREDIT_CARD
	} else {
		config.DefaultConfig.Log.Yellow("Unknown payment method, picking paypal...")
		config.DefaultConfig.Profile.Payment = PAYMENT_PAYPAL
	}

	config.DefaultConfig.MonitorNetwork.Client.Jar = nil

	return config, nil
}
