# Get Pin Analytics

## Description

Retrieves analytics data for a specific Pinterest pin, including metrics such as impressions, saves, clicks, and engagement rates over a specified time period.

## Details

- **Type**: core.ActionTypeAction

## Properties

| Name          | Type   | Required | Description                                                                  |
| ------------- | ------ | -------- | ---------------------------------------------------------------------------- |
| board_id      | String | Yes      | The unique identifier of the board containing the pin.                       |
| pin_id        | String | Yes      | The unique identifier of the Pinterest pin to analyze.                       |
| metric_type   | String | Yes      | The metric to retrieve (e.g., IMPRESSION, SAVE, PIN_CLICK).                  |
| start_date    | String | Yes      | Start date for analytics data in YYYY-MM-DD format.                          |
| end_date      | String | Yes      | End date for analytics data in YYYY-MM-DD format.                            |
| app_types     | String | No       | Filter by app type: ALL, MOBILE, TABLET, or WEB.                             |
| split_field   | String | No       | Split data by: NO_SPLIT, APP_TYPE, OWNED_PIN, SOURCE, AGE_BUCKET, or GENDER. |
| ad_account_id | String | No       | Ad account ID for retrieving advertising-related metrics.                    |

## Available Metrics

- **IMPRESSION**: Number of times the pin was shown
- **SAVE**: Number of times the pin was saved
- **PIN_CLICK**: Number of clicks on the pin
- **OUTBOUND_CLICK**: Number of clicks leading away from Pinterest
- **TOTAL_COMMENTS**: Total number of comments
- **TOTAL_REACTIONS**: Total number of reactions
- **VIDEO_MRC_VIEW**: Media Rating Council accredited video views
- **VIDEO_AVG_WATCH_TIME**: Average time spent watching video
- **VIDEO_V50_WATCH_TIME**: Time to reach 50% video completion
- **QUARTILE_95_PERCENT_VIEW**: Views reaching 95% completion
- **VIDEO_10S_VIEW**: 10-second video views
- **VIDEO_START**: Number of video starts

## Notes

- Date range cannot exceed 90 days
- Start date must be before end date
- The pin must belong to the authenticated user or be accessible via the provided ad account
- Analytics data may have a delay of up to 24 hours
