Questions....

# Pagination, why are message not wrapped
```
"Messages": [
    {
      "DataSetWriterId": 4,
      "Timestamp": "2022-09-17T14:57:33.701Z",
      "subResource": "acme.com/rock_solid_weather_sensor/OEC-ACME-UTILITY/F12SN894",
      "Payload": {
        "health": "NORMAL_0",
        "healthScore": 100
      }
    },
    {
      "DataSetWriterId": 40,
      "Timestamp": "2022-09-17T14:57:33.701Z",
      "subResource": "acme.com/OEC%20Demo%20OT%20Connector/OEC-ACME-WEATHER-CONNECTOR/mojo",
      "Payload": {
        "totalCount": 1,
        "perPage": 1,
        "page": 1,
        "hasNext": false
      }
    }
  ],
  ```
  instead of 
  ```
  "Messages": [
    {
      "DataSetWriterId": 40,
      "Timestamp": "2022-09-17T14:57:33.701Z",
      "Source": "acme.com/rock_solid_weather_sensor/OEC-ACME-UTILITY/F12SN894",
      "Payload": {
        Pages: [ {
          "health": "NORMAL_0",
          "healthScore": 100
        }],
        "totalCount": 1,
        "perPage": 20,
        "page": 1,
        "hasNext": false
      }
    }
  ],
  ```

# Last Will
-> shouldn't the last will be used to change health (what happens with assets -> no will on their topics because will can only be registered for one topic, also timestamp is wrong -> time series store does not work (timestamp from start of application)

# DataSetClassId
some are missing GUIDs e.g. data, metadata

# Retain messages
Does it make sense to set the retain flag for MAM and health? So if the registry crashes and restarts it gets all information?

# MesageId
why not timestamp GUID? -> source info is encoded in "Source" or "Publisher" so no need for oi4identifier

# Webhook support for "bigger messages"
Maybe some kind of webhook request/info for other application to query blobs or more data if available

# Publication Mode
immediate publish also on Interval? what happens if application and asset both  "application" mode for data?