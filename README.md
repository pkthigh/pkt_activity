# pkt_activity

|版本|说明|
|:--|:--:|
|v1.0|初版|

---

## HTTP接口

**WEB端需要从客户端获取玩家Token, 用于在WEB端对服务器HTTP接口的请求, 请求时键为token, 值为token.**

---

### 当前活动列表

| **Request**
|字段|类型|说明|必须|
|:---|:--:|:--:|:--:|

- Method: **GET**
- Path: ```/v1/activity/list```
- Header: {"tk": "xxxx"}
- Body:

```json
```

| **Response**
|字段|类型|说明|
|:---|:--:|:---:|
|online_time|int|上线时间|
|offline_time|int|下线时间|
|ac_type|int|活动类型|1:手数活动|
|status|int|活动状态|0:未开启 / 1:排期中 / 2:进行中 / 3:提前下线 4: 已结束|
|hand_num|int|活动门槛手数|
|bonus|int|奖金|

- Body:

```json
{
    "code": 0,
    "data":[
        {
            "info":{
                "id": 1,
                "online_time": 1234567890,
                "offline_time": 1234567890,
                "name_zh": "xxx",
                "content_zh": "xxx",
                "pic_zh_url": "xxx",
                "name_en": "xxx",
                "content_en": "xxx",
                "pic_en_url": "xxx",
                "ac_type": 1,
                "status" : 2,
                "hand_num": 100
            },
            "details":[
                {
                    "id": 10,
                    "activity_id" : 1,
                    "index": 1,
                    "bonus": 10
                }
            ]
        }
    ],
    "msg": "successful"
}
```

---

### 用户活动信息

| **Request**
|字段|类型|说明|必须|
|:---|:--:|:--:|:--:|
|aid|int|活动id|Y|

- Method: **GET**
- Path: ```/v1/activity/info?aid=1```
- Header: {"tk": "xxxx"}
- Body:

```json
```

| **Response**
|字段|类型|说明|
|:---|:--:|:---:|
|player_hands|int|用户当前手数|
|palyer_status|int|0: 不满足抽奖条件 1: 满足抽奖条件未抽奖 2: 已抽奖|

- Body:

```json
{
    "code": 0,
    "data": {
        "id": 1,
        "online_time": 1234567890,
        "offline_time": 1234567890,
        "name_zh": "xxx",
        "content_zh": "xxx",
        "pic_zh_url": "xxx",
        "name_en": "xxx",
        "content_en": "xxx",
        "pic_en_url": "xxx",
        "ac_type": 1,
        "status" : 2,
        "hand_num": 100,
        "player_hands": 120,
        "palyer_status": 1
    },
    "msg": "successful"
}
```

---

### 用户中奖记录

| **Request**
|字段|类型|说明|必须|
|:---|:--:|:--:|:--:|

- Method: **GET**
- Path: ```/v1/activity/records```
- Header: {"tk": "xxxx"}
- Body:

```json
```

| **Response**
|字段|类型|说明|
|:---|:--:|:---:|
|raffle_time|int|抽奖时间|
|bonus|int|奖金|
|status|int|0:未派发 1:已派发|

- Body:

```json
{
    "code": 0,
    "data": [
        {
            "id": 1,
            "create_at": 1234567890,
            "update_at": 1234567890,
            "player_id": 10000,
            "raffle_time": 1234567890,
            "activity_id": 1,
            "activity_detail_id": 10,
            "bonus": 10,
            "status": 1,

        },
    ],
    "msg": "successful"
}
```

---

### 用户参与活动

| **Request**
|字段|类型|说明|必须|
|:---|:--:|:--:|:--:|
|did|string|设备ID|Y|
|lang|bool|语言 f:zh/ t:en|Y|

- Method: **POST**
- Path: ```/v1/activity/participate```
- Header: {"tk": "xxxx"}
- Body:

```json
{
    "aid": 1,
    "atype": 1,
    "did": "xxxx",
    "lang": 0,
}
```

| **Response**
|字段|类型|说明|
|:---|:--:|:---:|
|detail_id|int|活动条件ID|
|bonus|int|中奖金额|

- Body:

```json
{
    "code": 0,
    "data": {
        "detail_id": 10,
        "bonus": 1000
    },
    "msg": "successful"
}
```

---

### 用户中奖跑马通知

| **Request**
|字段|类型|说明|必须|
|:---|:--:|:--:|:--:|
|aid|int|活动id|Y|
|count|int|数量, 选填, 默认50|N|

- Method: **GET**
- Path: ```/v1/activity/results```
  - ```/v1/activity/results?aid=18&count=2```

- Header:
- Body:

```json
```

| **Response**
|字段|类型|说明|
|:---|:--:|:---:|

- Body:

```json
{
    "code": 0,
    "data": [
        {"player":"xxx", "bonus": 1000},
    ],
    "msg": "Successful"
}
```

---
