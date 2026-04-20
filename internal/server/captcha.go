package server

import (
	"context"
	"errors"
	"fmt"

	recaptcha "cloud.google.com/go/recaptchaenterprise/v2/apiv1"
	recaptchapb "cloud.google.com/go/recaptchaenterprise/v2/apiv1/recaptchaenterprisepb"
	"google.golang.org/api/option"
)

type ReCaptchaParams struct {
	Key         string
	ProjectID   string
	Action      string
	Credentials string
}

/**
 * Create an assessment to analyze the risk of a UI action.
 *
 * @param projectID: Your Google Cloud Project ID.
 * @param recaptchaKey: The reCAPTCHA key associated with the site/app
 * @param token: The generated token obtained from the client.
 * @param recaptchaAction: Action name corresponding to the token.
 */
func CreateAssessment(
	projectID string,
	recaptchaKey string,
	token string,
	recaptchaAction string,
	credentialsFile string,
) (float32, error) {

	// Create the reCAPTCHA client.
	ctx := context.Background()
	client, err := recaptcha.NewClient(
		ctx,
		option.WithCredentialsFile(credentialsFile),
	)
	if err != nil {
		return 0, errors.New("Error creating reCAPTCHA client")
	}
	defer client.Close()

	// Set the properties of the event to be tracked.
	event := &recaptchapb.Event{
		Token:   token,
		SiteKey: recaptchaKey,
	}

	assessment := &recaptchapb.Assessment{
		Event: event,
	}

	// Build the assessment request.
	request := &recaptchapb.CreateAssessmentRequest{
		Assessment: assessment,
		Parent:     fmt.Sprintf("projects/%s", projectID),
	}

	response, err := client.CreateAssessment(
		ctx,
		request)

	if err != nil {
		return 0, fmt.Errorf("Error calling CreateAssessment: %v", err)
	}

	// Check if the token is valid.
	if !response.TokenProperties.Valid {
		return 0, fmt.Errorf(
			"Invalid Token: %v",
			response.TokenProperties.InvalidReason,
		)
	}

	// Check if the expected action was executed.
	if response.TokenProperties.Action != recaptchaAction {
		return 0, fmt.Errorf(
			"Mismatched action attribute. Expected [%s]. Actual [%s]",
			recaptchaAction,
			response.TokenProperties.Action,
		)
	}

	return response.RiskAnalysis.Score, nil

}
