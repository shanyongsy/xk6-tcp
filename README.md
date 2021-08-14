# xk6-tcp
* TCP extension components used in K6
* Link: https://k6.io/blog/extending-k6-with-xk6/#creating-the-k6-extension


## Building the extension with xk6
```
First install xk6 with:
$ go install github.com/k6io/xk6/cmd/xk6@latest

Next, let's build the k6 binary. To use the published version of the extension run:
$ xk6 build v0.32.0 --with github.com/shanyongsy/xk6-tcp

Or if you're working with a local directory run the following, replacing the path as needed:
$ xk6 build v0.32.0 --with github.com/shanyongsy/xk6-tcp="/home/shanyong/code/k6code/tcp-extension/xk6-tcp
```

## K6 Script
[k6-tcp script example](https://github.com/shanyongsy/xk6-tcp/blob/main/example/loadtest/test_tcp.js)

## Tcp Server
[Test server](https://github.com/shanyongsy/tcp-server-client-go)

```mermaid
gantt
       dateFormat                YYYY-MM-DD
       title                     Adding GANTT diagram functionality to mermaid
       excludes                  :excludes the named dates/days from being included in a charted task..
       section A section
       Completed task            :done,    des1, 2014-01-06,2014-01-08
       Active task               :active,  des2, 2014-01-09, 3d
       Future task               :         des3, after des2, 5d
       Future task2              :         des4, after des3, 5d

       section Critical tasks
       Completed task in the critical line :crit, done, 2014-01-06,24h
       Implement parser and jison          :crit, done, after des1, 2d
       Create tests for parser             :crit, active, 3d
       Future task in critical line        :crit, 5d
       Create tests for renderer           :2d
       Add to mermaid                      :1d

       section Documentation
       Describe gantt syntax               :active, a1, after des1, 3d
       Add gantt diagram to demo page      :after a1  , 20h
       Add another diagram to demo page    :doc1, after a1  , 48h

       section Last section
       Describe gantt syntax               :after doc1, 3d
       Add gantt diagram to demo page      :20h
       Add another diagram to demo page    :48h
```
