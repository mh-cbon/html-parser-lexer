# Changelog - html-parser-lexer

### 1.0.0

__Changes__

- removed TagOpenStartToken and TagCloseStartToken tokens
  to improve usability of the parsing result.
  For a tag such "<html>" is parsed it now returns 2 tokens :
  TagOpenToken    "<html"
  TagOpenEndToken ">"
  For a tag such "</html>" is parsed it now returns 2 tokens :
  TagCloseToken     "</html"
  TagCloseEndToken  ">"
- fix test for removal of TagOpenStartToken and TagCloseStartToken

__Contributors__

- mh-cbon

Released by mh-cbon, Mon 17 Oct 2016 -
[see the diff](https://github.com/mh-cbon/html-parser-lexer/compare/0.0.1...1.0.0#diff)
______________

### 0.0.1

__Changes__

- init

__Contributors__

- mh-cbon

Released by mh-cbon, Sun 16 Oct 2016 -
[see the diff](https://github.com/mh-cbon/html-parser-lexer/compare/487ec1079a4708e9b6801cebc3d390c2e0e52e84...0.0.1#diff)
______________


