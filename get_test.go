// Copyright (c) 2018 CA. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	//	"fmt"
	graphql "github.com/neelance/graphql-go"
	"io/ioutil"
	"log"
	"testing"
)

var testschema *graphql.Schema = graphql.MustParseSchema(Schema, &Resolver{})

func simpletest(t *testing.T, query string, result string) {
	var stest Test
	stest.Schema = testschema
	stest.Query = query
	stest.ExpectedResult = result
	RunTest(t, &stest)
}

func init() {
	cache := make(map[string]interface{})
	ctx := context.WithValue(context.Background(), "queryCache", &cache)
	setTestContext(&ctx)
	for _, fname := range []string{
		"deployment.json",
		"replicaset.json",
		"pod1.json"} {
		addToCache(&cache, "testdata/"+fname)
	}
}

func addToCache(cacheref *map[string]interface{}, fname string) {
	bytes, err := ioutil.ReadFile(fname)
	if err != nil {
		log.Fatal(err)
	}
	data := fromJson(bytes).(JsonObject)
	kind := data["kind"].(string)
	ns := getNamespace(data)
	name := getName(data)
	res := mapToResource(*getTestContext(), data)
	(*cacheref)[rawCacheKey(kind, ns, name)] = data
	(*cacheref)[cacheKey(kind, ns, name)] = res
	(*cacheref)[kind] = []resource{res}
}

func TestPods(t *testing.T) {
	simpletest(
		t,
		`{allDeployments() {
            metadata {
              creationTimestamp
              generation
              labels { name value }
            }
            spec {
              minReadySeconds
              paused
              progressDeadlineSeconds
              replicas
              revisionHistoryLimit
              selector {
                matchLabels { name value }
                matchExpressions {
                  key
                  operator
                  values
                }
              }
              strategy {
                type
                rollingUpdate {
                  maxSurgeInt
                  maxSurgeString
                  maxUnavailableInt
                  maxUnavailableString
                }
              }
              template {
                metadata {
                  creationTimestamp
                  labels { name value }
                },
                spec {
                  dnsPolicy
                  restartPolicy
                  schedulerName
                  terminationGracePeriodSeconds
                  volumes {
                    name
                    persistentVolumeClaim { claimName readOnly }
                  }
                }
              }
            }
         }}`,
		`{"allDeployments": [
            {
              "metadata": {
                "creationTimestamp": "2018-07-02T14:53:53Z",
                "generation": 1,
                "labels": [
                  {"name": "app", "value": "clunky-sabertooth-joomla"},
                  {"name": "chart", "value": "joomla-2.0.2"},
                  {"name": "heritage", "value": "Tiller"},
                  {"name": "release", "value": "clunky-sabertooth"}
                 ]
              },
              "spec": {
                "minReadySeconds": 0,
                "paused": false,
                "progressDeadlineSeconds": 600,
                "replicas": 1,
                "revisionHistoryLimit": 10,
                "selector": {
                  "matchExpressions": [],
                  "matchLabels": [
                    {"name": "app", "value": "clunky-sabertooth-joomla"}
                  ]
                },
                "strategy": {
                  "rollingUpdate": {
                    "maxSurgeInt": 1,
                    "maxSurgeString": null,
                    "maxUnavailableInt": 1,
                    "maxUnavailableString": null
                  },
                  "type": "RollingUpdate"
                },
                "template": {
                  "metadata": {
                    "creationTimestamp": null,
                    "labels": [
                      {"name": "app", "value": "clunky-sabertooth-joomla"}
                    ]
                  },
                  "spec": {
                    "dnsPolicy": "ClusterFirst",
                    "restartPolicy": "Always",
                    "schedulerName": "default-scheduler",
                    "terminationGracePeriodSeconds": 30,
                    "volumes": [
                      {
                        "name": "joomla-data",
                        "persistentVolumeClaim": {
                          "claimName": "clunky-sabertooth-joomla-joomla",
                          "readOnly": false
                        }
                      },
                      {
                        "name": "apache-data",
                        "persistentVolumeClaim": {
                          "claimName": "clunky-sabertooth-joomla-apache",
                          "readOnly": false
                        }
                      }
                    ]
                  }
                }
              }
            }
           ]
         }`)
}
