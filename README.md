# Google Cloud Build Pipeline Visualizer

[![build](https://github.com/RyanSiu1995/cloudbuild-visualizer/workflows/Go/badge.svg)](https://github.com/RyanSiu1995/cloudbuild-visualizer/workflows/Go/badge.svg)

For the current version of Google cloud build, it supports the async process with the variable `waitFor`. With the growth of complexity of your pipeline, it will be hard to maintain the async flow. Unlike Jenkins and CircleCI, there is no visualizer for your pipeline. This application aims at visualize the pipeline and help the developers to debug their cloud build.

## Rule of cloud build async process
From the Google [docs](https://cloud.google.com/cloud-build/docs/configuring-builds/configure-build-step-order), there are a few rules for the async process.
1. If no values are provided for waitFor, the build step waits for all prior build steps in the build request to complete successfully before running.
1. A step is dependent on every id in its waitFor and will not launch until each dependency has completed successfully.
1. By declaring that a step depends only on `-`, the step runs immediately when the build starts.

## How to install
1. Through `go get`
You can install the binary through the following command.
```
go get -u github.com/RyanSiu1995/gcb-visualizer
```
1. Install the pre-built binary
You can download the pre-built binary in [release page](https://github.com/RyanSiu1995/gcb-visualizer/releases) of this repo.

## How to use
You can visualize your pipeline with the following command.
```bash
gcb-visualizer <your-cloudbuild-yaml>
```
If you want to output the graph into other formats, you can use the output flag as the following.
```bash
gcb-visualizer --output my-pipeline.jpg <your-cloudbuild-yaml>
```
The current supported output formats are jpg, jpeg, dot and png.
