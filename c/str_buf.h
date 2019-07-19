#ifndef __STR_BUF_H__
#define __STR_BUF_H__

typedef struct {
  char *buf;
  int count, size;
} str_buf;

str_buf *str_buf_new(void);
void str_buf_add(str_buf *sb, char c);
char *str_buf_str(str_buf *sb);
void str_buf_free(str_buf *sb);

#endif // __STR_BUF_H__
