# \ObjectApi

All URIs are relative to *https://www.openmobilealliance.org*

Method | HTTP request | Description
------------- | ------------- | -------------
[**FindObject**](ObjectApi.md#FindObject) | **Get** /{objectURI} | 



## FindObject

> Lwm2M FindObject(ctx, objectURI)



get details of a single object

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**objectURI** | **string**| version of object to retrieve | 

### Return type

[**Lwm2M**](LWM2M.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: text/xml

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

