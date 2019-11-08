# UniQuery Examples

|     Query     | Description                                              |
| :-----------: | :------------------------------------------------------- |
|      `.`      | Root element.                                            |
|      `*`      | All elements.                                            |
|     `*.*`     | All children of all elements (all elements but root).    |
|      `a`      | Child `a` of root element.                               |
|     `a.b`     | Child `b` of root element's `a` child.                   |
|     `a.*`     | All children of root element's `a` child.                |
|     `*.a`     | Children `a` of every element with such child.           |
|     `a..`     | Parent of root's `a` child (root if root has `a` child). |
|    `*.a..`    | Every element with `a` child (parents of `a` elements).  |
|      `0`      | First item in root (root element must be array).         |
|     `*.0`     | First item of every array.                               |
|   `*.0.a..`   | First item with `a` child of every array.                |
|     `*=`      | All empty strings.                                       |
|   `*=abcd`    | All `abcd` strings.                                      |
| `*="*.a\"b="` | All `*.a"b=` strings.                                    |

## Special Characters

| Character | Description                              |
| :-------: | :--------------------------------------- |
|    `.`    | Child/parent accessor.                   |
|    `*`    | All elements/children.                   |
|    `=`    | Value equality filter.                   |
|    `~`    | Regular expression match filter.         |
|    `\`    | Escape character for special characters. |
