# simple-douyin-backend
[![build](https://img.shields.io/badge/build-0.1.0-brightgreen)](https://github.com/StellarisW/douyin)[![go-version](https://img.shields.io/badge/go-~%3D1.18-30dff3?logo=go)](https://github.com/StellarisW/douyin)[![OpenTracing Badge](https://img.shields.io/badge/OpenTracing-disabled-blue.svg)](http://opentracing.io)

## 1. 项目简介
本项目实现了极简版抖音后端的互动方向和社交方向全部功能，项目采用gRPC进行RPC通信，
基于DDD思想划分微服务架构，合理拆解各项业务，便于团队协作开发，建立开发规范约束，
使用Kafka消息队列对其中的一些高频写操作进行异步处理，削峰处理用于提高系统的并发和可用性；
(待实现)使用Jaeger实现链路追踪；采用MySQL、Redis缓存与存储配合的方式来提高性能，减少MySQL压力；
通过Nginx限流、参数校验、存储加密、访问控制和边界情况处理等设计提高安全性；
并且充分编写了单元测试，完成了集成测试和压力测试。

## 2. 项目特点
### 2.1易维护
- **架构高度工程化**

​ 设计了最合适的代码架构，对业务开发友好，封装复杂度，以便代码的快速迭代
- **代码生成工具**

  通过IDL自定义代码生成模板，以适配项目需求
- **全局错误码定义**

  为代码的每个逻辑定义的特定的错误码，可以在前端获取错误码即可定位错误来源
- **日志组件**

  通过zlog库，实现日志全覆盖
- **项目配置(待实现Apollo配置中心)**
  
  对于本地配置，复制/configs/config.ini.sample为/configs/config.ini，并将配置修改为本地自身配置

  (待实现)使用 Apollo 配置管理平台，能进行灰度发布，分环境、分集群管理配置
- **链路追踪(待实现)**

  使用Jaeger进行链路追踪，查看请求的全链路，确定其各个部分的耗时以及错误来源

### 2.2 高可用

- **弹性设计，面向故障编程**

  以面向故障的思维方式编程，考虑到不同的级别的故障或者问题，有效保障了业务的连续性，持续性

- **限流，提供过载保护**

  通过判断特定条件来确定是否执行某个策略，只处理自己能力范围之内的请求，超量的请求会被限流。

- **熔断降级，保证服务稳定(待实现)**

  负载超过系统的承载能力时，系统会自动采取保护措施，立即中断服务，确保自身不被压垮。

### 2.3 高并发

- **服务负载均衡**

  将rpc服务集群化部署，通过ETCD进行服务发现与负载均衡，当业务流量高峰时，可以实现轻松水平扩展

- **消息队列**

  使用消息队列，将部分业务进行异步处理，从而实现流量削峰

- **通过协程共享调用**

  相同的多个请求只需要发起一次拿结果的调用，其他请求可以"坐享其成"，有效减少了资源服务的并发压力，可以有效防止缓存击穿。

### 2.4 高性能

- **缓存支持**

  service层使用缓存，提高数据读取速度

- **并发优化**

  如用户信息、视频信息的拼装使用协程加速拼装速度

- **特定业务逻辑使用算法优化**

  如用户关系计算使用集合操作

### 2.5 安全性

- **用户密码加密**

  将密码进行哈希操作后存储到数据库

- **JWT鉴权**

  不同的合法用户，拥有不同的唯一JWT，有效防止水平越权

- **严密的的边界情况处理**

  考虑的用户可能输入的边界情况，进行特殊的处理

## 3. 项目实现
### 3.1 项目的技术选型与开发文档
#### 项目场景分析
本项目的业务场景主要包括用户、视频、点赞、评论与消息功能，针对不同的业务场景进行不同的技术选型策略。
从大方向上，视频点赞、用户关注等需要能够及时响应，操作流畅，因此需要尽可能地降低耗时。而用户的基本信息、
视频的基本信息、评论与消息的内容需要进行持久化存储。

### 3.2 项目目录结构
项目的目录结构采用了[golang-standards/project-layout](https://github.com/golang-standards/project-layout)
中所推荐的项目目录结构，方便项目的后续迭代与开发。每个目录下都有对应的package_info文件用于介绍该目录的
作用。

在项目的Internal目录下，将项目分解为标准的Controller,Service和Dao的三层架构，项目的Controller层
用于处理用户的API输入并调用Service层微服务；项目的service层拥有不同的领域(一个领域对应一个微服务)
，领域之间的信息通过RPC同步通信或消息异步通信的方式传递信息；
项目的Dao层用于操作数据库，使用ORM模型一定程度上防止SQL注入。

### 3.3 技术选型
- **主框架**

| 名称    | 说明       |
|-------|----------|
| Hertz | HTTP框架   |
| gRPC  | RPC框架    |
| Gorm  | ORM框架    |

- **存储层**

| 名称      | 说明       |
|---------|----------|
| MySQL   | 关系型数据库   |
| Redis   | 缓存       |
| Kafka   | 消息队列     |
| OSS     | 阿里云对象存储  |
| Apollo  | 配置中心     |

- **其他工具**

| 名称         | 说明                 |
|------------|--------------------|
| Nginx      | 高性能web代理服务器        |
| Hertz-JWT  | token生成，鉴权         |
| SnowFlake  | google开源的雪花UUID生成库 |
| ETCD       | 服务注册与发现            |
| GoConvey   | 单元测试               |
| Jaeger     | 链路追踪               |
| Docker     | 容器部署               |
| Kubernetes | 容器自动管理平台           |
| fx         | 依赖注入框架             |
| bcrypt     | 对密码进行随机加盐Hash      |
| CI         | Github Action      |
| Zlog       | 日志工具               |

- HTTP框架：Hertz框架
  字节跳动内部开源团队CloudWeGo开源的高性能HTTP框架，具有开箱即用的API注册功能。基于Fast HTTP
,Gin等框架的优点，进行融合并实现。
- RPC框架：gRPC框架
  由Google开源的高性能RPC框架，拥有易拓展，易部署的特点。在全球范围内都有广泛的引用。
- ORM框架：Gorm框架
  全功能 ORM，golang中最主流的orm框架，作者jinzhu现就职于字节跳动。对 SQL 注入问题，Gorm 使用 database/sql 的参数占位符来构造 SQL 语句，可以自动转义参数，避免 SQL 注入。然后我们限制使用 db.Raw 来构建查询语句，保证数据库安全。
- 消息队列：Kafka
  Apache Kafka 是分布式发布-订阅消息系统，在 kafka 官网上对 kafka 的定义：一个分布式发布-订阅消息传递系统。
Kafka 最初由 LinkedIn 公司开发，Linkedin 于 2010 年贡献给了 Apache 基金会并成为顶级开源项目。
Kafka 的主要应用场景有：日志收集系统和消息系统。
- 配置中心：Apollo
  使用apollo配置中心方便集中化管理服务的配置，配置修改后能够实时推送到应用端，还能进行灰度发布，分环境、分集群管理配置。适用于微服务配置管理场景。
- 代理服务器：Nginx
  反向代理 Web 服务器，内存占用少，启动极快，高并发能力强。我们用nginx配置ip级别的限流和转发。
- 鉴权方案：JWT
  是一个开放标准，它定义了一种以紧凑和自包含的方法，用于在双方之间安全地传输编码为JSON 对象的信息。我们在项目中使用 JWT 进行鉴权服务，将其注入 Gin 路由，所有 GET 和 POST 请求都会优先执行鉴权，对注册的新用户，用id生成 token，对其他请求传入的 token，会对其进行解析判断是否正确。对于越权问题，在需求中并未要求设置不同的权限角色，我们默认所有用户权限相同，所以不存在垂直越权的问题；对于水平越权，不同的合法用户有对应且唯一的 token，可以很好地防止水平越权问题。
- 服务发现：Etcd
  基于 Raft 算法的 etcd 天生是一个强一致性高可用的服务存储目录，用户可以在 etcd 中注册服务，并且对注册的服务设置key TTL，定时保持服务的心跳以达到监控健康状态的效果。通过在 etcd 指定的主题下注册的服务也能在对应的主题下查找到。我们使用Etcd去实现服务的注册和发现。
- 链路追踪：Jaeger
  分布式链路跟踪系统，有web页面进行可视化分析。
- 对象存储：OSS
  腾讯云对象存储存放用户上传的视频。对象存储具有横向扩容能力，存储成本低，而且能自动截取视频封面。并且能帮业务服务器分摊请求压力。

### 3.4 开发规范
#### 3.4.1 安全规范
基于[腾讯的Go安全指南](https://github.com/Tencent/secguide/blob/main/Go%E5%AE%89%E5%85%A8%E6%8C%87%E5%8D%97.md)

#### 3.4.2 建表规范
例如
```sql
CREATE TABLE `tablename` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
  `created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
  `updated_at` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间，用于软删除',
  PRIMARY KEY (`id`),
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='表用途';
```
以上字段是每张表必须的，将id设为主键，字段的默认值和注释必须填写。
部分表的deleted_at类型进行了调整，和规范不同，在数据库表设计部分有说明。

#### 3.4.3 Controller层分层规范
```
├── gateway       //gateway网关微服务
├── api           //用户接口层
├── application   //应用层
├── cmd           //启动以及初始化依赖注入
└── rpc           //rpc调用函数
```

## 4. 更多信息
关于数据库表设计、redis设计、消息队列设计、不同业务逻辑的设计、安全性设计、项目的
总结与反思等都能够在/docs文件夹中找到
