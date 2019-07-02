- A strict superset of JSON with the addition of syntactically significant newlines and indentation
- Comment: `#`
- Strings doesn't have to be quoted
- Multiple-line strings can be written either as a 'literal block' (using `|`), or a 'folded block' (using `>`).
- Nesting uses indentation
    - 2 space indent is preferred
- Maps doesn't have to have string keys
- The `<<` merge key is used to indicate that all the keys of one or more specified maps should be inserted into the current map
    - If the value associated with the key is a single mapping node, each of its key/value pairs is inserted into the current mapping, unless the key already exists in it
    - If the value associated with the merge key is a sequence, then this sequence is expected to contain mapping nodes and each of these nodes is merged in turn according to its order in the sequence. Keys in mapping nodes earlier in the sequence override keys specified in later mapping nodes
- > http://yaml.org/type/merge.html
- `*`, `&`
    - Repeated nodes (objects) are first identified by an anchor (marked with the ampersand - `&`), and are then aliased (referenced with an asterisk - `*`) thereafter
        - http://yaml.org/spec/1.2/spec.html
# [YAML](https://yaml.org/spec/1.2/spec.html)
- The primary objective of this revision (YAML 1.2) is to bring YAML into compliance with JSON as an official subset
- > Chapter 2 provides a short preview of the main YAML features. Chapter 3 describes the YAML information model, and the processes for converting from and to this model and the YAML text format. The bulk of the document, chapters 4 through 9, formally define this text format. Finally, chapter 10 recommends basic YAML schemas
- Repeated nodes (objects) are first identified by an anchor (marked with the ampersand - `“&”`), and are then aliased (referenced with an asterisk - `“*”`) thereafter

    ```yaml
    ---
    hr:
      - Mark McGwire
      # Following node labeled SS
      - &SS Sammy Sosa
    rbi:
      - *SS # Subsequent occurrence
      - Ken Griffey
    ```

- [Merge key](https://yaml.org/type/merge.html)

    ```yaml
    ---
    - &CENTER { x: 1, y: 2 }
    - &LEFT { x: 0, y: 2 }
    - &BIG { r: 10 }
    - &SMALL { r: 1 }

    # All the following maps are equal:

    - # Explicit keys
    x: 1
    y: 2
    r: 10
    label: center/big

    - # Merge one map
    << : *CENTER
    r: 10
    label: center/big

    - # Merge multiple maps
    << : [ *CENTER, *BIG ]
    label: center/big

    - # Override
    << : [ *BIG, *LEFT, *SMALL ]
    x: 1
    label: center/big
    ```

