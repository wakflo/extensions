// triggers/order_completed.md

# Order Completed

## Description

Automatically trigger workflows when a new order is marked as completed in your SendOwl account. This allows you to automate post-purchase processes such as customer onboarding, fulfillment, or marketing follow-ups.

## Properties

This trigger does not require any configuration properties.

## Details

- **Type**: sdkcore.TriggerTypePolling

This trigger polls the SendOwl API every 5 minutes to check for orders that have been marked as completed since the last check. When a new completed order is detected, it will trigger your workflow with the order details.
