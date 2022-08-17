[English](./README.md) | 简体中文

## 简介

katalyst 致力于解决云原生场景下的资源不合理利用问题，为资源管理和成本优化提供解决方案：
- QoS-Based 资源模型抽象：提供与业务场景匹配的资源 QoS 模型选择；
- 资源弹性管理：提供灵活可扩展的 HPA/VPA 资源弹性策略；
- 微拓扑及异构设备的调度、摆放：资源整体微拓扑感知调度、摆放，以及动态调整能力；
- 精细化资源分配、隔离：根据业务服务画像提供资源的精细化分配、出让和隔离

Katalyst 分为三个主要 Project：
- [Katalyst-API](https://github.com/kubewharf/katalyst-api.git) ：Katalyst 相关核心 API，包括 CRD、Protocol、QoS 定义等；
- [Katalyst-Core](https://github.com/kubewharf/katalyst-core.git) ：Katalyst 主体管控逻辑；
- [Charts](https://github.com/kubewharf/charts.git) ：Kubewharf 相关 Projects 的部署 helm charts；


更详细的介绍请参考 [Katalyst-Core](https://github.com/kubewharf/katalyst-core.git)

## 社区

### 贡献

若您期望成为 Katalyst 的贡献者，请参考 [CONTRIBUTING](CONTRIBUTING.md) 文档。

### 联系方式

如果您有任何疑问，欢迎提交 GitHub issues 或者 pull requests，或者联系我们的 [Maintainers](./MAINTAINERS.md)。

## 协议

Katalyst 采用 Apache 2.0 协议，协议详情请参考 [LICENSE](LICENSE)，另外 Katalyst 中的某些实现依赖于 Kubernetes 代码，此部分版权归属于 Kubernetes Authors。
