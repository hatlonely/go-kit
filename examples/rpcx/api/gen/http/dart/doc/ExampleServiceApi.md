# openapi.api.ExampleServiceApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**exampleServiceAdd**](ExampleServiceApi.md#exampleServiceAdd) | **POST** /v1/add | 
[**exampleServiceEcho**](ExampleServiceApi.md#exampleServiceEcho) | **GET** /v1/echo | 


# **exampleServiceAdd**
> ApiAddRes exampleServiceAdd(body)



### Example 
```dart
import 'package:openapi/api.dart';

var api_instance = ExampleServiceApi();
var body = ApiAddReq(); // ApiAddReq | 

try { 
    var result = api_instance.exampleServiceAdd(body);
    print(result);
} catch (e) {
    print("Exception when calling ExampleServiceApi->exampleServiceAdd: $e\n");
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**ApiAddReq**](ApiAddReq.md)|  | 

### Return type

[**ApiAddRes**](ApiAddRes.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **exampleServiceEcho**
> ApiEchoRes exampleServiceEcho(message)



### Example 
```dart
import 'package:openapi/api.dart';

var api_instance = ExampleServiceApi();
var message = message_example; // String | 

try { 
    var result = api_instance.exampleServiceEcho(message);
    print(result);
} catch (e) {
    print("Exception when calling ExampleServiceApi->exampleServiceEcho: $e\n");
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **message** | **String**|  | [optional] [default to null]

### Return type

[**ApiEchoRes**](ApiEchoRes.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

