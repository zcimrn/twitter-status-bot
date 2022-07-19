package tools

import (
  "strings"
)

func Escape(text string) string {
  for _, char := range []string{"\\", "_", "*", "[", "]", "(", ")", "~", "`", ">", "#", "+", "-", "=", "|", "{", "}", ".", "!"} {
    text = strings.Replace(text, char, "\\" + char, -1)
  }
  return text
}

func EscapeCode(code string) string {
  for _, char := range []string{"\\", "`"} {
    code = strings.Replace(code, char, "\\" + char, -1)
  }
  return code
}

func EscapeLink(link string) string {
  for _, char := range []string{"\\", "(", ")"} {
    link = strings.Replace(link, char, "\\" + char, -1)
  }
  return link
}
