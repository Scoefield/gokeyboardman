---
highlight: a11y-dark
---


## ES 简介及历史背景

### ES 简介

- 开源分布式搜索分析引擎
    - 近实时
    - 分布式存储/搜索/分析引擎
- 提供 restful api，可以让任何编程语言调用
- 支持多种方式集成接入
    - 多种编程语言类库
    - JDBC & ODBC（新版本支持）
- ES 软件下载量，超 3.5 亿次

网上统计 ES 在搜索引擎中的使用排名，当然国内的一些大厂也有使用，比如：滴滴、小米、360等。

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/1e518900817842549ac680bdaed94807~tplv-k3u1fbpfcp-watermark.image)


### ES 历史背景

- ES 起源 - Lucene
    - 创建于 1999 年，2005 年成为 Apache 顶级开源项目
    - 基于 Java 语言开发
    - Lucene 具有高性能、易扩展的优点
- ES 诞生 - Shay Banon
    - 2004 年 Shay Banon 基于 Lucene 开发了 Compass
    - 2010 年 Shay Banon 重写 Compass 取名 Elasticsearch

### ES 主要作用

- 海量数据的分布式存储以及集群管理
    - 服务与数据的高可用，水平扩张
- 近实时搜索，性能卓越
    - 结构化/全文/地理位置...
- 数据分析
    - 聚合功能
    
### ES 应用场景

主要应用场景分两大类：

- 搜索类(带上聚合)，与现有数据库进行同步，通过ES进行查询聚合
- 日志类，包括日志收集，指标性收集，通过 beats 等工具收集到 kafka 中，通过 logstash 进行转换，输送到 ES 中，然后通过 Kibana 进行展示。
    
## ES 生态圈

![image.png](https://p3-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/08d3206daa9e475db2b4c89ba67a3c05~tplv-k3u1fbpfcp-watermark.image)

### Logstash：数据处理管道

- 开源的服务端数据处理管道，支持从不同的来源采集数据，然后转换处理和聚合数据，并将数据发送到不同的存储中
- 诞生于 2009 年，最初用来做日志的采集与处理
- Logstash 创始人 Jordan Sisel
- 2013 年被 Elasticsearch 收购
- 特性...

### Beats：轻量的数据采集器

- go 语言开发，速度快
- 集成丰富的套件，比如：filebeat，实时的日志文件抓取收集

### Kibana：可视化分析利器

- 创始人- Rashid Khan
- 数据可视化工具，帮助用户查看和分析数据
- 也是 2013 年加入 Elasticsearch 公司
- 特性...

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/ca2c85d7cb6b4fb2ac02803b976ccea6~tplv-k3u1fbpfcp-watermark.image)

![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/a037129575274117b5da265855dd6a42~tplv-k3u1fbpfcp-watermark.image)

## Elasticsearch 与数据库的集成

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/c01522f24cb445619790b173ad2a7862~tplv-k3u1fbpfcp-watermark.image)

### 单独使用 ES 作为存储

### 与数据库集成使用

- 与现有系统的集成
- 数据更新频繁

## ES 指标分析/日志分析

![image.png](https://p1-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/8237c5dfc4ef4a689f6e1e69faf18c82~tplv-k3u1fbpfcp-watermark.image)


## ES 索引、Mapping、文档和 Restful Api

### 索引 - Index

- 索引是文档的容器，是一类文档的集合
    - Index 体现逻辑空间的概念，每个索引都有自己的 Mapping 定义，用于定义包含文档的字段名和字段类型
- 索引的 Mapping 与 Settings
    - Mapping 定义文档字段的类型
    - Settings 定义不同的数据分布
- Type
    - 同一个 Type 下可以定义多个相同结构的文档（类似一个表）
    - ES 7.0 之前，一个 Index 可以设置多个 Type
    - 7.0 开始，一个索引只能创建一个 Type -> "_doc"，（默认自动创建）
- 抽象与类比

| 关系型数据库 | Elasticsearch |
| --- | --- |
| Table | Index(Type) |
| Row | Document |
| Column | Field |
| Schema | Mapping |

### Mapping

- Mapping 类似数据库中的 Schema 定义，作用如下：
    - 定义索引中字段的名称
    - 定义字段的数据类型，例如：字符串、数字、布尔...
    - 定义倒排索引的相关配置
- Mapping 会把 json 文档映射成 Lucene 所需要的扁平格式

#### 动态 Mapping

- 在写入文档的时候，如果索引不存在，会自动创建索引
- ES 会自动根据文档信息，推算出字段的类型
- 默认允许自动新增字段
- 支持的类型，如下：

| json类型 | es类型 |
| --- | --- |
| 布尔值 | boolean |
| 浮点数 | float |
| 整数 | long |
| 对象 | object |
| 数组 | 由第一个非空数值的类型所决定 |
| 空值 | 忽略 |

- 需要注意的是
    - mapping中的字段类型一旦设定后，禁止修改
    - 原因：Lucene实现的倒排索引生成后不允许修改(提高效率)，如果要修改字段的类型，需要重新建立索引

#### 自定义（显示） Mapping

有时动态 Mapping 自动推算不是那么准确，因此可以通过自定义 Mapping 来创建。

- 预先定义好字段名称和类型，然后再调 API 创建

自定一创建 Mapping 的建议：

- 参考 API 手册，纯手写
- 但是为了减少输入的工作量，减少出错概率，可以依照以下步骤
    - 创建临时的index， 写入一条样本数据
    - 通过访问 Mapping API 获得该临时的动态 Mapping 定义
    - 然后进行相应的修改，再重新创建修改后的 Mapping
    - 删除临时索引



### 文档 - doc

- ES 是面向文档的，是所有可搜索数据的最小单位
- 文档会被序列化成 json 格式，保存在 ES 中
    - json 对象有相关字段组成
    - 每个字段都有对应的字段类型（字符串/数值/布尔/日期/二进制...）
- 每个文档都有一个 Unique ID
    - 可以自己定义（指定 ID）
    - 或者通过 ES 主动生成
    
- 文档元数据

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/e14aaf7937a44b55b7f58d50b35484e6~tplv-k3u1fbpfcp-watermark.image)

### Restful Api - 各种编程语言调用

![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/bab87c3f0d8c4d78bdbf927af37672b3~tplv-k3u1fbpfcp-watermark.image)

### 在 kibana dev tools CRUD 数据（Demo）

```bash
// 查看mapping
GET deal/_mapping

// 查看索引文档的总数
GET deal/_count

// 查看索引前10条文档（默认只返回10条）
GET deal/_search

// 查看索引信息
GET /_cat/indices?v

// 查看节点情况
GET /_cat/nodes?v

// 查看分片情况
GET /_cat/shardes?v

// 查看索引信息，安装文档个数排序
GET /_cat/indices?v&s=docs.count:desc

// 查看每个索引占用的内存情况
GET /_cat/indices?v&h=i,tm&s=tm:desc

// 搜索，返回所有字段
GET deal/_search
{
  "query": {
    "match": {
      "Fdeal_id": 11104514117
    }
  }
}

// 搜索，返回指定的字段
GET deal/_search?_source=Fdeal_id,Fdeal_gen_time,Frecv_name,Frecv_phone
{
  "query": {
    "match": {
      "Fdeal_id": 11104514117
    }
  }
}

************* story 测试 *************
// 插入一条数据，自动创建索引 story，指定id
PUT story/_doc/1
{
  "title": "守株待兔",
  "author": "Jack",
  "create_time": "2000-08-11"
}

// 插入一条数据，自动 id
POST story/_doc
{
  "title": "掩耳盗铃",
  "author": "Mike",
  "create_time": "2000-06-16"
}

// 查看所有文档（返回10条）
GET story/_search

// 查看指定文档
GET story/_doc/1

// 更新文档，文档必须已经存在，更新只会对相应字段做增量修改
POST story/_update/1
{
  "doc": {
    "create_time": "2000-08-12"
  }
}

// 删除文档
DELETE story/_doc/1
```

### 复合查询（bool query）

ES 经常用来做聚合、统计、分类等各种复杂的多条件查询，实际上，bool query用得非常多，因为查询条件个数不定，所以处理的逻辑思路时，外层用一个大的bool query来进行承载。

bool query 可以组合任意多个简单查询，各个简单查询之间的逻辑表示如下：

| 属性 | 说明 |
| --- | --- |
| must | 文档必须匹配must选项下的查询条件，相当于逻辑运算的 AND |
| should | 文档可以匹配 should 选项下的查询条件，也可以不匹配，相当于逻辑运算的 OR |
| must_not | 与 must 相反，匹配该选项下的查询条件的文档不会被返回 |
| filter | 和 mus t一样，匹配 filter 选项下的查询条件的文档才会被返回，但是 filter 不评分，只起到过滤功能 |

例子如下：

```bash
{
  "query": {
    "bool": {
      "must": {
        "match": {
          "content": "路由器"
        }
      },
      "must_not": {
        "match": {
          "content": "小米"
        }
      }
    }
  }
}
```

    需要注意的是，同一个bool下，只能有一个must、must_not、should 或 filter。

如果希望有多个 must 时，比如希望同时匹配"路由器"和"小米"，但是又故意分开这两个关键词（事实上，一个 must，然后使用 match，并且 operator 为 and 就可以达到目的），注意must下使用数组，然后里面多个 match 对象就可以了：

```bash
{
  "query": {
    "bool": {
      "must": [
        {
          "match": {
            "content": "路由器"
          }
        },
        {
          "match": {
            "content": "小米"
          }
        }
      ]
    }
  }
}
```

### 在命令行，curl 方式请求

```bash
  curl -XGET http://ip:9200/_cat/indices?pretty
  curl -XGET http://ip:9200/_cat/nodes?pretty
  curl -XGET http://ip:9200/_cluster/health?pretty
  curl -XGET http://ip:9200/_cat/shards?v
  curl -XGET http://ip:9200/mytest/_search?pretty
  curl -H"Content-Type:application/json" -XPOST http://ip:9200/deal/_search?pretty -d '{"query":{"match":{"name":"Jack"}}}'
  ......
```

### 聚合操作

相当于MySQL的聚合函数。

- max
```bash
{
  "size": 0,
  "aggs": {
    "max_id": {
      "max": {
        "field": "id"
      }
    }
  }
}
```
    size不设置为0，除了返回聚合结果外，还会返回其它所有的数据。
    
- min
```bash
{
  "size": 0,
  "aggs": {
    "min_id": {
      "min": {
        "field": "id"
      }
    }
  }
}
```

桶聚合，相当于 Mysql 的 group by.

......


### 批量操作 - Bulk API

- 批量操作可以减少网络连接产生的开销，提高性能
- 支持在一次 API 调用中，对不同的索引进行操作
- 支持四种类型操作
    - Index
    - Create
    - Update
    - Delete
- 操作中单条操作失败，不会影响其他操作
- 返回结果包括了每一条操作执行的结果

#### Bulk 批量操作
```bash
************* Bulk 批量操作 *************
GET test/_search
GET test2/_search

POST _bulk
{"index": {"_index": "test", "_id": "1"}}
{"field1": "value1"}
{"delete": {"_index": "test", "_id": "2"}}
{"create": {"_index": "test2", "_id": "3"}}
{"field1": "value1"}
{"update": {"_index": "test", "_id": "1"}}
{"doc": {"field2": "value2"}}
```

#### 批量读取 - mget

```bash
************* mget 批量操作 *************
GET _mget
{
  "docs": [
      {
        "_index": "test",
        "_id": "1"
      },
      {
        "_index": "test2",
        "_id": "3"
      },
      {
        "_index": "story",
        "_id": "1"
      }
    ]
}
```

#### msearch 批量操作

```bash
************* msearch 批量操作 *************
POST story/_msearch
{}
{"query": {"match_all": {}}, "size": 3}
{"index": "test"}
{"query": {"match_all": {}}, "size": 3}
{"index": "deal"}
{"query": {"match_all": {}}, "size": 3}
```

## 正排和倒排索引

### 图书和搜索引擎的类比

- 图书
    - 正排索引 - 目录页
    - 倒排索引 - 索引页
- 搜索引擎
    - 正排索引 - 文档 id 到文档内容和单词的关联
    - 倒排索引 - 单词到文档 id 的关系

![image.png](https://p6-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/6949c0f64b234e4f980e1eda7956cdcd~tplv-k3u1fbpfcp-watermark.image)

### 倒排索引的核心组成

- 单词词典，记录所有文档的单词，记录单词到倒排索引的关联关系
- 倒排列表 - 记录单词对应的文档相结合，由倒排索引组成
    - 倒排索引项
        - 文档 ID
        - 词频 - 该单词在文档中出现的次数，用于相关性评分
        - 位置 - 单词在文档中分词的位置，一般用于语句搜索
        - 偏移 - 记录单词的开始和结束位置，实现高亮显示

### ES 的倒排索引

- ES 的 json 文档中的每个字段，都有自己的倒排索引
- 可以指定对某些字段不做索引
    - 优点：节省存储空间
    - 缺点：字段无法被索引

## 分词器

分词：把全文本转换一系列的单词的过程。

分词是通过分词器（Analyzer）来实现的。

![image.png](https://p9-juejin.byteimg.com/tos-cn-i-k3u1fbpfcp/c6ce93a4b7fc40738941f87d46dbcab0~tplv-k3u1fbpfcp-watermark.image)

### ES 的内置分词器

- Standard Analyzer - 默认分词器，按词切分，小写处理
- Simple Analyzer - 按照非字母切分，非字母的都被去除，小写处理
- Stop Analyzer - 停用词过滤（the、a、is），小写处理
- Keyword Analyzer - 不分词，直接将输入当作输出
- Customer Analyzer - 自定义分词器

### 使用 _analyzer API

```bash
// Standard Analyzer - 默认分词器，按词切分，小写处理
GET _analyze
{
  "analyzer": "standard",
  "text": ["hello 2 world", "How are you."]
}

// Simple Analyzer 按照非字母切分，非字母的都被去除，小写处理
GET _analyze
{
  "analyzer": "simple",
  "text": ["hello 2 world", "Hi-you."]
}

// Stop Analyzer 按照非字母切分，非字母的都被去除，小写处理
GET _analyze
{
  "analyzer": "stop",
  "text": ["hello 2 world", "Hi-you."]
}
```

### 中文分词器 IK

- 支持自定义词库，支持热更新分词字典
- [https://github.com/medcl/elasticsearch-analysis-ik](https://links.jianshu.com/go?to=https%3A%2F%2Fgithub.com%2Fmedcl%2Felasticsearch-analysis-ik%2Freleases)

```bash
GET _analyze
{
  "analyzer": "ik_smart",
  "text": ["商品中心"]
}

curl -H"Content-Type:application/json" -XGET http://ip:9200/_analyze?pretty -d '{"analyzer":"ik_smart","text": ["Midea the shop.商品中心"]}'
```

### 搜索的相关性（打分score）

相关性打分是指文档与查询语句间的相关度，其本质就是搜索结果的文档返回排序的问题。

- 是否可以找到所有相关的内容
- 有多少不相关的内容被返回
- 文档打分是否合理
- 结合业务需求，平衡结果排名

比如：Google 的 Web 搜索，用的是 Page Rank 算法进行相关性计算和打分排名。不仅仅是内容，更重要的是内容的可信度。


| 相关性指标 | 解释 |
| --- | --- |
| 词频（TF） | 单词在该文档中出现的次数。词频越高，相关性越强 |
| 文档频率（DF） | 出现单词的文档数|




