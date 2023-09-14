#  文件监控

用来监控文件以及其子目录下的所有文件的增删改

此脚本[基于](https://gitee.com/piri47/AWD-1/blob/master/%E6%96%87%E4%BB%B6%E7%9B%91%E6%8E%A7.py)

做出的改进：
    1. 不会删除原文件，避免因为某些意外情况而导致的原文件被删除
    2. 前端打印所有详细信息，包括文件路径、被修改的时间、内容等

+ 用法

在需要监控的目录下直接启动即可
```go
monitor-Go.exe
```

+ 效果图

![](https://gallery-1304405887.cos.ap-nanjing.myqcloud.com/markdown33-14-23-133303.png)