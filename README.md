# SHUVolunteer

上海大学志愿者网站手机版后端

## API文档

### 登录

#### Request
+ Method: ```POST```
+ Url: ```/login```
+ Headers:
    ```Content-Type: application/json```
+ Body:
    ```json
    {
        "username": "string",
        "password": "string"
    }
    ```

#### Response
+ Code: ```200```
+ Body:
    ```json
    {
        "token": "string"
    }
    ```

### 报名活动
#### Request
+ Method: ```POST```
+ Url: ```/apply```
+ Authorization: ``` Bearer Token ```
+ Headers:
    ```
    Content-Type: application/json
    ```
+ Body:
    ```json
    {
        "activity_name": "string"
    }
    ```
#### Response
+ Code: ```200```

    表示操作成功执行

- Code: ```403```

    错误：没有登录

### 取消活动
#### Request
+ Method: ```POST```
+ Url: ```/resign```
+ Authorization: ``` Bearer Token ```
+ Headers:
    ```
    Content-Type: application/json
    ```
+ Body:
    ```json
    {
        "activity_name": "string"
    }
    ```
#### Response
+ Code: ```200```

    表示操作成功执行

- Code: ```403```

    错误：没有登录


### 获取活动列表
#### Request
+ Method: ```GET```
+ Url: ```/activities```
+ Authorization: ``` Bearer Token ```
#### Response
+ Code: ```200```
+ Body:

    返回所有的活动
    ```json
        [
            {
                "date": "date",
                "name": "string",
                "team": "string"
            },
            ...
        ]
    ```

- Code: ```403```

    错误：没有登录
