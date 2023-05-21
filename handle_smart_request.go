package adidas_backend

func HandleSmartModeRequests(config *Config) error {
	var err error

	if config.Resources.RequiredRequests.Payment {
		err = Payment(config)

		if err != nil {
			return err
		}

	}

	go AsyncDiscount(config)

	if config.Resources.RequiredRequests.Shipping {
		err = Address(config)

		if err != nil {
			return err
		}
	}

	return nil
}
