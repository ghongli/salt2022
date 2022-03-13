---
title: “Dapr”
date: 2022-03-04 16:35:39
---

### 集成模式

#### 基于消息的集成模式

一种把多个相互不相干的系统集成在一起的有效方法。

#### Saga 模式

1. 采取以事件为中心的视角，工作流响应不同的事件，并执行相应的动作。
   - 与工作流在未来状态与已经发生了的状态(表示为事件)之间进行协调。
2. 另一种对工作流建模的方式是将其视为一个`有限状态机`。简而言之，有限状态机具备一个所有可能状态的有限列表，并且根据输入以及当前状态，在状态之间进行转换。
   - Actor 可以为相关的状态切换，暴露出必要的方法，决定是否允许某个特定的切换。
   - 有限状态机，可以随时关闭并根据需要从其状态中恢复。除了要示状态的操作具备事务性之外，只能在预先定义的状态之间进行切换，要不能使用可能导致混乱的是间状态。

​		> 有限状态机更详细的讨论 ??
