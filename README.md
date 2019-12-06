# UniQuery

Unified query language for various commonly used data representation formats (JSON, XML, YAML, etc.)

## Usage

Install [Go](https://golang.org/) and run `go run cmd/uniquery/main.go -h` to get information about available flags and their meaning.

## Query Syntax

Please see [query examples](examples.md) for rough query syntax explanation.

## Data Format Support

| Format | Support                | Notes                                |
| :----: | :--------------------- | :----------------------------------- |
|  JSON  | :heavy_check_mark: Yes | Works according to tests.            |
|  YAML  | :question: Partial     | Not very well tested yet.            |
|  XML   | :x: No                 | More complicated than JSON and YAML. |
|  CSV   | :x: No                 | Support is not currently planned.    |


## Example (JSON)

Consider the following JSON file `users.json`, which maps real names to nicknames.

```json
{
    "alice80": "Alice Yang",
    "bob12": "Bob Jacobs",
    "tank": "Charlie Peterson"
}
```

|       Query | Result                                             |
| ----------: | :------------------------------------------------- |
| empty query | root element / whole data structure                |
|   `alice80` | `"Alice Yang"`                                     |
|      `tank` | `"Charlie Peterson"`                               |
|         `*` | `["Alice Yang", "Bob Jacobs", "Charlie Peterson"]` |
|  `alice80.` | root element (parent of `alice80`)                 |
