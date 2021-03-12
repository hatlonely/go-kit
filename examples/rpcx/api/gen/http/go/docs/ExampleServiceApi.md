# \ExampleServiceApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ExampleServiceAdd**](ExampleServiceApi.md#ExampleServiceAdd) | **Post** /v1/add | 
[**ExampleServiceEcho**](ExampleServiceApi.md#ExampleServiceEcho) | **Get** /v1/echo | 



## ExampleServiceAdd

> ApiAddRes ExampleServiceAdd(ctx, body)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**body** | [**ApiAddReq**](ApiAddReq.md)|  | 

### Return type

[**ApiAddRes**](apiAddRes.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ExampleServiceEcho

> ApiEchoRes ExampleServiceEcho(ctx, optional)



### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***ExampleServiceEchoOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a ExampleServiceEchoOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **message** | **optional.String**|  | 

### Return type

[**ApiEchoRes**](apiEchoRes.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

