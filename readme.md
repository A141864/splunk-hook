
# Splunk Hook
This hooks forwards custom JSON formatted events to the Splunk HTTP collector.

Features
--------
* Default fields

Installation
------------
```sh
go get -u github.com/flynnhandley/splunk-hook
```

Examples
--------

```go
package main

import (
	"github.com/flynnhandley/slack-hook"
)

func main() {
	
    token := "xxxxxxx-xxxxx-xxxxx-xxxxx-xxxxxxx"

	defaultValues := log.Fields{
		"CustomField1": "Some Value",
		"CustomField2": "Another Value",
	}
	log.AddHook(NewHook(
		"https://127.0.0.1:8088/services/collector",
		token,
		"SourceApp",
		"json",
		"default",
		defaultValues,
	))

	log.WithFields(log.Fields{
		"CustomField": "it works!",
	}).Error("This is the error message")
}
```


Documentation
-------------


License
-------
