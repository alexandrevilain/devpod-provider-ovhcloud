package options

import (
	"fmt"
	"os"
)

type Authentication struct {
	Endpoint    string
	AppKey      string
	AppSecret   string
	ConsumerKey string
	Region      string
	ServiceName string
}

type Options struct {
	MachineID     string
	MachineFolder string

	Authentication *Authentication

	Flavor string
}

func FromEnv(skipMachine bool) (*Options, error) {
	retOptions := &Options{}

	retOptions.Authentication = &Authentication{}

	var err error
	if !skipMachine {
		retOptions.MachineID, err = fromEnvOrError("MACHINE_ID")
		if err != nil {
			return nil, err
		}
		// prefix with devpod-
		retOptions.MachineID = "devpod-" + retOptions.MachineID

		retOptions.MachineFolder, err = fromEnvOrError("MACHINE_FOLDER")
		if err != nil {
			return nil, err
		}
	}

	retOptions.Authentication = &Authentication{}
	retOptions.Authentication.Endpoint, err = fromEnvOrError("OVHCLOUD_ENDPOINT")
	if err != nil {
		return nil, err
	}

	retOptions.Authentication.AppKey, err = fromEnvOrError("OVHCLOUD_APP_KEY")
	if err != nil {
		return nil, err
	}

	retOptions.Authentication.AppSecret, err = fromEnvOrError("OVHCLOUD_APP_SECRET")
	if err != nil {
		return nil, err
	}

	retOptions.Authentication.ConsumerKey, err = fromEnvOrError("OVHCLOUD_CONSUMER_KEY")
	if err != nil {
		return nil, err
	}

	retOptions.Authentication.ServiceName, err = fromEnvOrError("OVHCLOUD_SERVICE_NAME")
	if err != nil {
		return nil, err
	}

	retOptions.Authentication.Region, err = fromEnvOrError("OVHCLOUD_REGION")
	if err != nil {
		return nil, err
	}

	retOptions.Flavor, err = fromEnvOrError("OVHCLOUD_FLAVOR")
	if err != nil {
		return nil, err
	}

	return retOptions, nil
}

func fromEnvOrError(name string) (string, error) {
	val := os.Getenv(name)
	if val == "" {
		return "", fmt.Errorf("couldn't find option %s in environment, please make sure %s is defined", name, name)
	}

	return val, nil
}
