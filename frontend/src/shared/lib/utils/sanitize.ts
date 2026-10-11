// 统一 HTML 清洗工具：所有 dangerouslySetInnerHTML 必须经此清洗，防止存储型 XSS。
import DOMPurify from "dompurify";

// 富文本允许的 HTML 白名单（与后台审核弹窗保持一致）
const ALLOWED_TAGS = [
  "b",
  "i",
  "em",
  "strong",
  "a",
  "p",
  "br",
  "ul",
  "ol",
  "li",
  "img",
  "h1",
  "h2",
  "h3",
  "h4",
  "blockquote",
  "code",
  "pre",
];

const ALLOWED_ATTR = ["href", "target", "src", "alt", "class", "id"];

export function sanitizeHtml(html: string): string {
  return DOMPurify.sanitize(html, { ALLOWED_TAGS, ALLOWED_ATTR });
}