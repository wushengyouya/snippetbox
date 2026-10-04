// 在编译期将 HTML、CSS 和图片写入二进制，部署时不再需要同步 UI 目录。
package ui

import "embed"

// Files 包含应用的 HTML 模板和静态资源。
//
//go:embed "html" "static"
var Files embed.FS
