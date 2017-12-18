# SNS Hook for [Logrus](https://github.com/Sirupsen/logrus)

### Install
> $ go get github.com/stvvan/logrus_sns

### Usage
```
package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/stvvan/logrus_sns"
)

func main() {
	snsHook, err := logrus_sns.NewSNSHook("topic_arn", "subject", "us-east-1")

	if err != nil {
		panic(err)
	}

	log.AddHook(snsHook)

	log.WithFields(log.Fields{
		"error": "errorMsg",
	}).Error("error!")
}
```