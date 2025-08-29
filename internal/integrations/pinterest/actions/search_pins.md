# Search Pins

## Description

Searches for pins on Pinterest matching specific criteria or keywords.

## Details

- **Type**: sdkcore.ActionTypeNormal

## Properties

| Name      | Type   | Required | Description                                                  |
| --------- | ------ | -------- | ------------------------------------------------------------ |
| query     | String | Yes      | The search term to look for pins.                            |
| page-size | Number | No       | Number of pins to return per page (maximum 250). Default: 25 |
| bookmark  | String | No       | Cursor for next set of items (pagination).                   |
| board-id  | String | No       | Optional board ID to limit search within a specific board.   |
