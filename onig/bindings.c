#include <bindings.h>

int goonig_error_code_to_str(UChar *err_buf, int err_code, OnigErrorInfo *err_info)
{
    return onig_error_code_to_str(err_buf, err_code, err_info);
}

int goonig_init_regex(
    regex_t *reg,
    const char *pattern,
    int pattern_len,
    OnigOptionType option,
    OnigSyntaxType *syntax,
    OnigErrorInfo *err_info)
{
    return onig_new_without_alloc(
        reg,
        pattern,
        pattern + pattern_len,
        option,
        ONIG_ENCODING_UTF8,
        syntax,
        err_info);
}

void goonig_free_regex(regex_t *reg)
{
    onig_free_body(reg);
}

int goonig_regex_match(
    regex_t *reg,
    const char *str,
    int str_len,
    OnigRegion *region,
    OnigOptionType option)
{
    return onig_match(reg, str, str + str_len, str, region, option);
}

int goonig_regex_search(
    regex_t *reg,
    const char *str,
    int str_len,
    int rev,
    OnigRegion *region,
    OnigOptionType option)
{
    if (rev) {
        return onig_search(reg, str, str, str + str_len, str, region, option);
    }
    return onig_search(
        reg, str, str + str_len, str, str + str_len, region, option);
}

void goonig_init_region(OnigRegion *reg)
{
    onig_region_init(reg);
}

void goonig_free_region(OnigRegion *reg)
{
    onig_region_free(reg, 0);
}

int goonig_region_resize(OnigRegion *reg, int size)
{
    return onig_region_resize(reg, size);
}
