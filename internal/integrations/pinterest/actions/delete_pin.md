# Delete Pin

## Description

Deletes an existing pin from Pinterest. This action permanently removes a pin from your account.

## Details

- **Type**: core.ActionTypeAction

## Properties

| Name     | Type   | Required | Description                                            |
| -------- | ------ | -------- | ------------------------------------------------------ |
| board_id | String | Yes      | The unique identifier of the board containing the pin. |
| pin_id   | String | Yes      | The unique identifier of the pin to delete.            |

## Notes

- You must have ownership or write access to the pin to delete it
- Deletion is permanent and cannot be undone
- The pin will be removed from all boards where it appears
- Any repins of this pin by other users will not be affected
- Comments and likes on the pin will also be permanently deleted

## Error Handling

The action may fail if:

- The pin ID is invalid or doesn't exist
- You don't have permission to delete the pin
- The board ID doesn't match the pin's current board
- Authentication token is missing or invalid
