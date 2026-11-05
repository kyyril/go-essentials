
| Verb | Function / Output Format | Example Input | Example Output |
| --- | --- | --- | --- |
| **`%q`** | Single/double-quoted string or rune with safe escaping | `"Hello\nWorld"`, `'A'` | `"Hello\nWorld"`, `'A'` |
| **`%s`** | Raw string (as-is) | `Hello\nWorld`, `\x41` | `HelloWorld` (newline trimmed), error parsing |
| **`%v`** | Default value format | `Hello\nWorld`, `65` (rune integer value) | `Hello\nWorld`, `65` |
| **`%d`** | Decimal integer | *(invalid)*, `65` | `65` |
| **`%U`** | Unicode format (`U+XXXX`) | *(invalid)*, `U+0041` | `U+0041` |
| **`%+q`** | Quoted string with ASCII-only escaping | `"Hello\nWorld"` (non-ASCII escaped to `\uXXXX`), `'A'` | `"Hello\nWorld"`, `'A'` |
