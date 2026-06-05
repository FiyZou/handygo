# AI Collaboration Guide / AI 协作使用指南

This document is a user-facing guide. It explains how to use the collaboration workflow and is not a current task, requirement, or handoff.

本文档是面向使用者的说明文档，仅解释协作流程的使用方式，不代表当前任务、需求或交接内容。

## English

This project includes a closed-loop agent collaboration workflow. As a user, you only need to describe the goal. Agents maintain the project memory and role handoffs.

Example:

```text
Implement registration and login
```

After receiving the goal, the collaboration runner automatically executes:

```text
PM -> Architect -> Developer -> Reviewer
```

The agents are responsible for:

- updating `docs/product/PRD.md`
- updating `docs/tech/ARCHITECTURE.md`
- updating `docs/tasks.md`
- recording decisions in `docs/decision-log.md`
- updating `docs/handoff.md`
- implementing code and tests
- writing `docs/review/quality-report.md`
- updating `docs/qa/` when QA notes are needed

You should not manually maintain collaboration documents unless you intentionally want to correct or override the agent's understanding.

If work is interrupted, just say:

```text
Continue
```

The agent should read `docs/handoff.md` and resume from the current role.

## 中文

本项目包含一个闭环的 agent 协作流程。作为使用者，你只需要描述目标，项目记忆、角色流转和协作文档都由 agent 自动维护。

示例：

```text
实现注册登录
```

收到目标后，协作 runner 会自动执行：

```text
PM -> Architect -> Developer -> Reviewer
```

agent 会负责：

- 更新 `docs/product/PRD.md`
- 更新 `docs/tech/ARCHITECTURE.md`
- 更新 `docs/tasks.md`
- 将技术决策记录到 `docs/decision-log.md`
- 更新 `docs/handoff.md`
- 实现代码和测试
- 写入 `docs/review/quality-report.md`
- 需要 QA 记录时更新 `docs/qa/`

除非你要纠正或覆盖 agent 的理解，否则不需要手动维护这些协作文档。

如果工作被中断，只需要说：

```text
继续
```

agent 应该读取 `docs/handoff.md`，并从当前角色继续推进。
