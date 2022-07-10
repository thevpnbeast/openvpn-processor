package options

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/spf13/viper"
	"os"
)

// getStringEnv gets the specific environment variables with default value, returns default value if variable not set
func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// unmarshalConfig creates a new *viper.Viper and unmarshalls the config into struct using *viper.Viper
func unmarshalConfig(key string, opts interface{}) error {
	sub := viper.Sub(key)
	return sub.Unmarshal(opts)
}

func fetchConfig(applicationId, configProfileId string) error {
	sess := session.Must(session.NewSession())
	svc := appconfig.New(sess)

	application, err := svc.GetApplication(&appconfig.GetApplicationInput{
		ApplicationId: aws.String(applicationId),
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case appconfig.ErrCodeResourceNotFoundException:
				fmt.Println(appconfig.ErrCodeResourceNotFoundException, aerr.Error())
			case appconfig.ErrCodeInternalServerException:
				fmt.Println(appconfig.ErrCodeInternalServerException, aerr.Error())
			case appconfig.ErrCodeBadRequestException:
				fmt.Println(appconfig.ErrCodeBadRequestException, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return err
	}

	confProfile, _ := svc.GetConfigurationProfile(&appconfig.GetConfigurationProfileInput{
		ApplicationId: application.Id,
		ConfigurationProfileId: aws.String(configProfileId),
	})

	versions, _ := svc.ListHostedConfigurationVersions(&appconfig.ListHostedConfigurationVersionsInput{
		ConfigurationProfileId: confProfile.Id,
		ApplicationId: application.Id,
	})

	configuration, _ := svc.GetHostedConfigurationVersion(&appconfig.GetHostedConfigurationVersionInput{
		ApplicationId: application.Id,
		ConfigurationProfileId: confProfile.Id,
		VersionNumber: versions.Items[0].VersionNumber,
	})

	viper.SetConfigName(TargetConfigFileName)
	viper.SetConfigType(TargetConfigFileType)
	if err = viper.ReadConfig(bytes.NewBuffer(configuration.Content)); err != nil {
		return err
	}

	if err := unmarshalConfig(ApplicationName, opts); err != nil {
		return err
	}

	return nil
}