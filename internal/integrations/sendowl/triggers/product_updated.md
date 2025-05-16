// triggers/product_updated.md

# Product Updated

## Description

Automatically trigger workflows when a product in your SendOwl account is updated. This allows you to maintain synchronized product information across systems or notify team members about product changes.

## Properties

This trigger does not require any configuration properties.

## Details

- **Type**: sdkcore.TriggerTypePolling

This trigger polls the SendOwl API every 5 minutes to check for products that have been updated since the last check. When a product update is detected, it will trigger your workflow with the product details.
