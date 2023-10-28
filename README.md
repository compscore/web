# Web

## Check Format

```yaml
- name:
  release:
    org: compscore
    repo: smb
    tag: latest
  credentials:
    username:
    password:
  target:
  command:
  expectedOutput:
  weight:
  options:
    status_code:
    match:
    substring_match:
    regex_match:
```

## Parameters

|     parameter     |            path            |   type   | default  | required | description                                               |
| :---------------: | :------------------------: | :------: | :------: | :------: | :-------------------------------------------------------- |
|      `name`       |          `.name`           | `string` |   `""`   |  `true`  | `name of check (must be unique)`                          |
|       `org`       |       `.release.org`       | `string` |   `""`   |  `true`  | `organization that check repository belongs to`           |
|      `repo`       |      `.release.repo`       | `string` |   `""`   |  `true`  | `repository of the check`                                 |
|       `tag`       |       `.release.tag`       | `string` | `latest` | `false`  | `tagged version of check`                                 |
|    `username`     |  `.credentials.username`   | `string` |   `""`   | `false`  | `username for basic auth`                                 |
|    `password`     |  `.credentials.password`   | `string` |   `""`   | `false`  | `default password for basic auth or Authorization header` |
|     `target`      |         `.target`          | `string` |   `""`   |  `true`  | `network target for smb server`                           |
|     `command`     |         `.command`         | `string` |   `""`   | `false`  | `HTTP verb to create request with`                        |
| `expectedOutput`  |     `.expectedOutput`      | `string` |   `""`   | `false`  | `expected output for check to measured against`           |
|     `weight`      |         `.weight`          |  `int`   |   `0`    |  `true`  | `amount of points a successful check is worth`            |
|   `status_code`   |   `.options.status_code`   |  `int`   |   `0`    | `false`  | `check status code of response matches provided code`     |
|      `match`      |      `.options.match`      |  `bool`  | `false`  | `false`  | `check contents of targeted file are exact match`         |
| `substring_match` | `.options.substring_match` |  `bool`  | `false`  | `false`  | `check contents of targeted file are substring match`     |
|   `regex_match`   |   `.options.regex_match`   |  `bool`  | `false`  | `false`  | `check contents of targeted file are regex match`         |

## Examples

```yaml
- name: google.com-web
  release:
    org: compscore
    repo: web
    tag: latest
  target: https://google.com
  command: GET
  weight: 2
  options:
    status_code: 200
```

```yaml
- name: host_a-web
  release:
    org: compscore
    repo: web
    tag: latest
  target: http://10.{ .Team }.1.1
  command: GET
  expectedOutput: According to all known laws of aviation
  weight: 2
  options:
    status_code: 200
    substring_match:
```
