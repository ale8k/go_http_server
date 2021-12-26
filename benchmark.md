This server:
BACKLOG: 0 
BLOCKING: ON
Running 10s test @ http://127.0.0.1:8000/1
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    16.19ms   69.91ms 564.35ms   96.38%
    Req/Sec     2.63k     1.45k    3.75k    72.13%
  16384 requests in 10.04s, 800.00KB read
  Socket errors: connect 0, read 16384, write 0, timeout 0
Requests/sec:   1632.31
Transfer/sec:     79.70KB

------------------------

Express:
Running 10s test @ http://127.0.0.1:8001/1
  4 threads and 100 connections
  Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    12.21ms    4.05ms  59.16ms   91.20%
    Req/Sec     2.08k   428.01     2.52k    85.25%
  82990 requests in 10.02s, 17.02MB read
Requests/sec:   8286.39
Transfer/sec:      1.70MB