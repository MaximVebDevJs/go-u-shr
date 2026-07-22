package http

// MaxRequestBodySize — верхняя граница тела входящего запроса (защита от OOM).
const MaxRequestBodySize = 1 << 20 // 1 MiB

// MaxGzipDecompressedSize — верхняя граница распакованного gzip-тела (защита от gzip bomb).
const MaxGzipDecompressedSize = 1 << 20 // 1 MiB
