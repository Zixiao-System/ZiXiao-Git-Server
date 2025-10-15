#ifndef GIT_C_API_H
#define GIT_C_API_H

#ifdef __cplusplus
extern "C" {
#endif

// Repository operations
void* git_repository_new(const char* path);
void git_repository_free(void* repo);
int git_repository_init(void* repo, int bare);
int git_repository_exists(void* repo);
int git_repository_is_valid(void* repo);

// Reference operations
int git_repository_create_ref(void* repo, const char* refName, const char* sha);
char* git_repository_get_ref(void* repo, const char* refName);
char** git_repository_list_refs(void* repo, int* count);
int git_repository_delete_ref(void* repo, const char* refName);

// Branch operations
int git_repository_create_branch(void* repo, const char* branchName, const char* sha);
char** git_repository_list_branches(void* repo, int* count);
int git_repository_delete_branch(void* repo, const char* branchName);

// Pack operations
int git_repository_receive_pack(void* repo, const char* packData, int packLen);
char* git_repository_upload_pack(void* repo, const char** wants, int wantCount,
                                  const char** haves, int haveCount, int* outLen);

// Protocol operations
char* git_protocol_create_ref_advertisement(const char** refs, const char** shas,
                                            int refCount, const char* service, int* outLen);
char* git_protocol_pkt_line(const char* data);
char* git_protocol_flush_pkt();

// Utility functions
void git_free_string(char* str);
void git_free_string_array(char** arr, int count);

#ifdef __cplusplus
}
#endif

#endif // GIT_C_API_H
