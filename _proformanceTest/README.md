# Proformance Test

> webserver: github.com/simba-fs/echoServer

> test connamd: ab -n 1000000 -c 100 <url>

# Only Webserver
```
Server Software:        
Server Hostname:        localhost
Server Port:            3000

Document Path:          /
Document Length:        12 bytes

Concurrency Level:      100
Time taken for tests:   53.510 seconds
Complete requests:      1000000
Failed requests:        0
Total transferred:      135000000 bytes
HTML transferred:       12000000 bytes
Requests per second:    18688.16 [#/sec] (mean)
Time per request:       5.351 [ms] (mean)
Time per request:       0.054 [ms] (mean, across all concurrent requests)
Transfer rate:          2463.77 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.5      2       6
Processing:     0    3   0.7      3      15
Waiting:        0    2   0.7      2      15
Total:          0    5   0.7      5      18

Percentage of the requests served within a certain time (ms)
  50%      5
  66%      5
  75%      6
  80%      6
  90%      6
  95%      7
  98%      8
  99%      8
 100%     18 (longest request)
```

# proxy webserver
```toml
address = '0.0.0.0:4000'

[host.test1]
from = 'test1.localhost:3000'
to = 'http://localhost:3000'
```
```
Server Software:        
Server Hostname:        localhost
Server Port:            4000

Document Path:          /
Document Length:        339 bytes

Concurrency Level:      100
Time taken for tests:   52.660 seconds
Complete requests:      1000000
Failed requests:        0
Non-2xx responses:      1000000
Total transferred:      450000000 bytes
HTML transferred:       339000000 bytes
Requests per second:    18989.67 [#/sec] (mean)
Time per request:       5.266 [ms] (mean)
Time per request:       0.053 [ms] (mean, across all concurrent requests)
Transfer rate:          8345.07 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.5      2       7
Processing:     0    3   0.7      3      21
Waiting:        0    2   0.7      2      19
Total:          0    5   0.7      5      23

Percentage of the requests served within a certain time (ms)
  50%      5
  66%      5
  75%      6
  80%      6
  90%      6
  95%      6
  98%      7
  99%      8
 100%     23 (longest request)
```

# Host static file
```toml
address = '0.0.0.0:4000'

[host.test2]
from = 'test2.localhost:4000'
to = 'static://test2'

[static.test2]
repo = 'https://github.com/simba-fs/simba-fs'
branch = 'master'
```
```
Server Software:        
Server Hostname:        localhost
Server Port:            4000

Document Path:          /
Document Length:        339 bytes

Concurrency Level:      100
Time taken for tests:   53.295 seconds
Complete requests:      1000000
Failed requests:        0
Non-2xx responses:      1000000
Total transferred:      450000000 bytes
HTML transferred:       339000000 bytes
Requests per second:    18763.32 [#/sec] (mean)
Time per request:       5.330 [ms] (mean)
Time per request:       0.053 [ms] (mean, across all concurrent requests)
Transfer rate:          8245.60 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    2   0.5      2       6
Processing:     0    3   0.8      3      27
Waiting:        0    2   0.7      2      26
Total:          0    5   0.8      5      29

Percentage of the requests served within a certain time (ms)
  50%      5
  66%      5
  75%      6
  80%      6
  90%      6
  95%      7
  98%      8
  99%      8
 100%     29 (longest request)
```
