# Changelog


## v0.1.0

The API for test cases does not cover all aspects of VirtualServices.

-   Supported [HTTPMatchRequests](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPMatchRequest) fields to match requests against are: `authority`, `method`, `headers` and `uri`.
    -   Not supported ones: `scheme`, `port`, `queryParams`, etc.

-   Supported assert against [HTTPRouteDestination](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPRouteDestination)
    -   Not supported ones: [HTTPRedirect](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPRedirect), [HTTPRewrite](https://istio.io/docs/reference/config/networking/virtual-service/#HTTPRewrite), etc.
