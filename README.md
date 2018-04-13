# sorry-generator

![Travis](https://travis-ci.org/Hentioe/sorry-generator.svg?branch=master)
![GitHub release](https://img.shields.io/github/release/Hentioe/sorry-generator.svg)
![Docker Automated build](https://img.shields.io/docker/build/bluerain/sorry-generator.svg)


### 说明

本项目为`Sorry-为所欲为`系列视频/GIF 生成器，配套前端：https://sorry.bluerain.io

PS：灵感和资源模板来自 [xtyxtyx/sorry](https://github.com/xtyxtyx/sorry) 感谢：）

#### 使用

在有 Docker 的系统上直接执行下列命令即可（注意端口映射和挂载目录）：

```` bash
docker run -ti --name sorry-gen \
-d -p 4008:8080 --restart=always \
-v /data/apps/sorry-generator/dist:/data/resources \
-v /data/apps/sorry-generator/dist:/data/dist \
bluerain/sorry-generator
````

POST 以下数据到 `http://localhost:4008/generate/sorry/mp4`:

````
{"sentences":["第一句","第二句","第三句","第四句","第五句","第六句","第七句","第八句","第九句"]}
````

成功会返回：
````
{
  "hash": "74c6157d5dec218191835252aabda749"
}
````


同时会在 /data/apps/sorry-generator/dist 目录下生成对应 hash 作为文件名的文件（没有后缀的为 ass 字幕文件）。

注：修改 generate API 的最后一个 path 参数 mp4 为 gif 即产生 gif 文件。修改 sorry 为其它资源（例如王境泽：wangjingze）则产生相对应的资源。

假设你这样配置 nginx:

````
server {
        listen       80;
        server_name  your.domain;

        location / {
                root /data/apps/sorry-generator;
                index index.html;
        }
}
````
那么就可以直接提供生成文件的直链了：http://your.domain/dist/$hash.[mp4|gif]


对模板资源的数据进行查询：

我的前端（或者其它程序）该怎么知道某个资源有多少条字幕句子？

GET 访问首页 `http://localhost:4008`:

````
{
  "res": [
   {
    "tpl_key": "dagong",
    "name": "窃格瓦拉-打工是不可能……",
    "sentences": [],
    "sentences_count": 6
   },
   {
    "tpl_key": "sorry",
    "name": "为所欲为",
    "sentences": [],
    "sentences_count": 9
   },
   {
    "tpl_key": "wangjingze",
    "name": "王境泽-真香",
    "sentences": [],
    "sentences_count": 4
   }
  ],
  "res_count": 3
 }
````
会得到一个 res 数组，其中 tpl_key 就是模板名称，也就是上面的 sorry。sentences_count 表示有多少条字幕（需要输入多少句子）。sentences 数组是预设在程序中的默认字幕（用处例如提供前端输入框默认的 plachholder 的值）。以上所有数据都是程序扫描资源目录产生的结果，没有任何数据库成分。所以只要添加新的资源模板，API 结果会自动变更。

也可以 GET 访问 `http://localhost:4008/info/{tpl_key}` 对单独的资源进行数据查询：

````
{
  "tpl_key": "sorry",
  "name": "为所欲为",
  "sentences": [],
  "sentences_count": 9
 }
````

资源目录结构说明（以 resources 为根）：

````
.
└── template
    ├── dagong              # 模板 key
        ├── name            # 模板显示名称（文本文件）
        ├── sentences       # 预设字幕（文本文件，每一行表示一句字幕）
        ├── template.ass    # 字幕模板
        └── template.mp4    # 视频素材模板（实际上就是无字幕的原视频）
````

附加说明：

* 为什么不加入前端？
  
  因为这种东西本来就没必要限制为 Web 前端啊…… 需要前端自己写个静态页面即可。实际上应该将它视作任何 Programmably 项目的后端，例如各种平台的 Bot

* 如何手动安装资源？
  
  执行命令：`./sorry-gen -i res.zip`即可完成对资源包的安装，资源包的结构见上述说明。资源包中的任何文件都不会对已存在的资源文件进行替换，如果要更新指定资源请先删除相关目录再执行安装。


  在手动编译运行的情况下，默认是没有资源包的，你可以拉取并安装我的资源包：

  ````
  wget https://dl.bluerain.io/res.zip
  ./sorry-gen -i res.zip
  ````

  Docker 中镜像默认拉取的即是我的资源包。同样的，使用 sorry-generator 也可以安装资源包：

  ````
  docker run -ti --rm -v $PWD:/data bluerain/sorry-generator -i res.zip
  ````

### 申请添加

自行上传的功能将在下一个版本进行实现，不过现在你仍然可以通过投稿的方式请求新增需要的模板资源。

欢迎来这里投稿，直接去开 Issue 即可：

1. 标题为「建议添加 XX」。内容附上视频链接（如果是下载链接更好）、开始-结束时间段。
2. 标题为「希望添加 XX」。内容为视频片段的简短描述，上传视频附件（尺寸无所谓，我会自行会压缩）

第一种 Issue 会根据视频片段的热门程度、下载复杂度来决定是否添加，而第二种视频资源已经准备好的 Issue 有极大的可能会直接添加（精力有限）。

### 版本功能计划

- [x] v0.1: 实现基本功能
- [x] v0.2: 添加基于对模板资源扫描产生数据的查询相关的 API
- [x] v0.3: 程序本体和模板资源分离
- [ ] v0.4: 提供上传接口并持久化储存新增的模板（固定结构的压缩包资源）
- [ ] v1.0: 异步支持，对资源的生成请求立即响应，并提供查询接口返回任务实时状态
- [ ] v1.1: 回调支持，异步生成请求的任务完成主动触发 HookUrl
- [ ] v1.2: 基于可控长度队列任务控制并发

___

更多视频梗期待添加中……
