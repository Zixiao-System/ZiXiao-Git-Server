#include "git_object.h"
#include <openssl/sha.h>
#include <sstream>
#include <iomanip>
#include <cstring>

namespace GitCore {

// GitObject implementation
GitObject::GitObject(GitObjectType type, const std::string& data)
    : type(type), data(data) {
    sha = calculateSHA(serialize());
}

GitObject::~GitObject() {
}

GitObjectType GitObject::getType() const {
    return type;
}

std::string GitObject::getSHA() const {
    return sha;
}

std::string GitObject::getData() const {
    return data;
}

size_t GitObject::getSize() const {
    return data.size();
}

std::string GitObject::serialize() const {
    std::string typeStr;
    switch (type) {
        case GitObjectType::BLOB:   typeStr = "blob"; break;
        case GitObjectType::TREE:   typeStr = "tree"; break;
        case GitObjectType::COMMIT: typeStr = "commit"; break;
        case GitObjectType::TAG:    typeStr = "tag"; break;
    }

    std::ostringstream oss;
    oss << typeStr << " " << data.size() << '\0' << data;
    return oss.str();
}

std::string GitObject::calculateSHA(const std::string& content) {
    unsigned char hash[SHA_DIGEST_LENGTH];
    SHA1(reinterpret_cast<const unsigned char*>(content.c_str()),
         content.length(), hash);

    std::ostringstream oss;
    for (int i = 0; i < SHA_DIGEST_LENGTH; i++) {
        oss << std::hex << std::setw(2) << std::setfill('0')
            << static_cast<int>(hash[i]);
    }
    return oss.str();
}

// GitBlob implementation
GitBlob::GitBlob(const std::string& content)
    : GitObject(GitObjectType::BLOB, content) {
}

std::string GitBlob::getContent() const {
    return data;
}

// GitTreeEntry implementation
GitTreeEntry::GitTreeEntry(const std::string& mode, const std::string& name,
                          const std::string& sha)
    : mode(mode), name(name), sha(sha) {
}

// GitTree implementation
GitTree::GitTree(const std::vector<GitTreeEntry>& entries)
    : GitObject(GitObjectType::TREE, ""), entries(entries) {
    data = buildTreeData();
    sha = calculateSHA(serialize());
}

void GitTree::addEntry(const GitTreeEntry& entry) {
    entries.push_back(entry);
    data = buildTreeData();
    sha = calculateSHA(serialize());
}

std::vector<GitTreeEntry> GitTree::getEntries() const {
    return entries;
}

std::string GitTree::buildTreeData() const {
    std::ostringstream oss;
    for (const auto& entry : entries) {
        oss << entry.mode << " " << entry.name << '\0';

        // Convert hex SHA to binary (20 bytes)
        for (size_t i = 0; i < entry.sha.length(); i += 2) {
            std::string byteStr = entry.sha.substr(i, 2);
            unsigned char byte = static_cast<unsigned char>(
                std::stoi(byteStr, nullptr, 16));
            oss << byte;
        }
    }
    return oss.str();
}

// GitCommit implementation
GitCommit::GitCommit(const std::string& treeSHA,
                    const std::vector<std::string>& parentSHAs,
                    const std::string& author,
                    const std::string& committer,
                    const std::string& message)
    : GitObject(GitObjectType::COMMIT, ""),
      treeSHA(treeSHA),
      parentSHAs(parentSHAs),
      author(author),
      committer(committer),
      message(message) {
    data = buildCommitData();
    sha = calculateSHA(serialize());
}

std::string GitCommit::getTreeSHA() const {
    return treeSHA;
}

std::vector<std::string> GitCommit::getParentSHAs() const {
    return parentSHAs;
}

std::string GitCommit::getMessage() const {
    return message;
}

std::string GitCommit::buildCommitData() const {
    std::ostringstream oss;
    oss << "tree " << treeSHA << "\n";

    for (const auto& parent : parentSHAs) {
        oss << "parent " << parent << "\n";
    }

    oss << "author " << author << "\n";
    oss << "committer " << committer << "\n";
    oss << "\n" << message;

    return oss.str();
}

} // namespace GitCore
