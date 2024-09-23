# TestCases Reference

This document describes how one can define test cases to be used by the testing framework. Refer to the [examples](https://github.com/getyourguide/istio-config-validator/tree/main/examples) to understand how to test different use cases.

## TestCases

Test cases objects are interpreted by the framework to build up the mock and run the tests agaisnt the configuration.

| Field     | Type                    | Description                |
|-----------|-------------------------|----------------------------|
| testCases | [testCase[]](#TestCase) | List of test cases to run. |

## TestCase

TestCase defines each test that will be run sequentially.

| Field       | Type                                                                                                            | Description                                                |
|-------------|-----------------------------------------------------------------------------------------------------------------|------------------------------------------------------------|
| description | string                                                                                                          | Short description of what the testing is about.            |
| wantMatch   | bool                                                                                                            | If the test case should assert `true` or `false`           |
| request     | [request](#Request)                                                                                                         | Crafted requests that will mocked against VirtualServices  |
| route       | [HTTPRouteDestination[]](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPRouteDestination) | Route destinations that will be asserted for each request. |
| redirect    | [HTTPRedirect](https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPRedirect)          | Any redirect logic to test
| rewrite     | [HTTPRewrite](https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPRewrite)            | Any rewrite logic to test
| fault       | [HTTPFaultInjection](https://istio.io/latest/docs/reference/config/networking/virtual-service/#HTTPFaultInjection)   | Add fault injection (delay, abort) to test to test for timeouts, auth, etc.
| headers     | [Headers](https://istio.io/latest/docs/reference/config/networking/virtual-service/#Headers)                    | Test header manipulation rules. These are different from `request.headers`, i.e. headers present in the test request.
| delegate    | [Delegate](https://istio.io/latest/docs/reference/config/networking/virtual-service/#Delegate)                  | Any delegation logic to test


## Request

Request can contain more than one host (authority), method, uri, etc. The framework will mock requests in all possible combination defined here.


| Field     | Type              | Description                                                        |
|-----------|-------------------|--------------------------------------------------------------------|
| authority | string[]          | List of authority (host) that will be used to craft HTTP requests. |
| method    | string[]          | List of methods to craft requests.                                 |
| uri       | string[]          | List of URIs to craft requests.                                    |
| headers   | map[string]string | Headers present in all crafted requests.                           |
