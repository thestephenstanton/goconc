# Concurrency Download Challenge

## Problem Description

Implement a function in Go that, given a list of URLs and an existing download function, performs the following:

- Downloads data from all the URLs **concurrently**.
- Merges the downloaded results into a single map of the form:
  ```go
  map[string]int // mapping URL to its downloaded integer data

