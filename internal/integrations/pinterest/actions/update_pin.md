# Update Pin

## Description

Updates the details of an existing Pinterest pin, allowing modifications to title, description, link, alt text, and note fields.

## Details

- **Type**: core.ActionTypeAction

## Properties

| Name        | Type   | Required | Description                                                     |
| ----------- | ------ | -------- | --------------------------------------------------------------- |
| board_id    | String | Yes      | The unique identifier of the board containing the pin.          |
| pin_id      | String | Yes      | The unique identifier of the Pinterest pin to update.           |
| title       | String | No       | New title for the pin (max 100 characters).                     |
| description | String | No       | New description for the pin (max 800 characters).               |
| link        | String | No       | New link URL for the pin (must start with http:// or https://). |
| alt_text    | String | No       | New alternative text for accessibility (max 500 characters).    |
| note        | String | No       | New note to add to the pin (max 500 characters).                |

## Notes

- At least one optional field must be provided for the update to proceed
- Only the fields provided will be updated; other fields will remain unchanged
- The pin must belong to the authenticated user to be updated
