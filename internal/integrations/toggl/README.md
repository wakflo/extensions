# Toggl Documentation

Workflow provide a way for applications to convey information to other applications in real-time. A webhook enables an app or server to provide other applications with real-time information, creating automated workflows. This document describes a simple Webhook integration.

## Requirements

Before you begin, make sure that you have the following:

* A running web server to receive the webhook payloads.
* A public URL for your web server - services like ngrok can help with this.

## Steps for Webhook Integration

1. **Setup Web Server:**
   To receive an HTTP POST request, your application needs to have an accessible route. Here is an example route configured in Go: