# UniQuery Examples

|   Query   | Description                                                   |
| :-------: | :------------------------------------------------------------ |
|    `.`    | Root element.                                                 |
|    `*`    | All elements.                                                 |
|   `*.*`   | All children of all elements (all elements but root).         |
|    `a`    | Child `a` of root element.                                    |
|   `a.b`   | Child `b` of root element's `a` child.                        |
|   `a.*`   | All children of root element's `a` child.                     |
|   `*.a`   | Children `a` of every element with such child.                |
|   `a..`   | Parent of root's `a` child (root if root has `a` child).      |
|  `*.a..`  | Every element with `a` child (parents of `a` child elements). |
|    `0`    | First item in root (root element must be array).              |
|   `*.0`   | First item of every array.                                    |
| `*.0.a..` | First item of every array, where first item has `a` child.    |
