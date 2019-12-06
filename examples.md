# UniQuery Examples

|     Query      | Description                                              |
| :------------: | :------------------------------------------------------- |
| [empty query]  | Root element.                                            |
|      `**`      | All elements (scanned recursively).                      |
|     `**.*`     | All children of all elements (all elements but root).    |
|      `a`       | Child `a` of root element.                               |
|     `a.b`      | Child `b` of root element's `a` child.                   |
|     `a.*`      | All children of root element's `a` child.                |
|     `**.a`     | Children `a` of every element with such child.           |
|      `a.`      | Parent of root's `a` child (root if root has `a` child). |
|    `**.a.`     | All elements with `a` child (parents of `a` elements).   |
|      `0`       | First item in root (root element must be array).         |
|     `**.0`     | First item of every array.                               |
|   `**.0.a.`    | First item with `a` child of every array.                |
|      `*=`      | All empty strings.                                       |
|    `*=abcd`    | All `abcd` strings.                                      |
| `**="*.a\"b="` | All `*.a"b=` strings.                                    |

## Special Characters

| Character | Description                                             |
| :-------: | :------------------------------------------------------ |
|    `.`    | Child accessor (parent if specifier is empty).          |
|    `*`    | Child wildcard (`**` for recursion).                    |
|    `=`    | Value equality filter.                                  |
|    `~`    | Regular expression match filter.                        |
|    `\`    | Escape character for special characters.                |
|    `"`    | Quoted values and names may contain special characters. |
