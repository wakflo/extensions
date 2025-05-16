# Facebook Documentation

Workflow provide a way for applications to convey information to other applications in real-time. A webhook enables an app or server to provide other applications with real-time information, creating automated workflows. This document describes a simple Webhook integration.

## Requirements

To Obtain a Client ID and Client Secret:

1. Go to https://developers.facebook.com/
2. Make a new app, Select Other for usecase.
3. Choose Business as the type of app.
4. Fill the App Domain with Domain in Redirect URL.
5. Add new Product -> Facebook Login.
6. Navigate to Facebook Login Settings
7. Copy **Redirect Url Below** to "Valid OAuth Redirect URIs" and "Allowed Domains for the JavaScript SDK"
8. Create a new App Secret, then put the App ID and App Secret into Client ID and Client Secret.
   `