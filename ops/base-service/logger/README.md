# 日志

## 日志最佳实践建议

1. The modern way to do logging is to just write to stdout or stderr (whichever you prefer) and let another specialized tool figure out what to do with the data.

It is okay to log in a separate goroutine if you want. It's probably better to wait until you actually have a bottleneck before doing that though.

You can wrap an http handler with another http handler that does logging for you. Check out https://medium.com/statuscode/how-i-write-go-http-services-after-seven-years-37c208122831 for a lot of ideas on how to organize it.

2. most of developers follow by ideas of cloudian integration avoiding of features by a platform (OS in this case).

3. logging in docker containers is best done to stdout so that the logs can be routed correctly to whatever system you use to process them.

4. https://geshan.com.np/blog/2019/03/follow-these-logging-best-practices-to-get-the-most-out-of-application-level-logging-slides/