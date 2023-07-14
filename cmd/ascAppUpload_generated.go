// Code generated by piper's step-generator. DO NOT EDIT.

package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/SAP/jenkins-library/pkg/config"
	"github.com/SAP/jenkins-library/pkg/log"
	"github.com/SAP/jenkins-library/pkg/splunk"
	"github.com/SAP/jenkins-library/pkg/telemetry"
	"github.com/SAP/jenkins-library/pkg/validation"
	"github.com/spf13/cobra"
)

type ascAppUploadOptions struct {
	ServerURL          string `json:"serverUrl,omitempty"`
	AppToken           string `json:"appToken,omitempty"`
	AppID              string `json:"appId,omitempty"`
	FilePath           string `json:"filePath,omitempty"`
	JamfTargetSystem   string `json:"jamfTargetSystem,omitempty"`
	ReleaseAppVersion  string `json:"releaseAppVersion,omitempty"`
	ReleaseDescription string `json:"releaseDescription,omitempty"`
	ReleaseDate        string `json:"releaseDate,omitempty"`
	ReleaseVisible     bool   `json:"releaseVisible,omitempty"`
}

// AscAppUploadCommand Upload an app to ASC
func AscAppUploadCommand() *cobra.Command {
	const STEP_NAME = "ascAppUpload"

	metadata := ascAppUploadMetadata()
	var stepConfig ascAppUploadOptions
	var startTime time.Time
	var logCollector *log.CollectorHook
	var splunkClient *splunk.Splunk
	telemetryClient := &telemetry.Telemetry{}

	var createAscAppUploadCmd = &cobra.Command{
		Use:   STEP_NAME,
		Short: "Upload an app to ASC",
		Long: `With this step you can upload an app to ASC.
It creates a new release note in ASC and uploads the binary to ASC and therewith to Jamf.
For more information about ASC, check out [Application Support Center](https://github.com/SAP/application-support-center).`,
		PreRunE: func(cmd *cobra.Command, _ []string) error {
			startTime = time.Now()
			log.SetStepName(STEP_NAME)
			log.SetVerbose(GeneralConfig.Verbose)

			GeneralConfig.GitHubAccessTokens = ResolveAccessTokens(GeneralConfig.GitHubTokens)

			path, _ := os.Getwd()
			fatalHook := &log.FatalHook{CorrelationID: GeneralConfig.CorrelationID, Path: path}
			log.RegisterHook(fatalHook)

			err := PrepareConfig(cmd, &metadata, STEP_NAME, &stepConfig, config.OpenPiperFile)
			if err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}
			log.RegisterSecret(stepConfig.AppToken)

			if len(GeneralConfig.HookConfig.SentryConfig.Dsn) > 0 {
				sentryHook := log.NewSentryHook(GeneralConfig.HookConfig.SentryConfig.Dsn, GeneralConfig.CorrelationID)
				log.RegisterHook(&sentryHook)
			}

			if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
				splunkClient = &splunk.Splunk{}
				logCollector = &log.CollectorHook{CorrelationID: GeneralConfig.CorrelationID}
				log.RegisterHook(logCollector)
			}

			if err = log.RegisterANSHookIfConfigured(GeneralConfig.CorrelationID); err != nil {
				log.Entry().WithError(err).Warn("failed to set up SAP Alert Notification Service log hook")
			}

			validation, err := validation.New(validation.WithJSONNamesForStructFields(), validation.WithPredefinedErrorMessages())
			if err != nil {
				return err
			}
			if err = validation.ValidateStruct(stepConfig); err != nil {
				log.SetErrorCategory(log.ErrorConfiguration)
				return err
			}

			return nil
		},
		Run: func(_ *cobra.Command, _ []string) {
			stepTelemetryData := telemetry.CustomData{}
			stepTelemetryData.ErrorCode = "1"
			handler := func() {
				config.RemoveVaultSecretFiles()
				stepTelemetryData.Duration = fmt.Sprintf("%v", time.Since(startTime).Milliseconds())
				stepTelemetryData.ErrorCategory = log.GetErrorCategory().String()
				stepTelemetryData.PiperCommitHash = GitCommit
				telemetryClient.SetData(&stepTelemetryData)
				telemetryClient.Send()
				if len(GeneralConfig.HookConfig.SplunkConfig.Dsn) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.Dsn,
						GeneralConfig.HookConfig.SplunkConfig.Token,
						GeneralConfig.HookConfig.SplunkConfig.Index,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
				if len(GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint) > 0 {
					splunkClient.Initialize(GeneralConfig.CorrelationID,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblEndpoint,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblToken,
						GeneralConfig.HookConfig.SplunkConfig.ProdCriblIndex,
						GeneralConfig.HookConfig.SplunkConfig.SendLogs)
					splunkClient.Send(telemetryClient.GetData(), logCollector)
				}
			}
			log.DeferExitHandler(handler)
			defer handler()
			telemetryClient.Initialize(GeneralConfig.NoTelemetry, STEP_NAME)
			ascAppUpload(stepConfig, &stepTelemetryData)
			stepTelemetryData.ErrorCode = "0"
			log.Entry().Info("SUCCESS")
		},
	}

	addAscAppUploadFlags(createAscAppUploadCmd, &stepConfig)
	return createAscAppUploadCmd
}

func addAscAppUploadFlags(cmd *cobra.Command, stepConfig *ascAppUploadOptions) {
	cmd.Flags().StringVar(&stepConfig.ServerURL, "serverUrl", os.Getenv("PIPER_serverUrl"), "The URL to the ASC backend")
	cmd.Flags().StringVar(&stepConfig.AppToken, "appToken", os.Getenv("PIPER_appToken"), "App token used to authenticate with the ASC backend")
	cmd.Flags().StringVar(&stepConfig.AppID, "appId", os.Getenv("PIPER_appId"), "The app ID in ASC")
	cmd.Flags().StringVar(&stepConfig.FilePath, "filePath", os.Getenv("PIPER_filePath"), "The path to the app binary")
	cmd.Flags().StringVar(&stepConfig.JamfTargetSystem, "jamfTargetSystem", os.Getenv("PIPER_jamfTargetSystem"), "The jamf target system")
	cmd.Flags().StringVar(&stepConfig.ReleaseAppVersion, "releaseAppVersion", `Pending Release`, "The new app version name to be created in ASC")
	cmd.Flags().StringVar(&stepConfig.ReleaseDescription, "releaseDescription", `<p>TBD</p>`, "The new release description")
	cmd.Flags().StringVar(&stepConfig.ReleaseDate, "releaseDate", os.Getenv("PIPER_releaseDate"), "The new release date (Format: MM/DD/YYYY) Default is the current date")
	cmd.Flags().BoolVar(&stepConfig.ReleaseVisible, "releaseVisible", false, "The new release visible flag")

	cmd.MarkFlagRequired("serverUrl")
	cmd.MarkFlagRequired("appId")
	cmd.MarkFlagRequired("filePath")
	cmd.MarkFlagRequired("jamfTargetSystem")
}

// retrieve step metadata
func ascAppUploadMetadata() config.StepData {
	var theMetaData = config.StepData{
		Metadata: config.StepMetadata{
			Name:        "ascAppUpload",
			Aliases:     []config.Alias{},
			Description: "Upload an app to ASC",
		},
		Spec: config.StepSpec{
			Inputs: config.StepInputs{
				Secrets: []config.StepSecrets{
					{Name: "ascAppTokenCredentialsId", Description: "Jenkins secret text credential ID containing the authentication token for the ASC app", Type: "jenkins"},
				},
				Parameters: []config.StepParameters{
					{
						Name:        "serverUrl",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{{Name: "ascServerUrl"}},
						Default:     os.Getenv("PIPER_serverUrl"),
					},
					{
						Name: "appToken",
						ResourceRef: []config.ResourceReference{
							{
								Name:    "ascVaultSecretName",
								Type:    "vaultSecret",
								Default: "asc",
							},

							{
								Name: "ascAppTokenCredentialsId",
								Type: "secret",
							},
						},
						Scope:     []string{"PARAMETERS"},
						Type:      "string",
						Mandatory: false,
						Aliases:   []config.Alias{{Name: "ascAppToken"}},
						Default:   os.Getenv("PIPER_appToken"),
					},
					{
						Name:        "appId",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_appId"),
					},
					{
						Name:        "filePath",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_filePath"),
					},
					{
						Name:        "jamfTargetSystem",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   true,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_jamfTargetSystem"),
					},
					{
						Name:        "releaseAppVersion",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `Pending Release`,
					},
					{
						Name:        "releaseDescription",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     `<p>TBD</p>`,
					},
					{
						Name:        "releaseDate",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "string",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     os.Getenv("PIPER_releaseDate"),
					},
					{
						Name:        "releaseVisible",
						ResourceRef: []config.ResourceReference{},
						Scope:       []string{"PARAMETERS", "STAGES", "STEPS"},
						Type:        "bool",
						Mandatory:   false,
						Aliases:     []config.Alias{},
						Default:     false,
					},
				},
			},
		},
	}
	return theMetaData
}
