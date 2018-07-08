#include "oniguruma.h"

// This file contains some helper wrappers around oniguruma APIs. Any global
// symbols defined here must be namespaced as "goonig".

int goonig_error_code_to_str(UChar *err_buf, int err_code, OnigErrorInfo *err_info);
int goonig_init_regex(
    regex_t *reg,
    const char *pattern,
    int pattern_len,
    OnigOptionType option,
    OnigSyntaxType *syntax,
    OnigErrorInfo *err_info);
void goonig_free_regex(regex_t *reg);
int goonig_regex_match(
    regex_t *reg,
    const char *str,
    int str_len,
    OnigRegion *region,
    OnigOptionType option);
int goonig_regex_search(
    regex_t *reg,
    const char *str,
    int str_len,
    int rev, // bool
    OnigRegion *region,
    OnigOptionType option);

void goonig_init_region(OnigRegion *reg);
void goonig_free_region(OnigRegion *reg);
int goonig_region_resize(OnigRegion *reg, int size);
