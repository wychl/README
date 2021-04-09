## 检查网格的策略检查状态。

```bash
kubectl -n istio-system get cm istio -o jsonpath="{@.data.mesh}" | grep disablePolicyChecks
disablePolicyChecks: false
```

# 配额管理算法

https://djlimiter.readthedocs.io/en/stable/strategy.html

- FIXED_WINDOW
固定窗口算法允许出现两倍的速率峰值，滑动窗口算法不会出现这样的情况

- ROLLING_WINDOW
滑动窗口算法提供更加精确的控制，但是也会提高对Redis 资源的消耗


上面两种情况，在超过时间窗口限制之后，都会自动恢复。



定义 quota 对象，在系统中会展现为一系列的计数器，计数器的维度就是quota 中的维度的笛卡尔积。如果在 validDuration 的时间窗口过期之前调用次数超过了 maxAmount 规定，Mixer 就会返回 RESOURCE_EXHAUSTED 给 Envoy，Envoy 则会反馈 429 代码给调用方。


- upstream connect error or disconnect/reset before headers

原因 Readiness probe failed: HTTP probe failed with statuscode: 503


# 白名单

curl  http://47.103.82.239:31380/ui

PERMISSION_DENIED:staticversion.istio-system:<your mesh source ip> is not whitelisted



curl -H'X-Forwarded-For:192.168.0.15' http://47.103.82.239:31380/ui
