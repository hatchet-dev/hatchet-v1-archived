# UpdateMonitorRequest

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CronSchedule** | **string** |  | [optional] [default to null]
**Description** | **string** |  | [optional] [default to null]
**Disabled** | **bool** | Whether the monitor is disabled. In order to turn the monitor off for all modules, set disabled&#x3D;true. Passing in an empty module list will trigger this monitor for all modules. | [optional] [default to null]
**Kind** | **string** |  | [optional] [default to null]
**Modules** | **[]string** | A list of module ids. If empty or omitted, this monitor targets all modules. | [optional] [default to null]
**Name** | **string** |  | [optional] [default to null]
**PolicyBytes** | **string** |  | [optional] [default to null]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)

