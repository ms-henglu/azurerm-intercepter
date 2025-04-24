# Introduction

The main purpose of this project is to intercept requests from `terraform-provider-azurerm` and dump the request body.

This project handles requests sent from `terraform-provider-azurerm`.  
It responds to the requests by returning 400 errors for all requests, and makes no requests to Azure.  
In the 400 error response, it includes the request body in the response. This allows you to see the request body that was sent from `terraform-provider-azurerm` to Azure.