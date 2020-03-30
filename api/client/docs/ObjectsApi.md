# \ObjectsApi

All URIs are relative to *https://www.openmobilealliance.org*

Method | HTTP request | Description
------------- | ------------- | -------------
[**FindObjects**](ObjectsApi.md#FindObjects) | **Get** /Object | 



## FindObjects

> []ObjectMeta FindObjects(ctx, optional)



Returns metadata of of objects. Can be filtered by object ID and version 

### Required Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
 **optional** | ***FindObjectsOpts** | optional parameters | nil if no parameters

### Optional Parameters

Optional parameters are passed through a pointer to a FindObjectsOpts struct


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **objectVersion** | **optional.String**| version of object to retrieve | 
 **objectID** | **optional.Int32**| ID of object to retrieve | 

### Return type

[**[]ObjectMeta**](ObjectMeta.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

