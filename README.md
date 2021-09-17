# Churro - ETL for Kubernetes - spacex extension example

churro is a cloud-native Extract-Transform-Load (ETL) application designed to build, scale, and manage data pipeline applications.  That project can
be found at the [github churro site](https://github.com/churrodata/churro).

This particular project, spacex-geo-extension, is an example of
a custom churro extension you can write to perform any sort of
custom transformation logic that your use case might require.

## How does this work?
A churro extension has a gRPC interface that it will implement.  churro
registers your extension within a churro pipeline's extract source.

As data is being processed, churro will call any registered extensions
providing it a copy of the data being processed along with primary
keys for that data.

This allows you as within an extension the ability to know what
data is being processed and know its primary key location with
the pipeline database.  With that set of values, you can make ny
transform logic you would want to do.

This particular extension is pretty simple to show you the basics
of how a churro extension works as a starting point to more
complex extension writing.

This extension's container image is found on [DockerHub](https://hub.docker.com/repository/docker/churrodata/spacex-geo-extension).

## Design
* churro extensions run as Pods on a churro Kubernetes cluster
* churro extensions implement a gRPC interface, so you could write extensions in any language that support gRPC (e.g. Java, javascript, golang, etc.)
* You deploy your extension into your cluster as you would any other Pod

For more details on the churro design, checkout out the documentation at the [churro github pages](https://churrodata.github.io/churro/design-guide.html).

## Docs
Detailed documentation is found at the [churro github pages](https://churrodata.github.io/churro/), additional content such as blogs can be found at the [churrodata.com](https://www.churrodata.com) web site.

## Contributing
Since churro is open source, you can view the source code and make contributions such as pull requests on our github.com site.   We encourage users to log any issues they find on our github issues [site](https://github.com/churrodata/churro/issues).

## Support
churro enterprise support and services are provided by [churrodata.com](https://churrodata.com).

