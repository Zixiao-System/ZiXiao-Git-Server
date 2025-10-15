#include "git_c_api.h"
#include "git_repository.h"
#include "git_protocol.h"
#include <cstring>
#include <cstdlib>

using namespace GitCore;

extern "C" {

void* git_repository_new(const char* path) {
    return new GitRepository(std::string(path));
}

void git_repository_free(void* repo) {
    delete static_cast<GitRepository*>(repo);
}

int git_repository_init(void* repo, int bare) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    return r->init(bare != 0) ? 1 : 0;
}

int git_repository_exists(void* repo) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    return r->exists() ? 1 : 0;
}

int git_repository_is_valid(void* repo) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    return r->isValid() ? 1 : 0;
}

int git_repository_create_ref(void* repo, const char* refName, const char* sha) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    return r->createRef(refName, sha) ? 1 : 0;
}

char* git_repository_get_ref(void* repo, const char* refName) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    std::string ref = r->getRef(refName);
    if (ref.empty()) {
        return nullptr;
    }
    char* result = (char*)malloc(ref.length() + 1);
    strcpy(result, ref.c_str());
    return result;
}

char** git_repository_list_refs(void* repo, int* count) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    std::vector<std::string> refs = r->listRefs();

    *count = refs.size();
    if (refs.empty()) {
        return nullptr;
    }

    char** result = (char**)malloc(refs.size() * sizeof(char*));
    for (size_t i = 0; i < refs.size(); i++) {
        result[i] = (char*)malloc(refs[i].length() + 1);
        strcpy(result[i], refs[i].c_str());
    }
    return result;
}

int git_repository_delete_ref(void* repo, const char* refName) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    return r->deleteRef(refName) ? 1 : 0;
}

int git_repository_create_branch(void* repo, const char* branchName, const char* sha) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    return r->createBranch(branchName, sha) ? 1 : 0;
}

char** git_repository_list_branches(void* repo, int* count) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    std::vector<std::string> branches = r->listBranches();

    *count = branches.size();
    if (branches.empty()) {
        return nullptr;
    }

    char** result = (char**)malloc(branches.size() * sizeof(char*));
    for (size_t i = 0; i < branches.size(); i++) {
        result[i] = (char*)malloc(branches[i].length() + 1);
        strcpy(result[i], branches[i].c_str());
    }
    return result;
}

int git_repository_delete_branch(void* repo, const char* branchName) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    return r->deleteBranch(branchName) ? 1 : 0;
}

int git_repository_receive_pack(void* repo, const char* packData, int packLen) {
    GitRepository* r = static_cast<GitRepository*>(repo);
    std::string data(packData, packLen);
    return r->receivePack(data) ? 1 : 0;
}

char* git_repository_upload_pack(void* repo, const char** wants, int wantCount,
                                  const char** haves, int haveCount, int* outLen) {
    GitRepository* r = static_cast<GitRepository*>(repo);

    std::vector<std::string> wantVec;
    for (int i = 0; i < wantCount; i++) {
        wantVec.push_back(wants[i]);
    }

    std::vector<std::string> haveVec;
    for (int i = 0; i < haveCount; i++) {
        haveVec.push_back(haves[i]);
    }

    std::string pack = r->uploadPack(wantVec, haveVec);
    *outLen = pack.length();

    char* result = (char*)malloc(pack.length());
    memcpy(result, pack.data(), pack.length());
    return result;
}

char* git_protocol_create_ref_advertisement(const char** refs, const char** shas,
                                            int refCount, const char* service, int* outLen) {
    std::vector<GitProtocol::RefAdvertisement> refAds;
    for (int i = 0; i < refCount; i++) {
        GitProtocol::RefAdvertisement ad;
        ad.sha = shas[i];
        ad.refName = refs[i];
        refAds.push_back(ad);
    }

    std::string adv = GitProtocol::createRefAdvertisement(refAds, service);
    *outLen = adv.length();

    char* result = (char*)malloc(adv.length());
    memcpy(result, adv.data(), adv.length());
    return result;
}

char* git_protocol_pkt_line(const char* data) {
    std::string pkt = GitProtocol::pktLine(data);
    char* result = (char*)malloc(pkt.length() + 1);
    strcpy(result, pkt.c_str());
    return result;
}

char* git_protocol_flush_pkt() {
    std::string pkt = GitProtocol::flushPkt();
    char* result = (char*)malloc(pkt.length() + 1);
    strcpy(result, pkt.c_str());
    return result;
}

void git_free_string(char* str) {
    free(str);
}

void git_free_string_array(char** arr, int count) {
    for (int i = 0; i < count; i++) {
        free(arr[i]);
    }
    free(arr);
}

} // extern "C"
