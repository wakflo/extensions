# SurveyMonkey Integration

## Description

Integrate SurveyMonkey with your workflow automation to streamline survey management and data collection. Connect your SurveyMonkey account to automatically create surveys, manage collectors, track responses, and analyze results within your existing workflows. Key features include:

- Create and manage surveys programmatically
- Set up automated response collection
- Trigger workflows based on new survey responses
- Extract and analyze survey data in real-time
- Integrate survey data with other business systems
- Send automated follow-ups based on specific survey answers

## Prerequisites

- A SurveyMonkey account (Professional plan or higher recommended for API access)
- API credentials from SurveyMonkey Developer portal
- Wakflo workflow automation account

## Setup Instructions

1. Log in to your SurveyMonkey account and navigate to the [Developer portal](https://developer.surveymonkey.com)
2. Create a new app in the Developer portal
3. Set the OAuth Redirect URL to the callback URL provided by Wakflo
4. Copy your Client ID and Client Secret
5. In Wakflo, navigate to the Integrations section and select SurveyMonkey
6. Enter your Client ID and Client Secret, then click "Connect"
7. Follow the OAuth authorization flow to grant Wakflo access to your SurveyMonkey account

## Scopes Required

The integration requires the following API scopes:

- `create_surveys`
- `view_surveys`
- `edit_surveys`
- `view_collectors`
- `create_collectors`
- `view_responses`
- `create_responses`

## Actions

| Name                  | Description                                              | Link                                     |
| --------------------- | -------------------------------------------------------- | ---------------------------------------- |
| Get Survey            | Retrieves detailed information about a specific survey   | [docs](actions/get_survey.md)            |
| List Surveys          | Retrieves a list of surveys in your SurveyMonkey account | [docs](actions/list_surveys.md)          |
| Get Survey Response   | Retrieves a specific response to a survey                | [docs](actions/get_survey_response.md)   |
| List Survey Responses | Retrieves all responses for a specific survey            | [docs](actions/list_survey_responses.md) |
| Create Survey         | Creates a new survey in your SurveyMonkey account        | [docs](actions/create_survey.md)         |
| Create Collector      | Creates a new collector for a specific survey            | [docs](actions/create_collector.md)      |

## Triggers

| Name               | Description                                                        | Link                                   |
| ------------------ | ------------------------------------------------------------------ | -------------------------------------- |
| Survey Created     | Triggers when a new survey is created in your SurveyMonkey account | [docs](triggers/survey_created.md)     |
| Response Completed | Triggers when a new response is completed for a survey             | [docs](triggers/response_completed.md) |

## Troubleshooting

### Common Issues

1. **API Rate Limits**: SurveyMonkey enforces rate limits on API calls. If you encounter rate limit errors, try implementing a back-off strategy or spreading out your automation workflows.

2. **Permissions**: Ensure your SurveyMonkey account has the appropriate permissions for the actions you're attempting to perform. Some features may require a higher-tier plan.

3. **Authentication Errors**: If you encounter authentication errors, try reconnecting your SurveyMonkey account in Wakflo.

### Getting Help

For additional support, please contact:

- Wakflo Support: support@wakflo.com
- [SurveyMonkey API Documentation](https://developer.surveymonkey.com/api/v3/)

## Categories

- survey
- data collection

## Authors

- Wakflo <integrations@wakflo.com>
