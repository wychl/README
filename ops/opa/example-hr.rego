package httpapi.authz

import input

# 允许HR查看所有人的工资
allow {
  input.method == "GET"
  input.path = ["finance", "salary", _]
  input.user == hr[_]
}

# David is the only member of HR.
hr = [
  "david",
]
