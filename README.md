# 此项目停止维护已久，请看此

1. 请访问重置版本：[Hentioe/zhenxiang](https://github.com/Hentioe/zhenxiang)
1. 重置版本将继续提供免费的公共服务：）

# sorry-generator

![Travis](https://travis-ci.org/Hentioe/sorry-generator.svg?branch=master)
![GitHub release](https://img.shields.io/github/release/Hentioe/sorry-generator.svg)
![Docker Automated build](https://img.shields.io/docker/build/bluerain/sorry-generator.svg)


### 说明

本项目为`Sorry-为所欲为`系列视频/GIF 生成器，配套前端：https://sorry.bluerain.io

PS：灵感和部分资源模板来自 [xtyxtyx/sorry](https://github.com/xtyxtyx/sorry) 感谢：）

#### 使用

在有 Docker 的系统上直接执行下列命令即可（注意端口映射和挂载目录）：

```` bash
# 创建 sorry-tenerator 容器的 VOLUME
docker volume create tmp-sorry-gen
# 启动容器
docker run -ti --name sorry-gen \
-d -p 8080:8080 --restart=always \
-v tmp-sorry-gen:/data/tmp
-v /data/apps/sorry-generator/resources:/data/resources \
-v /data/apps/sorry-generator/dist:/data/dist \
bluerain/sorry-generator
````

附加解释：容器在启动时会持久化 `/data/tmp` 中的文件到 VOLUME，当前这个目录会存放通过上传接口上传的资源包。


程序默认绑定到 `:8080`，以 `test` 模式启动，若要更改需要手动添加 CLI 参数：

````
./sorry-gen -bind :80 -mode release
````

容器启动同样的直接将参数加在镜像后面。


注意：从 0.3 版本开始模板资源不会集成在项目或者 Docker 镜像中，需要自行安装【看[这里](#安装资源)】。

POST 以下数据到 `http://localhost:8080/generate/sorry/mp4`:

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
那么就可以直接提供生成文件的直链了：http://your.domain/dist/{hash}.[mp4|gif]


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
    ├── dagong              # 模板 KEY（API 中 tlp_key 的参数即是目录的名称）
        ├── name            # 模板显示名称（文本文件），自动生成
        ├── sentences       # 预设字幕（文本文件，每一行表示一句字幕），自动生成
        ├── template.ass    # 字幕模板，由原始字幕文件自动转换而成
        └── template.mp4    # 视频素材模板（实际上就是无字幕的原视频）
````

上传资源包 API（将 res.zip 放置在 ./assets 目录中）:

````
curl -X POST http://localhost:8080/upload/res \
  -F "file=@./assets/res.zip" \
  -H "Content-Type: multipart/form-data"
````

上传完成后会自动进行资源包的安装。安装成功会返回安装生成的文件列表，例如：

````
{
    "make_files": [
        "resources/template",
        "resources/template/sorry",
        "resources/template/sorry/template.ass",
        "resources/template/sorry/template.mp4",
        "resources/template/lese",
        "resources/template/lese/template.ass",
        "resources/template/lese/template.mp4",
        "resources/template/wangjingze",
        "resources/template/wangjingze/template.ass",
        "resources/template/wangjingze/template.mp4",
        "resources/template/dagong",
        "resources/template/dagong/template.ass",
        "resources/template/dagong/template.mp4"
    ]
}
````

如果安装的资源包中的资源已经存在，则不会生成任何文件（在安装资源包章节有详细描述）。假设上传的资源包中仅仅只有一个 sorry 资源，
在已经存在 sorry 的情况下，API 会返回一个空的 make_files。

异步任务和并发限制：

````
./sorry-gen -cl <number>
````
cl 参数即 Concurrency limits（并发限制），默认限制为 CPU 数量。需要注意的是，此限制并不对 `/generate/{tpl_key}/{res_type}` API 生效。这个参数影响的是生成异步任务的 API: `/task/generate/{tpl_key}`。

使用方式：

POST `http://localhost:8080/task/generate/sorry`

````
{"sentences":["第一句","第二句","第三句","第四句","第五句","第六句","第七句","第八句","第九句"]}
````

会立即返回：

````
{
    "hash": "776168419d55d4fe68792a73f6450791",
    "state": "waiting"
}
````

hash 表示产生的任务 ID（同时也表示产生的资源名称），state 表示当前任务状态。一般来讲调用此 API 的 state 状态都是 waiting，因为此 API 只会创建生成任务并立即返回资源 ID 并不会等到任务执行完毕才返回。

根据上面产生的 ID 来获取最新的任务状态：

GET `http://localhost:8080/task/generate/776168419d55d4fe68792a73f6450791`

````
{
    "hash": "776168419d55d4fe68792a73f6450791",
    "state": "completed"
}
````

当 state 为 completed 时，任务已经执行结束了，而不再是创建时的等待状态，这时候相应资源已经生成（包含 .gif 和 .mp4）完成，可以直接根据 ID 下载。

所有状态常量：

```` golang
const (
	// StateWaiting 等待状态（添加后默认）
	StateWaiting = "waiting"
	// StateCompleted 完成状态
	StateCompleted = "completed"
	// StateError 失败状态
	StateError = "failed"
	// StateNone 空状态（没有构建任务）
	StateNone = "none"
)
````

添加对异步任务和对其的并发限制支持的目的是，当遇到大量用户并发使用的场景的时候可以通过异步任务 API 解决响应延迟以及服务器资源紧张问题。


附加说明：

* 为什么不加入前端？
  
  因为这种东西本来就没必要限制为 Web 前端啊…… 需要前端自己写个静态页面即可。实际上应该将它视作任何 Programmably 项目的后端，例如各种平台的 Bot

### 安装资源

除了通过 `/upload/res` API 上传以外，还可以在本地执行命令：`./sorry-gen -i res.zip`完成对资源包的安装，资源包的结构见上述说明。资源包中的任何文件都不会对已存在的资源文件进行替换，如果要更新指定资源请先删除相关目录再执行安装。


在手动编译运行的情况下，默认是没有资源包的，你可以拉取并安装我的资源包：

````
wget https://dl.bluerain.io/res.zip
./sorry-gen -i res.zip
````

同样的，使用 sorry-generator 的 Docker 容器也可以这样安装资源包：

````
docker run -ti --rm -v $PWD/res.zip:/data/res.zip \
-v $PWD/resources:/data/resources bluerain/sorry-generator \
-i res.zip
````

如果你要创建可安装的资源包，需要遵循与以下标准：

1. 以 template 目录为根
2. 必须存在 template.mp4 和 template.ass 文件

假设你创建的安装包目录结构是这样的：

````
.
└── template
    └── sorry
        ├── template.ass
        └── template.mp4
````
template.ass 内容为：

````
[Script Info]
; Script generated by Aegisub 3.2.2
; http://www.aegisub.org/
Title: 为所欲为
ScriptType: v4.00+
WrapStyle: 0
ScaledBorderAndShadow: yes
YCbCr Matrix: TV.601
PlayResX: 300
PlayResY: 168

[Aegisub Project Garbage]
Audio File: template.mp4
Video File: template.mp4
Video AR Mode: 4
Video AR Value: 1.781250
Video Zoom Percent: 2.000000
Active Line: 8
Video Position: 25

[V4+ Styles]
Format: Name, Fontname, Fontsize, PrimaryColour, SecondaryColour, OutlineColour, BackColour, Bold, Italic, Underline, StrikeOut, ScaleX, ScaleY, Spacing, Angle, BorderStyle, Outline, Shadow, Alignment, MarginL, MarginR, MarginV, Encoding
Style: sorry,WenQuanYi Micro Hei,23,&H00FFFFFF,&H000000FF,&H00000000,&H00000000,0,0,0,0,100,100,0,0,1,1.1,0.5,2,5,5,5,1

[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:01.18,0:00:01.56,sorry,,0,0,0,,好啊
Dialogue: 0,0:00:03.18,0:00:04.43,sorry,,0,0,0,,就算你是一流程序员
Dialogue: 0,0:00:05.31,0:00:07.43,sorry,,0,0,0,,写出来的代码再完美
Dialogue: 0,0:00:07.56,0:00:09.93,sorry,,0,0,0,,我说这是 BUG 它就是 BUG
Dialogue: 0,0:00:10.06,0:00:11.56,sorry,,0,0,0,,毕竟我是用户
Dialogue: 0,0:00:11.93,0:00:13.06,sorry,,0,0,0,,你害我加班啊
Dialogue: 0,0:00:13.81,0:00:16.31,sorry,,0,0,0,,sorry 我就喜欢看程序猿加班
Dialogue: 0,0:00:18.06,0:00:19.56,sorry,,0,0,0,,以后天天找他 BUG
Dialogue: 0,0:00:19.60,0:00:21.60,sorry,,0,0,0,,天天找 天天找
````

将上述目录打包以后进行安装，会在 resources/template 中产生这样的文件结构（以 sorry 为根的视角）：

````
.
└── sorry
    ├── name
    ├── sentences
    ├── template.ass
    └── template.mp4
````

template.ass 的内容为：

````
# 上面的内容省略……
[Events]
Format: Layer, Start, End, Style, Name, MarginL, MarginR, MarginV, Effect, Text
Dialogue: 0,0:00:01.18,0:00:01.56,sorry,,0,0,0,,{{ index .sentences 0 }}
Dialogue: 0,0:00:03.18,0:00:04.43,sorry,,0,0,0,,{{ index .sentences 1 }}
Dialogue: 0,0:00:05.31,0:00:07.43,sorry,,0,0,0,,{{ index .sentences 2 }}
Dialogue: 0,0:00:07.56,0:00:09.93,sorry,,0,0,0,,{{ index .sentences 3 }}
Dialogue: 0,0:00:10.06,0:00:11.56,sorry,,0,0,0,,{{ index .sentences 4 }}
Dialogue: 0,0:00:11.93,0:00:13.06,sorry,,0,0,0,,{{ index .sentences 5 }}
Dialogue: 0,0:00:13.81,0:00:16.31,sorry,,0,0,0,,{{ index .sentences 6 }}
Dialogue: 0,0:00:18.06,0:00:19.56,sorry,,0,0,0,,{{ index .sentences 7 }}
Dialogue: 0,0:00:19.60,0:00:21.60,sorry,,0,0,0,,{{ index .sentences 8 }}
````
sentences 的内容为：

````
好啊
就算你是一流程序员
写出来的代码再完美
我说这是 BUG 它就是 BUG
毕竟我是用户
你害我加班啊
sorry 我就喜欢看程序猿加班
以后天天找他 BUG
天天找 天天找
````

name 的内容为：

````
为所欲为
````

可以发现，安装后的资源和原始资源包解压的区别在于：

1. template.ass 文件从原始字幕文件转换为模板字幕文件
2. 从原始字幕内容中提取的每一条字幕内容被持久化存储在了 sentences 文件中
3. 从原始字幕文件内容中提取的 Title 属性的值被持久化储存在了 name 文件中

只有经过安装的原始资源才能被程序正确的读取，原始资源是无法直接解压使用的。这样做的目的是方便对资源的创建，
因为在经过安装步骤之前需要手动创建字幕模板，是很别扭的。还要手动创建 name 和 sentences 文件这些跟资源无关的内容。
而安装功能可以直接使用最原始的资源（原始视频 + 原始字幕）。

PS: 有关视频字幕的制作建议了解一下 [Aegisub](http://www.aegisub.org/) 软件。


### 申请添加

理论上已经不需要申请加入新的模板资源，因为你可以自行上传。不过如果下载和剪切视频以及制作字幕对你而言仍然十分有难度的话

你可以用 Issue 投稿：

1. 标题为「建议添加 XX」。内容附上视频链接（如果是下载链接更好）、开始-结束时间段。
2. 标题为「希望添加 XX」。内容为视频片段的简短描述，上传视频附件（尺寸无所谓，我会自行会压缩）

第一种 Issue 会根据视频片段的热门程度、下载复杂度来决定是否添加，而第二种视频资源已经准备好的 Issue 有极大的可能会直接添加（精力有限）。

### 版本功能计划

- [x] v0.1: 实现基本功能
- [x] v0.2: 添加基于对模板资源扫描产生数据的查询相关的 API
- [x] v0.3: 程序本体和模板资源分离
- [x] v0.4: 提供上传接口并持久化储存新增的模板（固定结构的压缩包资源）
- [x] v1.0: 异步和并发限制支持，对资源的生成请求立即响应，提供查询接口返回任务实时状态
- [ ] v1.1: 回调支持，异步生成请求的任务完成主动触发 HookUrl

___

更多视频梗期待添加中……
