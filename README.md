# go-icinga2-client

Icinga2 API client.

## Getting started

```
import "github.com/Nexinto/go-icinga2-client/icinga2"

icinga, err := icinga2.New(icinga2.WebClient{
		URL:         "https://icinga.somewhere.com:5665,
		Username:    "icinga",
		Password:    "secret",
		Debug:       true,
		InsecureTLS: false})
```

### List hostgroups

```
hostGroups, err := icinga.ListHostGroups()
```

### Create a hostgroup

```
icinga.CreateHostGroup(icinga2.HostGroup{"mygroup"})
```

### Delete a hostgroup

```
icinga.DeleteHostGroup("mygroup")
```

## Supported Icinga objects

So far, supported are hostgroups, hosts and services.
