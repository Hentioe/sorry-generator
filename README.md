# sorry-generator

### 说明

本项目为`Sorry-为所欲为`系列视频/GIF 生成器，一时兴起写的代码，GIF 部分暂时没做。配套前端：https://sorry.bluerain.io

PS：灵感和资源模板来自 https://github.com/xtyxtyx/sorry 感谢：）

#### 使用

在有 Docker 的系统上直接执行下列命令即可（注意端口映射和挂载目录）：

```` bash
docker run -ti --name sorry-gen \
-d -p 4008:8080 --restart=always \
-v /data/apps/sorry-generator/dist:/data/dist \
bluerain/sorry-generator
````

POST 以下数据到 `http://localhost:4008/generate/sorry`:

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
那么就可以直接提供生成文件的直链了：http://your.domain/dist/${hash}.mp4

附加说明：

* 为什么不加入前端？
  
  因为这种东西本来就没必要限制为 Web 前端啊…… 需要前端自己写个静态页面即可。实际上应该将它视作任何 Programmably 项目的后端，例如各种平台的 Bot
___

更多视频梗，以及 GIF 支持中……
