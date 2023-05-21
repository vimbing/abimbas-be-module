package adidas_backend

import customjar "umbrella/internal/http_client/custom_jar"

func GetSensorAsCookie(config *Config) customjar.Cookie {
	return customjar.Cookie{
		Value: config.Resources.SensorData,
		Name:  "sensor_data",
	}
}
