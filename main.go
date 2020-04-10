package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mitchellh/mapstructure"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/xgfone/gconf/v4"
	"io/ioutil"
	"strings"
)

var (
	credentialsFile string
	profile         string
)

type AWSProfile struct {
	AccessKeyID     string `mapstructure:"aws_access_key_id" json:"AWS_ACCESS_KEY_ID,omitempty"`
	SecretAccessKey string `mapstructure:"aws_secret_access_key" json:"AWS_SECRET_ACCESS_KEY,omitempty"`
	SessionToken    string `mapstructure:"aws_session_token" json:"AWS_SESSION_TOKEN,omitempty"`
	Region          string `mapstructure:"region" json:"AWS_REGION,omitempty"`

	RoleArn         string `mapstructure:"role_arn" json:"AWS_ROLE_ARN,omitempty"`
	RoleSessionName string `mapstructure:"role_session_name" json:"AWS_ROLE_SESSION_NAME,omitempty"`
	SourceProfile   string `mapstructure:"source_profile" json:"AWS_SOURCE_PROFILE,omitempty"`
}

func (p AWSProfile) Validate() error {
	noAccessKeyID := p.AccessKeyID == ""
	noSecretAccessKey := p.SecretAccessKey == ""
	noSourceProfile := p.SourceProfile == ""
	noRoleArn := p.RoleArn == ""

	if noSourceProfile && noRoleArn {
		// expect an access key ID and secret access key
		if noAccessKeyID {
			return errors.New("no `aws_access_key_id` found in profile")
		}
		if noSecretAccessKey {
			return errors.New("no `aws_secret_access_key` found in profile")
		}
	}

	if noAccessKeyID && noSecretAccessKey {
		// expect a source profile and role arn
		if noSourceProfile {
			return errors.New("no `source_profile` found in profile")
		}
		if noRoleArn {
			return errors.New("no `role_arn` found in profile")
		}
	}

	return nil
}

func ParseAWSProfile(credentialsData []byte, profile string) (AWSProfile, error) {
	awsProfile := AWSProfile{}

	// check profile exists
	if !strings.Contains(string(credentialsData), fmt.Sprintf("[%s]", profile)) {
		return awsProfile, fmt.Errorf("no profile found with name: `%s`", profile)
	}

	// INI -> generic map
	rawData := make(map[string]interface{})
	err := gconf.NewIniDecoder(profile).Decode(credentialsData, rawData)
	if err != nil {
		return awsProfile, fmt.Errorf("failed to parse AWS profile from source: %s", err)
	}

	// generic map -> awsProfile
	err = mapstructure.Decode(rawData, &awsProfile)
	if err != nil {
		return awsProfile, fmt.Errorf("failed to parse AWS profile from loaded source: %s", err)
	}
	if err := awsProfile.Validate(); err != nil {
		return awsProfile, fmt.Errorf("invalid AWS profile [%s]: %s", profile, err)
	}

	return awsProfile, nil
}

var awsProfileCmd = &cobra.Command{
	Use:     "aw-profile-parser",
	Short:   "AWS profile parser",
	Long:    `AWS profile parser reads and validates a profile.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// read credentials file
		data, err := ioutil.ReadFile(credentialsFile)
		if err != nil {
			return fmt.Errorf("error reading file: %s", err)
		}

		awsProfile, err := ParseAWSProfile(data, profile)
		if err != nil {
			return err
		}

		result, err := json.MarshalIndent(awsProfile, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(result))
		return nil
	},
}

func main() {
	err := awsProfileCmd.Execute()
	if err != nil {
		log.Fatalf("Error executing command: %s", err)
	}
}

func init() {
	awsProfileCmd.PersistentFlags().StringVarP(&credentialsFile, "credentials", "c", "", "AWS credentials file")
	awsProfileCmd.PersistentFlags().StringVarP(&profile, "profile", "p", "default", "AWS profile name")
}
