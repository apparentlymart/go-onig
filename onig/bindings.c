#include <bindings.h>

int goonig_error_code_to_str(
    UChar *err_buf, int err_code, OnigErrorInfo *err_info)
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

int goonig_regex_capture_count(regex_t *reg)
{
    return onig_number_of_captures(reg);
}

typedef struct {
    goonig_name_table_entry *next;
    int count;
} goonig_regex_name_table_state;

int goonig_regex_name_table_cb(
    const UChar *start,
    const UChar *end,
    int num,
    int *groups,
    regex_t *reg,
    void *stateP)
{
    goonig_regex_name_table_state *state =
        (goonig_regex_name_table_state *)(stateP);

    for (int i = 0; i < num; i++) {
        state->count++;
        state->next->start = (UChar *)start;
        state->next->len = end - start;
        state->next->idx = groups[i];
        state->next++;
    }
}

int goonig_regex_name_table(regex_t *reg, goonig_name_table_entry *next)
{
    goonig_regex_name_table_state state;
    state.next = next;
    state.count = 0;
    onig_foreach_name(reg, goonig_regex_name_table_cb, &state);
    return state.count;
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
